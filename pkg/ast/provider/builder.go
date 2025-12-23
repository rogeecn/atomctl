package provider

import (
	"fmt"
	"go/ast"
	"strings"
)

// ProviderBuilder handles the construction of Provider objects from parsed AST components
type ProviderBuilder struct {
	config         *BuilderConfig
	commentParser  *CommentParser
	importResolver *ImportResolver
	astWalker      *ASTWalker
}

// BuilderConfig configures the provider builder behavior
type BuilderConfig struct {
	EnableValidation          bool
	StrictMode                bool
	DefaultProviderMode       ProviderMode
	DefaultInjectionMode      InjectionMode
	AutoGenerateReturnTypes   bool
	ResolveImportDependencies bool
}

// BuilderContext maintains context during provider building
type BuilderContext struct {
	FilePath       string
	PackageName    string
	ImportContext  *ImportContext
	ASTFile        *ast.File
	ProcessedTypes map[string]bool
	Errors         []error
	Warnings       []string
}

// NewProviderBuilder creates a new ProviderBuilder with default configuration
func NewProviderBuilder() *ProviderBuilder {
	return &ProviderBuilder{
		config: &BuilderConfig{
			EnableValidation:          true,
			StrictMode:                false,
			DefaultProviderMode:       ProviderModeBasic,
			DefaultInjectionMode:      InjectionModeAuto,
			AutoGenerateReturnTypes:   true,
			ResolveImportDependencies: true,
		},
		commentParser:  NewCommentParser(),
		importResolver: NewImportResolver(),
		astWalker:      NewASTWalker(),
	}
}

// NewProviderBuilderWithConfig creates a new ProviderBuilder with custom configuration
func NewProviderBuilderWithConfig(config *BuilderConfig) *ProviderBuilder {
	if config == nil {
		return NewProviderBuilder()
	}

	return &ProviderBuilder{
		config:         config,
		commentParser:  NewCommentParser(),
		importResolver: NewImportResolver(),
		astWalker: NewASTWalkerWithConfig(&WalkerConfig{
			StrictMode: config.StrictMode,
		}),
	}
}

// BuildFromTypeSpec builds a Provider from a type specification and its declaration
func (pb *ProviderBuilder) BuildFromTypeSpec(typeSpec *ast.TypeSpec, decl *ast.GenDecl, context *BuilderContext) (Provider, error) {
	if typeSpec == nil {
		return Provider{}, fmt.Errorf("type specification cannot be nil")
	}

	// Initialize builder context if not provided
	if context == nil {
		context = &BuilderContext{
			ProcessedTypes: make(map[string]bool),
			Errors:         make([]error, 0),
			Warnings:       make([]string, 0),
		}
	}

	// Check if type has already been processed
	if context.ProcessedTypes[typeSpec.Name.Name] {
		return Provider{}, fmt.Errorf("type %s has already been processed", typeSpec.Name.Name)
	}

	// Parse provider comment
	providerComment, err := pb.parseProviderComment(decl)
	if err != nil {
		return Provider{}, fmt.Errorf("failed to parse provider comment: %w", err)
	}

	// Create basic provider structure
	provider := Provider{
		StructName:    typeSpec.Name.Name,
		Mode:          pb.determineProviderMode(providerComment),
		ProviderGroup: pb.determineProviderGroup(providerComment),
		InjectParams:  make(map[string]InjectParam),
		Imports:       make(map[string]string),
		PkgName:       context.PackageName,
		ProviderFile:  context.FilePath,
		Location: SourceLocation{
			File: context.FilePath,
		},
	}

	// Set return type
	if err := pb.setReturnType(&provider, providerComment, typeSpec); err != nil {
		return Provider{}, err
	}

	// Process struct fields if it's a struct type
	if structType, ok := typeSpec.Type.(*ast.StructType); ok {
		if err := pb.processStructFields(&provider, structType, context); err != nil {
			return Provider{}, err
		}
	}

	// Resolve import dependencies
	if pb.config.ResolveImportDependencies {
		if err := pb.resolveImportDependencies(&provider, context); err != nil {
			return Provider{}, err
		}
	}

	// Apply mode-specific configurations
	if err := pb.applyModeSpecificConfig(&provider, providerComment); err != nil {
		return Provider{}, err
	}

	// Validate the built provider
	if pb.config.EnableValidation {
		if err := pb.validateProvider(&provider); err != nil {
			return Provider{}, err
		}
	}

	// Mark type as processed
	context.ProcessedTypes[typeSpec.Name.Name] = true

	return provider, nil
}

// BuildFromComment builds a Provider from a provider comment string
func (pb *ProviderBuilder) BuildFromComment(comment string, context *BuilderContext) (Provider, error) {
	// Parse the provider comment
	providerComment, err := pb.commentParser.ParseProviderComment(comment)
	if err != nil {
		return Provider{}, fmt.Errorf("failed to parse provider comment: %w", err)
	}

	// Create basic provider structure
	provider := Provider{
		Mode:          pb.determineProviderMode(providerComment),
		ProviderGroup: pb.determineProviderGroup(providerComment),
		InjectParams:  make(map[string]InjectParam),
		Imports:       make(map[string]string),
	}

	// Set return type from comment
	if providerComment.ReturnType != "" {
		provider.ReturnType = providerComment.ReturnType
	} else if pb.config.AutoGenerateReturnTypes {
		// Generate a default return type based on mode
		provider.ReturnType = pb.generateDefaultReturnType(provider.Mode)
	}

	// Apply mode-specific configurations
	if err := pb.applyModeSpecificConfig(&provider, providerComment); err != nil {
		return Provider{}, err
	}

	// Validate the built provider
	if pb.config.EnableValidation {
		if err := pb.validateProvider(&provider); err != nil {
			return Provider{}, err
		}
	}

	return provider, nil
}

// parseProviderComment parses the provider comment from a declaration
func (pb *ProviderBuilder) parseProviderComment(decl *ast.GenDecl) (*ProviderComment, error) {
	if decl.Doc == nil || len(decl.Doc.List) == 0 {
		return nil, fmt.Errorf("no documentation found for declaration")
	}

	// Extract comment lines
	commentLines := make([]string, len(decl.Doc.List))
	for i, comment := range decl.Doc.List {
		commentLines[i] = comment.Text
	}

	// Parse provider annotation
	return pb.commentParser.ParseCommentBlock(commentLines)
}

// determineProviderMode determines the provider mode from the comment
func (pb *ProviderBuilder) determineProviderMode(comment *ProviderComment) ProviderMode {
	if comment != nil && comment.Mode != "" {
		return comment.Mode
	}
	return pb.config.DefaultProviderMode
}

// determineProviderGroup determines the provider group from the comment
func (pb *ProviderBuilder) determineProviderGroup(comment *ProviderComment) string {
	if comment != nil && comment.Group != "" {
		return comment.Group
	}
	return ""
}

// setReturnType sets the return type for the provider
func (pb *ProviderBuilder) setReturnType(provider *Provider, comment *ProviderComment, typeSpec *ast.TypeSpec) error {
	if comment != nil && comment.ReturnType != "" {
		provider.ReturnType = comment.ReturnType
		return nil
	}

	if pb.config.AutoGenerateReturnTypes {
		provider.ReturnType = pb.generateDefaultReturnType(provider.Mode)
		return nil
	}

	// Default to pointer type
	provider.ReturnType = "*" + typeSpec.Name.Name
	return nil
}

// generateDefaultReturnType generates a default return type based on provider mode
func (pb *ProviderBuilder) generateDefaultReturnType(mode ProviderMode) string {
	switch mode {
	case ProviderModeBasic:
		return "interface{}"
	case ProviderModeGrpc:
		return "interface{}"
	case ProviderModeEvent:
		return "func() error"
	case ProviderModeJob:
		return "func(ctx context.Context) error"
	case ProviderModeCronJob:
		return "func(ctx context.Context) error"
	case ProviderModeModel:
		return "interface{}"
	default:
		return "interface{}"
	}
}

// processStructFields processes struct fields to extract injection parameters
func (pb *ProviderBuilder) processStructFields(provider *Provider, structType *ast.StructType, context *BuilderContext) error {
	if structType.Fields == nil {
		return nil
	}

	for _, field := range structType.Fields.List {
		if len(field.Names) == 0 {
			// Skip anonymous fields
			continue
		}

		for _, fieldName := range field.Names {
			injectParam, err := pb.createInjectParam(fieldName.Name, field, context)
			if err != nil {
				return fmt.Errorf("failed to create inject param for field %s: %w", fieldName.Name, err)
			}

			provider.InjectParams[fieldName.Name] = *injectParam
		}
	}

	return nil
}

// createInjectParam creates an InjectParam from a field
func (pb *ProviderBuilder) createInjectParam(fieldName string, field *ast.Field, context *BuilderContext) (*InjectParam, error) {
	// Extract field type
	typeStr := pb.extractFieldType(field)
	if typeStr == "" {
		return nil, fmt.Errorf("cannot determine type for field %s", fieldName)
	}

	// Check for import dependencies
	packagePath, packageAlias := pb.extractImportInfo(typeStr, context)

	return &InjectParam{
		Type:         typeStr,
		Package:      packagePath,
		PackageAlias: packageAlias,
	}, nil
}

// extractFieldType extracts the type string from a field
func (pb *ProviderBuilder) extractFieldType(field *ast.Field) string {
	if field.Type == nil {
		return ""
	}

	// Handle different type representations
	switch t := field.Type.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		// Handle qualified identifiers like "package.Type"
		if x, ok := t.X.(*ast.Ident); ok {
			return x.Name + "." + t.Sel.Name
		}
	case *ast.StarExpr:
		// Handle pointer types
		if x, ok := t.X.(*ast.Ident); ok {
			return "*" + x.Name
		}
	case *ast.ArrayType:
		// Handle array/slice types
		if x, ok := t.Elt.(*ast.Ident); ok {
			return "[]" + x.Name
		}
	case *ast.MapType:
		// Handle map types
		keyType := pb.extractTypeFromExpr(t.Key)
		valueType := pb.extractTypeFromExpr(t.Value)
		if keyType != "" && valueType != "" {
			return fmt.Sprintf("map[%s]%s", keyType, valueType)
		}
	}

	return ""
}

// extractTypeFromExpr extracts type string from an expression
func (pb *ProviderBuilder) extractTypeFromExpr(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		if x, ok := t.X.(*ast.Ident); ok {
			return x.Name + "." + t.Sel.Name
		}
	}
	return ""
}

// extractImportInfo extracts import information from a type string
func (pb *ProviderBuilder) extractImportInfo(typeStr string, context *BuilderContext) (packagePath, packageAlias string) {
	if !strings.Contains(typeStr, ".") {
		return "", ""
	}

	// Extract the package part
	parts := strings.Split(typeStr, ".")
	if len(parts) < 2 {
		return "", ""
	}

	packageAlias = parts[0]

	// Look up the import path from the import context
	if context.ImportContext != nil {
		if path, exists := context.ImportContext.ImportPaths[packageAlias]; exists {
			return path, packageAlias
		}
	}

	return "", packageAlias
}

// resolveImportDependencies resolves import dependencies for the provider
func (pb *ProviderBuilder) resolveImportDependencies(provider *Provider, context *BuilderContext) error {
	if context.ImportContext == nil {
		return nil
	}

	// Add imports from injection parameters
	for _, param := range provider.InjectParams {
		if param.Package != "" && param.PackageAlias != "" {
			provider.Imports[param.PackageAlias] = param.Package
		}
	}

	// For gRPC mode, extract and add imports from the original file's imports
	if provider.Mode == ProviderModeGrpc && provider.GrpcRegisterFunc != "" {
		// Extract package alias from gRPC register function name (e.g., "userv1" from "userv1.RegisterUserServiceServer")
		if pkgAlias := getTypePkgName(provider.GrpcRegisterFunc); pkgAlias != "" {
			// Look for this package in the original file's imports
			if importResolution, exists := context.ImportContext.FileImports[pkgAlias]; exists {
				// Add the import from the original file
				provider.Imports[importResolution.Path] = pkgAlias
			}
		}
	}

	// Add mode-specific imports
	modeImports := pb.getModeSpecificImports(provider.Mode)
	for alias, path := range modeImports {
		// Check for conflicts
		if existingPath, exists := provider.Imports[alias]; exists && existingPath != path {
			// Handle conflict by generating unique alias
			uniqueAlias := pb.generateUniqueAlias(alias, provider.Imports)
			provider.Imports[uniqueAlias] = path
		} else {
			provider.Imports[alias] = path
		}
	}

	return nil
}

// getModeSpecificImports returns mode-specific import requirements
func (pb *ProviderBuilder) getModeSpecificImports(mode ProviderMode) map[string]string {
	imports := make(map[string]string)

	switch mode {
	case ProviderModeGrpc:
		imports["grpc"] = "google.golang.org/grpc"
	case ProviderModeEvent:
		imports["context"] = "context"
	case ProviderModeJob:
		imports["context"] = "context"
	case ProviderModeCronJob:
		imports["context"] = "context"
	case ProviderModeModel:
		imports["encoding/json"] = "encoding/json"
	}

	return imports
}

// generateUniqueAlias generates a unique alias to avoid conflicts
func (pb *ProviderBuilder) generateUniqueAlias(baseAlias string, existingImports map[string]string) string {
	for i := 1; i < 1000; i++ {
		candidate := fmt.Sprintf("%s%d", baseAlias, i)
		if _, exists := existingImports[candidate]; !exists {
			return candidate
		}
	}
	return baseAlias
}

// applyModeSpecificConfig applies mode-specific configurations to the provider
func (pb *ProviderBuilder) applyModeSpecificConfig(provider *Provider, comment *ProviderComment) error {
	switch provider.Mode {
	case ProviderModeGrpc:
		provider.GrpcRegisterFunc = pb.generateGrpcRegisterFuncName(provider.StructName)
		provider.NeedPrepareFunc = true
	case ProviderModeEvent:
		provider.NeedPrepareFunc = true
	case ProviderModeJob:
		provider.NeedPrepareFunc = true
	case ProviderModeCronJob:
		provider.NeedPrepareFunc = true
	case ProviderModeModel:
		provider.NeedPrepareFunc = true
	default:
		// Basic mode - no special configuration
	}

	return nil
}

// generateGrpcRegisterFuncName generates a gRPC register function name
func (pb *ProviderBuilder) generateGrpcRegisterFuncName(structName string) string {
	// Convert struct name to register function name
	// Example: UserService -> RegisterUserServiceServer
	return "Register" + structName + "Server"
}

// validateProvider validates the constructed provider
func (pb *ProviderBuilder) validateProvider(provider *Provider) error {
	// Basic validation
	if provider.StructName == "" {
		return fmt.Errorf("provider struct name cannot be empty")
	}

	if provider.ReturnType == "" {
		return fmt.Errorf("provider return type cannot be empty")
	}

	if !IsValidProviderMode(string(provider.Mode)) {
		return fmt.Errorf("invalid provider mode: %s", provider.Mode)
	}

	// Validate injection parameters
	for name, param := range provider.InjectParams {
		if param.Type == "" {
			return fmt.Errorf("injection parameter type cannot be empty for %s", name)
		}
	}

	return nil
}

// GetConfig returns the builder configuration
func (pb *ProviderBuilder) GetConfig() *BuilderConfig {
	return pb.config
}

// SetConfig updates the builder configuration
func (pb *ProviderBuilder) SetConfig(config *BuilderConfig) {
	pb.config = config
}

// GetCommentParser returns the comment parser used by the builder
func (pb *ProviderBuilder) GetCommentParser() *CommentParser {
	return pb.commentParser
}

// GetImportResolver returns the import resolver used by the builder
func (pb *ProviderBuilder) GetImportResolver() *ImportResolver {
	return pb.importResolver
}

// GetASTWalker returns the AST walker used by the builder
func (pb *ProviderBuilder) GetASTWalker() *ASTWalker {
	return pb.astWalker
}
