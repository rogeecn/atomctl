package otel

import (
	"context"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

var (
	tracer trace.Tracer
	meter  metric.Meter
)

func Start(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return tracer.Start(ctx, spanName, opts...)
}

func Int64Counter(name string, options ...metric.Int64CounterOption) (metric.Int64Counter, error) {
	return meter.Int64Counter(name, options...)
}

// Int64UpDownCounter
func Int64UpDownCounter(name string, options ...metric.Int64UpDownCounterOption) (metric.Int64UpDownCounter, error) {
	return meter.Int64UpDownCounter(name, options...)
}

// Int64Histogram
func Int64Histogram(name string, options ...metric.Int64HistogramOption) (metric.Int64Histogram, error) {
	return meter.Int64Histogram(name, options...)
}

// Int64Gauge
func Int64Gauge(name string, options ...metric.Int64GaugeOption) (metric.Int64Gauge, error) {
	return meter.Int64Gauge(name, options...)
}

// Int64ObservableCounter
func Int64ObservableCounter(name string, options ...metric.Int64ObservableCounterOption) (metric.Int64ObservableCounter, error) {
	return meter.Int64ObservableCounter(name, options...)
}

// Int64ObservableUpDownCounter
func Int64ObservableUpDownCounter(name string, options ...metric.Int64ObservableUpDownCounterOption) (metric.Int64ObservableUpDownCounter, error) {
	return meter.Int64ObservableUpDownCounter(name, options...)
}

// Int64ObservableGauge
func Int64ObservableGauge(name string, options ...metric.Int64ObservableGaugeOption) (metric.Int64ObservableGauge, error) {
	return meter.Int64ObservableGauge(name, options...)
}

// Float64Counter
func Float64Counter(name string, options ...metric.Float64CounterOption) (metric.Float64Counter, error) {
	return meter.Float64Counter(name, options...)
}

// Float64UpDownCounter
func Float64UpDownCounter(name string, options ...metric.Float64UpDownCounterOption) (metric.Float64UpDownCounter, error) {
	return meter.Float64UpDownCounter(name, options...)
}

// Float64Histogram
func Float64Histogram(name string, options ...metric.Float64HistogramOption) (metric.Float64Histogram, error) {
	return meter.Float64Histogram(name, options...)
}

// Float64Gauge
func Float64Gauge(name string, options ...metric.Float64GaugeOption) (metric.Float64Gauge, error) {
	return meter.Float64Gauge(name, options...)
}

// Float64ObservableCounter
func Float64ObservableCounter(name string, options ...metric.Float64ObservableCounterOption) (metric.Float64ObservableCounter, error) {
	return meter.Float64ObservableCounter(name, options...)
}

// Float64ObservableUpDownCounter
func Float64ObservableUpDownCounter(name string, options ...metric.Float64ObservableUpDownCounterOption) (metric.Float64ObservableUpDownCounter, error) {
	return meter.Float64ObservableUpDownCounter(name, options...)
}

// Float64ObservableGauge
func Float64ObservableGauge(name string, options ...metric.Float64ObservableGaugeOption) (metric.Float64ObservableGauge, error) {
	return meter.Float64ObservableGauge(name, options...)
}

// RegisterCallback
func RegisterCallback(f metric.Callback, instruments ...metric.Observable) (metric.Registration, error) {
	return meter.RegisterCallback(f, instruments...)
}
