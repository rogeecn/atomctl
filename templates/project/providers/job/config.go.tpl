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

type Config struct {
	// Optional per-queue worker concurrency. If empty, defaults apply.
	QueueWorkers QueueWorkersConfig
}

// QueueWorkers allows configuring worker concurrency per queue.
// Key is the queue name, value is MaxWorkers. If empty, defaults are used.
// Example TOML:
//
//	[Job]
//	# high=20, default=10, low=5
//	# QueueWorkers = { high = 20, default = 10, low = 5 }
type QueueWorkersConfig map[string]int

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

// queueConfig returns a river.QueueConfig map built from QueueWorkers or defaults.
func (c *Config) queueConfig() map[string]river.QueueConfig {
	cfg := map[string]river.QueueConfig{}
	if c == nil || len(c.QueueWorkers) == 0 {
		cfg[QueueHigh] = river.QueueConfig{MaxWorkers: 10}
		cfg[QueueDefault] = river.QueueConfig{MaxWorkers: 10}
		cfg[QueueLow] = river.QueueConfig{MaxWorkers: 10}
		return cfg
	}
	for name, n := range c.QueueWorkers {
		if n <= 0 {
			n = 1
		}
		cfg[name] = river.QueueConfig{MaxWorkers: n}
	}

	if _, ok := cfg[QueueDefault]; !ok {
		cfg[QueueDefault] = river.QueueConfig{MaxWorkers: 10}
	}
	return cfg
}
