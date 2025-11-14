package datadog

import (
	"context"

	"github.com/gieart87/gohexaclean/internal/port/outbound/telemetry"
	ddtrace "gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// TracingServiceDatadog implements telemetry.TracingService using Datadog APM
type TracingServiceDatadog struct {
	serviceName string
}

// DatadogSpan wraps Datadog's span to implement telemetry.Span interface
type DatadogSpan struct {
	span ddtrace.Span
}

// NewTracingServiceDatadog creates a new Datadog tracing service
func NewTracingServiceDatadog(serviceName, agentHost string, agentPort string, env string) telemetry.TracingService {
	tracer.Start(
		tracer.WithService(serviceName),
		tracer.WithEnv(env),
		tracer.WithAgentAddr(agentHost + ":" + agentPort),
		tracer.WithAnalytics(true),
	)

	return &TracingServiceDatadog{
		serviceName: serviceName,
	}
}

// StartSpan starts a new root span
func (t *TracingServiceDatadog) StartSpan(ctx context.Context, operationName string, opts ...interface{}) (telemetry.Span, context.Context) {
	span, ctx := tracer.StartSpanFromContext(ctx, operationName)
	return &DatadogSpan{span: span}, ctx
}

// StartChildSpan starts a child span from a parent context
func (t *TracingServiceDatadog) StartChildSpan(ctx context.Context, operationName string) (telemetry.Span, context.Context) {
	span, ctx := tracer.StartSpanFromContext(ctx, operationName)
	return &DatadogSpan{span: span}, ctx
}

// Close stops the tracer
func (t *TracingServiceDatadog) Close() error {
	tracer.Stop()
	return nil
}

// SetTag sets a tag on the span
func (s *DatadogSpan) SetTag(key string, value interface{}) {
	s.span.SetTag(key, value)
}

// SetError marks the span as having an error
func (s *DatadogSpan) SetError(err error) {
	s.span.SetTag(ext.Error, err)
}

// Finish completes the span
func (s *DatadogSpan) Finish() {
	s.span.Finish()
}
