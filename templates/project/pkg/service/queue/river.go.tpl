package queue

import (
	"context"

	"demo01/app/jobs"
	"demo01/pkg/service"
	"demo01/providers/app"
	"demo01/providers/job"
	"demo01/providers/postgres"

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
		job.DefaultProvider(),
	}...)
}

func Command() atom.Option {
	return atom.Command(
		atom.Name("queue"),
		atom.Short("start queue processor"),
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
	Initials []contracts.Initial `group:"initials"`
}

func Serve(cmd *cobra.Command, args []string) error {
	return container.Container.Invoke(func(ctx context.Context, svc Service) error {
		log.SetFormatter(&log.JSONFormatter{})

		if svc.App.IsDevMode() {
			log.SetLevel(log.DebugLevel)
		}

		client, err := svc.Job.Client()
		if err != nil {
			return err
		}

		if err := client.Start(ctx); err != nil {
			return err
		}
		defer client.StopAndCancel(ctx)

		<-ctx.Done()
		return nil
	})
}
