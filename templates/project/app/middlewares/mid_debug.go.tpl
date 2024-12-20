package middlewares

import (
	"github.com/gofiber/fiber/v3"
)

func (f *Middlewares) DebugMode(c fiber.Ctx) error {
	return c.Next()
}
