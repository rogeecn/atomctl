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

		err = f(ctx, p)
		if err != nil {
			return err
		}
		return nil
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

		err = f(ctx, p1, p2)
		if err != nil {
			return err
		}
		return nil
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
		err = f(ctx, p1, p2, p3)
		if err != nil {
			return nil
		}
		return nil
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

		err = f(ctx, p1, p2, p3, p4)
		if err != nil {
			return nil
		}
		return nil
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
		err = f(ctx, p1, p2, p3, p4, p5)
		if err != nil {
			return nil
		}
		return nil
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
		err = f(ctx, p1, p2, p3, p4, p5, p6)
		if err != nil {
			return nil
		}
		return nil
	}
}
