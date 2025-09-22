package route

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
)

type RouteDefinition struct {
	FilePath string
	Path     string
	Name     string
	Imports  []string
	Actions  []ActionDefinition
}

type ActionDefinition struct {
	Route   string
	Method  string
	Name    string
	HasData bool
	Params  []ParamDefinition
}

type ParamDefinition struct {
	Name     string
	Type     string
	Key      string
	Model    string
	Position Position
}

type Position string

func positionFromString(v string) Position {
	switch v {
	case "path":
		return PositionPath
	case "query":
		return PositionQuery
	case "body":
		return PositionBody
	case "header":
		return PositionHeader
	case "cookie":
		return PositionCookie
	case "local":
		return PositionLocal
	case "file":
		return PositionFile
	}
	panic("invalid position: " + v)
}

const (
	PositionPath   Position = "path"
	PositionQuery  Position = "query"
	PositionBody   Position = "body"
	PositionHeader Position = "header"
	PositionCookie Position = "cookie"
	PositionLocal  Position = "local"
	PositionFile   Position = "file"
)

func ParseFile(file string) []RouteDefinition {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		log.WithError(err).Error("Failed to parse file")
		return nil
	}

	imports := extractImports(node)

	parser := &routeParser{
		file:        file,
		node:        node,
		imports:     imports,
		routes:      make(map[string]RouteDefinition),
		actions:     make(map[string][]ActionDefinition),
		usedImports: make(map[string][]string),
	}

	return parser.parse()
}

type routeParser struct {
	file        string
	node        *ast.File
	imports     map[string]string
	routes      map[string]RouteDefinition
	actions     map[string][]ActionDefinition
	usedImports map[string][]string
}

func (p *routeParser) parse() []RouteDefinition {
	p.parseFunctionDeclarations()
	return p.buildResult()
}

func (p *routeParser) parseFunctionDeclarations() {
	for _, decl := range p.node.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !p.isValidFunctionDeclaration(funcDecl, ok) {
			continue
		}

		recvType := p.extractReceiverType(funcDecl)
		p.initializeRoute(recvType)

		routeInfo := p.extractRouteInfo(funcDecl)
		if routeInfo == nil {
			continue
		}

		action := p.buildAction(funcDecl, routeInfo)
		p.actions[recvType] = append(p.actions[recvType], action)
	}
}

func (p *routeParser) isValidFunctionDeclaration(decl *ast.FuncDecl, ok bool) bool {
	return ok &&
		decl.Recv != nil &&
		decl.Doc != nil
}

func (p *routeParser) extractReceiverType(decl *ast.FuncDecl) string {
	return decl.Recv.List[0].Type.(*ast.StarExpr).X.(*ast.Ident).Name
}

func (p *routeParser) initializeRoute(recvType string) {
	if _, exists := p.routes[recvType]; !exists {
		p.routes[recvType] = RouteDefinition{
			Name:     recvType,
			FilePath: p.file,
			Actions:  []ActionDefinition{},
		}
		p.actions[recvType] = []ActionDefinition{}
	}
}

type routeInfo struct {
	path       string
	method     string
	bindParams []ParamDefinition
}

func (p *routeParser) extractRouteInfo(decl *ast.FuncDecl) *routeInfo {
	var info routeInfo
	var err error

	for _, comment := range decl.Doc.List {
		line := normalizeCommentLine(comment.Text)

		if strings.HasPrefix(line, "@Router") {
			info.path, info.method, err = parseRouteComment(line)
			if err != nil {
				log.WithError(err).WithFields(log.Fields{
					"file":   p.file,
					"action": decl.Name.Name,
				}).Error("Invalid route definition")
				return nil
			}
		}

		if strings.HasPrefix(line, "@Bind") {
			info.bindParams = append(info.bindParams, parseRouteBind(line))
		}
	}

	if info.path == "" || info.method == "" {
		return nil
	}

	log.WithFields(log.Fields{
		"file":   p.file,
		"action": decl.Name.Name,
		"path":   info.path,
		"method": info.method,
	}).Info("Found route")

	return &info
}

func (p *routeParser) buildAction(decl *ast.FuncDecl, routeInfo *routeInfo) ActionDefinition {
	orderBindParams := p.processFunctionParameters(decl, routeInfo.bindParams)

	hasData := false
	if decl.Type != nil && decl.Type.Results != nil {
		hasData = len(decl.Type.Results.List) > 1
	}

	return ActionDefinition{
		Route:   routeInfo.path,
		Method:  strings.ToUpper(routeInfo.method),
		Name:    decl.Name.Name,
		HasData: hasData,
		Params:  orderBindParams,
	}
}

func (p *routeParser) processFunctionParameters(decl *ast.FuncDecl, bindParams []ParamDefinition) []ParamDefinition {
	var orderBindParams []ParamDefinition

	if decl.Type == nil || decl.Type.Params == nil {
		return orderBindParams
	}

	for _, param := range decl.Type.Params.List {
		paramType := extractParameterType(param.Type)

		if isContextParameter(paramType) {
			continue
		}

		p.trackUsedImports(decl.Recv.List[0].Type.(*ast.StarExpr).X.(*ast.Ident).Name, paramType)

		for _, paramName := range param.Names {
			for i, bindParam := range bindParams {
				if bindParam.Name == paramName.Name {
					bindParams[i].Type = p.normalizeParameterType(paramType, bindParam.Position)
					orderBindParams = append(orderBindParams, bindParams[i])
					break
				}
			}
		}
	}

	return orderBindParams
}

func (p *routeParser) trackUsedImports(recvType, paramType string) {
	pkgParts := strings.Split(strings.Trim(paramType, "*"), ".")
	if len(pkgParts) == 2 {
		if importPath, exists := p.imports[pkgParts[0]]; exists {
			p.usedImports[recvType] = append(p.usedImports[recvType], importPath)
		}
	}
}

func (p *routeParser) normalizeParameterType(paramType string, position Position) string {
	if position != PositionLocal {
		return strings.TrimPrefix(paramType, "*")
	}
	return paramType
}

func (p *routeParser) buildResult() []RouteDefinition {
	var items []RouteDefinition
	for k, route := range p.routes {
		if actions, exists := p.actions[k]; exists {
			route.Actions = actions
			route.Imports = p.getUniqueImports(k)

			// Set the route path from the first action for backward compatibility
			if len(actions) > 0 {
				route.Path = actions[0].Route
			}

			items = append(items, route)
		}
	}
	return items
}

func (p *routeParser) getUniqueImports(recvType string) []string {
	if imports, exists := p.usedImports[recvType]; exists {
		return lo.Uniq(imports)
	}
	return []string{}
}

func extractImports(node *ast.File) map[string]string {
	imports := make(map[string]string)
	for _, imp := range node.Imports {
		pkg := strings.Trim(imp.Path.Value, "\"")

		// Handle empty or invalid package paths
		if pkg == "" {
			continue
		}

		// Use the last part of the package path as the name
		// This avoids calling gomod.GetPackageModuleName which can cause panics
		name := pkg
		if lastSlash := strings.LastIndex(pkg, "/"); lastSlash >= 0 {
			name = pkg[lastSlash+1:]
		}

		if imp.Name != nil {
			name = imp.Name.Name
			pkg = fmt.Sprintf(`%s %q`, name, pkg)
			imports[name] = pkg
			continue
		}
		imports[name] = fmt.Sprintf("%q", pkg)
	}
	return imports
}

func normalizeCommentLine(line string) string {
	return strings.TrimSpace(strings.TrimLeft(line, "/ \t"))
}

func extractParameterType(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + extractParameterType(t.X)
	case *ast.SelectorExpr:
		return fmt.Sprintf("%s.%s", extractParameterType(t.X), t.Sel.Name)
	default:
		return ""
	}
}

func isContextParameter(paramType string) bool {
	return strings.HasSuffix(paramType, "Context") || strings.HasSuffix(paramType, "Ctx")
}

func parseRouteComment(line string) (string, string, error) {
	parts := strings.FieldsFunc(line, func(r rune) bool {
		return r == ' ' || r == '\t' || r == '[' || r == ']'
	})
	parts = lo.Filter(parts, func(item string, idx int) bool {
		return item != ""
	})

	if len(parts) != 3 {
		return "", "", errors.New("invalid route definition")
	}

	return parts[1], parts[2], nil
}

func parseRouteBind(bind string) ParamDefinition {
	var param ParamDefinition
	parts := strings.FieldsFunc(bind, func(r rune) bool {
		return r == ' ' || r == '(' || r == ')' || r == '\t'
	})
	parts = lo.Filter(parts, func(item string, idx int) bool {
		return item != ""
	})

	for i, part := range parts {
		switch part {
		case "@Bind":
			if i+2 < len(parts) {
				param.Name = parts[i+1]
				param.Position = positionFromString(parts[i+2])
			}
		case "key":
			if i+1 < len(parts) {
				param.Key = parts[i+1]
			}
		case "model":
			if i+1 < len(parts) {
				// Supported formats:
				// - model(field:field_type) -> only specify model field/column;
				mv := parts[i+1]
				// if mv contains no dot, treat as field name directly
				if mv == "" {
					param.Model = "id"
					break
				}
				param.Model = mv
			}
		}
	}
	return param
}
