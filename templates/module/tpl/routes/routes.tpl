package routes

import (
	"github.com/gin-gonic/gin"
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
	engine := svc.GetEngine().(*gin.Engine)
	group := engine.Group("{{.NamePlural}}")
	log.Info("register route group: %s", group)

	return nil
}