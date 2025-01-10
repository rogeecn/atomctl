package event

import (
	"context"

	"{{.ModuleName}}/app/events"
	"{{.ModuleName}}/app/service"
	"{{.ModuleName}}/providers/app"
	"{{.ModuleName}}/providers/event"
	"{{.ModuleName}}/providers/postgres"

	"git.ipao.vip/rogeecn/atom"
	"git.ipao.vip/rogeecn/atom/container"
	"git.ipao.vip/rogeecn/atom/contracts"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go.uber.org/dig"
)

func defaultProviders() container.Providers {
	return service.Default(container.Providers{
		postgres.DefaultProvider(),
	}...)
}

func Command() atom.Option {
	return atom.Command(
		atom.Name("event"),
		atom.Short("start event processor"),
		atom.RunE(Serve),
		atom.Providers(
			defaultProviders().
				With(
					events.Provide,
				),
		),
	)
}

type Service struct {
	dig.In

	App      *app.Config
	PubSub   *event.PubSub
	Initials []contracts.Initial `group:"initials"`
}

func Serve(cmd *cobra.Command, args []string) error {
	return container.Container.Invoke(func(ctx context.Context, svc Service) error {
		log.SetFormatter(&log.JSONFormatter{})

		if svc.App.IsDevMode() {
			log.SetLevel(log.DebugLevel)
		}

		return svc.PubSub.Serve(ctx)
	})
}
