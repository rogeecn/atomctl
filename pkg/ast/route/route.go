package route

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"regexp"
	"strings"

	"git.ipao.vip/rogeecn/atomctl/pkg/utils/gomod"
	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
)

type RouteDefinition struct {
	Path    string
	Name    string
	Imports []string
	Actions []ActionDefinition
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
	Table    string
	Model    string
	Position Position
}

type Position string

func positionFromString(v string) Position {
	switch v {
	case "path":
		return PositionPath
	case "uri":
		return PositionURI
	case "query":
		return PositionQuery
	case "body":
		return PositionBody
	case "header":
		return PositionHeader
	case "cookie":
		return PositionCookie
	}
	panic("invalid position: " + v)
}

const (
	PositionPath   Position = "path"
	PositionURI    Position = "uri"
	PositionQuery  Position = "query"
	PositionBody   Position = "body"
	PositionHeader Position = "header"
	PositionCookie Position = "cookie"
)

func ParseFile(file string) []RouteDefinition {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		log.Println("ERR: ", err)
		return nil
	}

	imports := make(map[string]string)
	for _, imp := range node.Imports {
		pkg := strings.Trim(imp.Path.Value, "\"")
		name := gomod.GetPackageModuleName(pkg)
		if imp.Name != nil {
			name = imp.Name.Name
			pkg = fmt.Sprintf(`%s %q`, name, pkg)
			imports[name] = pkg
			continue
		}
		imports[name] = fmt.Sprintf("%q", pkg)
	}

	routes := make(map[string]RouteDefinition)
	actions := make(map[string][]ActionDefinition)
	usedImports := make(map[string][]string)

	// 再去遍历 struct 的方法去
	for _, decl := range node.Decls {
		decl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		// 普通方法不要
		if decl.Recv == nil {
			continue
		}

		// 没有Doc不要
		if decl.Doc == nil {
			continue
		}

		recvType := decl.Recv.List[0].Type.(*ast.StarExpr).X.(*ast.Ident).Name
		if _, ok := routes[recvType]; !ok {
			routes[recvType] = RouteDefinition{
				Name:    recvType,
				Path:    file,
				Actions: []ActionDefinition{},
			}
			actions[recvType] = []ActionDefinition{}
		}

		bindParams := []ParamDefinition{}

		// Doc 中把 @Router 的定义拿出来， Route 格式为 /user/:id [get] 两部分，表示路径和请求方法
		var path, method string
		var err error
		for _, l := range decl.Doc.List {
			line := strings.TrimLeft(l.Text, "/ \t")
			line = strings.TrimSpace(line)

			// 路由需要一些切换
			if strings.HasPrefix(line, "@Router") {
				path, method, err = parseRouteComment(line)
				if err != nil {
					log.Fatal(errors.Wrapf(err, "file: %s, action: %s", file, decl.Name.Name))
				}
			}

			if strings.HasPrefix(line, "@Bind") {
				//@Bind name query key() table() model()
				//@Bind name query
				bindParams = append(bindParams, parseRouteBind(line))
			}
		}

		if path == "" || method == "" {
			continue
		}
		log.WithField("file", file).WithField("action", decl.Name.Name).WithField("path", path).WithField("method", method).Info("get router")

		// 拿参数列表去, 忽略 context.Context 参数
		for _, param := range decl.Type.Params.List {
			// paramsType, ok := param.Type.(*ast.SelectorExpr)

			var typ string
			switch param.Type.(type) {
			case *ast.Ident:
				typ = param.Type.(*ast.Ident).Name
			case *ast.StarExpr:
				paramsType := param.Type.(*ast.StarExpr)
				switch paramsType.X.(type) {
				case *ast.SelectorExpr:
					X := paramsType.X.(*ast.SelectorExpr)
					typ = fmt.Sprintf("*%s.%s", X.X.(*ast.Ident).Name, X.Sel.Name)
				default:
					typ = fmt.Sprintf("*%s", paramsType.X.(*ast.Ident).Name)
				}
			case *ast.SelectorExpr:
				typ = fmt.Sprintf("%s.%s", param.Type.(*ast.SelectorExpr).X.(*ast.Ident).Name, param.Type.(*ast.SelectorExpr).Sel.Name)
			}

			if strings.HasSuffix(typ, "Context") || strings.HasSuffix(typ, "Ctx") {
				continue
			}
			pkgName := strings.Split(strings.Trim(typ, "*"), ".")
			if len(pkgName) == 2 {
				usedImports[recvType] = append(usedImports[recvType], imports[pkgName[0]])
			}

			typ = strings.TrimPrefix(typ, "*")

			for _, name := range param.Names {
				for i, bindParam := range bindParams {
					if bindParam.Name == name.Name {
						bindParams[i].Type = typ
						break
					}
				}
			}
		}

		actions[recvType] = append(actions[recvType], ActionDefinition{
			Route:   path,
			Method:  strings.ToUpper(method),
			Name:    decl.Name.Name,
			HasData: len(decl.Type.Results.List) > 1,
			Params:  bindParams,
		})
	}

	var items []RouteDefinition
	for k, item := range routes {
		a, ok := actions[k]
		if !ok {
			continue
		}
		item.Actions = a
		item.Imports = []string{}
		if im, ok := usedImports[k]; ok {
			item.Imports = lo.Uniq(im)
		}
		items = append(items, item)
	}
	return items
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

func getPackageRoute(mod, path string) string {
	paths := strings.SplitN(path, "modules", 2)
	pkg := paths[1]
	// path可能值为
	// /test/user_controller.go
	// /test/modules/user_controller.go

	return strings.TrimLeft(filepath.Dir(pkg), "/")
}

func formatRoute(route string) string {
	pattern := regexp.MustCompile(`(?mi)\{(.*?)\}`)
	if !pattern.MatchString(route) {
		return route
	}

	items := pattern.FindAllStringSubmatch(route, -1)
	for _, item := range items {
		param := strcase.ToLowerCamel(item[1])
		route = strings.ReplaceAll(route, item[0], fmt.Sprintf("{%s}", param))
	}

	route = pattern.ReplaceAllString(route, ":$1")
	route = strings.ReplaceAll(route, "/:id", "/:id<int>")
	route = strings.ReplaceAll(route, "Id/", "Id<int>/")
	return route
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
			param.Name = parts[i+1]
			param.Position = positionFromString(parts[i+2])
		case "key":
			param.Key = parts[i+1]
		case "table":
			param.Table = parts[i+1]
		case "model":
			param.Model = parts[i+1]
		}
	}
	return param
}
