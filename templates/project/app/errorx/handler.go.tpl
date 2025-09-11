package errorx

import (
    "errors"
    "net/http"

    "github.com/gofiber/fiber/v3"
    "gorm.io/gorm"
)

// ErrorHandler 错误处理器
type ErrorHandler struct{}

// NewErrorHandler 创建错误处理器
func NewErrorHandler() *ErrorHandler {
    return &ErrorHandler{}
}

// Handle 处理错误并返回统一格式
func (h *ErrorHandler) Handle(err error) *AppError {
    if appErr, ok := err.(*AppError); ok {
        return appErr
    }

    // 处理 Fiber 错误
    if fiberErr, ok := err.(*fiber.Error); ok {
        return h.handleFiberError(fiberErr)
    }

    // 处理 GORM 错误
    if appErr := h.handleGormError(err); appErr != nil {
        return appErr
    }

    // 默认内部错误
    return &AppError{
        Code:        ErrInternalError.Code,
        Message:     err.Error(),
        StatusCode:  http.StatusInternalServerError,
        originalErr: err,
    }
}

// handleFiberError 处理 Fiber 错误
func (h *ErrorHandler) handleFiberError(fiberErr *fiber.Error) *AppError {
    var appErr *AppError

    switch fiberErr.Code {
    case http.StatusBadRequest:
        appErr = ErrBadRequest
    case http.StatusUnauthorized:
        appErr = ErrUnauthorized
    case http.StatusForbidden:
        appErr = ErrForbidden
    case http.StatusNotFound:
        appErr = ErrRecordNotFound
    case http.StatusMethodNotAllowed:
        appErr = ErrUnsupportedMethod
    case http.StatusRequestEntityTooLarge:
        appErr = ErrRequestTooLarge
    case http.StatusTooManyRequests:
        appErr = ErrTooManyRequests
    default:
        appErr = ErrInternalError
    }

    return &AppError{
        Code:        appErr.Code,
        Message:     fiberErr.Message,
        StatusCode:  fiberErr.Code,
        originalErr: fiberErr,
    }
}

// handleGormError 处理 GORM 错误
func (h *ErrorHandler) handleGormError(err error) *AppError {
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return &AppError{
            Code:        ErrRecordNotFound.Code,
            Message:     ErrRecordNotFound.Message,
            StatusCode:  ErrRecordNotFound.StatusCode,
            originalErr: err,
        }
    }

    if errors.Is(err, gorm.ErrDuplicatedKey) {
        return &AppError{
            Code:        ErrRecordDuplicated.Code,
            Message:     ErrRecordDuplicated.Message,
            StatusCode:  ErrRecordDuplicated.StatusCode,
            originalErr: err,
        }
    }

    if errors.Is(err, gorm.ErrInvalidTransaction) {
        return &AppError{
            Code:        ErrConcurrencyError.Code,
            Message:     "事务无效",
            StatusCode:  ErrConcurrencyError.StatusCode,
            originalErr: err,
        }
    }

    return nil
}

