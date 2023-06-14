package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rogeecn/atom"
	"github.com/rogeecn/atom/container"
	"github.com/rogeecn/atom-addons/providers/http"
	"github.com/rogeecn/atom/utils/opt"
	"github.com/rogeecn/gen"
)

func Provide(opts ...opt.Option) error {
	newRoute := func (svc http.Service) http.Route {
		engine := svc.GetEngine().(*gin.Engine)
		// engine.GET("/", gen.DataFunc(controller.Show))
		return nil
	}
	return container.Container.Provide(newRoute, atom.GroupRoutes)
}