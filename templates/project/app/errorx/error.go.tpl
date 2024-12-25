package errorx

import (
	"errors"
	"fmt"
	"net/http"
	"runtime"

	"github.com/go-jet/jet/v2/qrm"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/binder"
	"github.com/gofiber/utils/v2"
	log "github.com/sirupsen/logrus"
)

func Middleware(c fiber.Ctx) error {
	err := c.Next()
	if err != nil {
		return Wrap(err).Response(c)
	}
	return err
}

type Response struct {
	isFormat bool
	err      error
	params   []any
	sql      string
	file     string

	StatusCode int    `json:"-" xml:"-"`
	Code       int    `json:"code" xml:"code"`
	Message    string `json:"message" xml:"message"`
}

func New(code, statusCode int, message string) *Response {
	return &Response{
		isFormat:   true,
		StatusCode: statusCode,
		Code:       code,
		Message:    message,
	}
}

func (r *Response) Sql(sql string) *Response {
	r.sql = sql
	return r
}

func (r *Response) Params(params ...any) *Response {
	r.params = params
	if _, file, line, ok := runtime.Caller(1); ok {
		r.file = fmt.Sprintf("%s:%d", file, line)
	}
	return r
}

func Wrap(err error) *Response {
	if e, ok := err.(*Response); ok {
		return e
	}
	return &Response{err: err}
}

func (r *Response) format() {
	r.isFormat = true
	if errors.Is(r.err, qrm.ErrNoRows) {
		r.Code = RecordNotExists.Code
		r.Message = RecordNotExists.Message
		r.StatusCode = RecordNotExists.StatusCode
		return
	}

	if e, ok := r.err.(*fiber.Error); ok {
		r.Code = e.Code
		r.Message = e.Message
		r.StatusCode = e.Code
		return
	}
}

func (r *Response) Error() string {
	if !r.isFormat {
		r.format()
	}

	return fmt.Sprintf("[%d] %s", r.Code, r.Message)
}

func (r *Response) Response(ctx fiber.Ctx) error {
	if !r.isFormat {
		r.format()
	}

	contentType := utils.ToLower(utils.UnsafeString(ctx.Context().Request.Header.ContentType()))
	contentType = binder.FilterFlags(utils.ParseVendorSpecificContentType(contentType))

	log.WithError(r.err).WithField("file", r.file).WithField("params", r.params).Errorf("response error: %+v", r)

	// Parse body accordingly
	switch contentType {
	case fiber.MIMETextXML, fiber.MIMEApplicationXML:
		return ctx.Status(r.StatusCode).XML(r)
	case fiber.MIMETextHTML, fiber.MIMETextPlain:
		return ctx.Status(r.StatusCode).SendString(r.Message)
	default:
		return ctx.Status(r.StatusCode).JSON(r.Message)
	}
}

var (
	RecordNotExists = New(http.StatusNotFound, http.StatusNotFound, "记录不存在")
	BadRequest      = New(http.StatusBadRequest, http.StatusBadRequest, "请求错误")
	Unauthorized    = New(http.StatusUnauthorized, http.StatusUnauthorized, "未授权")
	InternalErr     = New(http.StatusInternalServerError, http.StatusInternalServerError, "内部错误")
)
