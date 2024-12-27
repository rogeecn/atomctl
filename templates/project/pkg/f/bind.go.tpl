package f

import (
	"github.com/gofiber/fiber/v3"
	"github.com/pkg/errors"
)

func Path[T fiber.GenericType](key string) func(fiber.Ctx) (T, error) {
	return func(ctx fiber.Ctx) (T, error) {
		v := fiber.Params[T](ctx, key)
		return v, nil
	}
}

func URI[T any](name string) func(fiber.Ctx) (*T, error) {
	return func(ctx fiber.Ctx) (*T, error) {
		p := new(T)
		if err := ctx.Bind().URI(p); err != nil {
			return nil, errors.Wrapf(err, "uri: %s", name)
		}

		return p, nil
	}
}

func Body[T any](name string) func(fiber.Ctx) (*T, error) {
	return func(ctx fiber.Ctx) (*T, error) {
		p := new(T)
		if err := ctx.Bind().Body(p); err != nil {
			return nil, errors.Wrapf(err, "body: %s", name)
		}

		return p, nil
	}
}

func Query[T any](name string) func(fiber.Ctx) (*T, error) {
	return func(ctx fiber.Ctx) (*T, error) {
		p := new(T)
		if err := ctx.Bind().Query(p); err != nil {
			return nil, errors.Wrapf(err, "query: %s", name)
		}

		return p, nil
	}
}

func Header[T any](name string) func(fiber.Ctx) (*T, error) {
	return func(ctx fiber.Ctx) (*T, error) {
		p := new(T)
		err := ctx.Bind().Header(p)
		if err != nil {
			return nil, errors.Wrapf(err, "header: %s", name)
		}

		return p, nil
	}
}
