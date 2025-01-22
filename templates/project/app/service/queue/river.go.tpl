package queue

import (
	"context"

	"{{.ModuleName}}/app/jobs"
	"{{.ModuleName}}/app/service"
	"{{.ModuleName}}/pkg/atom"
	"{{.ModuleName}}/pkg/atom/container"
	"{{.ModuleName}}/pkg/atom/contracts"
	"{{.ModuleName}}/providers/app"
	"{{.ModuleName}}/providers/job"
	"{{.ModuleName}}/providers/postgres"

	"github.com/riverqueue/river"
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

		client, err := svc.Job.Client()
		if err != nil {
			return err
		}

		for _, cronJob := range svc.CronJobs {
			log.
				WithField("module", "cron").
				WithField("name", cronJob.Description()).
				WithField("duration", cronJob.Periodic().Seconds()).
				Info("registering cron job")

			for _, jobArgs := range cronJob.JobArgs() {
				client.PeriodicJobs().Add(
					river.NewPeriodicJob(
						river.PeriodicInterval(cronJob.Periodic()),
						func() (river.JobArgs, *river.InsertOpts) {
							return jobArgs, cronJob.InsertOpts()
						},
						&river.PeriodicJobOpts{
							RunOnStart: cronJob.RunOnStart(),
						},
					),
				)
			}
		}

		if err := client.Start(ctx); err != nil {
			return err
		}
		defer client.StopAndCancel(ctx)

		<-ctx.Done()
		return nil
	})
}
