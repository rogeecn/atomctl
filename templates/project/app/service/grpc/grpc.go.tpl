package grpc

import (
	"go.ipao.vip/atom"
	"go.ipao.vip/atom/container"
	"go.ipao.vip/atom/contracts"
	"{{.ModuleName}}/app/grpc/users"
	"{{.ModuleName}}/app/service"
	"{{.ModuleName}}/providers/app"
	"{{.ModuleName}}/providers/grpc"
	"{{.ModuleName}}/providers/postgres"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go.uber.org/dig"
)

func defaultProviders() container.Providers {
	return service.Default(container.Providers{
		postgres.DefaultProvider(),
		grpc.DefaultProvider(),
	}...)
}

func Command() atom.Option {
	return atom.Command(
		atom.Name("grpc"),
		atom.Short("run grpc server"),
		atom.RunE(Serve),
		atom.Providers(
			defaultProviders().
				With(
					users.Provide,
				),
		),
	)
}

type Service struct {
	dig.In

	App      *app.Config
	Grpc     *grpc.Grpc
	Initials []contracts.Initial `group:"initials"`
}

func Serve(cmd *cobra.Command, args []string) error {
	return container.Container.Invoke(func(svc Service) error {
		log.SetFormatter(&log.JSONFormatter{})

		if svc.App.IsDevMode() {
			log.SetLevel(log.DebugLevel)
		}

		return svc.Grpc.Serve()
	})
}
