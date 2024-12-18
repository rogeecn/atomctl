package http

import (
	"{{.ModuleName}}/pkg/service"
	"{{.ModuleName}}/providers/app"
	"{{.ModuleName}}/providers/hashids"
	"{{.ModuleName}}/providers/http"
	"{{.ModuleName}}/providers/jwt"
	"{{.ModuleName}}/providers/postgres"

	"git.ipao.vip/rogeecn/atom"
	"git.ipao.vip/rogeecn/atom/container"
	"git.ipao.vip/rogeecn/atom/contracts"
	"github.com/gofiber/fiber/v3/middleware/favicon"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go.uber.org/dig"
)

func defaultProviders() container.Providers {
	return service.Default(container.Providers{
		http.DefaultProvider(),
		postgres.DefaultProvider(),
		jwt.DefaultProvider(),
		hashids.DefaultProvider(),
	}...)
}

func Command() atom.Option {
	return atom.Command(
		atom.Name("serve"),
		atom.Short("run http server"),
		atom.RunE(Serve),
		atom.Providers(defaultProviders()),
	)
}

type Http struct {
	dig.In

	App      *app.Config
	Service  *http.Service
	Initials []contracts.Initial   `group:"initials"`
	Routes   []contracts.HttpRoute `group:"routes"`
}

func Serve(cmd *cobra.Command, args []string) error {
	return container.Container.Invoke(func(http Http) error {
		if http.App.Mode == app.AppModeDevelopment {
			log.SetLevel(log.DebugLevel)
		}

		http.Service.Engine.Use(favicon.New(favicon.Config{
			Data: []byte{},
		}))

		group := http.Service.Engine.Group("/v1")

		for _, route := range http.Routes {
			route.Register(group)
		}

		return http.Service.Serve()
	})
}
