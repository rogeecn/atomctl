package redis

import (
	"fmt"

	"git.ipao.vip/rogeecn/atom/container"
	"git.ipao.vip/rogeecn/atom/utils/opt"
)

const DefaultPrefix = "Redis"

func DefaultProvider() container.ProviderContainer {
	return container.ProviderContainer{
		Provider: Provide,
		Options: []opt.Option{
			opt.Prefix(DefaultPrefix),
		},
	}
}

type Config struct {
	Host     string
	Port     uint
	Password string
	DB       uint
}

func (c *Config) format() {
	if c.Host == "" {
		c.Host = "localhost"
	}

	if c.Port == 0 {
		c.Port = 6379
	}

	if c.DB == 0 {
		c.DB = 0
	}
}

func (c *Config) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
