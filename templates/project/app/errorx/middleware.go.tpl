package errorx

import "github.com/gofiber/fiber/v3"

// 全局实例
var DefaultSender = NewResponseSender()

// Middleware 错误处理中间件
func Middleware(c fiber.Ctx) error {
	err := c.Next()
	if err != nil {
		return DefaultSender.SendError(c, err)
	}
	return nil
}

// 便捷函数
func Wrap(err error) *AppError {
    if err == nil {
        return nil
    }

    if appErr, ok := err.(*AppError); ok {
        return &AppError{
            Code:        appErr.Code,
            Message:     appErr.Message,
            StatusCode:  appErr.StatusCode,
            Data:        appErr.Data,
            originalErr: appErr,
        }
    }

    return DefaultSender.handler.Handle(err)
}

func SendError(ctx fiber.Ctx, err error) error {
	return DefaultSender.SendError(ctx, err)
}
