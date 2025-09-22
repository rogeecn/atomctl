package route

import (
	"fmt"
	"sort"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/samber/lo"
)

type RenderBuildOpts struct {
	PackageName    string
	ProjectPackage string
	Routes         []RouteDefinition
}

func buildRenderData(opts RenderBuildOpts) (RenderData, error) {
	builder := &renderDataBuilder{
		opts: opts,
		data: RenderData{
			PackageName:    opts.PackageName,
			ProjectPackage: opts.ProjectPackage,
			Imports:        []string{},
			Controllers:    []string{},
			Routes:         make(map[string][]Router),
			RouteGroups:    []string{},
		},
		imports:          []string{},
		controllers:      []string{},
		needsFieldImport: false,
	}

	return builder.build()
}

type renderDataBuilder struct {
	opts             RenderBuildOpts
	data             RenderData
	imports          []string
	controllers      []string
	needsFieldImport bool
}

func (b *renderDataBuilder) build() (RenderData, error) {
	b.processRoutes()
	b.addRequiredImports()
	b.dedupeAndSortImports()
	b.dedupeAndSortControllers()
	b.sortRouteGroups()

	return b.data, nil
}

func (b *renderDataBuilder) processRoutes() {
	for _, route := range b.opts.Routes {
		b.collectRouteMetadata(route)
		b.buildRouteActions(route)
	}
}

func (b *renderDataBuilder) collectRouteMetadata(route RouteDefinition) {
	b.imports = append(b.imports, route.Imports...)
	b.controllers = append(b.controllers, fmt.Sprintf("%s *%s", strcase.ToLowerCamel(route.Name), route.Name))
}

func (b *renderDataBuilder) buildRouteActions(route RouteDefinition) {
	for _, action := range route.Actions {
		router := b.buildRouter(route, action)
		b.data.Routes[route.Name] = append(b.data.Routes[route.Name], router)
	}
}

func (b *renderDataBuilder) buildRouter(route RouteDefinition, action ActionDefinition) Router {
	funcName := b.generateFunctionName(action)
	params := b.buildParameters(action.Params)

	return Router{
		Method:     strcase.ToCamel(action.Method),
		Route:      action.Route,
		Controller: strcase.ToLowerCamel(route.Name),
		Action:     action.Name,
		Func:       funcName,
		Params:     params,
	}
}

func (b *renderDataBuilder) generateFunctionName(action ActionDefinition) string {
	funcName := fmt.Sprintf("Func%d", len(action.Params))
	if action.HasData {
		funcName = "Data" + funcName
	}
	return funcName
}

func (b *renderDataBuilder) buildParameters(params []ParamDefinition) []string {
	return lo.FilterMap(params, func(item ParamDefinition, _ int) (string, bool) {
		token := buildParamToken(item)
		if token == "" {
			return "", false
		}
		if item.Model != "" {
			b.needsFieldImport = true
		}
		return token, true
	})
}

func (b *renderDataBuilder) addRequiredImports() {
	if b.needsFieldImport {
		b.imports = append(b.imports, `field "go.ipao.vip/gen/field"`)
	}
}

func (b *renderDataBuilder) dedupeAndSortImports() {
	b.data.Imports = lo.Uniq(b.imports)
	sort.Strings(b.data.Imports)
}

func (b *renderDataBuilder) dedupeAndSortControllers() {
	b.data.Controllers = lo.Uniq(b.controllers)
	sort.Strings(b.data.Controllers)
}

func (b *renderDataBuilder) sortRouteGroups() {
	// Collect route groups
	for k := range b.data.Routes {
		b.data.RouteGroups = append(b.data.RouteGroups, k)
	}
	sort.Strings(b.data.RouteGroups)

	// Sort routes within each group
	for _, groupName := range b.data.RouteGroups {
		items := b.data.Routes[groupName]
		b.sortRouteItems(items)
		b.data.Routes[groupName] = items
	}
}

func (b *renderDataBuilder) sortRouteItems(items []Router) {
	sort.Slice(items, func(i, j int) bool {
		if items[i].Method != items[j].Method {
			return items[i].Method < items[j].Method
		}
		if items[i].Route != items[j].Route {
			return items[i].Route < items[j].Route
		}
		return items[i].Action < items[j].Action
	})
}

func buildParamToken(item ParamDefinition) string {
	key := item.getKey()
	builder := &paramTokenBuilder{item: item, key: key}
	return builder.build()
}

func (item ParamDefinition) getKey() string {
	if item.Key != "" {
		return item.Key
	}
	return item.Name
}

type paramTokenBuilder struct {
	item ParamDefinition
	key  string
}

func (b *paramTokenBuilder) build() string {
	switch b.item.Position {
	case PositionQuery:
		return b.buildQueryParam()
	case PositionHeader:
		return b.buildHeaderParam()
	case PositionFile:
		return b.buildFileParam()
	case PositionCookie:
		return b.buildCookieParam()
	case PositionBody:
		return b.buildBodyParam()
	case PositionPath:
		return b.buildPathParam()
	case PositionLocal:
		return b.buildLocalParam()
	default:
		return ""
	}
}

func (b *paramTokenBuilder) buildQueryParam() string {
	return fmt.Sprintf(`Query%s[%s]("%s")`, scalarSuffix(b.item.Type), b.item.Type, b.key)
}

func (b *paramTokenBuilder) buildHeaderParam() string {
	return fmt.Sprintf(`Header[%s]("%s")`, b.item.Type, b.key)
}

func (b *paramTokenBuilder) buildFileParam() string {
	return fmt.Sprintf(`File[multipart.FileHeader]("%s")`, b.key)
}

func (b *paramTokenBuilder) buildCookieParam() string {
	if b.item.Type == "string" {
		return fmt.Sprintf(`CookieParam("%s")`, b.key)
	}
	return fmt.Sprintf(`Cookie[%s]("%s")`, b.item.Type, b.key)
}

func (b *paramTokenBuilder) buildBodyParam() string {
	return fmt.Sprintf(`Body[%s]("%s")`, b.item.Type, b.key)
}

func (b *paramTokenBuilder) buildPathParam() string {
	if b.item.Model != "" {
		return b.buildModelLookupPath()
	}
	return fmt.Sprintf(`Path%s[%s]("%s")`, scalarSuffix(b.item.Type), b.item.Type, b.key)
}

func (b *paramTokenBuilder) buildModelLookupPath() string {
	field, fieldType := b.parseModelField()

	tpl := `func(ctx fiber.Ctx) (*%s, error) {
		v := fiber.Params[%s](ctx, "%s")
		return %sQuery.WithContext(ctx).Where(field.NewUnsafeFieldRaw("%s = ?", v)).First()
	}`
	return fmt.Sprintf(tpl, b.item.Type, fieldType, b.key, b.item.Type, field)
}

func (b *paramTokenBuilder) parseModelField() (string, string) {
	field := "id"
	fieldType := "int"

	if strings.Contains(b.item.Model, ":") {
		parts := strings.SplitN(b.item.Model, ":", 2)
		if len(parts) == 2 {
			field = parts[0]
			fieldType = parts[1]
		}
	} else {
		field = b.item.Model
	}

	return field, fieldType
}

func (b *paramTokenBuilder) buildLocalParam() string {
	return fmt.Sprintf(`Local[%s]("%s")`, b.item.Type, b.key)
}

func scalarSuffix(t string) string {
	switch t {
	case "string", "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"float32", "float64", "bool":
		return "Param"
	}
	return ""
}
