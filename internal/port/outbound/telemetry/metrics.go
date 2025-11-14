package telemetry

import "time"

// MetricsService defines the interface for metrics collection
type MetricsService interface {
	// Counter metrics
	IncrementCounter(name string, tags map[string]string, value float64)

	// Gauge metrics
	SetGauge(name string, tags map[string]string, value float64)

	// Histogram metrics
	RecordHistogram(name string, tags map[string]string, value float64)

	// Distribution metrics
	RecordDistribution(name string, tags map[string]string, value float64)

	// Timing metrics
	RecordTiming(name string, tags map[string]string, duration time.Duration)

	// Close closes the metrics service
	Close() error
}
