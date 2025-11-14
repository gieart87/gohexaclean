package datadog

import (
	"time"

	"github.com/DataDog/datadog-go/v5/statsd"
	"github.com/gieart87/gohexaclean/internal/port/outbound/telemetry"
)

// MetricsServiceDatadog implements telemetry.MetricsService using Datadog StatsD
type MetricsServiceDatadog struct {
	client *statsd.Client
}

// NewMetricsServiceDatadog creates a new Datadog metrics service
func NewMetricsServiceDatadog(address string, namespace string, tags []string) (telemetry.MetricsService, error) {
	client, err := statsd.New(address,
		statsd.WithNamespace(namespace),
		statsd.WithTags(tags),
	)
	if err != nil {
		return nil, err
	}

	return &MetricsServiceDatadog{
		client: client,
	}, nil
}

// IncrementCounter increments a counter metric
func (m *MetricsServiceDatadog) IncrementCounter(name string, tags map[string]string, value float64) {
	tagsList := convertTagsToList(tags)
	_ = m.client.Count(name, int64(value), tagsList, 1)
}

// SetGauge sets a gauge metric
func (m *MetricsServiceDatadog) SetGauge(name string, tags map[string]string, value float64) {
	tagsList := convertTagsToList(tags)
	_ = m.client.Gauge(name, value, tagsList, 1)
}

// RecordHistogram records a histogram metric
func (m *MetricsServiceDatadog) RecordHistogram(name string, tags map[string]string, value float64) {
	tagsList := convertTagsToList(tags)
	_ = m.client.Histogram(name, value, tagsList, 1)
}

// RecordDistribution records a distribution metric
func (m *MetricsServiceDatadog) RecordDistribution(name string, tags map[string]string, value float64) {
	tagsList := convertTagsToList(tags)
	_ = m.client.Distribution(name, value, tagsList, 1)
}

// RecordTiming records a timing metric
func (m *MetricsServiceDatadog) RecordTiming(name string, tags map[string]string, duration time.Duration) {
	tagsList := convertTagsToList(tags)
	_ = m.client.Timing(name, duration, tagsList, 1)
}

// Close closes the metrics service
func (m *MetricsServiceDatadog) Close() error {
	return m.client.Close()
}

// convertTagsToList converts a map of tags to a list of "key:value" strings
func convertTagsToList(tags map[string]string) []string {
	if len(tags) == 0 {
		return nil
	}

	tagsList := make([]string, 0, len(tags))
	for key, value := range tags {
		tagsList = append(tagsList, key+":"+value)
	}
	return tagsList
}
