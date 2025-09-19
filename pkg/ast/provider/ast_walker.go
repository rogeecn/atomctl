package provider

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// ASTWalker handles traversal of Go AST nodes to find provider-related structures
type ASTWalker struct {
	fileSet       *token.FileSet
	commentParser *CommentParser
	config        *WalkerConfig
	visitors      []NodeVisitor
}

// WalkerConfig configures the AST walker behavior
type WalkerConfig struct {
	IncludeTestFiles      bool
	IncludeGeneratedFiles bool
	MaxFileSize           int64
	StrictMode            bool
}

// NodeVisitor defines the interface for visiting AST nodes
type NodeVisitor interface {
	// VisitFile is called when a new file is processed
	VisitFile(filePath string, node *ast.File) error

	// VisitGenDecl is called for each generic declaration (type, var, const)
	VisitGenDecl(filePath string, decl *ast.GenDecl) error

	// VisitTypeSpec is called for each type specification
	VisitTypeSpec(filePath string, typeSpec *ast.TypeSpec, decl *ast.GenDecl) error

	// VisitStructType is called for each struct type
	VisitStructType(filePath string, structType *ast.StructType, typeSpec *ast.TypeSpec, decl *ast.GenDecl) error

	// VisitStructField is called for each field in a struct
	VisitStructField(filePath string, field *ast.Field, structType *ast.StructType) error

	// Complete is called when file processing is complete
	Complete(filePath string) error
}

// NewASTWalker creates a new ASTWalker with default configuration
func NewASTWalker() *ASTWalker {
	return &ASTWalker{
		fileSet:       token.NewFileSet(),
		commentParser: NewCommentParser(),
		config: &WalkerConfig{
			IncludeTestFiles:      false,
			IncludeGeneratedFiles: false,
			MaxFileSize:           10 * 1024 * 1024, // 10MB
			StrictMode:            false,
		},
		visitors: make([]NodeVisitor, 0),
	}
}

// NewASTWalkerWithConfig creates a new ASTWalker with custom configuration
func NewASTWalkerWithConfig(config *WalkerConfig) *ASTWalker {
	if config == nil {
		return NewASTWalker()
	}

	return &ASTWalker{
		fileSet:       token.NewFileSet(),
		commentParser: NewCommentParserWithStrictMode(config.StrictMode),
		config:        config,
		visitors:      make([]NodeVisitor, 0),
	}
}

// AddVisitor adds a node visitor to the walker
func (aw *ASTWalker) AddVisitor(visitor NodeVisitor) {
	aw.visitors = append(aw.visitors, visitor)
}

// RemoveVisitor removes a node visitor from the walker
func (aw *ASTWalker) RemoveVisitor(visitor NodeVisitor) {
	for i, v := range aw.visitors {
		if v == visitor {
			aw.visitors = append(aw.visitors[:i], aw.visitors[i+1:]...)
			break
		}
	}
}

// WalkFile traverses a single Go file
func (aw *ASTWalker) WalkFile(filePath string) error {
	// Check if file should be processed
	if !aw.shouldProcessFile(filePath) {
		return nil
	}

	// Parse the file
	node, err := parser.ParseFile(aw.fileSet, filePath, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("failed to parse file %s: %w", filePath, err)
	}

	// Notify visitors of file start
	for _, visitor := range aw.visitors {
		if err := visitor.VisitFile(filePath, node); err != nil {
			return err
		}
	}

	// Traverse the AST
	if err := aw.traverseFile(filePath, node); err != nil {
		return err
	}

	// Notify visitors of file completion
	for _, visitor := range aw.visitors {
		if err := visitor.Complete(filePath); err != nil {
			return err
		}
	}

	return nil
}

// WalkDir traverses all Go files in a directory
func (aw *ASTWalker) WalkDir(dirPath string) error {
	return filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
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

		// Process Go files
		if filepath.Ext(path) == ".go" && aw.shouldProcessFile(path) {
			if err := aw.WalkFile(path); err != nil {
				// Continue with other files, but log the error
				fmt.Printf("Warning: failed to process file %s: %v\n", path, err)
			}
		}

		return nil
	})
}

// traverseFile traverses the AST of a parsed file
func (aw *ASTWalker) traverseFile(filePath string, node *ast.File) error {
	// Traverse all declarations
	for _, decl := range node.Decls {
		if err := aw.traverseDeclaration(filePath, decl); err != nil {
			return err
		}
	}

	return nil
}

// traverseDeclaration traverses a single declaration
func (aw *ASTWalker) traverseDeclaration(filePath string, decl ast.Decl) error {
	genDecl, ok := decl.(*ast.GenDecl)
	if !ok {
		// Skip function declarations and other non-generic declarations
		return nil
	}

	// Notify visitors of generic declaration
	for _, visitor := range aw.visitors {
		if err := visitor.VisitGenDecl(filePath, genDecl); err != nil {
			return err
		}
	}

	// Traverse specs within the declaration
	for _, spec := range genDecl.Specs {
		if err := aw.traverseSpec(filePath, spec, genDecl); err != nil {
			return err
		}
	}

	return nil
}

// traverseSpec traverses a specification within a declaration
func (aw *ASTWalker) traverseSpec(filePath string, spec ast.Spec, decl *ast.GenDecl) error {
	typeSpec, ok := spec.(*ast.TypeSpec)
	if !ok {
		// Skip non-type specifications
		return nil
	}

	// Notify visitors of type specification
	for _, visitor := range aw.visitors {
		if err := visitor.VisitTypeSpec(filePath, typeSpec, decl); err != nil {
			return err
		}
	}

	// Check if it's a struct type
	structType, ok := typeSpec.Type.(*ast.StructType)
	if ok {
		// Notify visitors of struct type
		for _, visitor := range aw.visitors {
			if err := visitor.VisitStructType(filePath, structType, typeSpec, decl); err != nil {
				return err
			}
		}

		// Traverse struct fields
		if err := aw.traverseStructFields(filePath, structType); err != nil {
			return err
		}
	}

	return nil
}

// traverseStructFields traverses fields within a struct type
func (aw *ASTWalker) traverseStructFields(filePath string, structType *ast.StructType) error {
	if structType.Fields == nil {
		return nil
	}

	for _, field := range structType.Fields.List {
		// Notify visitors of struct field
		for _, visitor := range aw.visitors {
			if err := visitor.VisitStructField(filePath, field, structType); err != nil {
				return err
			}
		}
	}

	return nil
}

// shouldProcessFile determines if a file should be processed
func (aw *ASTWalker) shouldProcessFile(filePath string) bool {
	// Check file extension
	if filepath.Ext(filePath) != ".go" {
		return false
	}

	// Skip test files if not allowed
	if !aw.config.IncludeTestFiles && strings.HasSuffix(filePath, "_test.go") {
		return false
	}

	// Skip generated files if not allowed
	if !aw.config.IncludeGeneratedFiles && strings.HasSuffix(filePath, ".gen.go") {
		return false
	}

	// TODO: Check file size if needed (requires os.Stat)

	return true
}

// GetFileSet returns the file set used by the walker
func (aw *ASTWalker) GetFileSet() *token.FileSet {
	return aw.fileSet
}

// GetCommentParser returns the comment parser used by the walker
func (aw *ASTWalker) GetCommentParser() *CommentParser {
	return aw.commentParser
}

// GetConfig returns the walker configuration
func (aw *ASTWalker) GetConfig() *WalkerConfig {
	return aw.config
}

// ProviderDiscoveryVisitor implements NodeVisitor for discovering provider annotations
type ProviderDiscoveryVisitor struct {
	commentParser *CommentParser
	providers     []Provider
	currentFile   string
}

// NewProviderDiscoveryVisitor creates a new ProviderDiscoveryVisitor
func NewProviderDiscoveryVisitor(commentParser *CommentParser) *ProviderDiscoveryVisitor {
	return &ProviderDiscoveryVisitor{
		commentParser: commentParser,
		providers:     make([]Provider, 0),
	}
}

// VisitFile implements NodeVisitor.VisitFile
func (pdv *ProviderDiscoveryVisitor) VisitFile(filePath string, node *ast.File) error {
	pdv.currentFile = filePath
	return nil
}

// VisitGenDecl implements NodeVisitor.VisitGenDecl
func (pdv *ProviderDiscoveryVisitor) VisitGenDecl(filePath string, decl *ast.GenDecl) error {
	return nil
}

// VisitTypeSpec implements NodeVisitor.VisitTypeSpec
func (pdv *ProviderDiscoveryVisitor) VisitTypeSpec(filePath string, typeSpec *ast.TypeSpec, decl *ast.GenDecl) error {
	return nil
}

// VisitStructType implements NodeVisitor.VisitStructType
func (pdv *ProviderDiscoveryVisitor) VisitStructType(filePath string, structType *ast.StructType, typeSpec *ast.TypeSpec, decl *ast.GenDecl) error {
	// Check if the struct has a provider annotation
	if decl.Doc != nil && len(decl.Doc.List) > 0 {
		// Extract comment lines
		commentLines := make([]string, len(decl.Doc.List))
		for i, comment := range decl.Doc.List {
			commentLines[i] = comment.Text
		}

		// Parse provider annotation
		providerComment, err := pdv.commentParser.ParseCommentBlock(commentLines)
		if err == nil && providerComment != nil {
			// Create provider structure
			provider := Provider{
				StructName:    typeSpec.Name.Name,
				Mode:          providerComment.Mode,
				ProviderGroup: providerComment.Group,
				ReturnType:    providerComment.ReturnType,
				InjectParams:  make(map[string]InjectParam),
				Imports:       make(map[string]string),
			}

			// Set default return type if not specified
			if provider.ReturnType == "" {
				provider.ReturnType = "*" + provider.StructName
			}

			pdv.providers = append(pdv.providers, provider)
		}
	}

	return nil
}

// VisitStructField implements NodeVisitor.VisitStructField
func (pdv *ProviderDiscoveryVisitor) VisitStructField(filePath string, field *ast.Field, structType *ast.StructType) error {
	// This is where field-level processing would happen
	// For example, extracting inject tags and field types
	return nil
}

// Complete implements NodeVisitor.Complete
func (pdv *ProviderDiscoveryVisitor) Complete(filePath string) error {
	return nil
}

// GetProviders returns the discovered providers
func (pdv *ProviderDiscoveryVisitor) GetProviders() []Provider {
	return pdv.providers
}

// Reset clears the discovered providers
func (pdv *ProviderDiscoveryVisitor) Reset() {
	pdv.providers = make([]Provider, 0)
	pdv.currentFile = ""
}
