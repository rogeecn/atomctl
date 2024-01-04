package test

import (
	"a/common"
	"strings"

	"github.com/atom-apps/auth/modules/auth/controller"
	"github.com/atom-providers/log"
	"github.com/gofiber/fiber/v2"
	"github.com/rogeecn/atom/contracts"
)

// @provider
type UserRoutes struct {
	svc            contracts.HttpService
	userController *UserController
}

func (r *UserRoutes) Register() contracts.HttpRoute {
	engine := r.svc.GetEngine().(*fiber.App)
	group := engine.Group("tests")
	log.Infof("register route group: %s", group.(*fiber.Group).Prefix)

	r.routeUserController(group, userController)
	return nil
}

func (r *UserRoutes) routeUserController(engine fiber.Router, controller *controller.UserController) {
	groupPrefix := "/" + strings.TrimLeft(engine.(*fiber.Group).Prefix, "/")
	engine.Get(strings.TrimPrefix("/users/:id<int>", groupPrefix), DataFunc1(controller.Show, Integer[uint64]("id", PathParamError)))
	engine.Get(strings.TrimPrefix("/users", groupPrefix), DataFunc3(controller.List, Query[UserListQueryFilter](QueryParamError), Query[common.PageQueryFilter](QueryParamError), Query[common.SortQueryFilter](QueryParamError)))
	engine.Post(strings.TrimPrefix("/users", groupPrefix), Func1(controller.Create, Body[UserForm](BodyParamError)))
	engine.Put(strings.TrimPrefix("/users/:id<int>", groupPrefix), Func2(controller.Update, Integer[uint64]("id", PathParamError), Body[UserForm](BodyParamError)))
	engine.Delete(strings.TrimPrefix("/users/:id<int>", groupPrefix), Func1(controller.Delete, Integer[uint64]("id", PathParamError)))
}
