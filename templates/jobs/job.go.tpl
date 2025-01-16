package jobs

import (
	"context"
	"time"

	_ "git.ipao.vip/rogeecn/atom"
	_ "git.ipao.vip/rogeecn/atom/contracts"
	. "github.com/riverqueue/river"
)

var (
	_ JobArgs               = (*{{.Name}}Job)(nil)
	_ JobArgsWithInsertOpts = (*{{.Name}}Job)(nil)
)

type {{.Name}}Job struct {
}

// InsertOpts implements JobArgsWithInsertOpts.
func (s {{.Name}}Job) InsertOpts() InsertOpts {
	return InsertOpts{
		Queue:    QueueDefault,
		Priority: PriorityDefault,
		// UniqueOpts: UniqueOpts{
		// 	ByArgs: true,
		// },
	}
}

func ({{.Name}}Job) Kind() string {
	return "{{.Name}}Job"
}

var _ Worker[{{.Name}}Job] = (*{{.Name}}JobWorker)(nil)

// @provider(job)
type {{.Name}}JobWorker struct {
	WorkerDefaults[{{.Name}}Job]
}

func (w *{{.Name}}JobWorker) NextRetry(job *Job[{{.Name}}Job]) time.Time {
	return time.Now().Add(5 * time.Second)
}

func (w *{{.Name}}JobWorker) Work(ctx context.Context, job *Job[{{.Name}}Job]) error {
	return nil
}
