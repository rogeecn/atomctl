package provider

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"go.ipao.vip/atomctl/v2/pkg/utils/gomod"
)

// Parser defines the interface for parsing provider annotations
type Parser interface {
	// ParseFile parses a single Go file and returns providers found
	ParseFile(filePath string) ([]Provider, error)

	// ParseDir parses all Go files in a directory and returns providers found
	ParseDir(dirPath string) ([]Provider, error)

	// ParseString parses Go code from a string and returns providers found
	ParseString(code string) ([]Provider, error)

	// SetConfig sets the parser configuration
	SetConfig(config *ParserConfig)

	// GetConfig returns the current parser configuration
	GetConfig() *ParserConfig

	// GetContext returns the current parser context
	GetContext() *ParserContext
}

// GoParser implements the Parser interface for Go source files
type GoParser struct {
	config  *ParserConfig
	context *ParserContext
	mu      sync.RWMutex
}

// NewGoParser creates a new GoParser with default configuration
func NewGoParser() *GoParser {
	config := NewParserConfig()
	context := NewParserContext(config)

	// Initialize file set if not provided
	if config.FileSet == nil {
		config.FileSet = token.NewFileSet()
		context.FileSet = config.FileSet
	}

	return &GoParser{
		config:  config,
		context: context,
	}
}

// NewGoParserWithConfig creates a new GoParser with custom configuration
func NewGoParserWithConfig(config *ParserConfig) *GoParser {
	if config == nil {
		return NewGoParser()
	}

	context := NewParserContext(config)

	// Initialize file set if not provided
	if config.FileSet == nil {
		config.FileSet = token.NewFileSet()
		context.FileSet = config.FileSet
	}

	return &GoParser{
		config:  config,
		context: context,
	}
}

// ParseFile implements Parser.ParseFile
func (p *GoParser) ParseFile(filePath string) ([]Provider, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	// Check if file should be included
	if !p.context.ShouldIncludeFile(filePath) {
		p.context.FilesSkipped++
		return []Provider{}, nil
	}

	// Check cache if enabled
	if p.config.CacheEnabled {
		if cached, found := p.context.Cache[filePath]; found {
			if providers, ok := cached.([]Provider); ok {
				return providers, nil
			}
		}
	}

	// Parse the file
	node, err := parser.ParseFile(p.context.FileSet, filePath, nil, p.config.Mode)
	if err != nil {
		p.context.AddError(filePath, 0, 0, err.Error(), "error")
		return nil, err
	}

	// Parse providers from the file
	providers, err := p.parseFileContent(filePath, node)
	if err != nil {
		return nil, err
	}

	// Cache the result
	if p.config.CacheEnabled {
		p.context.Cache[filePath] = providers
	}

	p.context.FilesProcessed++
	p.context.ProvidersFound += len(providers)

	return providers, nil
}

// ParseDir implements Parser.ParseDir
func (p *GoParser) ParseDir(dirPath string) ([]Provider, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	var allProviders []Provider

	// Walk through directory
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			// Skip hidden directories and common build/dependency directories
			if strings.HasPrefix(info.Name(), ".") ||
				info.Name() == "node_modules" ||
				info.Name() == "vendor" ||
				info.Name() == "testdata" {
				return filepath.SkipDir
			}
			return nil
		}

		// Parse Go files
		if filepath.Ext(path) == ".go" && p.context.ShouldIncludeFile(path) {
			providers, err := p.ParseFile(path)
			if err != nil {
				log.Warnf("Failed to parse file %s: %v", path, err)
				// Continue with other files
				return nil
			}
			allProviders = append(allProviders, providers...)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return allProviders, nil
}

// ParseString implements Parser.ParseString
func (p *GoParser) ParseString(code string) ([]Provider, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	// Parse the code string
	node, err := parser.ParseFile(p.context.FileSet, "", strings.NewReader(code), p.config.Mode)
	if err != nil {
		return nil, err
	}

	// Parse providers from the AST
	return p.parseFileContent("<string>", node)
}

// SetConfig implements Parser.SetConfig
func (p *GoParser) SetConfig(config *ParserConfig) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.config = config
	p.context = NewParserContext(config)
}

// GetConfig implements Parser.GetConfig
func (p *GoParser) GetConfig() *ParserConfig {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.config
}

// GetContext implements Parser.GetContext
func (p *GoParser) GetContext() *ParserContext {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.context
}

// parseFileContent parses providers from a parsed AST node
func (p *GoParser) parseFileContent(filePath string, node *ast.File) ([]Provider, error) {
	// Extract imports
	imports := make(map[string]string)
	for _, imp := range node.Imports {
		name := ""
		pkgPath := strings.Trim(imp.Path.Value, "\"")

		if imp.Name != nil {
			name = imp.Name.Name
		} else {
			name = gomod.GetPackageModuleName(pkgPath)
		}

		// Handle anonymous imports
		if name == "_" {
			name = gomod.GetPackageModuleName(pkgPath)
			// Handle duplicates
			if _, ok := imports[name]; ok {
				name = fmt.Sprintf("%s%d", name, rand.Intn(100))
			}
		}

		imports[name] = pkgPath
		p.context.AddImport(name, pkgPath)
	}

	var providers []Provider

	// Parse providers from declarations
	for _, decl := range node.Decls {
		provider, err := p.parseProviderDecl(filePath, node, decl, imports)
		if err != nil {
			p.context.AddError(filePath, 0, 0, err.Error(), "warning")
			continue
		}

		if provider != nil {
			providers = append(providers, *provider)
		}
	}

	return providers, nil
}

// parseProviderDecl parses a provider from an AST declaration
func (p *GoParser) parseProviderDecl(filePath string, fileNode *ast.File, decl ast.Decl, imports map[string]string) (*Provider, error) {
	genDecl, ok := decl.(*ast.GenDecl)
	if !ok {
		return nil, nil
	}

	if len(genDecl.Specs) == 0 {
		return nil, nil
	}

	typeSpec, ok := genDecl.Specs[0].(*ast.TypeSpec)
	if !ok {
		return nil, nil
	}

	// Check if it's a struct type
	structType, ok := typeSpec.Type.(*ast.StructType)
	if !ok {
		return nil, nil
	}

	// Check for provider annotation
	if genDecl.Doc == nil || len(genDecl.Doc.List) == 0 {
		return nil, nil
	}

	docText := strings.TrimLeft(genDecl.Doc.List[len(genDecl.Doc.List)-1].Text, "/ \t")
	if !strings.HasPrefix(docText, "@provider") {
		return nil, nil
	}

	// Parse provider annotation
	providerDoc := parseProvider(docText)

	// Create provider struct
	provider := &Provider{
		StructName:    typeSpec.Name.Name,
		ReturnType:    providerDoc.ReturnType,
		Mode:          ProviderModeBasic,
		ProviderGroup: providerDoc.Group,
		InjectParams:  make(map[string]InjectParam),
		Imports:       make(map[string]string),
		PkgName:       fileNode.Name.Name,
		ProviderFile:  filepath.Join(filepath.Dir(filePath), "provider.gen.go"),
	}

	// Set default return type if not specified
	if provider.ReturnType == "" {
		provider.ReturnType = "*" + provider.StructName
	}

	// Parse provider mode
	if providerDoc.Mode != "" {
		if IsValidProviderMode(providerDoc.Mode) {
			provider.Mode = ProviderMode(providerDoc.Mode)
		} else {
			return nil, fmt.Errorf("invalid provider mode: %s", providerDoc.Mode)
		}
	}

	// Parse struct fields for injection
	if err := p.parseStructFields(structType, imports, provider, providerDoc.IsOnly); err != nil {
		return nil, err
	}

	// Handle special provider modes
	p.handleProviderModes(provider, providerDoc.Mode)

	// Add source location if enabled
	if p.config.SourceLocations {
		if genDecl.Doc != nil && len(genDecl.Doc.List) > 0 {
			position := p.context.FileSet.Position(genDecl.Doc.List[0].Pos())
			provider.Location = SourceLocation{
				File:   position.Filename,
				Line:   position.Line,
				Column: position.Column,
			}
		}
	}

	return provider, nil
}

// parseStructFields parses struct fields for injection parameters
func (p *GoParser) parseStructFields(structType *ast.StructType, imports map[string]string, provider *Provider, onlyMode bool) error {
	for _, field := range structType.Fields.List {
		if field.Names == nil {
			continue
		}

		// Check for struct tags
		if field.Tag != nil {
			provider.NeedPrepareFunc = true
		}

		// Check injection mode
		shouldInject := true
		if onlyMode {
			shouldInject = field.Tag != nil && strings.Contains(field.Tag.Value, `inject:"true"`)
		} else {
			shouldInject = field.Tag == nil || !strings.Contains(field.Tag.Value, `inject:"false"`)
		}

		if !shouldInject {
			continue
		}

		// Parse field type
		star, pkg, pkgAlias, typ, err := p.parseFieldType(field.Type, imports)
		if err != nil {
			continue
		}

		// Skip scalar types
		if lo.Contains(scalarTypes, typ) {
			continue
		}

		// Add injection parameter
		for _, name := range field.Names {
			provider.InjectParams[name.Name] = InjectParam{
				Star:         star,
				Type:         typ,
				Package:      pkg,
				PackageAlias: pkgAlias,
			}

			// Add to imports
			if pkg != "" && pkgAlias != "" {
				provider.Imports[pkg] = pkgAlias
			}
		}
	}

	return nil
}

// parseFieldType parses a field type and returns its components
func (p *GoParser) parseFieldType(expr ast.Expr, imports map[string]string) (star, pkg, pkgAlias, typ string, err error) {
	switch t := expr.(type) {
	case *ast.Ident:
		typ = t.Name
	case *ast.StarExpr:
		star = "*"
		return p.parseFieldType(t.X, imports)
	case *ast.SelectorExpr:
		if x, ok := t.X.(*ast.Ident); ok {
			pkgAlias = x.Name
			if path, ok := imports[pkgAlias]; ok {
				pkg = path
			}
			typ = t.Sel.Name
		}
	default:
		return "", "", "", "", errors.New("unsupported field type")
	}

	return star, pkg, pkgAlias, typ, nil
}

// handleProviderModes applies special handling for different provider modes
func (p *GoParser) handleProviderModes(provider *Provider, mode string) {
	moduleName := gomod.GetModuleName()

	switch provider.Mode {
	case ProviderModeGrpc:
		modePkg := moduleName + "/providers/grpc"
		provider.ProviderGroup = "atom.GroupInitial"
		provider.GrpcRegisterFunc = provider.ReturnType
		provider.ReturnType = "contracts.Initial"

		provider.Imports[atomPackage("")] = ""
		provider.Imports[atomPackage("contracts")] = ""
		provider.Imports[modePkg] = ""

		provider.InjectParams["__grpc"] = InjectParam{
			Star:         "*",
			Type:         "Grpc",
			Package:      modePkg,
			PackageAlias: "grpc",
		}

	case ProviderModeEvent:
		modePkg := moduleName + "/providers/event"
		provider.ProviderGroup = "atom.GroupInitial"
		provider.ReturnType = "contracts.Initial"

		provider.Imports[atomPackage("")] = ""
		provider.Imports[atomPackage("contracts")] = ""
		provider.Imports[modePkg] = ""

		provider.InjectParams["__event"] = InjectParam{
			Star:         "*",
			Type:         "PubSub",
			Package:      modePkg,
			PackageAlias: "event",
		}

	case ProviderModeJob, ProviderModeCronJob:
		modePkg := moduleName + "/providers/job"
		provider.ProviderGroup = "atom.GroupInitial"
		provider.ReturnType = "contracts.Initial"

		provider.Imports[atomPackage("")] = ""
		provider.Imports[atomPackage("contracts")] = ""
		provider.Imports["github.com/riverqueue/river"] = ""
		provider.Imports[modePkg] = ""

		provider.InjectParams["__job"] = InjectParam{
			Star:         "*",
			Type:         "Job",
			Package:      modePkg,
			PackageAlias: "job",
		}

	case ProviderModeModel:
		provider.ProviderGroup = "atom.GroupInitial"
		provider.ReturnType = "contracts.Initial"
		provider.NeedPrepareFunc = true
	}

	// Handle return type and group package imports
	if pkgAlias := getTypePkgName(provider.ReturnType); pkgAlias != "" {
		if importPkg, ok := p.context.Imports[pkgAlias]; ok {
			provider.Imports[importPkg] = pkgAlias
		}
	}

	if pkgAlias := getTypePkgName(provider.ProviderGroup); pkgAlias != "" {
		if importPkg, ok := p.context.Imports[pkgAlias]; ok {
			provider.Imports[importPkg] = pkgAlias
		}
	}
}
