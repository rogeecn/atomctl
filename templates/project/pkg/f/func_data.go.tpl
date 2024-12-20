package f

import (
	"github.com/gofiber/fiber/v3"
)

func DataFunc[T any](
	f func(fiber.Ctx) (T, error),
) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		data, err := f(ctx)
		if err != nil {
			return nil
		}
		return ctx.JSON(data)
	}
}

func DataFunc1[T any, P1 any](
	f func(fiber.Ctx, P1) (T, error),
	pf1 func(fiber.Ctx) (P1, error),
) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		p, err := pf1(ctx)
		if err != nil {
			return nil
		}

		data, err := f(ctx, p)
		if err != nil {
			return nil
		}
		return ctx.JSON(data)
	}
}

func DataFunc2[T any, P1 any, P2 any](
	f func(fiber.Ctx, P1, P2) (T, error),
	pf1 func(fiber.Ctx) (P1, error),
	pf2 func(fiber.Ctx) (P2, error),
) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		p1, err := pf1(ctx)
		if err != nil {
			return nil
		}
		p2, err := pf2(ctx)
		if err != nil {
			return nil
		}
		data, err := f(ctx, p1, p2)
		if err != nil {
			return nil
		}
		return ctx.JSON(data)
	}
}

func DataFunc3[T any, P1 any, P2 any, P3 any](
	f func(fiber.Ctx, P1, P2, P3) (T, error),
	pf1 func(fiber.Ctx) (P1, error),
	pf2 func(fiber.Ctx) (P2, error),
	pf3 func(fiber.Ctx) (P3, error),
) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		p1, err := pf1(ctx)
		if err != nil {
			return nil
		}
		p2, err := pf2(ctx)
		if err != nil {
			return nil
		}
		p3, err := pf3(ctx)
		if err != nil {
			return nil
		}
		data, err := f(ctx, p1, p2, p3)
		if err != nil {
			return nil
		}
		return ctx.JSON(data)
	}
}

func DataFunc4[T any, P1 any, P2 any, P3 any, P4 any](
	f func(fiber.Ctx, P1, P2, P3, P4) (T, error),
	pf1 func(fiber.Ctx) (P1, error),
	pf2 func(fiber.Ctx) (P2, error),
	pf3 func(fiber.Ctx) (P3, error),
	pf4 func(fiber.Ctx) (P4, error),
) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		p1, err := pf1(ctx)
		if err != nil {
			return nil
		}
		p2, err := pf2(ctx)
		if err != nil {
			return nil
		}
		p3, err := pf3(ctx)
		if err != nil {
			return nil
		}
		p4, err := pf4(ctx)
		if err != nil {
			return nil
		}
		data, err := f(ctx, p1, p2, p3, p4)
		if err != nil {
			return nil
		}
		return ctx.JSON(data)
	}
}

func DataFunc5[T any, P1 any, P2 any, P3 any, P4 any, P5 any](
	f func(fiber.Ctx, P1, P2, P3, P4, P5) (T, error),
	pf1 func(fiber.Ctx) (P1, error),
	pf2 func(fiber.Ctx) (P2, error),
	pf3 func(fiber.Ctx) (P3, error),
	pf4 func(fiber.Ctx) (P4, error),
	pf5 func(fiber.Ctx) (P5, error),
) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		p1, err := pf1(ctx)
		if err != nil {
			return nil
		}
		p2, err := pf2(ctx)
		if err != nil {
			return nil
		}
		p3, err := pf3(ctx)
		if err != nil {
			return nil
		}
		p4, err := pf4(ctx)
		if err != nil {
			return nil
		}
		p5, err := pf5(ctx)
		if err != nil {
			return nil
		}
		data, err := f(ctx, p1, p2, p3, p4, p5)
		if err != nil {
			return nil
		}
		return ctx.JSON(data)
	}
}

func DataFunc6[T any, P1 any, P2 any, P3 any, P4 any, P5 any, P6 any](
	f func(fiber.Ctx, P1, P2, P3, P4, P5, P6) (T, error),
	pf1 func(fiber.Ctx) (P1, error),
	pf2 func(fiber.Ctx) (P2, error),
	pf3 func(fiber.Ctx) (P3, error),
	pf4 func(fiber.Ctx) (P4, error),
	pf5 func(fiber.Ctx) (P5, error),
	pf6 func(fiber.Ctx) (P6, error),
) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		p1, err := pf1(ctx)
		if err != nil {
			return nil
		}
		p2, err := pf2(ctx)
		if err != nil {
			return nil
		}
		p3, err := pf3(ctx)
		if err != nil {
			return nil
		}
		p4, err := pf4(ctx)
		if err != nil {
			return nil
		}
		p5, err := pf5(ctx)
		if err != nil {
			return nil
		}
		p6, err := pf6(ctx)
		if err != nil {
			return nil
		}
		data, err := f(ctx, p1, p2, p3, p4, p5, p6)
		if err != nil {
			return nil
		}
		return ctx.JSON(data)
	}
}

func DataFunc7[T any, P1 any, P2 any, P3 any, P4 any, P5 any, P6 any, P7 any](
	f func(fiber.Ctx, P1, P2, P3, P4, P5, P6, P7) (T, error),
	pf1 func(fiber.Ctx) (P1, error),
	pf2 func(fiber.Ctx) (P2, error),
	pf3 func(fiber.Ctx) (P3, error),
	pf4 func(fiber.Ctx) (P4, error),
	pf5 func(fiber.Ctx) (P5, error),
	pf6 func(fiber.Ctx) (P6, error),
	pf7 func(fiber.Ctx) (P7, error),
) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		p1, err := pf1(ctx)
		if err != nil {
			return nil
		}
		p2, err := pf2(ctx)
		if err != nil {
			return nil
		}
		p3, err := pf3(ctx)
		if err != nil {
			return nil
		}
		p4, err := pf4(ctx)
		if err != nil {
			return nil
		}
		p5, err := pf5(ctx)
		if err != nil {
			return nil
		}
		p6, err := pf6(ctx)
		if err != nil {
			return nil
		}
		p7, err := pf7(ctx)
		if err != nil {
			return nil
		}
		data, err := f(ctx, p1, p2, p3, p4, p5, p6, p7)
		if err != nil {
			return nil
		}
		return ctx.JSON(data)
	}
}

func DataFunc8[T any, P1 any, P2 any, P3 any, P4 any, P5 any, P6 any, P7 any, P8 any](
	f func(fiber.Ctx, P1, P2, P3, P4, P5, P6, P7, P8) (T, error),
	pf1 func(fiber.Ctx) (P1, error),
	pf2 func(fiber.Ctx) (P2, error),
	pf3 func(fiber.Ctx) (P3, error),
	pf4 func(fiber.Ctx) (P4, error),
	pf5 func(fiber.Ctx) (P5, error),
	pf6 func(fiber.Ctx) (P6, error),
	pf7 func(fiber.Ctx) (P7, error),
	pf8 func(fiber.Ctx) (P8, error),
) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		p1, err := pf1(ctx)
		if err != nil {
			return nil
		}
		p2, err := pf2(ctx)
		if err != nil {
			return nil
		}
		p3, err := pf3(ctx)
		if err != nil {
			return nil
		}
		p4, err := pf4(ctx)
		if err != nil {
			return nil
		}
		p5, err := pf5(ctx)
		if err != nil {
			return nil
		}
		p6, err := pf6(ctx)
		if err != nil {
			return nil
		}
		p7, err := pf7(ctx)
		if err != nil {
			return nil
		}
		p8, err := pf8(ctx)
		if err != nil {
			return nil
		}
		data, err := f(ctx, p1, p2, p3, p4, p5, p6, p7, p8)
		if err != nil {
			return nil
		}
		return ctx.JSON(data)
	}
}

func DataFunc9[T any, P1 any, P2 any, P3 any, P4 any, P5 any, P6 any, P7 any, P8 any, P9 any](
	f func(fiber.Ctx, P1, P2, P3, P4, P5, P6, P7, P8, P9) (T, error),
	pf1 func(fiber.Ctx) (P1, error),
	pf2 func(fiber.Ctx) (P2, error),
	pf3 func(fiber.Ctx) (P3, error),
	pf4 func(fiber.Ctx) (P4, error),
	pf5 func(fiber.Ctx) (P5, error),
	pf6 func(fiber.Ctx) (P6, error),
	pf7 func(fiber.Ctx) (P7, error),
	pf8 func(fiber.Ctx) (P8, error),
	pf9 func(fiber.Ctx) (P9, error),
) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		p1, err := pf1(ctx)
		if err != nil {
			return nil
		}
		p2, err := pf2(ctx)
		if err != nil {
			return nil
		}
		p3, err := pf3(ctx)
		if err != nil {
			return nil
		}
		p4, err := pf4(ctx)
		if err != nil {
			return nil
		}
		p5, err := pf5(ctx)
		if err != nil {
			return nil
		}
		p6, err := pf6(ctx)
		if err != nil {
			return nil
		}
		p7, err := pf7(ctx)
		if err != nil {
			return nil
		}
		p8, err := pf8(ctx)
		if err != nil {
			return nil
		}
		p9, err := pf9(ctx)
		if err != nil {
			return nil
		}
		data, err := f(ctx, p1, p2, p3, p4, p5, p6, p7, p8, p9)
		if err != nil {
			return nil
		}
		return ctx.JSON(data)
	}
}

func DataFunc10[T any, P1 any, P2 any, P3 any, P4 any, P5 any, P6 any, P7 any, P8 any, P9 any, P10 any](
	f func(fiber.Ctx, P1, P2, P3, P4, P5, P6, P7, P8, P9, P10) (T, error),
	pf1 func(fiber.Ctx) (P1, error),
	pf2 func(fiber.Ctx) (P2, error),
	pf3 func(fiber.Ctx) (P3, error),
	pf4 func(fiber.Ctx) (P4, error),
	pf5 func(fiber.Ctx) (P5, error),
	pf6 func(fiber.Ctx) (P6, error),
	pf7 func(fiber.Ctx) (P7, error),
	pf8 func(fiber.Ctx) (P8, error),
	pf9 func(fiber.Ctx) (P9, error),
	pf10 func(fiber.Ctx) (P10, error),
) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		p1, err := pf1(ctx)
		if err != nil {
			return nil
		}
		p2, err := pf2(ctx)
		if err != nil {
			return nil
		}
		p3, err := pf3(ctx)
		if err != nil {
			return nil
		}
		p4, err := pf4(ctx)
		if err != nil {
			return nil
		}
		p5, err := pf5(ctx)
		if err != nil {
			return nil
		}
		p6, err := pf6(ctx)
		if err != nil {
			return nil
		}
		p7, err := pf7(ctx)
		if err != nil {
			return nil
		}
		p8, err := pf8(ctx)
		if err != nil {
			return nil
		}
		p9, err := pf9(ctx)
		if err != nil {
			return nil
		}
		p10, err := pf10(ctx)
		if err != nil {
			return nil
		}
		data, err := f(ctx, p1, p2, p3, p4, p5, p6, p7, p8, p9, p10)
		if err != nil {
			return nil
		}
		return ctx.JSON(data)
	}
}
