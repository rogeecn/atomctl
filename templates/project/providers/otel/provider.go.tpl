package otel

import (
	"context"
	"os"
	"time"

	"git.ipao.vip/rogeecn/atom/container"
	"git.ipao.vip/rogeecn/atom/utils/opt"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.15.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/encoding/gzip"
)

func Provide(opts ...opt.Option) error {
	o := opt.New(opts...)
	var config Config
	if err := o.UnmarshalConfig(&config); err != nil {
		return err
	}
	return container.Container.Provide(func(ctx context.Context) (*OTEL, error) {
		o := &OTEL{
			Tracer: otel.Tracer(config.ServiceName),
		}

		var err error
		o.Resource, err = initResource(ctx, &config)
		if err != nil {
			return o, errors.Wrapf(err, "Failed to create OpenTelemetry resource")
		}

		o.Exporter, o.SpanProcessor, err = initGrpcExporterAndSpanProcessor(ctx, &config)
		if err != nil {
			return o, errors.Wrapf(err, "Failed to create OpenTelemetry trace exporter")
		}

		traceProvider := sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithResource(o.Resource),
			sdktrace.WithSpanProcessor(o.SpanProcessor),
		)

		otel.SetTracerProvider(traceProvider)
		otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		))

		container.AddCloseAble(func() {
			cxt, cancel := context.WithTimeout(ctx, time.Second)
			defer cancel()
			if err := o.Exporter.Shutdown(cxt); err != nil {
				otel.Handle(err)
			}
		})

		return o, nil
	}, o.DiOptions()...)
}

type OTEL struct {
	Tracer        trace.Tracer
	Resource      *resource.Resource
	Exporter      *otlptrace.Exporter
	SpanProcessor sdktrace.SpanProcessor
}

func initResource(ctx context.Context, conf *Config) (*resource.Resource, error) {
	hostName, _ := os.Hostname()

	r, err := resource.New(
		ctx,
		resource.WithFromEnv(),
		resource.WithProcess(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(conf.ServiceName),   // 应用名
			semconv.ServiceVersionKey.String(conf.Version),    // 应用版本
			semconv.DeploymentEnvironmentKey.String(conf.Env), // 部署环境
			semconv.HostNameKey.String(hostName),              // 主机名
		),
	)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func initGrpcExporterAndSpanProcessor(ctx context.Context, conf *Config) (*otlptrace.Exporter, sdktrace.SpanProcessor, error) {
	opts := []otlptracegrpc.Option{
		otlptracegrpc.WithCompressor(gzip.Name),
	}

	if conf.Token != "" {
		headers := map[string]string{"Authentication": conf.Token}
		opts = append(opts, otlptracegrpc.WithHeaders(headers))
	}

	if conf.EndpointGRPC != "" {
		opts = append(opts, otlptracegrpc.WithEndpoint(conf.EndpointGRPC))
	}

	traceExporter, err := otlptrace.New(ctx, otlptracegrpc.NewClient(opts...))
	if err != nil {
		return nil, nil, err
	}

	batchSpanProcessor := sdktrace.NewBatchSpanProcessor(traceExporter)

	return traceExporter, batchSpanProcessor, nil
}
