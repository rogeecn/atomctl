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

// WithData 添加数据
func (e *AppError) WithData(data any) *AppError {
    e.Data = data
    return e
}

// WithMsg 设置消息
func (e *AppError) WithMsg(msg string) *AppError {
    e.Message = msg
    return e
}

// WithSQL 记录SQL信息
func (e *AppError) WithSQL(sql string) *AppError {
    e.sql = sql
    return e
}

// WithParams 记录参数信息，并自动获取调用位置
func (e *AppError) WithParams(params ...any) *AppError {
    e.params = params
    if _, file, line, ok := runtime.Caller(1); ok {
        e.file = fmt.Sprintf("%s:%d", file, line)
    }
    return e
}

// NewError 创建应用错误
func NewError(code ErrorCode, statusCode int, message string) *AppError {
    return &AppError{
        Code:       code,
        Message:    message,
        StatusCode: statusCode,
    }
}
