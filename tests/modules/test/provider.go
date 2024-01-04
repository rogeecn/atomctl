package test

import (
	"github.com/rogeecn/atom"
	"github.com/rogeecn/atom/container"
	"github.com/rogeecn/atom/utils/opt"
)

func Provide(opts ...opt.Option) error {
	if err := container.Container.Provide(newRoute, atom.GroupRoutes); err != nil {
		return err
	}

	if err := container.Container.Provide(func(
		userSvc *UserService,
	) (*UserController, error) {
		obj := &UserController{
			userSvc: userSvc,
		}
		return obj, nil
	}); err != nil {
		return err
	}

	if err := container.Container.Provide(func() (*UserService, error) {
		obj := &UserService{}
		return obj, nil
	}); err != nil {
		return err
	}

	return nil
}
