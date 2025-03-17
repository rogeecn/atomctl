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

		queue := &Job{Workers: workers, driver: pool, ctx: ctx, periodicJobs: make(map[string]rivertype.PeriodicJobHandle), jobs: make(map[string]*rivertype.JobInsertResult)}
		container.AddCloseAble(queue.Close)

		return queue, nil
	}, o.DiOptions()...)
}

type Job struct {
	Workers *river.Workers

	l sync.Mutex

	ctx    context.Context
	driver *riverpgxv5.Driver

	client       *river.Client[pgx.Tx]
	periodicJobs map[string]rivertype.PeriodicJobHandle
	jobs         map[string]*rivertype.JobInsertResult
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

func (q *Job) AddPeriodicJobs(job contracts.CronJob) error {
	for _, job := range job.Args() {
		if err := q.AddPeriodicJob(job); err != nil {
			return err
		}
	}
	return nil
}

func (q *Job) AddPeriodicJob(job contracts.CronJobArg) error {
	client, err := q.Client()
	if err != nil {
		return err
	}

	q.periodicJobs[job.Arg.Kind()] = client.PeriodicJobs().Add(river.NewPeriodicJob(
		job.PeriodicInterval,
		func() (river.JobArgs, *river.InsertOpts) {
			return job.Arg, lo.ToPtr(job.Arg.InsertOpts())
		},
		&river.PeriodicJobOpts{
			RunOnStart: job.RunOnStart,
		},
	))

	return nil
}

func (q *Job) Cancel(kind string) error {
	client, err := q.Client()
	if err != nil {
		return err
	}

	if h, ok := q.periodicJobs[kind]; ok {
		client.PeriodicJobs().Remove(h)
		delete(q.periodicJobs, kind)
		return nil
	}

	if r, ok := q.jobs[kind]; ok {
		_, err = client.JobCancel(q.ctx, r.Job.ID)
		if err != nil {
			return err
		}
		delete(q.jobs, kind)
		return nil
	}

	return errors.New("job by kind(" + kind + ") not found")
}

func (q *Job) Add(job contracts.JobArgs) error {
	client, err := q.Client()
	if err != nil {
		return err
	}

	q.jobs[job.Kind()], err = client.Insert(q.ctx, job, lo.ToPtr(job.InsertOpts()))
	return err
}
