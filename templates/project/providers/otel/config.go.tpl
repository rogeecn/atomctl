package otel

import (
	"os"

	"go.ipao.vip/atom"
	"go.ipao.vip/atom/container"
	"go.ipao.vip/atom/opt"
)

const DefaultPrefix = "OTEL"

func DefaultProvider() container.ProviderContainer {
	return container.ProviderContainer{
		Provider: Provide,
		Options: []opt.Option{
			opt.Prefix(DefaultPrefix),
			opt.Group(atom.GroupInitialName),
		},
	}
}

type Config struct {
	ServiceName string
	Version     string
	Env         string

	EndpointGRPC string
	EndpointHTTP string
	Token        string
}

func (c *Config) format() {
	if c.ServiceName == "" {
		c.ServiceName = os.Getenv("SERVICE_NAME")
		if c.ServiceName == "" {
			c.ServiceName = "unknown"
		}
	}

	if c.Version == "" {
		c.Version = os.Getenv("SERVICE_VERSION")
		if c.Version == "" {
			c.Version = "unknown"
		}
	}

	if c.Env == "" {
		c.Env = os.Getenv("DEPLOY_ENVIRONMENT")
		if c.Env == "" {
			c.Env = "unknown"
		}
	}
}
