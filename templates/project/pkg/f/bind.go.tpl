package f

import (
	"github.com/gofiber/fiber/v3"
	"github.com/pkg/errors"
)

func Local[T any](key string) func(fiber.Ctx) (T, error) {
	return func(ctx fiber.Ctx) (T, error) {
		v := fiber.Locals[T](ctx, key)
		return v, nil
	}
}

func Path[T fiber.GenericType](key string) func(fiber.Ctx) (T, error) {
	return func(ctx fiber.Ctx) (T, error) {
		v := fiber.Params[T](ctx, key)
		return v, nil
	}
}

func PathParam[T fiber.GenericType](name string) func(fiber.Ctx) (T, error) {
	return func(ctx fiber.Ctx) (T, error) {
		v := fiber.Params[T](ctx, name)
		return v, nil
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

func QueryParam[T fiber.GenericType](key string) func(fiber.Ctx) (T, error) {
	return func(ctx fiber.Ctx) (T, error) {
		v := fiber.Query[T](ctx, key)
		return v, nil
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

func Cookie[T any](name string) func(fiber.Ctx) (*T, error) {
	return func(ctx fiber.Ctx) (*T, error) {
		p := new(T)
		if err := ctx.Bind().Cookie(p); err != nil {
			return nil, errors.Wrapf(err, "cookie: %s", name)
		}

		return p, nil
	}
}

func CookieParam(name string) func(fiber.Ctx) (string, error) {
	return func(ctx fiber.Ctx) (string, error) {
		return ctx.Cookies(name), nil
	}
}
