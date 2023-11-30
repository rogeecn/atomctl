package common

import (
	"strings"

	"github.com/atom-providers/jwt"
	"github.com/gofiber/fiber/v2"
)

func GetJwtToken(ctx *fiber.Ctx) (string, error) {
	headers, ok := ctx.GetReqHeaders()[jwt.HttpHeader]
	if !ok {
		return "", ctx.SendStatus(fiber.StatusUnauthorized)
	}
	if len(headers) == 0 {
		return "", ctx.SendStatus(fiber.StatusUnauthorized)
	}
	token := headers[0]

	token = token[len(jwt.TokenPrefix):]
	return token, nil
}
