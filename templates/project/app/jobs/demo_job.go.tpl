package jobs

import (
	"context"
	"sort"
	"time"

	. "github.com/riverqueue/river"
	log "github.com/sirupsen/logrus"
	_ "go.ipao.vip/atom"
	"go.ipao.vip/atom/contracts"
	_ "go.ipao.vip/atom/contracts"
)

var _ contracts.JobArgs = DemoJob{}

type DemoJob struct {
	Strings []string `json:"strings"`
}

func (s DemoJob) InsertOpts() InsertOpts {
	return InsertOpts{
		Queue:    QueueDefault,
		Priority: PriorityDefault,
	}
}

func (DemoJob) Kind() string       { return "demo_job" }
func (a DemoJob) UniqueID() string { return a.Kind() }

var _ Worker[DemoJob] = (*DemoJobWorker)(nil)

// @provider(job)
type DemoJobWorker struct {
	WorkerDefaults[DemoJob]
}

func (w *DemoJobWorker) NextRetry(job *Job[DemoJob]) time.Time {
	return time.Now().Add(30 * time.Second)
}

func (w *DemoJobWorker) Work(ctx context.Context, job *Job[DemoJob]) error {
	logger := log.WithField("job", job.Args.Kind())

	logger.Infof("[START] %s args: %v", job.Args.Kind(), job.Args.Strings)
	defer logger.Infof("[END] %s", job.Args.Kind())

	// modify below
	sort.Strings(job.Args.Strings)
	logger.Infof("[%s] Sorted strings: %v\n", time.Now().Format(time.TimeOnly), job.Args.Strings)

	return nil
}
