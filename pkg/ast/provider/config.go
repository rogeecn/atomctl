package provider

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// ParserConfig represents the configuration for the parser
type ParserConfig struct {
	// File parsing options
	ParseComments bool           // Whether to parse comments (default: true)
	Mode          parser.Mode    // Parser mode
	FileSet       *token.FileSet // File set for position information

	// Include/exclude options
	IncludePatterns []string // Glob patterns for files to include
	ExcludePatterns []string // Glob patterns for files to exclude

	// Provider parsing options
	StrictMode     bool   // Whether to use strict validation (default: false)
	DefaultMode    string // Default provider mode for simple @provider annotations
	AllowTestFiles bool   // Whether to parse test files (default: false)
	AllowGenFiles  bool   // Whether to parse generated files (default: false)

	// Performance options
	MaxFileSize  int64 // Maximum file size to parse (bytes, default: 10MB)
	Concurrency  int   // Number of concurrent parsers (default: 1)
	CacheEnabled bool  // Whether to enable caching (default: true)

	// Output options
	OutputDir       string // Output directory for generated files
	OutputFileName  string // Output file name (default: provider.gen.go)
	SourceLocations bool   // Whether to include source location info (default: false)
}

// ParserContext represents the context for parsing operations
type ParserContext struct {
	// Configuration
	Config *ParserConfig

	// Parsing state
	FileSet    *token.FileSet
	WorkingDir string
	ModuleName string

	// Import resolution
	Imports    map[string]string // Package alias -> package path
	ModuleInfo map[string]string // Module path -> module name

	// Statistics and metrics
	FilesProcessed int
	FilesSkipped   int
	ProvidersFound int
	ParseErrors    []ParseError

	// Caching
	Cache map[string]interface{} // File path -> parsed content
}

// ParseError represents a parsing error with location information
type ParseError struct {
	File     string `json:"file"`
	Line     int    `json:"line"`
	Column   int    `json:"column"`
	Message  string `json:"message"`
	Severity string `json:"severity"` // "error", "warning", "info"
}

// NewParserConfig creates a new ParserConfig with default values
func NewParserConfig() *ParserConfig {
	return &ParserConfig{
		ParseComments:   true,
		Mode:            parser.ParseComments,
		StrictMode:      false,
		DefaultMode:     "basic",
		AllowTestFiles:  false,
		AllowGenFiles:   false,
		MaxFileSize:     10 * 1024 * 1024, // 10MB
		Concurrency:     1,
		CacheEnabled:    true,
		OutputFileName:  "provider.gen.go",
		SourceLocations: false,
	}
}

// NewParserContext creates a new ParserContext with the given configuration
func NewParserContext(config *ParserConfig) *ParserContext {
	if config == nil {
		config = NewParserConfig()
	}

	return &ParserContext{
		Config:      config,
		FileSet:     config.FileSet,
		Imports:     make(map[string]string),
		ModuleInfo:  make(map[string]string),
		ParseErrors: make([]ParseError, 0),
		Cache:       make(map[string]interface{}),
	}
}

// ShouldIncludeFile determines if a file should be included in parsing
func (c *ParserContext) ShouldIncludeFile(filePath string) bool {
	// Check file extension
	if filepath.Ext(filePath) != ".go" {
		return false
	}

	// Skip test files if not allowed
	if !c.Config.AllowTestFiles && strings.HasSuffix(filePath, "_test.go") {
		return false
	}

	// Skip generated files if not allowed, but allow routes.gen.go and services.gen.go since they may contain providers
	if !c.Config.AllowGenFiles {
		return false
	}
	// if !c.Config.AllowGenFiles && strings.HasSuffix(filePath, ".gen.go") && !strings.HasSuffix(filePath, "routes.gen.go") && !strings.HasSuffix(filePath, "services.gen.go") {
	// 	return false
	// }

	// Check file size
	if info, err := os.Stat(filePath); err == nil {
		if info.Size() > c.Config.MaxFileSize {
			c.AddError(filePath, 0, 0, "file exceeds maximum size", "warning")
			return false
		}
	}

	// TODO: Implement include/exclude pattern matching
	// For now, include all Go files that pass the basic checks
	return true
}

// AddError adds a parsing error to the context
func (c *ParserContext) AddError(file string, line, column int, message, severity string) {
	c.ParseErrors = append(c.ParseErrors, ParseError{
		File:     file,
		Line:     line,
		Column:   column,
		Message:  message,
		Severity: severity,
	})
}

// HasErrors returns true if there are any errors in the context
func (c *ParserContext) HasErrors() bool {
	for _, err := range c.ParseErrors {
		if err.Severity == "error" {
			return true
		}
	}
	return false
}

// GetErrors returns all errors of a specific severity
func (c *ParserContext) GetErrors(severity string) []ParseError {
	var errors []ParseError
	for _, err := range c.ParseErrors {
		if err.Severity == severity {
			errors = append(errors, err)
		}
	}
	return errors
}

// AddImport adds an import to the context
func (c *ParserContext) AddImport(alias, path string) {
	c.Imports[alias] = path
}

// GetImportPath returns the import path for a given alias
func (c *ParserContext) GetImportPath(alias string) (string, bool) {
	path, ok := c.Imports[alias]
	return path, ok
}
