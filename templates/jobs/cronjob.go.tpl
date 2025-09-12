package jobs

import (
	"context"
	"time"

	. "github.com/riverqueue/river"
	_ "go.ipao.vip/atom"
	"go.ipao.vip/atom/contracts"
)

var _ contracts.CronJob = (*Cron{{.Name}})(nil)

// @provider(cronjob)
type Cron{{.Name}} struct{}

func (c *Cron{{.Name}}) Args() []contracts.CronJobArg {
	jobs := []contracts.CronJobArg{
		{
			RunOnStart:       true, // 启动时运行一次
			PeriodicInterval: cron.Every(time.Hour * 2), // 每2小时运行一次
			// PeriodicInterval: PeriodicInterval("* * * * *"), // 每天凌晨1点运行一次
			Arg:              {{.Name}}{},
		},
	}
	return jobs
}
