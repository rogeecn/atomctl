package jobs

import (
	"time"

	"github.com/riverqueue/river"
	"github.com/sirupsen/logrus"
	_ "go.ipao.vip/atom"
	"go.ipao.vip/atom/contracts"
)

var _ contracts.CronJob = (*CronJob)(nil)

// @provider contracts.CronJob atom.GroupCronJob
type CronJob struct {
	log *logrus.Entry `inject:"false"`
}

func (cron *CronJob) Prepare() error {
	cron.log = logrus.WithField("module", "cron")
	return nil
}

func (cron *CronJob) Description() string {
	return "hello world cron job"
}

// InsertOpts implements contracts.CronJob.
func (cron *CronJob) InsertOpts() *river.InsertOpts {
	return nil
}

// JobArgs implements contracts.CronJob.
func (cron *CronJob) JobArgs() []river.JobArgs {
	return []river.JobArgs{
		SortArgs{
			Strings: []string{"a", "c", "b", "d"},
		},
	}
}

// Periodic implements contracts.CronJob.
func (cron *CronJob) Periodic() time.Duration {
	return time.Second * 10
}

// RunOnStart implements contracts.CronJob.
func (cron *CronJob) RunOnStart() bool {
	return true
}
