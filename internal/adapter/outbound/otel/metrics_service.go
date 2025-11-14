package otel

import (
	"context"
	"time"

	"github.com/gieart87/gohexaclean/internal/port/outbound/telemetry"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

// MetricsServiceOTEL implements telemetry.MetricsService using OpenTelemetry
type MetricsServiceOTEL struct {
	meterProvider *sdkmetric.MeterProvider
	meter         metric.Meter
}

// NewMetricsServiceOTEL creates a new OpenTelemetry metrics service
func NewMetricsServiceOTEL(ctx context.Context, serviceName, collectorEndpoint string) (telemetry.MetricsService, error) {
	// Create OTLP exporter
	exporter, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithEndpoint(collectorEndpoint),
		otlpmetricgrpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	// Create resource
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
		),
	)
	if err != nil {
		return nil, err
	}

	// Create meter provider
	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter)),
		sdkmetric.WithResource(res),
	)

	// Set global meter provider
	otel.SetMeterProvider(meterProvider)

	// Get meter
	meter := meterProvider.Meter(serviceName)

	return &MetricsServiceOTEL{
		meterProvider: meterProvider,
		meter:         meter,
	}, nil
}

// IncrementCounter increments a counter metric
func (m *MetricsServiceOTEL) IncrementCounter(name string, tags map[string]string, value float64) {
	counter, err := m.meter.Float64Counter(name)
	if err != nil {
		return
	}

	attrs := convertTagsToAttributes(tags)
	counter.Add(context.Background(), value, metric.WithAttributes(attrs...))
}

// SetGauge sets a gauge metric
func (m *MetricsServiceOTEL) SetGauge(name string, tags map[string]string, value float64) {
	// OTEL doesn't have direct gauge support, use observable gauge
	attrs := convertTagsToAttributes(tags)
	_, _ = m.meter.Float64ObservableGauge(name,
		metric.WithFloat64Callback(func(_ context.Context, observer metric.Float64Observer) error {
			observer.Observe(value, metric.WithAttributes(attrs...))
			return nil
		}),
	)
}

// RecordHistogram records a histogram metric
func (m *MetricsServiceOTEL) RecordHistogram(name string, tags map[string]string, value float64) {
	histogram, err := m.meter.Float64Histogram(name)
	if err != nil {
		return
	}

	attrs := convertTagsToAttributes(tags)
	histogram.Record(context.Background(), value, metric.WithAttributes(attrs...))
}

// RecordDistribution records a distribution metric (using histogram in OTEL)
func (m *MetricsServiceOTEL) RecordDistribution(name string, tags map[string]string, value float64) {
	m.RecordHistogram(name, tags, value)
}

// RecordTiming records a timing metric
func (m *MetricsServiceOTEL) RecordTiming(name string, tags map[string]string, duration time.Duration) {
	histogram, err := m.meter.Float64Histogram(name)
	if err != nil {
		return
	}

	attrs := convertTagsToAttributes(tags)
	histogram.Record(context.Background(), duration.Seconds(), metric.WithAttributes(attrs...))
}

// Close closes the metrics service
func (m *MetricsServiceOTEL) Close() error {
	if m.meterProvider != nil {
		return m.meterProvider.Shutdown(context.Background())
	}
	return nil
}

// convertTagsToAttributes converts a map of tags to OTEL attributes
func convertTagsToAttributes(tags map[string]string) []attribute.KeyValue {
	if len(tags) == 0 {
		return nil
	}

	attrs := make([]attribute.KeyValue, 0, len(tags))
	for key, value := range tags {
		attrs = append(attrs, attribute.String(key, value))
	}
	return attrs
}
