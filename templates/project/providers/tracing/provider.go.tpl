package tracing

import (
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	config "github.com/uber/jaeger-client-go/config"
	"go.ipao.vip/atom/container"
	"go.ipao.vip/atom/opt"
)

func Provide(opts ...opt.Option) error {
	o := opt.New(opts...)
	var conf Config
	if err := o.UnmarshalConfig(&conf); err != nil {
		return err
	}
	conf.format()
	return container.Container.Provide(func() (opentracing.Tracer, error) {
		log := logrus.New()
		log.SetFormatter(&logrus.JSONFormatter{TimestampFormat: time.RFC3339})

		cfg := &config.Configuration{
			ServiceName: conf.Name,
			Disabled:    conf.Disabled,
			Sampler: &config.SamplerConfig{
				Type:                    conf.Sampler_Type,
				Param:                   conf.Sampler_Param,
				SamplingServerURL:       conf.Sampler_SamplingServerURL,
				MaxOperations:           conf.Sampler_MaxOperations,
				SamplingRefreshInterval: time.Duration(conf.Sampler_RefreshIntervalSec) * time.Second,
			},
			Reporter: &config.ReporterConfig{
				LogSpans: func() bool {
					if conf.Reporter_LogSpans == nil {
						return true
					}
					return *conf.Reporter_LogSpans
				}(),
				LocalAgentHostPort:  conf.Reporter_LocalAgentHostPort,
				CollectorEndpoint:   conf.Reporter_CollectorEndpoint,
				BufferFlushInterval: time.Duration(conf.Reporter_BufferFlushMs) * time.Millisecond,
				QueueSize:           conf.Reporter_QueueSize,
			},
			RPCMetrics: conf.RPCMetrics,
			Gen128Bit:  conf.Gen128Bit,
		}

		if len(conf.Tags) > 0 {
			cfg.Tags = make([]opentracing.Tag, 0, len(conf.Tags))
			for k, v := range conf.Tags {
				cfg.Tags = append(cfg.Tags, opentracing.Tag{Key: k, Value: v})
			}
		}

		jLogger := &jaegerLogrus{logger: log}
		zipkinShared := conf.ZipkinSharedRPCSpan
		if !zipkinShared {
			zipkinShared = true
		}
		tracer, closer, err := cfg.NewTracer(
			config.Logger(jLogger),
			config.ZipkinSharedRPCSpan(zipkinShared),
		)
		if err != nil {
			log.Errorf("无法初始化 Jaeger: %v", err)
			return opentracing.NoopTracer{}, nil
		}
		opentracing.SetGlobalTracer(tracer)
		container.AddCloseAble(func() { _ = closer.Close() })

		return tracer, nil
	}, o.DiOptions()...)
}
