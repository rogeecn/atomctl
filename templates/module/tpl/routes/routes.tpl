package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rogeecn/atom"
	"github.com/rogeecn/atom/container"
	"github.com/rogeecn/atom-addons/providers/http"
	"github.com/rogeecn/atom-addons/providers/log"
	"github.com/rogeecn/atom/utils/opt"
)

func Provide(opts ...opt.Option) error {
	return container.Container.Provide(newRoute, atom.GroupRoutes)
}

func newRoute (svc http.Service) http.Route {
	engine := svc.GetEngine().(*fiber.App)
	group := engine.Group("{{.NamePlural}}")
	log.Infof("register route group: %s", group.(*fiber.Group).Prefix)

	return nil
}