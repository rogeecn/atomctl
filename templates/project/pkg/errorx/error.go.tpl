package errorx

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v3"
)

type Response struct {
	StatusCode int    `json:"-"`
	Code       int    `json:"code"`
	Message    string `json:"message"`
}

func Wrap(err error) Response {
	return Response{http.StatusInternalServerError, http.StatusInternalServerError, err.Error()}
}

func (r Response) Error() string {
	return fmt.Sprintf("[%d] %s", r.Code, r.Message)
}

func (r Response) Response(ctx fiber.Ctx) error {
	return ctx.Status(r.StatusCode).JSON(r)
}

var (
	RequestParseError = Response{http.StatusBadRequest, http.StatusBadRequest, "请求解析错误"}
	InternalError     = Response{http.StatusInternalServerError, http.StatusInternalServerError, "内部错误"}
)
