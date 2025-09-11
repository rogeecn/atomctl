package tracing

import (
	"github.com/sirupsen/logrus"
	"go.ipao.vip/atom/container"
	"go.ipao.vip/atom/opt"
)

const DefaultPrefix = "Tracing"

func DefaultProvider() container.ProviderContainer {
	return container.ProviderContainer{
		Provider: Provide,
		Options: []opt.Option{
			opt.Prefix(DefaultPrefix),
		},
	}
}

// 自定义的 Logger 实现
type jaegerLogrus struct {
	logger *logrus.Logger
}

func (l *jaegerLogrus) Error(msg string) {
	l.logger.Error(msg)
}

func (l *jaegerLogrus) Infof(msg string, args ...interface{}) {
	l.logger.Infof(msg, args...)
}

type Config struct {
	Name                        string
	Reporter_LocalAgentHostPort string //:  "127.0.0.1:6831",
	Reporter_CollectorEndpoint  string //:   "http://127.0.0.1:14268/api/traces",
	Disabled                    bool
	Gen128Bit                   bool
	ZipkinSharedRPCSpan         bool
	RPCMetrics                  bool

	// Sampler configuration
	Sampler_Type               string
	Sampler_Param              float64
	Sampler_SamplingServerURL  string
	Sampler_MaxOperations      int
	Sampler_RefreshIntervalSec uint

	// Reporter configuration
	Reporter_LogSpans      *bool
	Reporter_BufferFlushMs uint
	Reporter_QueueSize     int

	// Process tags
	Tags map[string]string
}

func (c *Config) format() {
	if c.Reporter_LocalAgentHostPort == "" {
		c.Reporter_LocalAgentHostPort = "127.0.0.1:6831"
	}

	if c.Reporter_CollectorEndpoint == "" {
		c.Reporter_CollectorEndpoint = "http://127.0.0.1:14268/api/traces"
	}

	if c.Name == "" {
		c.Name = "default"
	}

	if c.Sampler_Type == "" {
		c.Sampler_Type = "const"
	}
	if c.Sampler_Param == 0 {
		c.Sampler_Param = 1
	}
	if c.Reporter_BufferFlushMs == 0 {
		c.Reporter_BufferFlushMs = 100
	}
	if c.Reporter_QueueSize == 0 {
		c.Reporter_QueueSize = 1000
	}
	if c.Reporter_LogSpans == nil {
		b := true
		c.Reporter_LogSpans = &b
	}
}
