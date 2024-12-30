package tracing

import (
	"io"
	"time"

	"git.ipao.vip/rogeecn/atom/container"
	"git.ipao.vip/rogeecn/atom/utils/opt"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	config "github.com/uber/jaeger-client-go/config"
)

func Provide(opts ...opt.Option) error {
	o := opt.New(opts...)
	var conf Config
	if err := o.UnmarshalConfig(&conf); err != nil {
		return err
	}
	conf.format()
	return container.Container.Provide(func() (opentracing.Tracer, io.Closer, error) {
		log := logrus.New()
		log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		})

		cfg := &config.Configuration{
			ServiceName: conf.Name,
			Sampler: &config.SamplerConfig{
				Type:  "const",
				Param: 1,
			},
			Reporter: &config.ReporterConfig{
				LogSpans:            true,
				LocalAgentHostPort:  conf.Reporter_LocalAgentHostPort,
				CollectorEndpoint:   conf.Reporter_CollectorEndpoint,
				BufferFlushInterval: 100 * time.Millisecond,
				QueueSize:           1000,
			},
		}

		// 使用自定义的 logger
		jLogger := &jaegerLogrus{logger: log}
		tracer, closer, err := cfg.NewTracer(
			config.Logger(jLogger),
			config.ZipkinSharedRPCSpan(true),
		)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "无法初始化 Jaeger: %v", err)
		}
		opentracing.SetGlobalTracer(tracer)

		return tracer, closer, nil
	}, o.DiOptions()...)
}
