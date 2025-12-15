package job

import (
	"context"
	"fmt"
	"sync"
	"time"

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
		// health check ping with timeout
		pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()
		if err := dbPool.Ping(pingCtx); err != nil {
			return nil, fmt.Errorf("job provider: db ping failed: %w", err)
		}
		container.AddCloseAble(dbPool.Close)
		pool := riverpgxv5.New(dbPool)

		queue := &Job{
			Workers:      workers,
			driver:       pool,
			ctx:          ctx,
			conf:         &config,
			periodicJobs: make(map[string]rivertype.PeriodicJobHandle),
		}
		container.AddCloseAble(queue.Close)

		return queue, nil
	}, o.DiOptions()...)
}

type Job struct {
	ctx     context.Context
	conf    *Config
	Workers *river.Workers
	driver  *riverpgxv5.Driver

	l      sync.Mutex
	client *river.Client[pgx.Tx]

	periodicJobs map[string]rivertype.PeriodicJobHandle
}

func (q *Job) Close() {
	if q.client == nil {
		return
	}

	if err := q.client.StopAndCancel(q.ctx); err != nil {
		log.Errorf("Failed to stop and cancel client: %s", err)
	}
	// clear references
	q.l.Lock()
	q.periodicJobs = map[string]rivertype.PeriodicJobHandle{}
	q.l.Unlock()
}

func (q *Job) Client() (*river.Client[pgx.Tx], error) {
	q.l.Lock()
	defer q.l.Unlock()

	if q.client == nil {
		var err error
		q.client, err = river.NewClient(q.driver, &river.Config{
			Workers: q.Workers,
			Queues:  q.conf.queueConfig(),
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
	q.l.Lock()
	defer q.l.Unlock()

	q.periodicJobs[job.Arg.UniqueID()] = client.PeriodicJobs().Add(river.NewPeriodicJob(
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

func (q *Job) Cancel(id string) error {
	client, err := q.Client()
	if err != nil {
		return err
	}

	q.l.Lock()
	defer q.l.Unlock()

	if h, ok := q.periodicJobs[id]; ok {
		client.PeriodicJobs().Remove(h)
		delete(q.periodicJobs, id)
	}
	return nil
}

// CancelContext is like Cancel but allows passing a context.
func (q *Job) CancelContext(ctx context.Context, id string) error {
	client, err := q.Client()
	if err != nil {
		return err
	}
	q.l.Lock()
	defer q.l.Unlock()
	if h, ok := q.periodicJobs[id]; ok {
		client.PeriodicJobs().Remove(h)
		delete(q.periodicJobs, id)
		return nil
	}

	return nil
}

func (q *Job) Add(job contracts.JobArgs) error {
	client, err := q.Client()
	if err != nil {
		return err
	}

	q.l.Lock()
	defer q.l.Unlock()

	_, err = client.Insert(q.ctx, job, lo.ToPtr(job.InsertOpts()))
	return err
}
