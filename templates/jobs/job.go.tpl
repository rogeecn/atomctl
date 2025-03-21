package jobs

import (
	"context"
	"time"

	. "github.com/riverqueue/river"
	_ "go.ipao.vip/atom"
	"go.ipao.vip/atom/contracts"
)

var _ contracts.JobArgs = (*{{.Name}})(nil)

type {{.Name}} struct {}

func (s {{.Name}}) InsertOpts() InsertOpts {
	return InsertOpts{
		Queue:    QueueDefault,
		Priority: PriorityDefault,
		// UniqueOpts: UniqueOpts{
		// 	ByArgs: true,
		// },
	}
}

func ({{.Name}}) Kind() string { return "{{.Name}}" }
func (arg {{.Name}}) UniqueID() string { return arg.Kind()}

var _ Worker[{{.Name}}] = (*{{.Name}}Worker)(nil)

// @provider(job)
type {{.Name}}Worker struct {
	WorkerDefaults[{{.Name}}]
}

func (w *{{.Name}}Worker) NextRetry(job *Job[{{.Name}}]) time.Time {
	return time.Now().Add(5 * time.Second)
}

func (w *{{.Name}}Worker) Work(ctx context.Context, job *Job[{{.Name}}]) error {
	return nil
}
