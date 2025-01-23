package contracts

import (
	"time"

	"github.com/riverqueue/river"
)

type CronJob interface {
	Description() string
	Periodic() time.Duration
	JobArgs() []river.JobArgs
	InsertOpts() *river.InsertOpts
	RunOnStart() bool
}
