package provider

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"go.ipao.vip/atomctl/v2/pkg/utils/gomod"
)

// MainParser represents the main parser that uses extracted components
type MainParser struct {
	commentParser  *CommentParser
	importResolver *ImportResolver
	astWalker      *ASTWalker
	builder        *ProviderBuilder
	validator      *GoValidator
	config         *ParserConfig
}

// NewParser creates a new MainParser with default configuration
func NewParser() *MainParser {
	return &MainParser{
		commentParser:  NewCommentParser(),
		importResolver: NewImportResolver(),
		astWalker:      NewASTWalker(),
		builder:        NewProviderBuilder(),
		validator:      NewGoValidator(),
		config:         NewParserConfig(),
	}
}

// NewParserWithConfig creates a new MainParser with custom configuration
func NewParserWithConfig(config *ParserConfig) *MainParser {
	if config == nil {
		return NewParser()
	}

	return &MainParser{
		commentParser:  NewCommentParser(),
		importResolver: NewImportResolver(),
		astWalker: NewASTWalkerWithConfig(&WalkerConfig{
			StrictMode: config.StrictMode,
		}),
		builder: NewProviderBuilderWithConfig(&BuilderConfig{
			EnableValidation:          config.StrictMode,
			StrictMode:                config.StrictMode,
			DefaultProviderMode:       ProviderModeBasic,
			DefaultInjectionMode:      InjectionModeAuto,
			AutoGenerateReturnTypes:   true,
			ResolveImportDependencies: true,
		}),
		validator: NewGoValidator(),
		config:    config,
	}
}

// Parse parses a Go source file and returns discovered providers
// This is the refactored version of the original Parse function
func ParseRefactored(source string) []Provider {
	parser := NewParser()
	providers, err := parser.ParseFile(source)
	if err != nil {
		log.Error("Parse error: ", err)
		return []Provider{}
	}
	return providers
}

// ParseFile parses a single Go source file and returns discovered providers
func (p *MainParser) ParseFile(source string) ([]Provider, error) {
	// Check if file should be processed
	if !p.shouldProcessFile(source) {
		return []Provider{}, nil
	}

	// Parse the AST
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, source, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file %s: %w", source, err)
	}

	// Create parser context
	context := NewParserContext(p.config)
	context.WorkingDir = filepath.Dir(source)
	context.ModuleName = gomod.GetModuleName()

	// Resolve imports
	importContext, err := p.importResolver.ResolveFileImports(node, source)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve imports: %w", err)
	}

	// Create builder context
	builderContext := &BuilderContext{
		FilePath:       source,
		PackageName:    node.Name.Name,
		ImportContext:  importContext,
		ASTFile:        node,
		ProcessedTypes: make(map[string]bool),
		Errors:         make([]error, 0),
		Warnings:       make([]string, 0),
	}

	// Use AST walker to find provider annotations
	visitor := NewProviderDiscoveryVisitor(p.commentParser)
	p.astWalker.AddVisitor(visitor)

	// Walk the AST
	if err := p.astWalker.WalkFile(source); err != nil {
		return nil, fmt.Errorf("failed to walk AST: %w", err)
	}

	// Build providers from discovered annotations
	providers := make([]Provider, 0)
	discoveredProviders := visitor.GetProviders()

	for _, discoveredProvider := range discoveredProviders {
		// Find the corresponding AST node for this provider
		provider, err := p.buildProviderFromDiscovery(discoveredProvider, node, builderContext)
		if err != nil {
			context.AddError(source, 0, 0, fmt.Sprintf("failed to build provider %s: %v", discoveredProvider.StructName, err), "error")
			continue
		}

		// Validate the provider if enabled
		if p.config.StrictMode {
			if err := p.validator.Validate(&provider); err != nil {
				context.AddError(source, 0, 0, fmt.Sprintf("validation failed for provider %s: %v", provider.StructName, err), "error")
				continue
			}
		}

		providers = append(providers, provider)
	}

	// Log any warnings or errors
	for _, parseErr := range context.GetErrors("warning") {
		log.Warnf("Warning while parsing %s: %s", source, parseErr.Message)
	}
	for _, parseErr := range context.GetErrors("error") {
		log.Errorf("Error while parsing %s: %s", source, parseErr.Message)
	}

	return providers, nil
}

// ParseDir parses all Go files in a directory and returns discovered providers
func (p *MainParser) ParseDir(dir string) ([]Provider, error) {
	var allProviders []Provider

	// Use AST walker to traverse the directory
	if err := p.astWalker.WalkDir(dir); err != nil {
		return nil, fmt.Errorf("failed to walk directory: %w", err)
	}

	// Note: This would need to be enhanced to collect providers from all files
	// For now, we'll return an empty slice and log a warning
	log.Warn("ParseDir not fully implemented yet")
	return allProviders, nil
}

// shouldProcessFile determines if a file should be processed
func (p *MainParser) shouldProcessFile(source string) bool {
	// Skip test files
	if strings.HasSuffix(source, "_test.go") {
		return false
	}

	// Skip generated provider files
	if strings.HasSuffix(source, "provider.gen.go") {
		return false
	}

	return true
}

// buildProviderFromDiscovery builds a complete Provider from a discovered provider annotation
func (p *MainParser) buildProviderFromDiscovery(discoveredProvider Provider, node *ast.File, context *BuilderContext) (Provider, error) {
	// Find the corresponding type specification in the AST
	var typeSpec *ast.TypeSpec
	var genDecl *ast.GenDecl

	for _, decl := range node.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}

		if len(gd.Specs) == 0 {
			continue
		}

		ts, ok := gd.Specs[0].(*ast.TypeSpec)
		if !ok {
			continue
		}

		if ts.Name.Name == discoveredProvider.StructName {
			typeSpec = ts
			genDecl = gd
			break
		}
	}

	if typeSpec == nil {
		return Provider{}, fmt.Errorf("type specification not found for %s", discoveredProvider.StructName)
	}

	// Use the builder to construct the complete provider
	provider, err := p.builder.BuildFromTypeSpec(typeSpec, genDecl, context)
	if err != nil {
		return Provider{}, fmt.Errorf("failed to build provider: %w", err)
	}

	// Apply legacy compatibility transformations
	result, err := p.applyLegacyCompatibility(&provider)
	if err != nil {
		return Provider{}, fmt.Errorf("failed to apply legacy compatibility: %w", err)
	}

	return result, nil
}

// applyLegacyCompatibility applies transformations to maintain backward compatibility
func (p *MainParser) applyLegacyCompatibility(provider *Provider) (Provider, error) {
	// Set provider file path
	provider.ProviderFile = filepath.Join(filepath.Dir(provider.ProviderFile), "provider.gen.go")

	// Apply mode-specific transformations based on the original logic
	switch provider.Mode {
	case ProviderModeGrpc:
		p.applyGrpcCompatibility(provider)
	case ProviderModeEvent:
		p.applyEventCompatibility(provider)
	case ProviderModeJob, ProviderModeCronJob:
		p.applyJobCompatibility(provider)
	case ProviderModeModel:
		p.applyModelCompatibility(provider)
	}

	return *provider, nil
}

// applyGrpcCompatibility applies gRPC-specific compatibility transformations
func (p *MainParser) applyGrpcCompatibility(provider *Provider) {
	modePkg := gomod.GetModuleName() + "/providers/grpc"

	// Add required imports
	provider.Imports[createAtomPackage("")] = ""
	provider.Imports[createAtomPackage("contracts")] = ""
	provider.Imports[modePkg] = ""

	// Set provider group
	if provider.ProviderGroup == "" {
		provider.ProviderGroup = "atom.GroupInitial"
	}

	// Set return type and register function
	if provider.GrpcRegisterFunc == "" {
		provider.GrpcRegisterFunc = provider.ReturnType
	}
	provider.ReturnType = "contracts.Initial"

	// Add gRPC injection parameter
	provider.InjectParams["__grpc"] = InjectParam{
		Star:         "*",
		Type:         "Grpc",
		Package:      modePkg,
		PackageAlias: "grpc",
	}
}

// applyEventCompatibility applies event-specific compatibility transformations
func (p *MainParser) applyEventCompatibility(provider *Provider) {
	modePkg := gomod.GetModuleName() + "/providers/event"

	// Add required imports
	provider.Imports[createAtomPackage("")] = ""
	provider.Imports[createAtomPackage("contracts")] = ""
	provider.Imports[modePkg] = ""

	// Set provider group
	if provider.ProviderGroup == "" {
		provider.ProviderGroup = "atom.GroupInitial"
	}

	// Set return type
	provider.ReturnType = "contracts.Initial"

	// Add event injection parameter
	provider.InjectParams["__event"] = InjectParam{
		Star:         "*",
		Type:         "PubSub",
		Package:      modePkg,
		PackageAlias: "event",
	}
}

// applyJobCompatibility applies job-specific compatibility transformations
func (p *MainParser) applyJobCompatibility(provider *Provider) {
	modePkg := gomod.GetModuleName() + "/providers/job"

	// Add required imports
	provider.Imports[createAtomPackage("")] = ""
	provider.Imports[createAtomPackage("contracts")] = ""
	provider.Imports["github.com/riverqueue/river"] = ""
	provider.Imports[modePkg] = ""

	// Set provider group
	if provider.ProviderGroup == "" {
		provider.ProviderGroup = "atom.GroupInitial"
	}

	// Set return type
	provider.ReturnType = "contracts.Initial"

	// Add job injection parameter
	provider.InjectParams["__job"] = InjectParam{
		Star:         "*",
		Type:         "Job",
		Package:      modePkg,
		PackageAlias: "job",
	}
}

// applyModelCompatibility applies model-specific compatibility transformations
func (p *MainParser) applyModelCompatibility(provider *Provider) {
	// Set provider group
	if provider.ProviderGroup == "" {
		provider.ProviderGroup = "atom.GroupInitial"
	}

	// Set return type
	provider.ReturnType = "contracts.Initial"

	// Ensure prepare function is needed
	provider.NeedPrepareFunc = true
}

// GetCommentParser returns the comment parser used by this parser
func (p *MainParser) GetCommentParser() *CommentParser {
	return p.commentParser
}

// GetImportResolver returns the import resolver used by this parser
func (p *MainParser) GetImportResolver() *ImportResolver {
	return p.importResolver
}

// GetASTWalker returns the AST walker used by this parser
func (p *MainParser) GetASTWalker() *ASTWalker {
	return p.astWalker
}

// GetBuilder returns the provider builder used by this parser
func (p *MainParser) GetBuilder() *ProviderBuilder {
	return p.builder
}

// GetValidator returns the validator used by this parser
func (p *MainParser) GetValidator() *GoValidator {
	return p.validator
}

// GetConfig returns the parser configuration
func (p *MainParser) GetConfig() *ParserConfig {
	return p.config
}

// SetConfig updates the parser configuration
func (p *MainParser) SetConfig(config *ParserConfig) {
	p.config = config
}

// Helper function to create atom package paths
func createAtomPackage(suffix string) string {
	root := "go.ipao.vip/atom"
	if suffix != "" {
		return fmt.Sprintf("%s/%s", root, suffix)
	}
	return root
}
