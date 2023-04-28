package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rogeecn/atom"
	"github.com/rogeecn/atom/container"
	"github.com/rogeecn/atom/providers/http"
	"github.com/rogeecn/atom/utils/opt"
	"github.com/rogeecn/gen"
)

func Provide(opts ...opt.Option) error {
	return container.Container.Provide(NewRoute, atom.GroupRoutes)
}


type Route struct {
	engine  *gin.Engine
}

func NewRoute(svc http.Service) http.Route {
	engine := svc.GetEngine().(*gin.Engine)
	return &Route{engine: engine}
}

func (r *Route) Register() {
	// r.engine.GET("/", gen.DataFunc(r.controller.Show))
}
