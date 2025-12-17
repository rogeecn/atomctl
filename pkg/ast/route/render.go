package route

import (
	_ "embed"
	"os"
	"path/filepath"

	"go.ipao.vip/atomctl/v2/pkg/utils/gomod"
)

//go:embed router.go.tpl
var routeTpl string

//go:embed manual.go.tpl
var routeManualTpl string

type RenderData struct {
	PackageName    string
	ModuleName     string
	ProjectPackage string
	Imports        []string
	Controllers    []string
	Routes         map[string][]Router
	RouteGroups    []string
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
	routeManualPath := filepath.Join(path, "routes.manual.go")

	data, err := buildRenderData(RenderBuildOpts{
		PackageName:    filepath.Base(path),
		ModuleName:     gomod.GetModuleName(),
		ProjectPackage: gomod.GetModuleName(),
		Routes:         routes,
	})
	if err != nil {
		return err
	}

	out, err := renderTemplate(data)
	if err != nil {
		return err
	}

	if err := os.WriteFile(routePath, out, 0o644); err != nil {
		return err
	}
	// if routes.manual.go not exists then create it
	if _, err := os.Stat(routeManualPath); os.IsNotExist(err) {
		manualOut, err := renderManualTemplate(data)
		if err != nil {
			return err
		}
		if err := os.WriteFile(routeManualPath, manualOut, 0o644); err != nil {
			return err
		}
	}
	return nil
}
