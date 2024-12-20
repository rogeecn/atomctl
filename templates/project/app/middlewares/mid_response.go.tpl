package middlewares

import (
	"{{.ModuleName}}/app/errorx"

	"github.com/gofiber/fiber/v3"
	log "github.com/sirupsen/logrus"
)

func (f *Middlewares) ProcessResponse(c fiber.Ctx) error {
	err := c.Next()
	if err != nil {
		log.WithError(err).Error("process response error")

		if e, ok := err.(errorx.Response); ok {
			return e.Response(c)
		}

		if e, ok := err.(*fiber.Error); ok {
			return errorx.Response{
				StatusCode: e.Code,
				Code:       e.Code,
				Message:    e.Message,
			}.Response(c)
		}

		return errorx.Wrap(err).Response(c)

	}
	return err
}
