package route

import (
	"fmt"
	"sort"

	"github.com/iancoleman/strcase"
	"github.com/samber/lo"
)

type RenderBuildOpts struct {
	PackageName    string
	ProjectPackage string
	Routes         []RouteDefinition
}

func buildRenderData(opts RenderBuildOpts) (RenderData, error) {
	rd := RenderData{
		PackageName:    opts.PackageName,
		ProjectPackage: opts.ProjectPackage,
		Imports:        []string{},
		Controllers:    []string{},
		Routes:         make(map[string][]Router),
		RouteGroups:    []string{},
	}

	imports := []string{}
	controllers := []string{}

	for _, route := range opts.Routes {
		imports = append(imports, route.Imports...)
		controllers = append(controllers, fmt.Sprintf("%s *%s", strcase.ToLowerCamel(route.Name), route.Name))

		for _, action := range route.Actions {
			funcName := fmt.Sprintf("Func%d", len(action.Params))
			if action.HasData {
				funcName = "Data" + funcName
			}

			params := lo.FilterMap(action.Params, func(item ParamDefinition, _ int) (string, bool) {
				tok := buildParamToken(item)
				if tok == "" {
					return "", false
				}
				return tok, true
			})

			rd.Routes[route.Name] = append(rd.Routes[route.Name], Router{
				Method:     strcase.ToCamel(action.Method),
				Route:      action.Route,
				Controller: strcase.ToLowerCamel(route.Name),
				Action:     action.Name,
				Func:       funcName,
				Params:     params,
			})
		}
	}

	// de-dup and sort imports/controllers for stable output
	rd.Imports = lo.Uniq(imports)
	sort.Strings(rd.Imports)
	rd.Controllers = lo.Uniq(controllers)
	sort.Strings(rd.Controllers)

	// stable order for route groups and entries
	for k := range rd.Routes {
		rd.RouteGroups = append(rd.RouteGroups, k)
	}
	sort.Strings(rd.RouteGroups)
	for _, k := range rd.RouteGroups {
		items := rd.Routes[k]
		sort.Slice(items, func(i, j int) bool {
			if items[i].Method != items[j].Method {
				return items[i].Method < items[j].Method
			}
			if items[i].Route != items[j].Route {
				return items[i].Route < items[j].Route
			}
			return items[i].Action < items[j].Action
		})
		rd.Routes[k] = items
	}

	return rd, nil
}

func buildParamToken(item ParamDefinition) string {
	key := item.Name
	if item.Key != "" {
		key = item.Key
	}

	switch item.Position {
	case PositionQuery:
		return fmt.Sprintf(`Query%s[%s]("%s")`, scalarSuffix(item.Type), item.Type, key)
	case PositionHeader:
		return fmt.Sprintf(`Header[%s]("%s")`, item.Type, key)
	case PositionFile:
		return fmt.Sprintf(`File[multipart.FileHeader]("%s")`, key)
	case PositionCookie:
		if item.Type == "string" {
			return fmt.Sprintf(`CookieParam("%s")`, key)
		}
		return fmt.Sprintf(`Cookie[%s]("%s")`, item.Type, key)
	case PositionBody:
		return fmt.Sprintf(`Body[%s]("%s")`, item.Type, key)
	case PositionPath:
		return fmt.Sprintf(`Path%s[%s]("%s")`, scalarSuffix(item.Type), item.Type, key)
	case PositionLocal:
		return fmt.Sprintf(`Local[%s]("%s")`, item.Type, key)
	}
	return ""
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
