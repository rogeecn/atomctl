package route

import (
	"bytes"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	log "github.com/sirupsen/logrus"
)

// TemplateRenderer defines the interface for template rendering operations
type TemplateRenderer interface {
	Render(data RenderData) ([]byte, error)
	Validate() error
	GetTemplateInfo() TemplateInfo
}

// TemplateInfo provides metadata about the template
type TemplateInfo struct {
	Name       string
	Version    string
	Functions  []string
	Options    []string
	Size       int
}

// RouteRenderer implements TemplateRenderer for route generation
type RouteRenderer struct {
	template *template.Template
	info     TemplateInfo
	logger   *log.Entry
}

// NewRouteRenderer creates a new RouteRenderer instance with proper initialization
func NewRouteRenderer() *RouteRenderer {
	renderer := &RouteRenderer{
		logger: log.WithField("module", "route-renderer"),
		info: TemplateInfo{
			Name:    "router",
			Version: "1.0.0",
			Functions: []string{
				"sprig",
				"template",
				"custom",
			},
			Options: []string{
				"missingkey=error",
			},
		},
	}

	// Initialize template with error handling
	if err := renderer.initializeTemplate(); err != nil {
		renderer.logger.WithError(err).Error("Failed to initialize template")
		return nil
	}

	renderer.info.Size = len(routeTpl)
	renderer.logger.WithFields(log.Fields{
		"template_size": renderer.info.Size,
		"version":       renderer.info.Version,
	}).Info("Route renderer initialized successfully")

	return renderer
}

// initializeTemplate sets up the template with proper functions and options
func (r *RouteRenderer) initializeTemplate() error {
	// Create template with sprig functions and custom options
	tmpl := template.New(r.info.Name).
		Funcs(sprig.FuncMap()).
		Option("missingkey=error")

	// Parse the template
	parsedTmpl, err := tmpl.Parse(routeTpl)
	if err != nil {
		return WrapError(err, "failed to parse route template")
	}

	r.template = parsedTmpl
	return nil
}

// Render renders the template with the provided data
func (r *RouteRenderer) Render(data RenderData) ([]byte, error) {
	// Validate input data
	if err := r.validateRenderData(data); err != nil {
		return nil, err
	}

	// Create buffer for rendering
	var buf bytes.Buffer
	buf.Grow(estimatedBufferSize(data)) // Pre-allocate buffer for better performance

	// Execute template with error handling
	if err := r.template.Execute(&buf, data); err != nil {
		r.logger.WithError(err).WithFields(log.Fields{
			"package_name": data.PackageName,
			"routes_count": len(data.Routes),
		}).Error("Template execution failed")
		return nil, WrapError(err, "template execution failed for package: %s", data.PackageName)
	}

	// Validate rendered content
	result := buf.Bytes()
	if len(result) == 0 {
		return nil, NewRouteError(ErrTemplateFailed, "rendered content is empty for package: %s", data.PackageName)
	}

	r.logger.WithFields(log.Fields{
		"package_name":  data.PackageName,
		"routes_count":  len(data.Routes),
		"content_length": len(result),
	}).Debug("Template rendered successfully")

	return result, nil
}

// Validate checks if the renderer is properly configured
func (r *RouteRenderer) Validate() error {
	if r.template == nil {
		return NewRouteError(ErrTemplateFailed, "template is not initialized")
	}

	if r.info.Name == "" {
		return NewRouteError(ErrTemplateFailed, "template name is not set")
	}

	return nil
}

// GetTemplateInfo returns metadata about the template
func (r *RouteRenderer) GetTemplateInfo() TemplateInfo {
	return r.info
}

// validateRenderData validates the input data before rendering
func (r *RouteRenderer) validateRenderData(data RenderData) error {
	if data.PackageName == "" {
		return NewRouteError(ErrInvalidInput, "package name cannot be empty")
	}

	if len(data.Routes) == 0 {
		return NewRouteError(ErrNoRoutes, "no routes to render for package: %s", data.PackageName)
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

// estimatedBufferSize calculates the estimated buffer size needed for rendering
func estimatedBufferSize(data RenderData) int {
	// Base size for package structure
	baseSize := 1024 // ~1KB for base package structure

	// Add size based on routes
	routesSize := len(data.Routes) * 256 // ~256 bytes per route

	// Add size based on imports
	importsSize := len(data.Imports) * 64 // ~64 bytes per import

	// Add size based on controllers
	controllersSize := len(data.Controllers) * 128 // ~128 bytes per controller

	return baseSize + routesSize + importsSize + controllersSize
}

// renderTemplate is the legacy function for backward compatibility
// Use NewRouteRenderer().Render() for new code
func renderTemplate(data RenderData) ([]byte, error) {
	renderer := NewRouteRenderer()
	if renderer == nil {
		return nil, NewRouteError(ErrTemplateFailed, "failed to create route renderer")
	}
	return renderer.Render(data)
}
