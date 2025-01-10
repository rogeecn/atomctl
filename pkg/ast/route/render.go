package route

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"git.ipao.vip/rogeecn/atomctl/pkg/utils/gomod"
	"github.com/Masterminds/sprig/v3"
	"github.com/iancoleman/strcase"
	"github.com/samber/lo"
)

//go:embed router.go.tpl
var routeTpl string

type RenderData struct {
	PackageName    string
	ProjectPackage string
	Imports        []string
	Controllers    []string
	Routes         map[string][]Router
}

type Router struct {
	Method     string
	Route      string
	Controller string
	Action     string
	Func       string
	Params     []string
}

func Render(path string, routes []RouteDefinition) error {
	routePath := filepath.Join(path, "routes.gen.go")

	tmpl, err := template.New("route").Funcs(sprig.FuncMap()).Parse(routeTpl)
	if err != nil {
		return err
	}

	renderData := RenderData{
		PackageName:    filepath.Base(path),
		ProjectPackage: gomod.GetModuleName(),
		Routes:         make(map[string][]Router),
	}

	// collect imports
	imports := []string{}
	controllers := []string{}
	for _, route := range routes {
		imports = append(imports, route.Imports...)
		controllers = append(controllers, fmt.Sprintf("%s *%s", strcase.ToLowerCamel(route.Name), route.Name))
		for _, action := range route.Actions {
			funcName := fmt.Sprintf("Func%d", len(action.Params))
			if action.HasData {
				funcName = "Data" + funcName
			}

			renderData.Routes[route.Name] = append(renderData.Routes[route.Name], Router{
				Method:     strcase.ToCamel(action.Method),
				Route:      action.Route,
				Controller: strcase.ToLowerCamel(route.Name),
				Action:     action.Name,
				Func:       funcName,
				Params: lo.FilterMap(action.Params, func(item ParamDefinition, _ int) (string, bool) {
					switch item.Position {
					case PositionQuery:
						return fmt.Sprintf(`Query%s[%s]("%s")`, isScalarType(item.Type), item.Type, item.Name), true
					case PositionHeader:
						return fmt.Sprintf(`Header[%s]("%s")`, item.Type, item.Name), true
					case PositionCookie:
						return fmt.Sprintf(`Cookie%s[%s]("%s")`, isScalarType(item.Type), item.Type, item.Name), true
					case PositionBody:
						return fmt.Sprintf(`Body[%s]("%s")`, item.Type, item.Name), true
					case PositionPath:
						return fmt.Sprintf(`Path%s[%s]("%s")`, isScalarType(item.Type), item.Type, item.Name), true
					case PositionLocal:
						key := item.Name
						if item.Key != "" {
							key = item.Key
						}
						return fmt.Sprintf(`Local[%s]("%s")`, item.Type, key), true
					}
					return "", false
				}),
			})
		}
	}

	renderData.Imports = lo.Uniq(imports)
	renderData.Controllers = lo.Uniq(controllers)

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, renderData)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(routePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(buf.Bytes())
	if err != nil {
		return err
	}
	return nil
}

func isScalarType(t string) string {
	switch t {
	case "string", "int", "int32", "int64", "float32", "float64", "bool":
		return "Param"
	}
	return ""
}
