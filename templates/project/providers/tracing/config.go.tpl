package tracing

import (
	"github.com/sirupsen/logrus"
	"{{.ModuleName}}/pkg/atom/container"
	"{{.ModuleName}}/pkg/atom/opt"
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
}
