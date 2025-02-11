package job

import (
	"github.com/riverqueue/river"
	"go.ipao.vip/atom/container"
	"go.ipao.vip/atom/opt"
)

const DefaultPrefix = "Job"

func DefaultProvider() container.ProviderContainer {
	return container.ProviderContainer{
		Provider: Provide,
		Options: []opt.Option{
			opt.Prefix(DefaultPrefix),
		},
	}
}

type Config struct{}

const (
	PriorityDefault = river.PriorityDefault
	PriorityLow     = 2
	PriorityMiddle  = 3
	PriorityHigh    = 3
)

const (
	QueueHigh    = "high"
	QueueDefault = river.QueueDefault
	QueueLow     = "low"
)
