package otel

import (
	"context"

	"github.com/gieart87/gohexaclean/internal/port/outbound/telemetry"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

// TracingServiceOTEL implements telemetry.TracingService using OpenTelemetry
type TracingServiceOTEL struct {
	tracerProvider *sdktrace.TracerProvider
	tracer         trace.Tracer
}

// OTELSpan wraps OTEL's span to implement telemetry.Span interface
type OTELSpan struct {
	span trace.Span
}

// NewTracingServiceOTEL creates a new OpenTelemetry tracing service
func NewTracingServiceOTEL(ctx context.Context, serviceName, collectorEndpoint string) (telemetry.TracingService, error) {
	// Create OTLP exporter
	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(collectorEndpoint),
		otlptracegrpc.WithInsecure(),
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

	// Create tracer provider
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	// Set global tracer provider
	otel.SetTracerProvider(tracerProvider)

	// Get tracer
	tracer := tracerProvider.Tracer(serviceName)

	return &TracingServiceOTEL{
		tracerProvider: tracerProvider,
		tracer:         tracer,
	}, nil
}

// StartSpan starts a new root span
func (t *TracingServiceOTEL) StartSpan(ctx context.Context, operationName string, opts ...interface{}) (telemetry.Span, context.Context) {
	ctx, span := t.tracer.Start(ctx, operationName)
	return &OTELSpan{span: span}, ctx
}

// StartChildSpan starts a child span from a parent context
func (t *TracingServiceOTEL) StartChildSpan(ctx context.Context, operationName string) (telemetry.Span, context.Context) {
	ctx, span := t.tracer.Start(ctx, operationName)
	return &OTELSpan{span: span}, ctx
}

// Close stops the tracer
func (t *TracingServiceOTEL) Close() error {
	if t.tracerProvider != nil {
		return t.tracerProvider.Shutdown(context.Background())
	}
	return nil
}

// SetTag sets a tag on the span
func (s *OTELSpan) SetTag(key string, value interface{}) {
	switch v := value.(type) {
	case string:
		s.span.SetAttributes(attribute.String(key, v))
	case int:
		s.span.SetAttributes(attribute.Int(key, v))
	case int64:
		s.span.SetAttributes(attribute.Int64(key, v))
	case float64:
		s.span.SetAttributes(attribute.Float64(key, v))
	case bool:
		s.span.SetAttributes(attribute.Bool(key, v))
	default:
		s.span.SetAttributes(attribute.String(key, "unknown"))
	}
}

// SetError marks the span as having an error
func (s *OTELSpan) SetError(err error) {
	s.span.RecordError(err)
	s.span.SetStatus(codes.Error, err.Error())
}

// Finish completes the span
func (s *OTELSpan) Finish() {
	s.span.End()
}
