# Phase 1: Quickstart Guide

## Overview
This quickstart guide demonstrates how to use the refactored AST route generation system. It covers the basic workflow for defining routes and generating route handler code.

## Prerequisites

- Go 1.21 or higher
- Existing atomctl project structure
- Basic understanding of Go annotations

## Basic Route Definition

### 1. Simple Route

Create a controller with basic route annotation:

```go
// app/http/user_controller.go
package http

// @Router /users [get]
type UserController struct {}
```

Generate routes:
```bash
atomctl gen route
```

### 2. Route with Parameters

Add parameter bindings using `@Bind` annotations:

```go
// app/http/user_controller.go
package http

// @Router /users/:id [get]
// @Bind id (path) model()
// @Bind limit (query) model(limit:int)
type UserController struct {}
```

## Parameter Binding Types

### Path Parameters
```go
// @Bind id (path) model()
// @Bind name (path) model(name:string)
```

### Query Parameters
```go
// @Bind limit (query) model(limit:int)
// @Bind offset (query) model(offset:int)
// @Bind filter (query)
```

### Body Parameters
```go
// @Bind user (body) model(User)
// @Bind data (body) model(CreateUserRequest)
```

### Header Parameters
```go
// @Bind authorization (header)
// @Bind x-api-key (header) model(APIKey)
```

## Generated Code Structure

The route generation will create a `routes.gen.go` file:

```go
// app/http/routes.gen.go
package http

import (
    "context"
    "net/http"

    "go.ipao.vip/atom/contracts"
    "go.ipao.vip/atom/http"
    "go.ipao.vip/gen/model"
)

type RouteProvider struct {
    userController *UserController
}

func (p *RouteProvider) Provide(opts ...contracts.Option) error {
    // Route registration logic here
    p.userController = &UserController{}

    // Register /users route
    http.Handle("/users", p.userController.GetUsers)

    return nil
}

// UserController method stubs
func (c *UserController) GetUsers(ctx context.Context, w http.ResponseWriter, r *http.Request) {
    // Generated method implementation
}
```

## Testing Your Routes

### 1. Unit Test
```go
// app/http/user_controller_test.go
package http

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestUserController_GetUsers(t *testing.T) {
    controller := &UserController{}

    req := httptest.NewRequest("GET", "/users", nil)
    w := httptest.NewRecorder()

    controller.GetUsers(context.Background(), w, req)

    assert.Equal(t, http.StatusOK, w.Code)
}
```

### 2. Integration Test
```go
// integration/user_routes_test.go
package integration

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestUserRoutes(t *testing.T) {
    // Setup router with generated routes
    router := setupRouter()

    tests := []struct {
        name       string
        path       string
        method     string
        wantStatus int
    }{
        {"Get Users", "/users", "GET", http.StatusOK},
        {"Get User by ID", "/users/123", "GET", http.StatusOK},
        {"Create User", "/users", "POST", http.StatusCreated},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            req := httptest.NewRequest(tt.method, tt.path, nil)
            w := httptest.NewRecorder()

            router.ServeHTTP(w, req)

            assert.Equal(t, tt.wantStatus, w.Code)
        })
    }
}
```

## Advanced Features

### 1. Route Groups
```go
// @Router /api/v1/users [get]
// @Router /api/v1/users/:id [get,put,delete]
type UserController struct {}
```

### 2. Middleware Integration
```go
// @Router /admin [get]
// @Middleware auth,admin
type AdminController struct {}
```

### 3. Custom Return Types
```go
// @Router /users [post]
// @ReturnType UserResponse
type UserController struct {}
```

## Configuration Options

### Parser Configuration
```go
config := &route.RouteParserConfig{
    StrictMode:      true,
    ParseComments:   true,
    SourceLocations: true,
    EnableValidation: true,
}
```

### Builder Configuration
```go
config := &route.BuilderConfig{
    EnableValidation:          true,
    StrictMode:                true,
    DefaultParamPosition:      route.ParamPositionQuery,
    AutoGenerateReturnTypes:  true,
    ResolveImportDependencies: true,
}
```

## Error Handling

### 1. Validation Errors
The refactored system provides detailed error messages:

```bash
$ atomctl gen route
Error: invalid route syntax in user_controller.go:15
  Expected: @Router /path [method]
  Found:    @Router /users
  Fix: Add HTTP methods in brackets
```

### 2. Parameter Binding Errors
```bash
$ atomctl gen route
Error: invalid parameter binding in user_controller.go:16
  Parameter 'id' has invalid position 'invalid'
  Valid positions: path, query, body, header, cookie, local, file
```

## Migration from Legacy System

### 1. Existing Code Compatibility
The refactored system maintains full backward compatibility:

```go
// This still works
// @Router /users [get]
// @Bind id (path) model()
type UserController struct {}
```

### 2. Gradual Migration
You can migrate files incrementally:

```bash
# Generate routes for specific files
atomctl gen route app/http/user_controller.go

# Or generate for entire directory
atomctl gen route app/http/
```

## Performance Considerations

### 1. Caching
The refactored system includes caching for improved performance:

```go
config := &route.RouteParserConfig{
    CacheEnabled: true,
}
```

### 2. Parallel Processing
Enable parallel processing for large projects:

```go
config := &route.RouteParserConfig{
    ParallelProcessing: true,
}
```

## Debugging and Diagnostics

### 1. Enable Detailed Logging
```go
config := &route.RouteParserConfig{
    SourceLocations: true,
}
```

### 2. Access Diagnostics
```go
parser := route.NewRouteParser()
routes, err := parser.ParseFile("controller.go")

// Get detailed diagnostics
diagnostics := parser.GetDiagnostics()
for _, diag := range diagnostics {
    fmt.Printf("%s: %s (%s:%d)\n", diag.Level, diag.Message, diag.File, diag.Location.Line)
}
```

## Best Practices

### 1. Route Definition
- Use descriptive route names
- Group related routes together
- Follow REST conventions where applicable

### 2. Parameter Binding
- Use appropriate parameter positions
- Provide clear parameter names
- Add validation for complex parameters

### 3. Error Handling
- Implement proper error handling in controllers
- Use appropriate HTTP status codes
- Provide meaningful error messages

### 4. Testing
- Write comprehensive tests for all routes
- Test both success and error scenarios
- Use contract tests for consistency

## Troubleshooting

### Common Issues

1. **Routes not generated**: Check file naming and location
2. **Parameters not parsed**: Verify annotation syntax
3. **Import errors**: Ensure all dependencies are available
4. **Compilation errors**: Check generated code syntax

### Getting Help

- Review the contract tests in `contracts/` directory
- Check the diagnostic output for detailed error information
- Run tests to verify implementation correctness

## Next Steps

1. Define your routes using `@Router` and `@Bind` annotations
2. Run `atomctl gen route` to generate route code
3. Implement the generated controller methods
4. Write tests to verify functionality
5. Configure options as needed for your project

The refactored system provides a solid foundation for route generation with improved maintainability, testability, and extensibility.