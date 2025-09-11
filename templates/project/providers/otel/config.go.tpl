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

    // Connection security
    InsecureGRPC bool // if true, use grpc insecure for OTLP gRPC
    InsecureHTTP bool // if true, use http insecure for OTLP HTTP

    // Tracing sampler
    // Sampler: "always" (default) or "ratio"
    Sampler      string
    SamplerRatio float64 // used when Sampler == "ratio"; 0..1

    // Tracing batcher options (milliseconds)
    BatchTimeoutMs      uint
    ExportTimeoutMs     uint
    MaxQueueSize        int
    MaxExportBatchSize  int

    // Metrics options (milliseconds)
    MetricReaderIntervalMs        uint // export interval for PeriodicReader
    RuntimeReadMemStatsIntervalMs uint // runtime metrics min read interval
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
