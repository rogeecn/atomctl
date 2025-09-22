package contracts

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// RouteParserContract defines the contract tests for RouteParser implementations
type RouteParserContract struct {
	parser RouteParser
}

// RouteParser interface definition for contract testing
type RouteParser interface {
	ParseFile(filePath string) ([]RouteDefinition, error)
	ParseDir(dirPath string) ([]RouteDefinition, error)
	ParseString(code string) ([]RouteDefinition, error)
	SetConfig(config *RouteParserConfig) error
	GetConfig() *RouteParserConfig
	GetContext() *ParserContext
	GetDiagnostics() []Diagnostic
}

// RouteDefinition represents a route definition (simplified for testing)
type RouteDefinition struct {
	StructName  string
	Path        string
	Methods     []string
	Parameters  []ParamDefinition
	Imports     map[string]string
}

// ParamDefinition represents a parameter definition (simplified for testing)
type ParamDefinition struct {
	Name     string
	Position string
	Type     string
}

// RouteParserConfig represents parser configuration (simplified for testing)
type RouteParserConfig struct {
	StrictMode      bool
	ParseComments   bool
	SourceLocations bool
}

// ParserContext represents parser context (simplified for testing)
type ParserContext struct {
	WorkingDir string
	ModuleName string
}

// Diagnostic represents diagnostic information (simplified for testing)
type Diagnostic struct {
	Level   string
	Message string
	File    string
}

// NewRouteParserContract creates a new contract test instance
func NewRouteParserContract(parser RouteParser) *RouteParserContract {
	return &RouteParserContract{parser: parser}
}

// TestSuite runs all contract tests for RouteParser
func (c *RouteParserContract) TestSuite(t *testing.T) {
	t.Run("RouteParser_ParseFile_BasicRoute", c.testParseFileBasicRoute)
	t.Run("RouteParser_ParseFile_WithParameters", c.testParseFileWithParameters)
	t.Run("RouteParser_ParseFile_InvalidSyntax", c.testParseFileInvalidSyntax)
	t.Run("RouteParser_ParseFile_NonexistentFile", c.testParseFileNonexistentFile)
	t.Run("RouteParser_ParseString_SimpleRoute", c.testParseStringSimpleRoute)
	t.Run("RouteParser_ParseString_MultipleRoutes", c.testParseStringMultipleRoutes)
	t.Run("RouteParser_ParseDir_EmptyDirectory", c.testParseDirEmptyDirectory)
	t.Run("RouteParser_Configuration", c.testConfiguration)
	t.Run("RouteParser_Diagnostics", c.testDiagnostics)
}

// Contract Tests

func (c *RouteParserContract) testParseFileBasicRoute(t *testing.T) {
	// GIVEN a Go file with basic route annotation
	filePath := "testdata/basic_route.go"

	// WHEN parsing the file
	routes, err := c.parser.ParseFile(filePath)

	// THEN it should succeed and return the route definition
	assert.NoError(t, err, "ParseFile should not error")
	assert.Len(t, routes, 1, "Should find exactly one route")

	route := routes[0]
	assert.Equal(t, "UserController", route.StructName, "Should extract correct struct name")
	assert.Equal(t, "/users", route.Path, "Should extract correct path")
	assert.Contains(t, route.Methods, "GET", "Should include GET method")
	assert.Empty(t, route.Parameters, "Basic route should have no parameters")
}

func (c *RouteParserContract) testParseFileWithParameters(t *testing.T) {
	// GIVEN a Go file with route containing parameters
	filePath := "testdata/params_route.go"

	// WHEN parsing the file
	routes, err := c.parser.ParseFile(filePath)

	// THEN it should succeed and extract parameters
	assert.NoError(t, err, "ParseFile should not error")
	assert.Len(t, routes, 1, "Should find exactly one route")

	route := routes[0]
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

func (c *RouteParserContract) testParseFileInvalidSyntax(t *testing.T) {
	// GIVEN a Go file with invalid syntax
	filePath := "testdata/invalid_syntax.go"

	// WHEN parsing the file
	routes, err := c.parser.ParseFile(filePath)

	// THEN it should fail with appropriate error
	assert.Error(t, err, "ParseFile should error on invalid syntax")
	assert.Empty(t, routes, "Should return no routes on error")

	// Verify diagnostic information
	diagnostics := c.parser.GetDiagnostics()
	assert.NotEmpty(t, diagnostics, "Should provide diagnostic information")
	assert.Contains(t, diagnostics[0].Message, "syntax", "Error message should mention syntax")
}

func (c *RouteParserContract) testParseFileNonexistentFile(t *testing.T) {
	// GIVEN a nonexistent file path
	filePath := "testdata/nonexistent.go"

	// WHEN parsing the file
	routes, err := c.parser.ParseFile(filePath)

	// THEN it should fail with file not found error
	assert.Error(t, err, "ParseFile should error on nonexistent file")
	assert.Empty(t, routes, "Should return no routes on error")
}

func (c *RouteParserContract) testParseStringSimpleRoute(t *testing.T) {
	// GIVEN Go code string with simple route
	code := `
package main

// @Router /users [get]
type UserController struct {}
`

	// WHEN parsing the string
	routes, err := c.parser.ParseString(code)

	// THEN it should succeed and return the route
	assert.NoError(t, err, "ParseString should not error")
	assert.Len(t, routes, 1, "Should find exactly one route")

	route := routes[0]
	assert.Equal(t, "UserController", route.StructName, "Should extract correct struct name")
	assert.Equal(t, "/users", route.Path, "Should extract correct path")
}

func (c *RouteParserContract) testParseStringMultipleRoutes(t *testing.T) {
	// GIVEN Go code string with multiple routes
	code := `
package main

// @Router /users [get]
type UserController struct {}

// @Router /products [post]
type ProductController struct {}
`

	// WHEN parsing the string
	routes, err := c.parser.ParseString(code)

	// THEN it should succeed and return all routes
	assert.NoError(t, err, "ParseString should not error")
	assert.Len(t, routes, 2, "Should find exactly two routes")

	// Verify both routes are found
	routeNames := []string{routes[0].StructName, routes[1].StructName}
	assert.Contains(t, routeNames, "UserController", "Should find UserController")
	assert.Contains(t, routeNames, "ProductController", "Should find ProductController")
}

func (c *RouteParserContract) testParseDirEmptyDirectory(t *testing.T) {
	// GIVEN an empty directory
	dirPath := "testdata/empty"

	// WHEN parsing the directory
	routes, err := c.parser.ParseDir(dirPath)

	// THEN it should succeed with no routes
	assert.NoError(t, err, "ParseDir should not error on empty directory")
	assert.Empty(t, routes, "Should return no routes from empty directory")
}

func (c *RouteParserContract) testConfiguration(t *testing.T) {
	// GIVEN default configuration
	config := c.parser.GetConfig()
	assert.NotNil(t, config, "Should have default configuration")

	// WHEN setting new configuration
	newConfig := &RouteParserConfig{
		StrictMode:      true,
		ParseComments:   false,
		SourceLocations: true,
	}

	err := c.parser.SetConfig(newConfig)

	// THEN configuration should be updated
	assert.NoError(t, err, "SetConfig should not error")

	updatedConfig := c.parser.GetConfig()
	assert.True(t, updatedConfig.StrictMode, "StrictMode should be updated")
	assert.False(t, updatedConfig.ParseComments, "ParseComments should be updated")
	assert.True(t, updatedConfig.SourceLocations, "SourceLocations should be updated")
}

func (c *RouteParserContract) testDiagnostics(t *testing.T) {
	// GIVEN a file with warnings
	filePath := "testdata/warnings.go"

	// WHEN parsing the file
	_, err := c.parser.ParseFile(filePath)

	// THEN diagnostics should be available
	diagnostics := c.parser.GetDiagnostics()
	assert.NotEmpty(t, diagnostics, "Should provide diagnostic information")

	// Verify diagnostic structure
	for _, diag := range diagnostics {
		assert.NotEmpty(t, diag.Level, "Diagnostic should have level")
		assert.NotEmpty(t, diag.Message, "Diagnostic should have message")
	}
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