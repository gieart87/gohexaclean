package telemetry

import "context"

// Span represents a tracing span
type Span interface {
	// SetTag sets a tag on the span
	SetTag(key string, value interface{})

	// SetError marks the span as having an error
	SetError(err error)

	// Finish completes the span
	Finish()
}

// TracingService defines the interface for distributed tracing
type TracingService interface {
	// StartSpan starts a new span
	StartSpan(ctx context.Context, operationName string, opts ...interface{}) (Span, context.Context)

	// StartChildSpan starts a child span from a parent context
	StartChildSpan(ctx context.Context, operationName string) (Span, context.Context)

	// Close closes the tracing service
	Close() error
}
