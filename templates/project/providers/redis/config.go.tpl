package redis

import (
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"go.ipao.vip/atom/container"
	"go.ipao.vip/atom/opt"
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
	Host       string
	Port       uint
	Password   string
	DB         uint
	Username   string
	ClientName string

	// Pool & retries
	PoolSize     int
	MinIdleConns int
	MaxRetries   int

	// Timeouts (seconds); 0 = library default
	DialTimeoutSeconds     uint
	ReadTimeoutSeconds     uint
	WriteTimeoutSeconds    uint
	ConnMaxIdleTimeSeconds uint
	ConnMaxLifetimeSeconds uint
	PingTimeoutSeconds     uint
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

	if c.PingTimeoutSeconds == 0 {
		c.PingTimeoutSeconds = 5
	}
}

func (c *Config) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// DialTimeout returns the dial timeout duration, 0 means library default
func (c *Config) DialTimeout() time.Duration {
	if c.DialTimeoutSeconds == 0 {
		return 0
	}
	return time.Duration(c.DialTimeoutSeconds) * time.Second
}

// ReadTimeout returns the read timeout duration, 0 means library default
func (c *Config) ReadTimeout() time.Duration {
	if c.ReadTimeoutSeconds == 0 {
		return 0
	}
	return time.Duration(c.ReadTimeoutSeconds) * time.Second
}

// WriteTimeout returns the write timeout duration, 0 means library default
func (c *Config) WriteTimeout() time.Duration {
	if c.WriteTimeoutSeconds == 0 {
		return 0
	}
	return time.Duration(c.WriteTimeoutSeconds) * time.Second
}

// ConnMaxIdleTime returns the max idle time for a connection, 0 means default
func (c *Config) ConnMaxIdleTime() time.Duration {
	if c.ConnMaxIdleTimeSeconds == 0 {
		return 0
	}
	return time.Duration(c.ConnMaxIdleTimeSeconds) * time.Second
}

// ConnMaxLifetime returns the max lifetime for a connection, 0 means default
func (c *Config) ConnMaxLifetime() time.Duration {
	if c.ConnMaxLifetimeSeconds == 0 {
		return 0
	}
	return time.Duration(c.ConnMaxLifetimeSeconds) * time.Second
}

// PingTimeout returns the ping timeout duration, default 5s
func (c *Config) PingTimeout() time.Duration {
	if c.PingTimeoutSeconds == 0 {
		return 5 * time.Second
	}
	return time.Duration(c.PingTimeoutSeconds) * time.Second
}

// Options builds a redis.Options based on the config values.
// Only non-zero/meaningful values are set, to preserve library defaults.
func (c *Config) Options() *redis.Options {
	c.format()
	ro := &redis.Options{
		Addr:       c.Addr(),
		Username:   c.Username,
		Password:   c.Password,
		DB:         int(c.DB),
		ClientName: c.ClientName,
	}
	if c.PoolSize > 0 {
		ro.PoolSize = c.PoolSize
	}
	if c.MinIdleConns > 0 {
		ro.MinIdleConns = c.MinIdleConns
	}
	if c.MaxRetries > 0 {
		ro.MaxRetries = c.MaxRetries
	}
	if dt := c.DialTimeout(); dt > 0 {
		ro.DialTimeout = dt
	}
	if rt := c.ReadTimeout(); rt > 0 {
		ro.ReadTimeout = rt
	}
	if wt := c.WriteTimeout(); wt > 0 {
		ro.WriteTimeout = wt
	}
	if it := c.ConnMaxIdleTime(); it > 0 {
		ro.ConnMaxIdleTime = it
	}
	if lt := c.ConnMaxLifetime(); lt > 0 {
		ro.ConnMaxLifetime = lt
	}
	return ro
}
