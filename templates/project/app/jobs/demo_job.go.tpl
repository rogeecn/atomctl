package jobs

import (
	"context"
	"sort"
	"time"

	. "github.com/riverqueue/river"
	log "github.com/sirupsen/logrus"
	_ "{{.ModuleName}}/pkg/atom"
	_ "{{.ModuleName}}/pkg/atom/contracts"
)

// provider:[except|only] [returnType] [group]

var (
	_ JobArgs               = SortArgs{}
	_ JobArgsWithInsertOpts = SortArgs{}
)

type SortArgs struct {
	Strings []string `json:"strings"`
}

// InsertOpts implements JobArgsWithInsertOpts.
func (s SortArgs) InsertOpts() InsertOpts {
	return InsertOpts{
		Queue:    QueueDefault,
		Priority: PriorityDefault,
	}
}

func (SortArgs) Kind() string {
	return "sort"
}

var _ Worker[SortArgs] = (*SortWorker)(nil)

// @provider(job)
type SortWorker struct {
	WorkerDefaults[SortArgs]
}

func (w *SortWorker) Work(ctx context.Context, job *Job[SortArgs]) error {
	sort.Strings(job.Args.Strings)

	log.Infof("[%s] Sorted strings: %v\n", time.Now().Format(time.TimeOnly), job.Args.Strings)
	return nil
}

func (w *SortWorker) NextRetry(job *Job[SortArgs]) time.Time {
	return time.Now().Add(5 * time.Second)
}
