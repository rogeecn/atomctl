package route

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRouteParsing(t *testing.T) {
	t.Run("BasicRouteParsing", func(t *testing.T) {
		// This test will fail until we improve the existing parsing
		// GIVEN a Go file with basic route annotation
		code := `
package main

// UserController defines user-related routes
type UserController struct {}

// @Router /users [get]
func (c *UserController) GetUser() error {
	return nil
}
`

		// Create a temporary file for testing
		tmpFile := "/tmp/test_route.go"
		err := os.WriteFile(tmpFile, []byte(code), 0o644)
		assert.NoError(t, err, "Should create temp file")
		defer os.Remove(tmpFile)

		// WHEN parsing the file using existing API
		routes := ParseFile(tmpFile)

		// THEN it should return the route definition
		// Note: Current implementation may not extract all info correctly
		assert.NotEmpty(t, routes, "Should find at least one route")

		route := routes[0]
		assert.Equal(t, "/users", route.Path, "Should extract correct path")

		// Find the GET action
		getAction := findActionByMethod(route.Actions, "GET")
		assert.NotNil(t, getAction, "Should find GET action")
	})
}

func TestParameterBinding(t *testing.T) {
	t.Run("ParameterBinding", func(t *testing.T) {
		// GIVEN a Go file with parameter bindings
		code := `
package main

type UserController struct {}

// @Router /users/:id [get]
// @Bind id (path) model()
// @Bind limit (query) model(limit:int)
func (c *UserController) GetUser(id string, limit int) {
}
`

		// Create a temporary file for testing
		tmpFile := "/tmp/test_params.go"
		err := os.WriteFile(tmpFile, []byte(code), 0o644)
		assert.NoError(t, err, "Should create temp file")
		defer os.Remove(tmpFile)

		// WHEN parsing the file using existing API
		routes := ParseFile(tmpFile)

		// THEN it should succeed and extract parameters
		assert.NotEmpty(t, routes, "Should find at least one route")

		route := routes[0]
		getAction := findActionByMethod(route.Actions, "GET")
		assert.NotNil(t, getAction, "Should find GET action")
		assert.NotEmpty(t, getAction.Params, "GET action should have parameters")

		// Verify path parameter
		pathParam := findParameterByPosition(getAction.Params, PositionPath)
		assert.NotNil(t, pathParam, "Should find path parameter")
		assert.Equal(t, "id", pathParam.Name, "Path parameter should be named 'id'")

		// Verify query parameter
		queryParam := findParameterByPosition(getAction.Params, PositionQuery)
		assert.NotNil(t, queryParam, "Should find query parameter")
		assert.Equal(t, "limit", queryParam.Name, "Query parameter should be named 'limit'")
	})
}

func TestErrorHandling(t *testing.T) {
	t.Run("InvalidSyntax", func(t *testing.T) {
		// GIVEN invalid Go code
		code := `
package main

type UserController struct {}

// @Router /users [get]  // Missing closing bracket
func (c *UserController) GetUser() {
`

		// Create a temporary file for testing
		tmpFile := "/tmp/test_invalid.go"
		err := os.WriteFile(tmpFile, []byte(code), 0o644)
		assert.NoError(t, err, "Should create temp file")
		defer os.Remove(tmpFile)

		// WHEN parsing the file using existing API
		routes := ParseFile(tmpFile)

		// THEN it should handle the error gracefully
		// Note: Current implementation may log errors but return nil/empty
		assert.Empty(t, routes, "Should return no routes on invalid syntax")
	})

	t.Run("EmptyRoute", func(t *testing.T) {
		// GIVEN a route with empty path
		code := `
package main

type UserController struct {}

// @Router  [get]
func (c *UserController) GetUser() {
}
`

		// Create a temporary file for testing
		tmpFile := "/tmp/test_empty.go"
		err := os.WriteFile(tmpFile, []byte(code), 0o644)
		assert.NoError(t, err, "Should create temp file")
		defer os.Remove(tmpFile)

		// WHEN parsing the file using existing API
		routes := ParseFile(tmpFile)

		// THEN it should handle the validation appropriately
		// Current implementation creates route entries but skips invalid route definitions
		// So we expect 1 route but with no actions
		assert.Equal(t, 1, len(routes), "Should create route entry for controller")
		assert.Equal(t, 0, len(routes[0].Actions), "Should have no actions for invalid route")
	})
}

// Helper functions (using existing data structures)
func findParameterByPosition(params []ParamDefinition, position Position) *ParamDefinition {
	for _, param := range params {
		if param.Position == position {
			return &param
		}
	}
	return nil
}

func findActionByMethod(actions []ActionDefinition, method string) *ActionDefinition {
	for _, action := range actions {
		if action.Method == method {
			return &action
		}
	}
	return nil
}

func TestBackwardCompatibility(t *testing.T) {
	t.Run("ExistingAnnotationFormats", func(t *testing.T) {
		// GIVEN existing annotation formats that must continue to work
		code := `
package main

type UserController struct {
}

// @Router /api/v1/users [get]
// @Bind id (path) model()
// @Bind name (query) model(name:string)
func (c *UserController) GetUser(id int, name string) error {
	return nil
}

// @Router /api/v1/users [post]
// @Bind user (body) model(*User)
func (c *UserController) CreateUser(user *User) error {
	return nil
}

type HealthController struct {
}

// @Router /health [get]
func (c *HealthController) Check() error {
	return nil
}
`
		// Create a temporary file for testing
		tmpFile := "/tmp/test_compatibility.go"
		err := os.WriteFile(tmpFile, []byte(code), 0o644)
		assert.NoError(t, err, "Should create temp file")
		defer os.Remove(tmpFile)

		// WHEN parsing the file using existing API
		routes := ParseFile(tmpFile)

		// THEN it should extract all routes without breaking
		assert.NotEmpty(t, routes, "Should find routes")
		assert.Equal(t, 2, len(routes), "Should find exactly 2 routes")

		// Verify user route with multiple actions
		userRoute := findRouteByPath(routes, "/api/v1/users")
		assert.NotNil(t, userRoute, "Should find user route")
		assert.Equal(t, 2, len(userRoute.Actions), "Should have GET and POST actions")

		// Verify HTTP methods
		getAction := findActionByMethod(userRoute.Actions, "GET")
		assert.NotNil(t, getAction, "Should find GET action")
		postAction := findActionByMethod(userRoute.Actions, "POST")
		assert.NotNil(t, postAction, "Should find POST action")

		// Verify parameter binding compatibility
		assert.GreaterOrEqual(t, len(getAction.Params), 2, "GET action should have parameters bound")

		// Verify health route
		healthRoute := findRouteByPath(routes, "/health")
		assert.NotNil(t, healthRoute, "Should find health route")
		assert.Equal(t, 1, len(healthRoute.Actions), "Should have 1 action")
	})

	t.Run("SpecialCharactersInPaths", func(t *testing.T) {
		// GIVEN routes with special characters that must work
		code := `
package main

type ApiController struct {
}

// @Router /api/v1/users/:id/profile [get]
func (c *ApiController) GetUserProfile(id string) error {
	return nil
}

// @Router /api/v1/orders/:order_id/items/:item_id [post]
func (c *ApiController) GetOrderItem(orderID, itemID string) error {
	return nil
}

// @Router /download/*filename [put]
func (c *ApiController) DownloadFile(filename string) error {
	return nil
}
`
		tmpFile := "/tmp/test_special_paths.go"
		err := os.WriteFile(tmpFile, []byte(code), 0o644)
		assert.NoError(t, err, "Should create temp file")
		defer os.Remove(tmpFile)

		// WHEN parsing
		routes := ParseFile(tmpFile)

		// THEN special paths should be preserved
		assert.Equal(t, 1, len(routes), "Should find 1 controller with multiple actions")

		// All routes belong to ApiController, check that all actions are present
		apiRoute := findRouteByPath(routes, "/api/v1/users/:id/profile")
		assert.NotNil(t, apiRoute, "Should find API controller route")
		assert.Equal(t, 3, len(apiRoute.Actions), "Should have 3 actions")

		// Verify all three methods are present
		actionMethods := make(map[string]bool)
		for _, action := range apiRoute.Actions {
			actionMethods[action.Method] = true
		}

		assert.True(t, actionMethods["GET"], "Should have GET method")
		assert.True(t, actionMethods["POST"], "Should have POST method")
		assert.True(t, actionMethods["PUT"], "Should have PUT method")
		assert.Equal(t, 3, len(actionMethods), "Should have 3 different methods")

		// Verify we can find actions by checking names
		var foundUserProfile, foundOrderItem, foundDownload bool
		for _, action := range apiRoute.Actions {
			switch action.Name {
			case "GetUserProfile":
				foundUserProfile = true
			case "GetOrderItem":
				foundOrderItem = true
			case "DownloadFile":
				foundDownload = true
			}
		}

		assert.True(t, foundUserProfile, "Should find GetUserProfile action")
		assert.True(t, foundOrderItem, "Should find GetOrderItem action")
		assert.True(t, foundDownload, "Should find DownloadFile action")
	})
}

// Helper function to find route by path
func findRouteByPath(routes []RouteDefinition, path string) *RouteDefinition {
	for _, route := range routes {
		if route.Path == path {
			return &route
		}
	}
	return nil
}

func TestCLIIntegration(t *testing.T) {
	t.Run("ParseFileIntegration", func(t *testing.T) {
		// GIVEN a realistic controller file structure
		code := `
package main

import (
	"net/http"
)

type UserController struct {
}

// @Router /api/v1/users [get]
// @Bind id (path) model()
// @Bind filter (query) model(filter:UserFilter)
func (c *UserController) GetUser(id int, filter string) error {
	return nil
}

// @Router /api/v1/users [post]
// @Bind user (body) model(*User)
func (c *UserController) CreateUser(user *User) error {
	return nil
}

type ProductController struct {
}

// @Router /api/v1/products [get]
// @Bind category (query) model(category:string)
func (c *ProductController) GetProducts(category string) error {
	return nil
}

type HealthController struct {
}

// @Router /health [get]
func (c *HealthController) Check() error {
	return nil
}
`
		// Create a realistic file structure
		tmpDir := "/tmp/test_app"
		httpDir := tmpDir + "/app/http"
		err := os.MkdirAll(httpDir, 0o755)
		assert.NoError(t, err, "Should create directory structure")
		defer os.RemoveAll(tmpDir)

		// Write controller file
		controllerFile := httpDir + "/user_controller.go"
		err = os.WriteFile(controllerFile, []byte(code), 0o644)
		assert.NoError(t, err, "Should write controller file")

		// WHEN parsing using the same method as CLI
		routes := ParseFile(controllerFile)

		// THEN it should work as expected by the CLI
		assert.NotEmpty(t, routes, "Should parse routes successfully")
		assert.Equal(t, 3, len(routes), "Should find 3 controllers")

		// Verify the data structure matches CLI expectations
		userRoute := findRouteByPath(routes, "/api/v1/users")
		assert.NotNil(t, userRoute, "Should find user route")
		assert.Equal(t, "UserController", userRoute.Name, "Should preserve controller name")
		assert.Equal(t, 2, len(userRoute.Actions), "Should have GET and POST actions")

		// Verify parameter structure is compatible
		getAction := findActionByMethod(userRoute.Actions, "GET")
		assert.NotNil(t, getAction, "Should have GET action")

		// Verify the data can be processed by the existing grouping logic
		// (similar to lo.GroupBy in cmd/gen_route.go:105)
		routeGroups := make(map[string][]RouteDefinition)
		for _, route := range routes {
			dirPath := filepath.Dir(route.Path)
			routeGroups[dirPath] = append(routeGroups[dirPath], route)
		}

		assert.Greater(t, len(routeGroups), 0, "Should be groupable by path")
	})

	t.Run("EmptyFileHandling", func(t *testing.T) {
		// GIVEN an empty Go file
		code := `
package main

type EmptyController struct {
}
`
		tmpFile := "/tmp/test_empty.go"
		err := os.WriteFile(tmpFile, []byte(code), 0o644)
		assert.NoError(t, err, "Should create temp file")
		defer os.Remove(tmpFile)

		// WHEN parsing
		routes := ParseFile(tmpFile)

		// THEN it should handle gracefully (no panic, empty result)
		assert.Empty(t, routes, "Should return empty routes for file without annotations")
	})

	t.Run("FileWithSyntaxErrors", func(t *testing.T) {
		// GIVEN a file with syntax errors
		code := `
package main

// @Router /test [get]
// This is not valid Go syntax
type BrokenController struct {
`
		tmpFile := "/tmp/test_broken.go"
		err := os.WriteFile(tmpFile, []byte(code), 0o644)
		assert.NoError(t, err, "Should create temp file")
		defer os.Remove(tmpFile)

		// WHEN parsing
		routes := ParseFile(tmpFile)

		// THEN it should handle gracefully
		// The current implementation may return empty or partial results
		// This test verifies the CLI won't crash
		assert.True(t, len(routes) == 0, "Should handle syntax errors gracefully")
	})
}

func TestRouteRenderer(t *testing.T) {
	t.Run("NewRouteRenderer", func(t *testing.T) {
		renderer := NewRouteRenderer()
		assert.NotNil(t, renderer)
		assert.NotNil(t, renderer.logger)

		info := renderer.GetTemplateInfo()
		assert.Equal(t, "router", info.Name)
		assert.Equal(t, "1.0.0", info.Version)
		assert.Contains(t, info.Functions, "sprig")
	})

	t.Run("RendererValidation", func(t *testing.T) {
		renderer := NewRouteRenderer()
		assert.NotNil(t, renderer)

		err := renderer.Validate()
		assert.NoError(t, err)
	})

	t.Run("RenderValidData", func(t *testing.T) {
		renderer := NewRouteRenderer()
		assert.NotNil(t, renderer)

		data := RenderData{
			PackageName:    "main",
			ProjectPackage: "test/project",
			Imports:        []string{`"fmt"`},
			Controllers:    []string{"userController *UserController"},
			Routes: map[string][]Router{
				"UserController": {
					{
						Method:     "Get",
						Route:      "/users",
						Controller: "userController",
						Action:     "GetUser",
						Func:       "Func0",
					},
				},
			},
			RouteGroups: []string{"UserController"},
		}

		content, err := renderer.Render(data)
		assert.NoError(t, err)
		assert.NotEmpty(t, content)
		assert.Contains(t, string(content), "package main")
		assert.Contains(t, string(content), "func (r *Routes) Register")
	})

	t.Run("RenderInvalidData", func(t *testing.T) {
		renderer := NewRouteRenderer()
		assert.NotNil(t, renderer)

		data := RenderData{
			PackageName: "", // Invalid empty package name
			Routes:      map[string][]Router{},
		}

		content, err := renderer.Render(data)
		assert.Error(t, err)
		assert.Nil(t, content)
		assert.Contains(t, err.Error(), "package name cannot be empty")
	})

	t.Run("RenderNoRoutes", func(t *testing.T) {
		renderer := NewRouteRenderer()
		assert.NotNil(t, renderer)

		data := RenderData{
			PackageName: "main",
			Routes:      map[string][]Router{}, // Empty routes
		}

		content, err := renderer.Render(data)
		assert.Error(t, err)
		assert.Nil(t, content)
		assert.Contains(t, err.Error(), "no routes to render")
	})

	t.Run("LegacyRenderTemplate", func(t *testing.T) {
		data := RenderData{
			PackageName:    "main",
			ProjectPackage: "test/project",
			Imports:        []string{`"fmt"`},
			Controllers:    []string{"userController *UserController"},
			Routes: map[string][]Router{
				"UserController": {
					{
						Method:     "Get",
						Route:      "/users",
						Controller: "userController",
						Action:     "GetUser",
						Func:       "Func0",
					},
				},
			},
			RouteGroups: []string{"UserController"},
		}

		content, err := renderTemplate(data)
		assert.NoError(t, err)
		assert.NotEmpty(t, content)
		assert.Contains(t, string(content), "package main")
	})
}
