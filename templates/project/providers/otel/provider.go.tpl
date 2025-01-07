package otel

import (
	"context"
	"os"
	"time"

	"git.ipao.vip/rogeecn/atom/container"
	"git.ipao.vip/rogeecn/atom/utils/opt"

	"github.com/pkg/errors"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
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
			config: &config,
		}

		err := o.initResource(ctx)
		if err != nil {
			return o, errors.Wrapf(err, "Failed to create OpenTelemetry resource")
		}

		if err := o.initMeterProvider(ctx); err != nil {
			return o, errors.Wrapf(err, "Failed to create OpenTelemetry metric provider")
		}

		if err := o.initTracerProvider(ctx); err != nil {
			return o, errors.Wrapf(err, "Failed to create OpenTelemetry tracer provider")
		}

		o.Tracer = otel.Tracer(config.ServiceName)
		o.Meter = otel.Meter(config.ServiceName)

		return o, nil
	}, o.DiOptions()...)
}

type OTEL struct {
	config *Config

	Tracer trace.Tracer
	Meter  metric.Meter

	resource *resource.Resource
}

func (o *OTEL) initResource(ctx context.Context) (err error) {
	hostName, _ := os.Hostname()

	o.resource, err = resource.New(
		ctx,
		resource.WithFromEnv(),
		resource.WithProcess(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(o.config.ServiceName),   // 应用名
			semconv.ServiceVersionKey.String(o.config.Version),    // 应用版本
			semconv.DeploymentEnvironmentKey.String(o.config.Env), // 部署环境
			semconv.HostNameKey.String(hostName),                  // 主机名
		),
	)
	return
}

func (o *OTEL) initMeterProvider(ctx context.Context) (err error) {
	exporterGrpcFunc := func(ctx context.Context) (sdkmetric.Exporter, error) {
		opts := []otlpmetricgrpc.Option{
			otlpmetricgrpc.WithEndpoint(o.config.EndpointGRPC),
			otlpmetricgrpc.WithCompressor(gzip.Name),
		}

		if o.config.Token != "" {
			headers := map[string]string{"Authentication": o.config.Token}
			opts = append(opts, otlpmetricgrpc.WithHeaders(headers))
		}

		exporter, err := otlpmetricgrpc.New(ctx, opts...)
		if err != nil {
			return nil, err
		}
		return exporter, nil
	}

	exporterHttpFunc := func(ctx context.Context) (sdkmetric.Exporter, error) {
		opts := []otlpmetrichttp.Option{
			otlpmetrichttp.WithEndpoint(o.config.EndpointHTTP),
			otlpmetrichttp.WithCompression(1),
		}

		if o.config.Token != "" {
			opts = append(opts, otlpmetrichttp.WithURLPath(o.config.Token))
		}

		exporter, err := otlpmetrichttp.New(ctx, opts...)
		if err != nil {
			return nil, err
		}
		return exporter, nil
	}

	var exporter sdkmetric.Exporter
	if o.config.EndpointHTTP != "" {
		exporter, err = exporterHttpFunc(ctx)
	} else {
		exporter, err = exporterGrpcFunc(ctx)
	}

	if err != nil {
		return
	}

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(
			sdkmetric.NewPeriodicReader(exporter),
		),
		sdkmetric.WithResource(o.resource),
	)
	otel.SetMeterProvider(meterProvider)

	err = runtime.Start(runtime.WithMinimumReadMemStatsInterval(time.Second * 5))
	if err != nil {
		return errors.Wrapf(err, "Failed to start runtime metrics")
	}

	container.AddCloseAble(func() {
		if err := meterProvider.Shutdown(ctx); err != nil {
			otel.Handle(err)
		}
	})

	return
}

func (o *OTEL) initTracerProvider(ctx context.Context) error {
	exporterGrpcFunc := func(ctx context.Context) (*otlptrace.Exporter, error) {
		opts := []otlptracegrpc.Option{
			otlptracegrpc.WithCompressor(gzip.Name),
			otlptracegrpc.WithEndpoint(o.config.EndpointGRPC),
		}

		if o.config.Token != "" {
			headers := map[string]string{"Authentication": o.config.Token}
			opts = append(opts, otlptracegrpc.WithHeaders(headers))
		}

		exporter, err := otlptrace.New(ctx, otlptracegrpc.NewClient(opts...))
		if err != nil {
			return nil, err
		}

		container.AddCloseAble(func() {
			cxt, cancel := context.WithTimeout(ctx, time.Second)
			defer cancel()
			if err := exporter.Shutdown(cxt); err != nil {
				otel.Handle(err)
			}
		})

		return exporter, nil
	}

	exporterHttpFunc := func(ctx context.Context) (*otlptrace.Exporter, error) {
		opts := []otlptracehttp.Option{
			otlptracehttp.WithInsecure(),
			otlptracehttp.WithCompression(1),
			otlptracehttp.WithEndpoint(o.config.EndpointHTTP),
		}

		if o.config.Token != "" {
			opts = append(opts, otlptracehttp.WithURLPath(o.config.Token))
		}

		exporter, err := otlptrace.New(ctx, otlptracehttp.NewClient(opts...))
		if err != nil {
			return nil, err
		}

		return exporter, nil
	}

	var exporter *otlptrace.Exporter
	var err error
	if o.config.EndpointHTTP != "" {
		exporter, err = exporterHttpFunc(ctx)
	} else {
		exporter, err = exporterGrpcFunc(ctx)
	}

	if err != nil {
		return err
	}

	traceProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(o.resource),
		sdktrace.WithSpanProcessor(sdktrace.NewBatchSpanProcessor(exporter)),
	)
	container.AddCloseAble(func() {
		if err := traceProvider.Shutdown(ctx); err != nil {
			otel.Handle(err)
		}
	})

	otel.SetTracerProvider(traceProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return err
}
