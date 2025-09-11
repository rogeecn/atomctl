package otel

import (
	"context"
	"os"
	"time"

	"go.ipao.vip/atom/container"
	"go.ipao.vip/atom/contracts"
	"go.ipao.vip/atom/opt"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.15.0"
	"google.golang.org/grpc/encoding/gzip"
)

// formatAuth formats token into an Authorization header value.
func formatAuth(token string) string {
    if token == "" {
        return ""
    }
    if len(token) > 7 && (token[:7] == "Bearer " || token[:7] == "bearer ") {
        return token
    }
    return "Bearer " + token
}

func Provide(opts ...opt.Option) error {
	o := opt.New(opts...)
	var config Config
	if err := o.UnmarshalConfig(&config); err != nil {
		return err
	}
	config.format()
	return container.Container.Provide(func(ctx context.Context) (contracts.Initial, error) {
		o := &builder{
			config: &config,
		}

		if err := o.initResource(ctx); err != nil {
			return o, errors.Wrapf(err, "Failed to create OpenTelemetry resource")
		}

		if err := o.initMeterProvider(ctx); err != nil {
			return o, errors.Wrapf(err, "Failed to create OpenTelemetry metric provider")
		}

		if err := o.initTracerProvider(ctx); err != nil {
			return o, errors.Wrapf(err, "Failed to create OpenTelemetry tracer provider")
		}

		tracer = otel.Tracer(config.ServiceName)
		meter = otel.Meter(config.ServiceName)

		log.Info("otel provider init success")
		return o, nil
	}, o.DiOptions()...)
}

type builder struct {
	config *Config

	resource *resource.Resource
}

func (o *builder) initResource(ctx context.Context) (err error) {
	hostName, _ := os.Hostname()

	o.resource, err = resource.New(
		ctx,
		resource.WithFromEnv(),
		resource.WithProcess(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithOS(),
		resource.WithContainer(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(o.config.ServiceName),   // 应用名
			semconv.ServiceVersionKey.String(o.config.Version),    // 应用版本
			semconv.DeploymentEnvironmentKey.String(o.config.Env), // 部署环境
			semconv.HostNameKey.String(hostName),                  // 主机名
		),
	)
	return err
}

func (o *builder) initMeterProvider(ctx context.Context) (err error) {
	exporterGrpcFunc := func(ctx context.Context) (sdkmetric.Exporter, error) {
		opts := []otlpmetricgrpc.Option{
			otlpmetricgrpc.WithEndpoint(o.config.EndpointGRPC),
			otlpmetricgrpc.WithCompressor(gzip.Name),
		}

		if h := formatAuth(o.config.Token); h != "" {
			opts = append(opts, otlpmetricgrpc.WithHeaders(map[string]string{"Authorization": h}))
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
        if o.config.InsecureHTTP {
            opts = append(opts, otlpmetrichttp.WithInsecure())
        }
        if h := formatAuth(o.config.Token); h != "" {
            opts = append(opts, otlpmetrichttp.WithHeaders(map[string]string{"Authorization": h}))
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
		return err
	}

    // periodic reader with optional custom interval
    var readerOpts []sdkmetric.PeriodicReaderOption
    if o.config.MetricReaderIntervalMs > 0 {
        readerOpts = append(readerOpts, sdkmetric.WithInterval(time.Duration(o.config.MetricReaderIntervalMs)*time.Millisecond))
    }
    meterProvider := sdkmetric.NewMeterProvider(
        sdkmetric.WithReader(
            sdkmetric.NewPeriodicReader(exporter, readerOpts...),
        ),
        sdkmetric.WithResource(o.resource),
    )
	otel.SetMeterProvider(meterProvider)

    interval := 5 * time.Second
    if o.config.RuntimeReadMemStatsIntervalMs > 0 {
        interval = time.Duration(o.config.RuntimeReadMemStatsIntervalMs) * time.Millisecond
    }
    err = runtime.Start(runtime.WithMinimumReadMemStatsInterval(interval))
	if err != nil {
		return errors.Wrapf(err, "Failed to start runtime metrics")
	}

	container.AddCloseAble(func() {
		if err := meterProvider.Shutdown(ctx); err != nil {
			otel.Handle(err)
		}
	})

	return err
}

func (o *builder) initTracerProvider(ctx context.Context) error {
	exporterGrpcFunc := func(ctx context.Context) (*otlptrace.Exporter, error) {
        opts := []otlptracegrpc.Option{
            otlptracegrpc.WithCompressor(gzip.Name),
            otlptracegrpc.WithEndpoint(o.config.EndpointGRPC),
        }
        if o.config.InsecureGRPC {
            opts = append(opts, otlptracegrpc.WithInsecure())
        }

		if h := formatAuth(o.config.Token); h != "" {
			opts = append(opts, otlptracegrpc.WithHeaders(map[string]string{"Authorization": h}))
		}

		log.Debugf("Creating GRPC trace exporter with endpoint: %s", o.config.EndpointGRPC)
		exporter, err := otlptrace.New(ctx, otlptracegrpc.NewClient(opts...))
		if err != nil {
			return nil, errors.Wrap(err, "failed to create GRPC trace exporter")
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
            otlptracehttp.WithCompression(1),
            otlptracehttp.WithEndpoint(o.config.EndpointHTTP),
        }
        if o.config.InsecureHTTP {
            opts = append(opts, otlptracehttp.WithInsecure())
        }

		if h := formatAuth(o.config.Token); h != "" {
			opts = append(opts, otlptracehttp.WithHeaders(map[string]string{"Authorization": h}))
		}

		log.Debugf("Creating HTTP trace exporter with endpoint: %s", o.config.EndpointHTTP)
		exporter, err := otlptrace.New(ctx, otlptracehttp.NewClient(opts...))
		if err != nil {
			return nil, errors.Wrap(err, "failed to create HTTP trace exporter")
		}

		return exporter, nil
	}

	var exporter *otlptrace.Exporter
	var err error
	if o.config.EndpointHTTP != "" {
		exporter, err = exporterHttpFunc(ctx)
		log.Infof("otel http exporter: %s", o.config.EndpointHTTP)
	} else {
		exporter, err = exporterGrpcFunc(ctx)
		log.Infof("otel grpc exporter: %s", o.config.EndpointGRPC)
	}

	if err != nil {
		return err
	}

    // Sampler
    sampler := sdktrace.AlwaysSample()
    if o.config.Sampler == "ratio" {
        ratio := o.config.SamplerRatio
        if ratio <= 0 {
            ratio = 0
        }
        if ratio > 1 {
            ratio = 1
        }
        sampler = sdktrace.ParentBased(sdktrace.TraceIDRatioBased(ratio))
    }

    // Batcher options
    var batchOpts []sdktrace.BatchSpanProcessorOption
    if o.config.BatchTimeoutMs > 0 {
        batchOpts = append(batchOpts, sdktrace.WithBatchTimeout(time.Duration(o.config.BatchTimeoutMs)*time.Millisecond))
    }
    if o.config.ExportTimeoutMs > 0 {
        batchOpts = append(batchOpts, sdktrace.WithExportTimeout(time.Duration(o.config.ExportTimeoutMs)*time.Millisecond))
    }
    if o.config.MaxQueueSize > 0 {
        batchOpts = append(batchOpts, sdktrace.WithMaxQueueSize(o.config.MaxQueueSize))
    }
    if o.config.MaxExportBatchSize > 0 {
        batchOpts = append(batchOpts, sdktrace.WithMaxExportBatchSize(o.config.MaxExportBatchSize))
    }

    traceProvider := sdktrace.NewTracerProvider(
        sdktrace.WithSampler(sampler),
        sdktrace.WithResource(o.resource),
        sdktrace.WithBatcher(exporter, batchOpts...),
    )
	container.AddCloseAble(func() {
		log.Error("shut down")
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
