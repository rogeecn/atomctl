package errorx

import (
	"fmt"
	"runtime"
)

// AppError 应用错误结构
type AppError struct {
	Code       ErrorCode `json:"code"`
	Message    string    `json:"message"`
	StatusCode int       `json:"-"`
	Data       any       `json:"data,omitempty"`
	ID         string    `json:"id,omitempty"`

	// 调试信息
	originalErr error
	file        string
	params      []any
	sql         string
}

// Error 实现 error 接口
func (e *AppError) Error() string {
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// Unwrap 允许通过 errors.Unwrap 遍历到原始错误
func (e *AppError) Unwrap() error { return e.originalErr }

// copy 返回 AppError 的副本，用于链式调用时的并发安全
func (e *AppError) copy() *AppError {
	newErr := *e
	return &newErr
}

// WithCause 携带原始错误并记录调用栈
func (e *AppError) WithCause(err error) *AppError {
	newErr := e.copy()
	newErr.originalErr = err

	// 记录调用者位置
	if _, file, line, ok := runtime.Caller(1); ok {
		newErr.file = fmt.Sprintf("%s:%d", file, line)
	}
	return newErr
}

// WithData 添加数据
func (e *AppError) WithData(data any) *AppError {
	newErr := e.copy()
	newErr.Data = data
	return newErr
}

// WithMsg 设置消息
func (e *AppError) WithMsg(msg string) *AppError {
	newErr := e.copy()
	newErr.Message = msg
	return newErr
}

func (e *AppError) WithMsgf(format string, args ...any) *AppError {
	newErr := e.copy()
	newErr.Message = fmt.Sprintf(format, args...)
	return newErr
}

// WithSQL 记录SQL信息
func (e *AppError) WithSQL(sql string) *AppError {
	newErr := e.copy()
	newErr.sql = sql
	return newErr
}

// WithParams 记录参数信息，并自动获取调用位置
func (e *AppError) WithParams(params ...any) *AppError {
	newErr := e.copy()
	newErr.params = params
	if _, file, line, ok := runtime.Caller(1); ok {
		newErr.file = fmt.Sprintf("%s:%d", file, line)
	}
	return newErr
}

// NewError 创建应用错误
func NewError(code ErrorCode, statusCode int, message string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
	}
}
