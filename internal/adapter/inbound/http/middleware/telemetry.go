package middleware

import (
	"context"
	"strconv"
	"time"

	"github.com/gieart87/gohexaclean/internal/port/outbound/telemetry"
	"github.com/gofiber/fiber/v2"
)

// TelemetryMiddleware creates middleware for collecting HTTP metrics and traces
func TelemetryMiddleware(metrics telemetry.MetricsService, tracing telemetry.TracingService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Start tracing span if tracing is enabled
		var span telemetry.Span
		if tracing != nil {
			operationName := c.Method() + " " + c.Route().Path
			var ctx context.Context
			span, ctx = tracing.StartSpan(c.UserContext(), operationName)
			c.SetUserContext(ctx)
			defer span.Finish()

			// Set span tags
			span.SetTag("http.method", c.Method())
			span.SetTag("http.url", c.OriginalURL())
			span.SetTag("http.route", c.Route().Path)
		}

		// Process request
		err := c.Next()

		// Record metrics if metrics service is enabled
		if metrics != nil {
			duration := time.Since(start)
			statusCode := c.Response().StatusCode()

			tags := map[string]string{
				"method": c.Method(),
				"route":  c.Route().Path,
				"status": strconv.Itoa(statusCode),
			}

			// Record request count
			metrics.IncrementCounter("http.requests.total", tags, 1)

			// Record request duration
			metrics.RecordTiming("http.request.duration", tags, duration)

			// Record status code counts
			if statusCode >= 500 {
				metrics.IncrementCounter("http.requests.errors", tags, 1)
			} else if statusCode >= 400 {
				metrics.IncrementCounter("http.requests.client_errors", tags, 1)
			} else if statusCode >= 200 && statusCode < 300 {
				metrics.IncrementCounter("http.requests.success", tags, 1)
			}
		}

		// Update span with status code if tracing is enabled
		if span != nil {
			statusCode := c.Response().StatusCode()
			span.SetTag("http.status_code", statusCode)

			if statusCode >= 400 {
				span.SetTag("error", true)
			}

			if err != nil {
				span.SetError(err)
			}
		}

		return err
	}
}
