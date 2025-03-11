package jobs

import (
	"time"

	. "github.com/riverqueue/river"
	"github.com/sirupsen/logrus"
	_ "go.ipao.vip/atom"
	"go.ipao.vip/atom/contracts"
)

var _ contracts.CronJob = (*CronJob)(nil)

// @provider(cronjob)
type CronJob struct {
	log *logrus.Entry `inject:"false"`
}

func (j *CronJob) Prepare() error {
	j.log = logrus.WithField("module", "cron")
	return nil
}

func (CronJob) Kind() string {
	return "cron_job"
}

// InsertOpts implements contracts.CronJob.
func (CronJob) InsertOpts() InsertOpts {
	return InsertOpts{
		MaxAttempts: 1,
	}
}

// JobArgs implements contracts.CronJob.
func (CronJob) JobArgs() JobArgs {
	return SortArgs{
		Strings: []string{"a", "c", "b", "d"},
	}
}

// Periodic implements contracts.CronJob.
func (cron *CronJob) Periodic() PeriodicSchedule {
	return PeriodicInterval(time.Minute)
}

// RunOnStart implements contracts.CronJob.
func (CronJob) RunOnStart() bool {
	return true
}
