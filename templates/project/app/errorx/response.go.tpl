package errorx

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/binder"
	"github.com/gofiber/utils/v2"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// ResponseSender 响应发送器
type ResponseSender struct {
	handler *ErrorHandler
}

// NewResponseSender 创建响应发送器
func NewResponseSender() *ResponseSender {
	return &ResponseSender{
		handler: NewErrorHandler(),
	}
}

// SendError 发送错误响应
func (s *ResponseSender) SendError(ctx fiber.Ctx, err error) error {
	appErr := s.handler.Handle(err)

	// 记录错误日志
	s.logError(appErr)

	// 根据 Content-Type 返回不同格式
	return s.sendResponse(ctx, appErr)
}

// logError 记录错误日志
func (s *ResponseSender) logError(appErr *AppError) {
	// 确保每个错误实例都有唯一ID，便于日志关联
	if appErr.ID == "" {
		appErr.ID = uuid.NewString()
	}

	// 构造详细的错误级联链路（包含类型、状态、定位等）
	chain := make([]map[string]any, 0, 4)
	var e error = appErr
	for e != nil {
		entry := map[string]any{
			"type":  fmt.Sprintf("%T", e),
			"error": e.Error(),
		}
		switch v := e.(type) {
		case *AppError:
			entry["code"] = v.Code
			entry["statusCode"] = v.StatusCode
			if v.file != "" {
				entry["file"] = v.file
			}
			if len(v.params) > 0 {
				entry["params"] = v.params
			}
			if v.sql != "" {
				entry["sql"] = v.sql
			}
			if v.ID != "" {
				entry["id"] = v.ID
			}
		case *fiber.Error:
			entry["statusCode"] = v.Code
			entry["message"] = v.Message
		}

		// GORM 常见错误归类标记
		if errors.Is(e, gorm.ErrRecordNotFound) {
			entry["gorm"] = "record_not_found"
		} else if errors.Is(e, gorm.ErrDuplicatedKey) {
			entry["gorm"] = "duplicated_key"
		} else if errors.Is(e, gorm.ErrInvalidTransaction) {
			entry["gorm"] = "invalid_transaction"
		}

		chain = append(chain, entry)
		e = errors.Unwrap(e)
	}

	root := chain[len(chain)-1]["error"]

	logEntry := log.WithFields(log.Fields{
		"id":          appErr.ID,
		"code":        appErr.Code,
		"statusCode":  appErr.StatusCode,
		"file":        appErr.file,
		"sql":         appErr.sql,
		"params":      appErr.params,
		"error_chain": chain,
		"root_error":  root,
	})

	if appErr.originalErr != nil {
		logEntry = logEntry.WithError(appErr.originalErr)
	}

	// 根据错误级别记录不同级别的日志
	if appErr.StatusCode >= 500 {
		logEntry.Error("系统错误: ", appErr.Message)
	} else if appErr.StatusCode >= 400 {
		logEntry.Warn("客户端错误: ", appErr.Message)
	} else {
		logEntry.Info("应用错误: ", appErr.Message)
	}
}

// sendResponse 发送响应
func (s *ResponseSender) sendResponse(ctx fiber.Ctx, appErr *AppError) error {
	contentType := utils.ToLower(utils.UnsafeString(ctx.Request().Header.ContentType()))
	contentType = binder.FilterFlags(utils.ParseVendorSpecificContentType(contentType))

	switch contentType {
	case fiber.MIMETextXML, fiber.MIMEApplicationXML:
		return ctx.Status(appErr.StatusCode).XML(appErr)
	case fiber.MIMETextHTML, fiber.MIMETextPlain:
		return ctx.Status(appErr.StatusCode).SendString(appErr.Message)
	default:
		return ctx.Status(appErr.StatusCode).JSON(appErr)
	}
}
