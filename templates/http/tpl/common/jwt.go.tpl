package common

import (
	"strings"

	"github.com/atom-providers/jwt"
	"github.com/gofiber/fiber/v2"
)

func GetJwtToken(ctx *fiber.Ctx) (string, error) {
	token, ok := ctx.GetReqHeaders()[jwt.HttpHeader]
	if !ok {
		return "", ctx.SendStatus(fiber.StatusUnauthorized)
	}

	if !strings.HasPrefix(token, jwt.TokenPrefix) {
		return "", ctx.SendStatus(fiber.StatusUnauthorized)
	}
	token = token[len(jwt.TokenPrefix):]
	return token, nil
}
