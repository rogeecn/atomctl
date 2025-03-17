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

// Prepare implements contracts.CronJob.
func (CronJob) Prepare() error {
	return nil
}

// JobArgs implements contracts.CronJob.
func (CronJob) Args() []contracts.CronJobArg {
	return []contracts.CronJobArg{
		{
			Arg: SortArgs{
				Strings: []string{"a", "b", "c", "d"},
			},

			PeriodicInterval: PeriodicInterval(time.Second * 10),
			RunOnStart:       false,
		},
	}
}
