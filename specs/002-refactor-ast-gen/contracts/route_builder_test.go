package contracts

import (
	"go/ast"
	"testing"

	"github.com/stretchr/testify/assert"
)

// RouteBuilderContract defines the contract tests for RouteBuilder implementations
type RouteBuilderContract struct {
	builder RouteBuilder
}

// RouteBuilder interface definition for contract testing
type RouteBuilder interface {
	BuildFromTypeSpec(typeSpec *ast.TypeSpec, decl *ast.GenDecl, context *BuilderContext) (RouteDefinition, error)
	BuildFromComment(comment string, context *BuilderContext) (RouteDefinition, error)
	ValidateDefinition(def *RouteDefinition) error
}

// BuilderContext represents builder context (simplified for testing)
type BuilderContext struct {
	FilePath      string
	PackageName   string
	ImportContext *ImportContext
	ASTFile       *ast.File
}

// ImportContext represents import context (simplified for testing)
type ImportContext struct {
	Imports map[string]string
}

// NewRouteBuilderContract creates a new contract test instance
func NewRouteBuilderContract(builder RouteBuilder) *RouteBuilderContract {
	return &RouteBuilderContract{builder: builder}
}

// TestSuite runs all contract tests for RouteBuilder
func (c *RouteBuilderContract) TestSuite(t *testing.T) {
	t.Run("RouteBuilder_BuildFromTypeSpec_BasicRoute", c.testBuildFromTypeSpecBasicRoute)
	t.Run("RouteBuilder_BuildFromTypeSpec_WithParameters", c.testBuildFromTypeSpecWithParameters)
	t.Run("RouteBuilder_BuildFromTypeSpec_InvalidInput", c.testBuildFromTypeSpecInvalidInput)
	t.Run("RouteBuilder_BuildFromComment_SimpleComment", c.testBuildFromCommentSimpleComment)
	t.Run("RouteBuilder_BuildFromComment_ComplexComment", c.testBuildFromCommentComplexComment)
	t.Run("RouteBuilder_ValidateDefinition_ValidRoute", c.testValidateDefinitionValidRoute)
	t.Run("RouteBuilder_ValidateDefinition_InvalidRoute", c.testValidateDefinitionInvalidRoute)
}

// Contract Tests

func (c *RouteBuilderContract) testBuildFromTypeSpecBasicRoute(t *testing.T) {
	// GIVEN a valid type specification and declaration
	typeSpec := &ast.TypeSpec{
		Name: &ast.Ident{Name: "UserController"},
	}

	decl := &ast.GenDecl{
		Doc: &ast.CommentGroup{
			List: []*ast.Comment{
				{Text: "// @Router /users [get]"},
			},
		},
		Specs: []ast.Spec{typeSpec},
	}

	context := &BuilderContext{
		FilePath:    "UserController.go",
		PackageName: "controllers",
		ImportContext: &ImportContext{
			Imports: make(map[string]string),
		},
	}

	// WHEN building route from type spec
	route, err := c.builder.BuildFromTypeSpec(typeSpec, decl, context)

	// THEN it should succeed and return valid route definition
	assert.NoError(t, err, "BuildFromTypeSpec should not error")
	assert.NotNil(t, route, "Should return route definition")

	assert.Equal(t, "UserController", route.StructName, "Should extract correct struct name")
	assert.Equal(t, "/users", route.Path, "Should extract correct path")
	assert.Contains(t, route.Methods, "GET", "Should include GET method")
}

func (c *RouteBuilderContract) testBuildFromTypeSpecWithParameters(t *testing.T) {
	// GIVEN a type specification with parameter bindings
	typeSpec := &ast.TypeSpec{
		Name: &ast.Ident{Name: "UserController"},
	}

	decl := &ast.GenDecl{
		Doc: &ast.CommentGroup{
			List: []*ast.Comment{
				{Text: "// @Router /users/:id [get]"},
				{Text: "// @Bind id (path) model()"},
				{Text: "// @Bind limit (query) model(limit:int)"},
			},
		},
		Specs: []ast.Spec{typeSpec},
	}

	context := &BuilderContext{
		FilePath:    "UserController.go",
		PackageName: "controllers",
		ImportContext: &ImportContext{
			Imports: map[string]string{
				"model": "go.ipao.vip/gen/model",
			},
		},
	}

	// WHEN building route from type spec
	route, err := c.builder.BuildFromTypeSpec(typeSpec, decl, context)

	// THEN it should succeed and extract parameters
	assert.NoError(t, err, "BuildFromTypeSpec should not error")
	assert.NotNil(t, route, "Should return route definition")

	assert.NotEmpty(t, route.Parameters, "Route should have parameters")

	// Verify path parameter
	pathParam := findParameterByPosition(route.Parameters, "path")
	assert.NotNil(t, pathParam, "Should find path parameter")
	assert.Equal(t, "id", pathParam.Name, "Path parameter should be named 'id'")

	// Verify query parameter
	queryParam := findParameterByPosition(route.Parameters, "query")
	assert.NotNil(t, queryParam, "Should find query parameter")
	assert.Equal(t, "limit", queryParam.Name, "Query parameter should be named 'limit'")
}

func (c *RouteBuilderContract) testBuildFromTypeSpecInvalidInput(t *testing.T) {
	// GIVEN an invalid type specification (nil)
	var typeSpec *ast.TypeSpec = nil
	decl := &ast.GenDecl{}
	context := &BuilderContext{}

	// WHEN building route from invalid type spec
	route, err := c.builder.BuildFromTypeSpec(typeSpec, decl, context)

	// THEN it should fail with appropriate error
	assert.Error(t, err, "BuildFromTypeSpec should error on invalid input")
	assert.Equal(t, RouteDefinition{}, route, "Should return empty route on error")
}

func (c *RouteBuilderContract) testBuildFromCommentSimpleComment(t *testing.T) {
	// GIVEN a simple route comment
	comment := `@Router /users [get]`
	context := &BuilderContext{
		FilePath:    "UserController.go",
		PackageName: "controllers",
	}

	// WHEN building route from comment
	route, err := c.builder.BuildFromComment(comment, context)

	// THEN it should succeed and return valid route
	assert.NoError(t, err, "BuildFromComment should not error")
	assert.NotNil(t, route, "Should return route definition")

	assert.Equal(t, "/users", route.Path, "Should extract correct path")
	assert.Contains(t, route.Methods, "GET", "Should include GET method")
}

func (c *RouteBuilderContract) testBuildFromCommentComplexComment(t *testing.T) {
	// GIVEN a complex route comment with parameters
	comment := `@Router /users/:id [get,put]
@Bind id (path) model()
@Bind user (body) model(User)`
	context := &BuilderContext{
		FilePath:    "UserController.go",
		PackageName: "controllers",
		ImportContext: &ImportContext{
			Imports: map[string]string{
				"model": "go.ipao.vip/gen/model",
			},
		},
	}

	// WHEN building route from comment
	route, err := c.builder.BuildFromComment(comment, context)

	// THEN it should succeed and extract all information
	assert.NoError(t, err, "BuildFromComment should not error")
	assert.NotNil(t, route, "Should return route definition")

	assert.Equal(t, "/users/:id", route.Path, "Should extract correct path")
	assert.Contains(t, route.Methods, "GET", "Should include GET method")
	assert.Contains(t, route.Methods, "PUT", "Should include PUT method")

	assert.NotEmpty(t, route.Parameters, "Route should have parameters")
}

func (c *RouteBuilderContract) testValidateDefinitionValidRoute(t *testing.T) {
	// GIVEN a valid route definition
	route := RouteDefinition{
		StructName: "UserController",
		Path:       "/users",
		Methods:    []string{"GET", "POST"},
		Parameters: []ParamDefinition{
			{Name: "id", Position: "path", Type: "int"},
		},
	}

	// WHEN validating the route definition
	err := c.builder.ValidateDefinition(&route)

	// THEN it should succeed
	assert.NoError(t, err, "ValidateDefinition should not error on valid route")
}

func (c *RouteBuilderContract) testValidateDefinitionInvalidRoute(t *testing.T) {
	// GIVEN an invalid route definition (empty path)
	route := RouteDefinition{
		StructName: "UserController",
		Path:       "", // Empty path is invalid
		Methods:    []string{"GET"},
	}

	// WHEN validating the route definition
	err := c.builder.ValidateDefinition(&route)

	// THEN it should fail
	assert.Error(t, err, "ValidateDefinition should error on invalid route")
}

// Helper functions

func findParameterByPosition(params []ParamDefinition, position string) *ParamDefinition {
	for _, param := range params {
		if param.Position == position {
			return &param
		}
	}
	return nil
}