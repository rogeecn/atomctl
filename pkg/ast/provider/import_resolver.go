package provider

import (
	"fmt"
	"go/ast"
	"math/rand"
	"path/filepath"
	"strings"

	"go.ipao.vip/atomctl/v2/pkg/utils/gomod"
)

// ImportResolver handles resolution of Go imports and package aliases
type ImportResolver struct {
	resolverConfig *ResolverConfig
	cache          map[string]*ImportResolution
}

// ResolverConfig configures the import resolver behavior
type ResolverConfig struct {
	EnableCache             bool
	StrictMode              bool
	DefaultAliasStrategy    AliasStrategy
	AnonymousImportHandling AnonymousImportPolicy
}

// AliasStrategy defines how to generate default aliases
type AliasStrategy int

const (
	AliasStrategyModuleName AliasStrategy = iota
	AliasStrategyLastPath
	AliasStrategyCustom
)

// AnonymousImportPolicy defines how to handle anonymous imports
type AnonymousImportPolicy int

const (
	AnonymousImportSkip AnonymousImportPolicy = iota
	AnonymousImportUseModuleName
	AnonymousImportGenerateUnique
)

// ImportResolution represents a resolved import
type ImportResolution struct {
	Path         string            // Import path
	Alias        string            // Package alias
	IsAnonymous  bool              // Is this an anonymous import (_)
	IsValid      bool              // Is the import valid
	Error        string            // Error message if invalid
	PackageName  string            // Actual package name
	Dependencies map[string]string // Dependencies of this import
}

// ImportContext maintains context for import resolution
type ImportContext struct {
	FileImports    map[string]*ImportResolution // Alias -> Resolution
	ImportPaths    map[string]string            // Path -> Alias
	ModuleInfo     map[string]string            // Module path -> module name
	WorkingDir     string                       // Current working directory
	ModuleName     string                       // Current module name
	ProcessedFiles map[string]bool              // Track processed files
}

// NewImportResolver creates a new ImportResolver
func NewImportResolver() *ImportResolver {
	return &ImportResolver{
		resolverConfig: &ResolverConfig{
			EnableCache:             true,
			StrictMode:              false,
			DefaultAliasStrategy:    AliasStrategyModuleName,
			AnonymousImportHandling: AnonymousImportUseModuleName,
		},
		cache: make(map[string]*ImportResolution),
	}
}

// NewImportResolverWithConfig creates a new ImportResolver with custom configuration
func NewImportResolverWithConfig(config *ResolverConfig) *ImportResolver {
	if config == nil {
		return NewImportResolver()
	}

	return &ImportResolver{
		resolverConfig: config,
		cache:          make(map[string]*ImportResolution),
	}
}

// ResolveFileImports resolves all imports for a given AST file
func (ir *ImportResolver) ResolveFileImports(file *ast.File, filePath string) (*ImportContext, error) {
	context := &ImportContext{
		FileImports:    make(map[string]*ImportResolution),
		ImportPaths:    make(map[string]string),
		ModuleInfo:     make(map[string]string),
		WorkingDir:     filepath.Dir(filePath),
		ProcessedFiles: make(map[string]bool),
	}

	// Resolve current module name
	moduleName := gomod.GetModuleName()
	context.ModuleName = moduleName

	// Process imports
	for _, imp := range file.Imports {
		resolution, err := ir.resolveImportSpec(imp, context)
		if err != nil {
			if ir.resolverConfig.StrictMode {
				return nil, err
			}
			// In non-strict mode, continue with other imports
			continue
		}

		if resolution != nil {
			context.FileImports[resolution.Alias] = resolution
			context.ImportPaths[resolution.Path] = resolution.Alias
		}
	}

	return context, nil
}

// resolveImportSpec resolves a single import specification
func (ir *ImportResolver) resolveImportSpec(imp *ast.ImportSpec, context *ImportContext) (*ImportResolution, error) {
	// Extract import path
	path := strings.Trim(imp.Path.Value, "\"")
	if path == "" {
		return nil, fmt.Errorf("empty import path")
	}

	// Check cache first
	if ir.resolverConfig.EnableCache {
		if cached, found := ir.cache[path]; found {
			return cached, nil
		}
	}

	// Determine alias
	alias := ir.determineAlias(imp, path, context)

	// Resolve package name
	packageName, err := ir.resolvePackageName(path, context)
	if err != nil {
		resolution := &ImportResolution{
			Path:        path,
			Alias:       alias,
			IsAnonymous: imp.Name != nil && imp.Name.Name == "_",
			IsValid:     false,
			Error:       err.Error(),
			PackageName: "",
		}

		if ir.resolverConfig.EnableCache {
			ir.cache[path] = resolution
		}

		return resolution, err
	}

	// Create resolution
	resolution := &ImportResolution{
		Path:         path,
		Alias:        alias,
		IsAnonymous:  imp.Name != nil && imp.Name.Name == "_",
		IsValid:      true,
		PackageName:  packageName,
		Dependencies: make(map[string]string),
	}

	// Resolve dependencies if needed
	if err := ir.resolveDependencies(resolution, context); err != nil {
		resolution.IsValid = false
		resolution.Error = err.Error()
	}

	// Cache the result
	if ir.resolverConfig.EnableCache {
		ir.cache[path] = resolution
	}

	return resolution, nil
}

// determineAlias determines the appropriate alias for an import
func (ir *ImportResolver) determineAlias(imp *ast.ImportSpec, path string, context *ImportContext) string {
	// If explicit alias is provided, use it
	if imp.Name != nil {
		if imp.Name.Name == "_" {
			// Handle anonymous import based on policy
			return ir.handleAnonymousImport(path, context)
		}
		return imp.Name.Name
	}

	// Generate default alias based on strategy
	switch ir.resolverConfig.DefaultAliasStrategy {
	case AliasStrategyModuleName:
		return gomod.GetPackageModuleName(path)
	case AliasStrategyLastPath:
		return ir.getLastPathComponent(path)
	case AliasStrategyCustom:
		return ir.generateCustomAlias(path, context)
	default:
		return gomod.GetPackageModuleName(path)
	}
}

// handleAnonymousImport handles anonymous imports based on policy
func (ir *ImportResolver) handleAnonymousImport(path string, context *ImportContext) string {
	switch ir.resolverConfig.AnonymousImportHandling {
	case AnonymousImportSkip:
		return "_"
	case AnonymousImportUseModuleName:
		alias := gomod.GetPackageModuleName(path)
		// Check for conflicts
		if _, exists := context.FileImports[alias]; exists {
			return ir.generateUniqueAlias(alias, context)
		}
		return alias
	case AnonymousImportGenerateUnique:
		baseAlias := gomod.GetPackageModuleName(path)
		return ir.generateUniqueAlias(baseAlias, context)
	default:
		return "_"
	}
}

// resolvePackageName resolves the actual package name for an import path
func (ir *ImportResolver) resolvePackageName(path string, context *ImportContext) (string, error) {
	// Handle standard library packages
	if !strings.Contains(path, ".") {
		// For standard library, the package name is typically the last component
		return ir.getLastPathComponent(path), nil
	}

	// Handle third-party packages
	packageName := gomod.GetPackageModuleName(path)
	if packageName == "" {
		return "", fmt.Errorf("could not resolve package name for %s", path)
	}

	return packageName, nil
}

// resolveDependencies resolves dependencies for an import
func (ir *ImportResolver) resolveDependencies(resolution *ImportResolution, context *ImportContext) error {
	// This is a placeholder for dependency resolution
	// In a more sophisticated implementation, this could:
	// - Parse the imported package to find its dependencies
	// - Check for version conflicts
	// - Validate import compatibility

	// For now, we'll just note that third-party packages might have dependencies
	if strings.Contains(resolution.Path, ".") {
		// Add some common dependencies as examples
		// This could be made configurable
	}
	return nil
}

// GetAlias returns the alias for a given import path
func (ir *ImportResolver) GetAlias(path string, context *ImportContext) (string, bool) {
	alias, exists := context.ImportPaths[path]
	return alias, exists
}

// GetPath returns the import path for a given alias
func (ir *ImportResolver) GetPath(alias string, context *ImportContext) (string, bool) {
	if resolution, exists := context.FileImports[alias]; exists {
		return resolution.Path, true
	}
	return "", false
}

// GetPackageName returns the package name for a given alias or path
func (ir *ImportResolver) GetPackageName(identifier string, context *ImportContext) (string, bool) {
	// First try as alias
	if resolution, exists := context.FileImports[identifier]; exists {
		return resolution.PackageName, true
	}

	// Then try as path
	if alias, exists := context.ImportPaths[identifier]; exists {
		if resolution, resExists := context.FileImports[alias]; resExists {
			return resolution.PackageName, true
		}
	}

	return "", false
}

// IsValidImport checks if an import path is valid
func (ir *ImportResolver) IsValidImport(path string) bool {
	// Basic validation
	if path == "" {
		return false
	}

	// Check for invalid characters
	if strings.ContainsAny(path, " \t\n\r\"'") {
		return false
	}

	// TODO: Add more sophisticated validation
	return true
}

// GetImportPathFromType extracts the import path from a qualified type name
func (ir *ImportResolver) GetImportPathFromType(typeName string, context *ImportContext) (string, bool) {
	if !strings.Contains(typeName, ".") {
		return "", false
	}

	alias := strings.Split(typeName, ".")[0]
	path, exists := ir.GetPath(alias, context)
	return path, exists
}

// Helper methods

func (ir *ImportResolver) getLastPathComponent(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) == 0 {
		return ""
	}
	return parts[len(parts)-1]
}

func (ir *ImportResolver) generateCustomAlias(path string, context *ImportContext) string {
	// Generate a meaningful alias based on the path
	parts := strings.Split(path, "/")
	if len(parts) == 0 {
		return "unknown"
	}

	// Use the last few parts to create a meaningful alias
	start := 0
	if len(parts) > 2 {
		start = len(parts) - 2
	}

	aliasParts := parts[start:]
	for i, part := range aliasParts {
		aliasParts[i] = strings.ToLower(part)
	}

	return strings.Join(aliasParts, "")
}

func (ir *ImportResolver) generateUniqueAlias(baseAlias string, context *ImportContext) string {
	// Check if base alias is available
	if _, exists := context.FileImports[baseAlias]; !exists {
		return baseAlias
	}

	// Generate unique alias by adding suffix
	for i := 1; i < 1000; i++ {
		candidate := fmt.Sprintf("%s%d", baseAlias, i)
		if _, exists := context.FileImports[candidate]; !exists {
			return candidate
		}
	}

	// Fallback to random suffix
	return fmt.Sprintf("%s%d", baseAlias, rand.Intn(10000))
}

// ClearCache clears the import resolution cache
func (ir *ImportResolver) ClearCache() {
	ir.cache = make(map[string]*ImportResolution)
}

// GetCacheSize returns the number of cached resolutions
func (ir *ImportResolver) GetCacheSize() int {
	return len(ir.cache)
}

// GetConfig returns the resolver configuration
func (ir *ImportResolver) GetConfig() *ResolverConfig {
	return ir.resolverConfig
}

// SetConfig updates the resolver configuration
func (ir *ImportResolver) SetConfig(config *ResolverConfig) {
	ir.resolverConfig = config
	// Clear cache when config changes
	ir.ClearCache()
}
