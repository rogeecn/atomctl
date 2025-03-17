package job

import (
	"context"
	"sync"

	"{{.ModuleName}}/providers/postgres"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"github.com/riverqueue/river/rivertype"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"go.ipao.vip/atom/container"
	"go.ipao.vip/atom/contracts"
	"go.ipao.vip/atom/opt"
)

func Provide(opts ...opt.Option) error {
	o := opt.New(opts...)
	var config Config
	if err := o.UnmarshalConfig(&config); err != nil {
		return err
	}
	return container.Container.Provide(func(ctx context.Context, dbConf *postgres.Config) (*Job, error) {
		workers := river.NewWorkers()

		dbPoolConfig, err := pgxpool.ParseConfig(dbConf.DSN())
		if err != nil {
			return nil, err
		}

		dbPool, err := pgxpool.NewWithConfig(ctx, dbPoolConfig)
		if err != nil {
			return nil, err
		}
		container.AddCloseAble(dbPool.Close)
		pool := riverpgxv5.New(dbPool)

		queue := &Job{Workers: workers, Driver: pool, ctx: ctx}
		container.AddCloseAble(queue.Close)

		return queue, nil
	}, o.DiOptions()...)
}

type Job struct {
	ctx     context.Context
	Workers *river.Workers
	Driver  *riverpgxv5.Driver

	l      sync.Mutex
	client *river.Client[pgx.Tx]
}

func (q *Job) Close() {
	if q.client == nil {
		return
	}

	if err := q.client.StopAndCancel(q.ctx); err != nil {
		log.Errorf("Failed to stop and cancel client: %s", err)
	}
}

func (q *Job) Client() (*river.Client[pgx.Tx], error) {
	q.l.Lock()
	defer q.l.Unlock()

	if q.client == nil {
		var err error
		q.client, err = river.NewClient(q.Driver, &river.Config{
			Workers: q.Workers,
			Queues: map[string]river.QueueConfig{
				QueueHigh:    {MaxWorkers: 10},
				QueueDefault: {MaxWorkers: 10},
				QueueLow:     {MaxWorkers: 10},
			},
		})
		if err != nil {
			return nil, err
		}
	}

	return q.client, nil
}

func (q *Job) Start(ctx context.Context) error {
	client, err := q.Client()
	if err != nil {
		return errors.Wrap(err, "get client failed")
	}

	if err := client.Start(ctx); err != nil {
		return err
	}
	defer client.StopAndCancel(ctx)

	<-ctx.Done()

	return nil
}

func (q *Job) StopAndCancel(ctx context.Context) error {
	client, err := q.Client()
	if err != nil {
		return errors.Wrap(err, "get client failed")
	}

	return client.StopAndCancel(ctx)
}

func (q *Job) AddPeriodicJobs(job contracts.CronJob) (map[string]rivertype.PeriodicJobHandle, error) {
	var err error

	handles := make(map[string]rivertype.PeriodicJobHandle)
	for _, job := range job.Args() {
		handles[job.Kind], err = q.AddPeriodicJob(job)
		if err != nil {
			return handles, err
		}
	}
	return handles, nil
}

func (q *Job) AddPeriodicJob(job contracts.CronJobArg) (rivertype.PeriodicJobHandle, error) {
	client, err := q.Client()
	if err != nil {
		return 0, err
	}

	return client.PeriodicJobs().Add(river.NewPeriodicJob(
		job.PeriodicInterval,
		func() (river.JobArgs, *river.InsertOpts) {
			return job.Arg, lo.ToPtr(job.Arg.InsertOpts())
		},
		&river.PeriodicJobOpts{
			RunOnStart: job.RunOnStart,
		},
	)), nil
}
