package f

import (
	"github.com/gofiber/fiber/v3"
)

func Func(f fiber.Handler) fiber.Handler {
	return f
}

func Func1[P1 any](
	f func(fiber.Ctx, P1) error,
	pf1 func(fiber.Ctx) (P1, error),
) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		p, err := pf1(ctx)
		if err != nil {
			return err
		}

		return f(ctx, p)
	}
}

func Func2[P1 any, P2 any](
	f func(fiber.Ctx, P1, P2) error,
	pf1 func(fiber.Ctx) (P1, error),
	pf2 func(fiber.Ctx) (P2, error),
) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		p1, err := pf1(ctx)
		if err != nil {
			return err
		}

		p2, err := pf2(ctx)
		if err != nil {
			return err
		}

		return f(ctx, p1, p2)
	}
}

func Func3[P1 any, P2 any, P3 any](
	f func(fiber.Ctx, P1, P2, P3) error,
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
		return f(ctx, p1, p2, p3)
	}
}

func Func4[P1 any, P2 any, P3 any, P4 any](
	f func(fiber.Ctx, P1, P2, P3, P4) error,
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

		return f(ctx, p1, p2, p3, p4)
	}
}

func Func5[P1 any, P2 any, P3 any, P4 any, P5 any](
	f func(fiber.Ctx, P1, P2, P3, P4, P5) error,
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
		return f(ctx, p1, p2, p3, p4, p5)
	}
}

func Func6[P1 any, P2 any, P3 any, P4 any, P5 any, P6 any](
	f func(fiber.Ctx, P1, P2, P3, P4, P5, P6) error,
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
		return f(ctx, p1, p2, p3, p4, p5, p6)
	}
}

func Func7[P1 any, P2 any, P3 any, P4 any, P5 any, P6 any, P7 any](
	f func(fiber.Ctx, P1, P2, P3, P4, P5, P6, P7) error,
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
		return f(ctx, p1, p2, p3, p4, p5, p6, p7)
	}
}

func Func8[P1 any, P2 any, P3 any, P4 any, P5 any, P6 any, P7 any, P8 any](
	f func(fiber.Ctx, P1, P2, P3, P4, P5, P6, P7, P8) error,
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
		return f(ctx, p1, p2, p3, p4, p5, p6, p7, p8)
	}
}

func Func9[P1 any, P2 any, P3 any, P4 any, P5 any, P6 any, P7 any, P8 any, P9 any](
	f func(fiber.Ctx, P1, P2, P3, P4, P5, P6, P7, P8, P9) error,
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
		return f(ctx, p1, p2, p3, p4, p5, p6, p7, p8, p9)
	}
}

func Func10[P1 any, P2 any, P3 any, P4 any, P5 any, P6 any, P7 any, P8 any, P9 any, P10 any](
	f func(fiber.Ctx, P1, P2, P3, P4, P5, P6, P7, P8, P9, P10) error,
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
		return f(ctx, p1, p2, p3, p4, p5, p6, p7, p8, p9, p10)
	}
}

