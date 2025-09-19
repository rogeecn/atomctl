package provider

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"runtime"
	"strings"
)

// ParserError represents errors that occur during parsing
type ParserError struct {
	Operation string `json:"operation"`
	File      string `json:"file"`
	Line      int    `json:"line"`
	Column    int    `json:"column"`
	Message   string `json:"message"`
	Code      string `json:"code,omitempty"`
	Severity  string `json:"severity"` // "error", "warning", "info"
	Cause     error  `json:"cause,omitempty"`
	Stack     string `json:"stack,omitempty"`
}

// Error implements the error interface
func (e *ParserError) Error() string {
	if e.File != "" {
		return fmt.Sprintf("%s: %s at %s:%d:%d", e.Operation, e.Message, e.File, e.Line, e.Column)
	}
	return fmt.Sprintf("%s: %s", e.Operation, e.Message)
}

// Unwrap implements the error unwrapping interface
func (e *ParserError) Unwrap() error {
	return e.Cause
}

// Is implements the error comparison interface
func (e *ParserError) Is(target error) bool {
	var other *ParserError
	if errors.As(target, &other) {
		return e.Operation == other.Operation && e.Code == other.Code
	}
	return false
}

// RendererError represents errors that occur during rendering
type RendererError struct {
	Operation string `json:"operation"`
	Template  string `json:"template"`
	Target    string `json:"target,omitempty"`
	Message   string `json:"message"`
	Code      string `json:"code,omitempty"`
	Cause     error  `json:"cause,omitempty"`
	Stack     string `json:"stack,omitempty"`
}

// Error implements the error interface
func (e *RendererError) Error() string {
	if e.Template != "" {
		return fmt.Sprintf("renderer %s (template %s): %s", e.Operation, e.Template, e.Message)
	}
	return fmt.Sprintf("renderer %s: %s", e.Operation, e.Message)
}

// Unwrap implements the error unwrapping interface
func (e *RendererError) Unwrap() error {
	return e.Cause
}

// Is implements the error comparison interface
func (e *RendererError) Is(target error) bool {
	var other *RendererError
	if errors.As(target, &other) {
		return e.Operation == other.Operation && e.Code == other.Code
	}
	return false
}

// FileSystemError represents file system related errors
type FileSystemError struct {
	Operation string `json:"operation"`
	Path      string `json:"path"`
	Message   string `json:"message"`
	Code      string `json:"code,omitempty"`
	Cause     error  `json:"cause,omitempty"`
}

// Error implements the error interface
func (e *FileSystemError) Error() string {
	return fmt.Sprintf("file system %s failed for %s: %s", e.Operation, e.Path, e.Message)
}

// Unwrap implements the error unwrapping interface
func (e *FileSystemError) Unwrap() error {
	return e.Cause
}

// ConfigurationError represents configuration related errors
type ConfigurationError struct {
	Field   string `json:"field"`
	Value   string `json:"value,omitempty"`
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
	Cause   error  `json:"cause,omitempty"`
}

// Error implements the error interface
func (e *ConfigurationError) Error() string {
	if e.Value != "" {
		return fmt.Sprintf("configuration error for field %s (value: %s): %s", e.Field, e.Value, e.Message)
	}
	return fmt.Sprintf("configuration error for field %s: %s", e.Field, e.Message)
}

// Unwrap implements the error unwrapping interface
func (e *ConfigurationError) Unwrap() error {
	return e.Cause
}

// Error codes
const (
	ErrCodeFileNotFound       = "FILE_NOT_FOUND"
	ErrCodePermissionDenied   = "PERMISSION_DENIED"
	ErrCodeInvalidSyntax      = "INVALID_SYNTAX"
	ErrCodeInvalidAnnotation  = "INVALID_ANNOTATION"
	ErrCodeInvalidMode        = "INVALID_MODE"
	ErrCodeInvalidType        = "INVALID_TYPE"
	ErrCodeTemplateNotFound   = "TEMPLATE_NOT_FOUND"
	ErrCodeTemplateError      = "TEMPLATE_ERROR"
	ErrCodeValidationFailed   = "VALIDATION_FAILED"
	ErrCodeConfigurationError = "CONFIGURATION_ERROR"
	ErrCodeFileSystemError    = "FILE_SYSTEM_ERROR"
	ErrCodeUnknownError       = "UNKNOWN_ERROR"
)

// Error severity levels
const (
	SeverityError   = "error"
	SeverityWarning = "warning"
	SeverityInfo    = "info"
)

// Error builder functions

// NewParserError creates a new ParserError
func NewParserError(operation, message string) *ParserError {
	return &ParserError{
		Operation: operation,
		Message:   message,
		Severity:  SeverityError,
	}
}

// NewParserErrorWithCause creates a new ParserError with a cause
func NewParserErrorWithCause(operation, message string, cause error) *ParserError {
	err := NewParserError(operation, message)
	err.Cause = cause
	err.Stack = captureStackTrace(2)
	return err
}

// NewParserErrorAtLocation creates a new ParserError with file location
func NewParserErrorAtLocation(operation, file string, line, column int, message string) *ParserError {
	return &ParserError{
		Operation: operation,
		File:      file,
		Line:      line,
		Column:    column,
		Message:   message,
		Severity:  SeverityError,
	}
}

// NewValidationError creates a new ValidationError
func NewValidationError(ruleName, message string) *ValidationError {
	return &ValidationError{
		RuleName: ruleName,
		Message:  message,
		Severity: SeverityError,
	}
}

// NewValidationErrorWithCause creates a new ValidationError with a cause
func NewValidationErrorWithCause(ruleName, message string, cause error) *ValidationError {
	err := NewValidationError(ruleName, message)
	err.Cause = cause
	return err
}

// NewRendererError creates a new RendererError
func NewRendererError(operation, message string) *RendererError {
	return &RendererError{
		Operation: operation,
		Message:   message,
	}
}

// NewRendererErrorWithCause creates a new RendererError with a cause
func NewRendererErrorWithCause(operation, message string, cause error) *RendererError {
	err := NewRendererError(operation, message)
	err.Cause = cause
	err.Stack = captureStackTrace(2)
	return err
}

// NewFileSystemError creates a new FileSystemError
func NewFileSystemError(operation, path, message string) *FileSystemError {
	return &FileSystemError{
		Operation: operation,
		Path:      path,
		Message:   message,
	}
}

// NewFileSystemErrorFromError creates a new FileSystemError from an existing error
func NewFileSystemErrorFromError(operation, path string, err error) *FileSystemError {
	return &FileSystemError{
		Operation: operation,
		Path:      path,
		Message:   err.Error(),
		Cause:     err,
	}
}

// NewConfigurationError creates a new ConfigurationError
func NewConfigurationError(field, message string) *ConfigurationError {
	return &ConfigurationError{
		Field:   field,
		Message: message,
	}
}

// WrapError wraps an error with additional context
func WrapError(err error, operation string) error {
	if err == nil {
		return nil
	}

	switch e := err.(type) {
	case *ParserError:
		return NewParserErrorWithCause(e.Operation, e.Message, err)
	case *ValidationError:
		return NewValidationErrorWithCause(e.RuleName, e.Message, err)
	case *RendererError:
		return NewRendererErrorWithCause(e.Operation, e.Message, err)
	case *FileSystemError:
		return NewFileSystemErrorFromError(e.Operation, e.Path, err)
	case *ConfigurationError:
		return NewConfigurationError(e.Field, e.Message)
	default:
		return fmt.Errorf("%s: %w", operation, err)
	}
}

// Error utility functions

// IsParserError checks if an error is a ParserError
func IsParserError(err error) bool {
	var target *ParserError
	return errors.As(err, &target)
}

// IsValidationError checks if an error is a ValidationError
func IsValidationError(err error) bool {
	var target *ValidationError
	return errors.As(err, &target)
}

// IsRendererError checks if an error is a RendererError
func IsRendererError(err error) bool {
	var target *RendererError
	return errors.As(err, &target)
}

// IsFileSystemError checks if an error is a FileSystemError
func IsFileSystemError(err error) bool {
	var target *FileSystemError
	return errors.As(err, &target)
}

// IsConfigurationError checks if an error is a ConfigurationError
func IsConfigurationError(err error) bool {
	var target *ConfigurationError
	return errors.As(err, &target)
}

// GetErrorCode returns the error code for a given error
func GetErrorCode(err error) string {
	if err == nil {
		return ""
	}

	switch e := err.(type) {
	case *ParserError:
		return e.Code
	case *RendererError:
		return e.Code
	case *FileSystemError:
		return e.Code
	case *ConfigurationError:
		return e.Code
	default:
		return ErrCodeUnknownError
	}
}

// IsFileNotFoundError checks if an error is a file not found error
func IsFileNotFoundError(err error) bool {
	if errors.Is(err, fs.ErrNotExist) {
		return true
	}

	if pathErr, ok := err.(*os.PathError); ok {
		return errors.Is(pathErr.Err, fs.ErrNotExist)
	}

	return false
}

// IsPermissionError checks if an error is a permission error
func IsPermissionError(err error) bool {
	if errors.Is(err, os.ErrPermission) {
		return true
	}

	if pathErr, ok := err.(*os.PathError); ok {
		return errors.Is(pathErr.Err, os.ErrPermission)
	}

	return false
}

// Error recovery functions

// RecoverFromParseError attempts to recover from parsing errors
func RecoverFromParseError(err error) error {
	if err == nil {
		return nil
	}

	// If it's a file system error, provide more helpful message
	if IsFileSystemError(err) {
		var fsErr *FileSystemError
		if errors.As(err, &fsErr) {
			if IsFileNotFoundError(fsErr.Cause) {
				return NewParserErrorWithCause(
					"parse",
					fmt.Sprintf("file not found: %s", fsErr.Path),
					err,
				)
			}
			if IsPermissionError(fsErr.Cause) {
				return NewParserErrorWithCause(
					"parse",
					fmt.Sprintf("permission denied: %s", fsErr.Path),
					err,
				)
			}
		}
	}

	// For syntax errors, try to provide location information
	if strings.Contains(err.Error(), "syntax") {
		return NewParserErrorWithCause("parse", "syntax error in source code", err)
	}

	// Default: wrap the error with parsing context
	return WrapError(err, "parse")
}

// captureStackTrace captures the current stack trace
func captureStackTrace(skip int) string {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(skip, pcs[:])
	frames := runtime.CallersFrames(pcs[:n])

	var builder strings.Builder
	for {
		frame, more := frames.Next()
		if !more || strings.Contains(frame.Function, "runtime.") {
			break
		}
		fmt.Fprintf(&builder, "%s\n\t%s:%d\n", frame.Function, frame.File, frame.Line)
	}

	return builder.String()
}

// Error aggregation

// ErrorGroup represents a group of related errors
type ErrorGroup struct {
	Errors []error `json:"errors"`
}

// Error implements the error interface
func (eg *ErrorGroup) Error() string {
	if len(eg.Errors) == 0 {
		return "no errors"
	}

	if len(eg.Errors) == 1 {
		return eg.Errors[0].Error()
	}

	var messages []string
	for _, err := range eg.Errors {
		messages = append(messages, err.Error())
	}
	return fmt.Sprintf("multiple errors occurred:\n\t%s", strings.Join(messages, "\n\t"))
}

// Unwrap implements the error unwrapping interface
func (eg *ErrorGroup) Unwrap() []error {
	return eg.Errors
}

// NewErrorGroup creates a new ErrorGroup
func NewErrorGroup(errors ...error) *ErrorGroup {
	var filteredErrors []error
	for _, err := range errors {
		if err != nil {
			filteredErrors = append(filteredErrors, err)
		}
	}
	return &ErrorGroup{Errors: filteredErrors}
}

// CollectErrors collects non-nil errors from a slice
func CollectErrors(errors ...error) *ErrorGroup {
	return NewErrorGroup(errors...)
}
