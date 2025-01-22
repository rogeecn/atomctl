package http

import (
	"{{.ModuleName}}/app/errorx"
	"{{.ModuleName}}/app/jobs"
	"{{.ModuleName}}/app/service"
	_ "{{.ModuleName}}/docs"
	"{{.ModuleName}}/pkg/atom"
	"{{.ModuleName}}/pkg/atom/container"
	"{{.ModuleName}}/pkg/atom/contracts"
	"{{.ModuleName}}/providers/app"
	"{{.ModuleName}}/providers/hashids"
	"{{.ModuleName}}/providers/http"
	"{{.ModuleName}}/providers/http/swagger"
	"{{.ModuleName}}/providers/job"
	"{{.ModuleName}}/providers/jwt"
	"{{.ModuleName}}/providers/postgres"

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
		job.DefaultProvider(),
	}...)
}

func Command() atom.Option {
	return atom.Command(
		atom.Name("serve"),
		atom.Short("run http server"),
		atom.RunE(Serve),
		atom.Providers(
			defaultProviders().
				With(
					jobs.Provide,
				),
		),
	)
}

type Http struct {
	dig.In

	App      *app.Config
	Job      *job.Job
	Service  *http.Service
	Initials []contracts.Initial   `group:"initials"`
	Routes   []contracts.HttpRoute `group:"routes"`
}

func Serve(cmd *cobra.Command, args []string) error {
	return container.Container.Invoke(func(http Http) error {
		log.SetFormatter(&log.JSONFormatter{})

		if http.App.Mode == app.AppModeDevelopment {
			log.SetLevel(log.DebugLevel)

			http.Service.Engine.Get("/swagger/*", swagger.HandlerDefault)
		}
		http.Service.Engine.Use(errorx.Middleware)
		http.Service.Engine.Use(favicon.New(favicon.Config{
			Data: []byte{},
		}))

		group := http.Service.Engine.Group("")
		for _, route := range http.Routes {
			route.Register(group)
		}

		return http.Service.Serve()
	})
}
