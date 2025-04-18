package jobs

import (
	"time"

	. "github.com/riverqueue/river"
	"github.com/sirupsen/logrus"
	_ "go.ipao.vip/atom"
	"go.ipao.vip/atom/contracts"
)

var _ contracts.CronJob = (*DemoCronJob)(nil)

// @provider(cronjob)
type DemoCronJob struct {
	log *logrus.Entry `inject:"false"`
}

// Prepare implements contracts.CronJob.
func (DemoCronJob) Prepare() error {
	return nil
}

// JobArgs implements contracts.CronJob.
func (DemoCronJob) Args() []contracts.CronJobArg {
	return []contracts.CronJobArg{
		{
			Arg: DemoJob{
				Strings: []string{"a", "b", "c", "d"},
			},

			PeriodicInterval: PeriodicInterval(time.Second * 10),
			RunOnStart:       false,
		},
	}
}
