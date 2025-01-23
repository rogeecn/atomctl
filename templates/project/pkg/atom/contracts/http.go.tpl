package contracts

import (
	"github.com/gofiber/fiber/v3"
)

type HttpRoute interface {
	Register(fiber.Router)
	Name() string
}
