package route

import (
	_ "embed"
	"os"
	"path/filepath"

	"go.ipao.vip/atomctl/v2/pkg/utils/gomod"
)


//go:embed router.go.tpl
var routeTpl string

type RenderData struct {
	PackageName    string
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
	// Validate input parameters
	if err := validateRenderInput(path, routes); err != nil {
		return err
	}

	renderer := &routeRenderer{
		path:   path,
		routes: routes,
	}

	return renderer.render()
}

type routeRenderer struct {
	path   string
	routes []RouteDefinition
}

func (r *routeRenderer) render() error {
	// Prepare render data
	data, err := r.prepareRenderData()
	if err != nil {
		return err
	}

	// Generate content
	content, err := r.generateContent(data)
	if err != nil {
		return err
	}

	// Write to file atomically
	return r.writeFileAtomically(content)
}

func (r *routeRenderer) prepareRenderData() (RenderData, error) {
	data, err := buildRenderData(RenderBuildOpts{
		PackageName:    filepath.Base(r.path),
		ProjectPackage: getProjectPackageWithFallback(),
		Routes:         r.routes,
	})
	if err != nil {
		return RenderData{}, WrapError(err, "failed to build render data for path: %s", r.path)
	}

	// Validate the generated data
	if err := r.validateRenderData(data); err != nil {
		return RenderData{}, err
	}

	return data, nil
}

func (r *routeRenderer) validateRenderData(data RenderData) error {
	if data.PackageName == "" {
		return NewRouteError(ErrInvalidInput, "package name cannot be empty")
	}

	if len(data.Routes) == 0 {
		return NewRouteError(ErrNoRoutes, "no routes to render")
	}

	// Validate that all routes have required fields
	for controllerName, routes := range data.Routes {
		if controllerName == "" {
			return NewRouteError(ErrInvalidInput, "controller name cannot be empty")
		}

		for i, route := range routes {
			if route.Method == "" {
				return NewRouteError(ErrInvalidInput, "route method cannot be empty for controller %s, route %d", controllerName, i)
			}
			if route.Route == "" {
				return NewRouteError(ErrInvalidInput, "route path cannot be empty for controller %s, route %d", controllerName, i)
			}
		}
	}

	return nil
}

func (r *routeRenderer) generateContent(data RenderData) ([]byte, error) {
	content, err := renderTemplate(data)
	if err != nil {
		return nil, WrapError(err, "failed to render template for path: %s", r.path)
	}

	// Validate generated content is not empty
	if len(content) == 0 {
		return nil, NewRouteError(ErrTemplateFailed, "generated content is empty")
	}

	return content, nil
}

func (r *routeRenderer) writeFileAtomically(content []byte) error {
	routePath := filepath.Join(r.path, "routes.gen.go")

	// Write to temporary file first for atomic operation
	tempPath := routePath + ".tmp"
	if err := os.WriteFile(tempPath, content, 0o644); err != nil {
		return WrapError(err, "failed to write temporary route file: %s", tempPath)
	}

	// Rename temporary file to final destination (atomic operation)
	if err := os.Rename(tempPath, routePath); err != nil {
		// Clean up temporary file if rename fails
		_ = os.Remove(tempPath)
		return WrapError(err, "failed to rename temporary file to final destination: %s -> %s", tempPath, routePath)
	}

	return nil
}

func validateRenderInput(path string, routes []RouteDefinition) error {
	if path == "" {
		return NewRouteError(ErrInvalidInput, "path cannot be empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return NewRouteError(ErrInvalidPath, "directory does not exist: %s", path)
	}

	if len(routes) == 0 {
		// This is not necessarily an error, but worth noting
		return NewRouteError(ErrNoRoutes, "no routes provided for rendering")
	}

	return nil
}

func getProjectPackageWithFallback() string {
	if moduleName := gomod.GetModuleName(); moduleName != "" {
		return moduleName
	}
	return "unknown" // fallback to prevent crashes
}
