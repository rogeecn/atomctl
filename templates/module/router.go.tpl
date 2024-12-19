package {{.ModuleName}}

import (
	_ "git.ipao.vip/rogeecn/atom"
	_ "git.ipao.vip/rogeecn/atom/contracts"
	"github.com/gofiber/fiber/v3"
	log "github.com/sirupsen/logrus"
)

// @provider:except	contracts.HttpRoute	 atom.GroupRoutes
type Router struct {
	log *log.Entry `inject:"false"`

	controller *Controller
}

func (r *Router) Name() string {
	return "{{.ModuleName}}"
}

func (r *Router) Prepare() error {
	r.log = log.WithField("http.group", r.Name())
	return nil
}

func (r *Router) Register(router fiber.Router) {
	group := router.Group(r.Name())
	_ = group
}
