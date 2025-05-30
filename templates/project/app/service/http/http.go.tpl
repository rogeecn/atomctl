package http

import (
	"context"

	"go.ipao.vip/atom"
	"go.ipao.vip/atom/container"
	"go.ipao.vip/atom/contracts"
	"{{.ModuleName}}/app/errorx"
	"{{.ModuleName}}/app/jobs"
	"{{.ModuleName}}/app/service"
	_ "{{.ModuleName}}/docs"
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

type Service struct {
	dig.In

	App      *app.Config
	Job      *job.Job
	Http     *http.Service
	Initials []contracts.Initial   `group:"initials"`
	Routes   []contracts.HttpRoute `group:"routes"`
}

func Serve(cmd *cobra.Command, args []string) error {
	return container.Container.Invoke(func(ctx context.Context, svc Service) error {
		log.SetFormatter(&log.JSONFormatter{})

		if svc.App.Mode == app.AppModeDevelopment {
			log.SetLevel(log.DebugLevel)

			svc.Http.Engine.Get("/swagger/*", swagger.HandlerDefault)
		}
		svc.Http.Engine.Use(errorx.Middleware)
		svc.Http.Engine.Use(favicon.New(favicon.Config{
			Data: []byte{},
		}))

		group := svc.Http.Engine.Group("")
		for _, route := range svc.Routes {
			route.Register(group)
		}

		return svc.Http.Serve()
	})
}
