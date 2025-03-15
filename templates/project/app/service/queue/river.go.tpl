package queue

import (
	"context"

	"go.ipao.vip/atom"
	"go.ipao.vip/atom/container"
	"go.ipao.vip/atom/contracts"
	"{{.ModuleName}}/app/jobs"
	"{{.ModuleName}}/app/service"
	"{{.ModuleName}}/providers/app"
	"{{.ModuleName}}/providers/job"
	"{{.ModuleName}}/providers/postgres"

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
	CronJobs []contracts.CronJob `group:"cron_jobs"`
}

func Serve(cmd *cobra.Command, args []string) error {
	return container.Container.Invoke(func(ctx context.Context, svc Service) error {
		log.SetFormatter(&log.JSONFormatter{})

		if svc.App.IsDevMode() {
			log.SetLevel(log.DebugLevel)
		}

		if err := svc.Job.Start(ctx); err != nil {
			return err
		}
		defer svc.Job.Close()

		<-ctx.Done()
		return nil
	})
}
