# Phase 1: Data Model Design

## Overview
This document defines the core data models and contracts for the refactored AST route generation system, following the component-based architecture pattern established in the provider module.

## Core Entities

### 1. RouteDefinition
Represents a complete route definition extracted from source code annotations.

```go
type RouteDefinition struct {
    // Basic Information
    StructName    string              // Name of the handler struct
    FilePath      string              // Source file path
    PackageName   string              // Go package name

    // Route Configuration
    Path          string              // HTTP path pattern
    Methods       []string           // HTTP methods (GET, POST, etc.)
    Name          string              // Route name for identification

    // Dependencies & Imports
    Imports       map[string]string  // Package path -> alias mapping
    Parameters    []ParamDefinition  // Route parameters
    Middleware    []string           // Middleware chain

    // Code Generation
    HandlerFunc   string              // Handler function name
    ReturnType    string              // Return type specification
    ProviderGroup string              // Dependency injection group

    // Metadata
    Location      SourceLocation      // Source location for error reporting
    Annotations   map[string]string   // Additional annotations
}
```

### 2. ParamDefinition
Represents a parameter binding from different sources.

```go
type ParamDefinition struct {
    // Parameter Identification
    Name          string              // Parameter name
    Position      ParamPosition       // Parameter location (path, query, etc.)
    Type          string              // Parameter type
    Source        string              // Source annotation text

    // Type Information
    BaseType      string              // Base type without pointer
    IsPointer     bool                // Is pointer type
    IsSlice       bool                // Is slice type
    IsMap         bool                // Is map type

    // Model Information (for structured parameters)
    ModelName     string              // Model name for model() binding
    ModelField    string              // Target field in model
    ModelType     string              // Model field type

    // Validation & Constraints
    Required      bool                // Is required parameter
    DefaultValue  string              // Default value
    Validation    []ValidationRule    // Validation rules

    // Code Generation
    ParamToken    string              // Template token for generation
    ImportPath    string              // Import path for type
}
```

### 3. ParamPosition
Parameter position enumeration.

```go
type ParamPosition string

const (
    ParamPositionPath   ParamPosition = "path"   // URL path parameters
    ParamPositionQuery  ParamPosition = "query"  // Query string parameters
    ParamPositionBody   ParamPosition = "body"   // Request body
    ParamPositionHeader ParamPosition = "header" // HTTP headers
    ParamPositionCookie ParamPosition = "cookie" // Cookies
    ParamPositionLocal  ParamPosition = "local"  // Local context values
    ParamPositionFile   ParamPosition = "file"   // File uploads
)
```

### 4. ValidationRule
Parameter validation rule definition.

```go
type ValidationRule struct {
    Type        ValidationType // Validation type (required, min, max, etc.)
    Value       string         // Validation value
    Message     string         // Error message
    Constraint  string         // Constraint expression
}
```

### 5. SourceLocation
Source code location information.

```go
type SourceLocation struct {
    File   string // Source file path
    Line   int    // Line number
    Column int    // Column position
}
```

## Component Interfaces

### 1. RouteParser Interface
Main coordinator interface for route parsing.

```go
type RouteParser interface {
    ParseFile(filePath string) ([]RouteDefinition, error)
    ParseDir(dirPath string) ([]RouteDefinition, error)
    ParseString(code string) ([]RouteDefinition, error)

    // Configuration
    SetConfig(config *RouteParserConfig) error
    GetConfig() *RouteParserConfig

    // Context & Diagnostics
    GetContext() *ParserContext
    GetDiagnostics() []Diagnostic
}
```

### 2. CommentParser Interface
Handles annotation parsing from source comments.

```go
type CommentParser interface {
    ParseRouteComment(comment string) (*RouteAnnotation, error)
    ParseBindComment(comment string) (*BindAnnotation, error)
    IsRouteComment(comment string) bool
    IsBindComment(comment string) bool
}
```

### 3. ImportResolver Interface
Manages import resolution and dependencies.

```go
type ImportResolver interface {
    ResolveFileImports(node *ast.File, filePath string) (*ImportContext, error)
    ResolveImportPath(alias string) (string, error)
    AddImport(imports map[string]string, path, alias string) error
}
```

### 4. RouteBuilder Interface
Constructs route definitions from parsed components.

```go
type RouteBuilder interface {
    BuildFromTypeSpec(typeSpec *ast.TypeSpec, decl *ast.GenDecl, context *BuilderContext) (RouteDefinition, error)
    BuildFromComment(comment string, context *BuilderContext) (RouteDefinition, error)
    ValidateDefinition(def *RouteDefinition) error
}
```

### 5. RouteValidator Interface
Validates route definitions and configurations.

```go
type RouteValidator interface {
    ValidateRoute(def *RouteDefinition) error
    ValidateParameters(params []ParamDefinition) error
    ValidateImports(imports map[string]string) error
    GetValidationRules() []ValidationRule
}
```

### 6. RouteRenderer Interface
Handles template rendering and code generation.

```go
type RouteRenderer interface {
    Render(routes []RouteDefinition, outputPath string) error
    RenderToFile(route RouteDefinition, outputPath string) error
    SetTemplate(template *template.Template) error
    GetTemplate() *template.Template
}
```

## Supporting Data Structures

### 1. RouteAnnotation
Parsed route annotation information.

```go
type RouteAnnotation struct {
    Path    string   // Route path
    Methods []string // HTTP methods
    Name    string   // Route name
    Options  map[string]string // Additional options
}
```

### 2. BindAnnotation
Parsed bind annotation information.

```go
type BindAnnotation struct {
    Name      string         // Parameter name
    Position  ParamPosition  // Parameter position
    Key       string         // Source key
    Model     *ModelInfo     // Model binding info (optional)
    Options   map[string]string // Additional options
}
```

### 3. ModelInfo
Model binding information.

```go
type ModelInfo struct {
    Name     string // Model name
    Field    string // Target field
    Type     string // Field type
    Required bool   // Is required
}
```

### 4. ImportContext
Import resolution context.

```go
type ImportContext struct {
    FileImports    map[string]*ImportResolution // Alias -> Resolution
    ImportPaths    map[string]string            // Path -> Alias
    ModuleInfo     map[string]string            // Module path -> module name
    WorkingDir     string                       // Current working directory
    ModuleName     string                       // Current module name
    ProcessedFiles map[string]bool              // Track processed files
}
```

### 5. ImportResolution
Individual import resolution information.

```go
type ImportResolution struct {
    Path     string // Import path
    Alias    string // Import alias
    Type     ImportType // Import type
    Used     bool   // Whether import is used
}
```

### 6. BuilderContext
Context for route building process.

```go
type BuilderContext struct {
    FilePath       string              // Current file path
    PackageName    string              // Package name
    ImportContext  *ImportContext      // Import information
    ASTFile        *ast.File           // AST node
    ProcessedTypes map[string]bool     // Processed types cache
    Errors         []error             // Error collection
    Warnings       []string            // Warning collection
    Config         *BuilderConfig      // Builder configuration
}
```

## Configuration Structures

### 1. RouteParserConfig
Configuration for route parser behavior.

```go
type RouteParserConfig struct {
    // Parsing Options
    ParseComments     bool           // Parse comments (default: true)
    StrictMode        bool           // Strict validation mode
    SourceLocations   bool           // Include source locations

    // File Processing
    SkipTestFiles     bool           // Skip test files (default: true)
    SkipGenerated     bool           // Skip generated files (default: true)
    AllowedPatterns   []string       // Allowed file patterns

    // Validation Options
    EnableValidation  bool           // Enable validation (default: true)
    ValidationLevel   ValidationLevel // Validation strictness

    // Performance Options
    CacheEnabled      bool           // Enable parsing cache
    ParallelProcessing bool           // Enable parallel processing
}
```

### 2. BuilderConfig
Configuration for route builder.

```go
type BuilderConfig struct {
    EnableValidation          bool          // Enable validation
    StrictMode                bool          // Strict validation mode
    DefaultParamPosition      ParamPosition // Default parameter position
    AutoGenerateReturnTypes  bool          // Auto-generate return types
    ResolveImportDependencies bool          // Resolve import dependencies
}
```

### 3. ValidationLevel
Validation strictness levels.

```go
type ValidationLevel int

const (
    ValidationLevelNone    ValidationLevel = iota // No validation
    ValidationLevelBasic                         // Basic validation
    ValidationLevelStrict                        // Strict validation
    ValidationLevelPedantic                      // Pedantic validation
)
```

## Error Handling

### 1. RouteError
Route-specific error type.

```go
type RouteError struct {
    Code      ErrorCode     // Error code
    Message   string        // Error message
    File      string        // File path
    Line      int           // Line number
    Column    int           // Column number
    Context   string        // Error context
    Severity  ErrorSeverity // Error severity
    Inner     error         // Inner error
}
```

### 2. ErrorCode
Error code enumeration.

```go
type ErrorCode string

const (
    ErrCodeInvalidSyntax      ErrorCode = "INVALID_SYNTAX"
    ErrCodeMissingAnnotation  ErrorCode = "MISSING_ANNOTATION"
    ErrCodeInvalidParameter   ErrorCode = "INVALID_PARAMETER"
    ErrCodeDuplicateRoute     ErrorCode = "DUPLICATE_ROUTE"
    ErrCodeImportResolution   ErrorCode = "IMPORT_RESOLUTION"
    ErrCodeValidation         ErrorCode = "VALIDATION_ERROR"
    ErrCodeTemplateError      ErrorCode = "TEMPLATE_ERROR"
)
```

### 3. Diagnostic
Rich diagnostic information.

```go
type Diagnostic struct {
    Level     DiagnosticLevel // Diagnostic level
    Code      ErrorCode       // Error code
    Message   string          // Diagnostic message
    File      string          // File path
    Location  SourceLocation  // Source location
    Context   string          // Additional context
    Suggestions []string      // Suggested fixes
}
```

## Summary

This data model design provides a comprehensive foundation for the refactored route generation system. Key features include:

1. **Clear separation of concerns**: Each component has well-defined interfaces and responsibilities
2. **Comprehensive error handling**: Structured error types with rich diagnostic information
3. **Extensible validation**: Configurable validation system with multiple levels
4. **Type safety**: Strong typing throughout the system
5. **Configuration management**: Flexible configuration system for different use cases
6. **Backward compatibility**: Designed to support existing annotation formats

The design follows SOLID principles and provides a solid foundation for implementing the refactored route generation system.