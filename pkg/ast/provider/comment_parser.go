package provider

import (
	"fmt"
	"strings"
)

// CommentParser handles parsing of provider annotations from Go comments
type CommentParser struct {
	strictMode bool
}

// NewCommentParser creates a new CommentParser
func NewCommentParser() *CommentParser {
	return &CommentParser{
		strictMode: false,
	}
}

// NewCommentParserWithStrictMode creates a new CommentParser with strict mode enabled
func NewCommentParserWithStrictMode(strictMode bool) *CommentParser {
	return &CommentParser{
		strictMode: strictMode,
	}
}

// ParseProviderComment parses a provider annotation from a comment line
func (cp *CommentParser) ParseProviderComment(comment string) (*ProviderComment, error) {
	// Trim the comment markers
	comment = strings.TrimSpace(comment)
	comment = strings.TrimPrefix(comment, "//")
	comment = strings.TrimPrefix(comment, "/*")
	comment = strings.TrimSuffix(comment, "*/")
	comment = strings.TrimSpace(comment)

	// Check if it's a provider annotation
	if !strings.HasPrefix(comment, "@provider") {
		return nil, fmt.Errorf("not a provider annotation")
	}

	// Parse the provider annotation
	return cp.parseProviderAnnotation(comment)
}

// parseProviderAnnotation parses the provider annotation structure
func (cp *CommentParser) parseProviderAnnotation(annotation string) (*ProviderComment, error) {
	result := &ProviderComment{
		RawText: annotation,
		IsValid: true,
		Errors:  make([]string, 0),
	}

	// Remove @provider prefix
	content := strings.TrimSpace(strings.TrimPrefix(annotation, "@provider"))

	// Handle empty case
	if content == "" {
		result.Mode = ProviderModeBasic
		result.Injection = InjectionModeAuto
		return result, nil
	}

	// Parse the annotation components
	return cp.parseAnnotationComponents(content, result)
}

// parseAnnotationComponents parses the components of the provider annotation
func (cp *CommentParser) parseAnnotationComponents(content string, result *ProviderComment) (*ProviderComment, error) {
	// Parse injection mode first (only/except)
	injectionMode, remaining := cp.parseInjectionMode(content)
	result.Injection = injectionMode

	// Parse provider mode (in parentheses)
	providerMode, remaining := cp.parseProviderMode(remaining)
	if providerMode != "" {
		if IsValidProviderMode(providerMode) {
			result.Mode = ProviderMode(providerMode)
		} else {
			result.IsValid = false
			result.Errors = append(result.Errors, fmt.Sprintf("invalid provider mode: %s", providerMode))
			if cp.strictMode {
				return result, fmt.Errorf("invalid provider mode: %s", providerMode)
			}
		}
	} else {
		result.Mode = ProviderModeBasic
	}

	// Parse return type and group
	returnType, group := cp.parseReturnTypeAndGroup(remaining)
	result.ReturnType = returnType
	result.Group = group

	return result, nil
}

// parseInjectionMode parses the injection mode (only/except)
func (cp *CommentParser) parseInjectionMode(content string) (InjectionMode, string) {
	if strings.Contains(content, ":only") {
		return InjectionModeOnly, strings.Replace(content, ":only", "", 1)
	} else if strings.Contains(content, ":except") {
		return InjectionModeExcept, strings.Replace(content, ":except", "", 1)
	}
	return InjectionModeAuto, content
}

// parseProviderMode parses the provider mode from parentheses
func (cp *CommentParser) parseProviderMode(content string) (string, string) {
	start := strings.Index(content, "(")
	end := strings.Index(content, ")")

	if start >= 0 && end > start {
		mode := strings.TrimSpace(content[start+1 : end])
		remaining := content[:start] + strings.TrimSpace(content[end+1:])
		return mode, strings.TrimSpace(remaining)
	}
	return "", strings.TrimSpace(content)
}

// parseReturnTypeAndGroup parses the return type and group from remaining content
func (cp *CommentParser) parseReturnTypeAndGroup(content string) (string, string) {
	parts := strings.Fields(content)

	if len(parts) == 0 {
		return "", ""
	}

	if len(parts) == 1 {
		return parts[0], ""
	}

	return parts[0], parts[1]
}

// IsProviderAnnotation checks if a comment line is a provider annotation
func (cp *CommentParser) IsProviderAnnotation(comment string) bool {
	comment = strings.TrimSpace(comment)
	comment = strings.TrimPrefix(comment, "//")
	comment = strings.TrimPrefix(comment, "/*")
	comment = strings.TrimSpace(comment)
	return strings.HasPrefix(comment, "@provider")
}

// ParseCommentBlock parses a block of comments to find provider annotations
func (cp *CommentParser) ParseCommentBlock(comments []string) (*ProviderComment, error) {
	if len(comments) == 0 {
		return nil, fmt.Errorf("empty comment block")
	}

	// Check each comment line for provider annotation (from bottom to top)
	for i := len(comments) - 1; i >= 0; i-- {
		comment := comments[i]
		if cp.IsProviderAnnotation(comment) {
			return cp.ParseProviderComment(comment)
		}
	}

	return nil, fmt.Errorf("no provider annotation found in comment block")
}

// ValidateProviderComment validates a parsed provider comment
func (cp *CommentParser) ValidateProviderComment(comment *ProviderComment) []string {
	var errors []string

	if comment == nil {
		errors = append(errors, "comment is nil")
		return errors
	}

	// Validate provider mode
	if comment.Mode != "" && !IsValidProviderMode(string(comment.Mode)) {
		errors = append(errors, fmt.Sprintf("invalid provider mode: %s", comment.Mode))
	}

	// Validate injection mode
	if comment.Injection != "" && !IsValidInjectionMode(string(comment.Injection)) {
		errors = append(errors, fmt.Sprintf("invalid injection mode: %s", comment.Injection))
	}

	// Validate return type format
	if comment.ReturnType != "" && !isValidGoType(comment.ReturnType) {
		errors = append(errors, fmt.Sprintf("invalid return type format: %s", comment.ReturnType))
	}

	// Validate group format
	if comment.Group != "" && !isValidGoIdentifier(comment.Group) {
		errors = append(errors, fmt.Sprintf("invalid group identifier: %s", comment.Group))
	}

	return errors
}

// ProviderComment represents a parsed provider annotation comment
type ProviderComment struct {
	RawText    string        // Original comment text
	Mode       ProviderMode  // Provider mode
	Injection  InjectionMode // Injection mode (only/except/auto)
	ReturnType string        // Return type specification
	Group      string        // Provider group
	IsValid    bool          // Whether the comment is valid
	Errors     []string      // Validation errors
}

// IsOnlyMode returns true if this is an "only" injection mode
func (pc *ProviderComment) IsOnlyMode() bool {
	return pc.Injection == InjectionModeOnly
}

// IsExceptMode returns true if this is an "except" injection mode
func (pc *ProviderComment) IsExceptMode() bool {
	return pc.Injection == InjectionModeExcept
}

// IsAutoMode returns true if this is an "auto" injection mode
func (pc *ProviderComment) IsAutoMode() bool {
	return pc.Injection == InjectionModeAuto
}

// HasMode returns true if a specific provider mode is set
func (pc *ProviderComment) HasMode(mode ProviderMode) bool {
	return pc.Mode == mode
}

// String returns a string representation of the provider comment
func (pc *ProviderComment) String() string {
	var builder strings.Builder

	builder.WriteString("@provider")

	if pc.Mode != ProviderModeBasic {
		builder.WriteString(fmt.Sprintf("(%s)", pc.Mode))
	}

	if pc.Injection == InjectionModeOnly {
		builder.WriteString(":only")
	} else if pc.Injection == InjectionModeExcept {
		builder.WriteString(":except")
	}

	if pc.ReturnType != "" {
		builder.WriteString(" ")
		builder.WriteString(pc.ReturnType)
	}

	if pc.Group != "" {
		builder.WriteString(" ")
		builder.WriteString(pc.Group)
	}

	return builder.String()
}
