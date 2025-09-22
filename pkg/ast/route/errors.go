package route

import "fmt"

// Error types for better error handling
type RouteErrorCode int

const (
	ErrInvalidInput RouteErrorCode = iota
	ErrInvalidPath
	ErrNoRoutes
	ErrParseFailed
	ErrTemplateFailed
	ErrFileWriteFailed
)

type RouteError struct {
	Code    RouteErrorCode
	Message string
	Cause   error
}

func (e *RouteError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("route error [%d]: %s (cause: %v)", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("route error [%d]: %s", e.Code, e.Message)
}

func (e *RouteError) Unwrap() error {
	return e.Cause
}

func (e *RouteError) WithCause(cause error) *RouteError {
	e.Cause = cause
	return e
}

func NewRouteError(code RouteErrorCode, format string, args ...interface{}) *RouteError {
	return &RouteError{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}

func WrapError(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	// If it's already a RouteError, just wrap it with more context
	if routeErr, ok := err.(*RouteError); ok {
		return NewRouteError(routeErr.Code, format, args...).WithCause(routeErr)
	}

	// Wrap other errors with a generic parse error
	return NewRouteError(ErrParseFailed, format, args...).WithCause(err)
}