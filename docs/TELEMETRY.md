# Telemetry & Datadog Integration

This document describes the telemetry implementation and Datadog integration in the GoHexaClean boilerplate.

## Overview

The application implements observability through two main pillars:
- **Metrics**: Application and business metrics using Datadog StatsD
- **Distributed Tracing**: Request tracing and APM using Datadog APM

## Architecture

### Port/Adapter Pattern

Telemetry follows the hexagonal architecture pattern:

```
internal/port/outbound/telemetry/     # Port interfaces
  ├── metrics.go                      # MetricsService interface
  └── tracing.go                      # TracingService interface

internal/adapter/outbound/datadog/    # Datadog adapter implementation
  ├── metrics_service.go              # Datadog StatsD implementation
  └── tracing_service.go              # Datadog APM implementation
```

### Key Components

1. **Telemetry Ports** (`internal/port/outbound/telemetry/`)
   - `MetricsService`: Interface for metrics collection
   - `TracingService`: Interface for distributed tracing
   - `Span`: Interface representing a trace span

2. **Datadog Adapter** (`internal/adapter/outbound/datadog/`)
   - `MetricsServiceDatadog`: Implements metrics using Datadog StatsD
   - `TracingServiceDatadog`: Implements tracing using Datadog APM
   - `DatadogSpan`: Wraps Datadog span implementation

3. **HTTP Middleware** (`internal/adapter/inbound/http/middleware/telemetry.go`)
   - Automatic metrics collection for all HTTP requests
   - Automatic trace creation for each request
   - Response time tracking
   - Status code metrics

## Configuration

### YAML Configuration (`config/app.yaml`)

```yaml
datadog:
  enabled: false              # Enable/disable Datadog integration
  agent_host: localhost       # Datadog agent host
  agent_port: "8125"         # Datadog agent port (StatsD)
  namespace: gohexaclean     # Metrics namespace prefix
  tags:                      # Global tags for all metrics
    - env:development
    - service:gohexaclean
  apm_enabled: false         # Enable/disable APM tracing
```

### Environment Variables

Override configuration using environment variables:

```bash
DD_ENABLED=true              # Enable Datadog
DD_AGENT_HOST=localhost      # Agent host
DD_AGENT_PORT=8125          # Agent port
DD_APM_ENABLED=true         # Enable APM tracing
```

## Metrics

### Automatically Collected Metrics

The telemetry middleware automatically collects the following metrics for HTTP requests:

1. **Request Count**
   - Metric: `http.requests.total`
   - Tags: `method`, `route`, `status`
   - Type: Counter

2. **Request Duration**
   - Metric: `http.request.duration`
   - Tags: `method`, `route`, `status`
   - Type: Timing

3. **Success Requests**
   - Metric: `http.requests.success`
   - Tags: `method`, `route`, `status`
   - Type: Counter
   - Condition: HTTP 2xx responses

4. **Client Errors**
   - Metric: `http.requests.client_errors`
   - Tags: `method`, `route`, `status`
   - Type: Counter
   - Condition: HTTP 4xx responses

5. **Server Errors**
   - Metric: `http.requests.errors`
   - Tags: `method`, `route`, `status`
   - Type: Counter
   - Condition: HTTP 5xx responses

### Custom Metrics

You can record custom metrics in your application code:

```go
// Inject MetricsService from container
func (s *YourService) YourMethod(metricsService telemetry.MetricsService) {
    // Counter
    metricsService.IncrementCounter("custom.counter", map[string]string{
        "type": "business_event",
    }, 1)

    // Gauge
    metricsService.SetGauge("custom.gauge", map[string]string{
        "resource": "database",
    }, 42.5)

    // Timing
    start := time.Now()
    // ... your operation ...
    metricsService.RecordTiming("custom.operation.duration", nil, time.Since(start))

    // Histogram
    metricsService.RecordHistogram("custom.histogram", map[string]string{
        "bucket": "large",
    }, 100.0)

    // Distribution
    metricsService.RecordDistribution("custom.distribution", nil, 75.3)
}
```

## Distributed Tracing (APM)

### Automatic Tracing

The telemetry middleware automatically creates traces for all HTTP requests with:
- Operation name: `{METHOD} {ROUTE}` (e.g., "GET /api/v1/users")
- HTTP method tag
- HTTP URL tag
- HTTP route tag
- HTTP status code tag
- Error tag (if status >= 400)

### Custom Spans

Create custom spans in your code:

```go
func (s *YourService) YourMethod(ctx context.Context, tracing telemetry.TracingService) error {
    // Start a child span
    span, ctx := tracing.StartChildSpan(ctx, "database.query")
    defer span.Finish()

    // Add tags to the span
    span.SetTag("query.type", "select")
    span.SetTag("table", "users")

    // Perform operation
    result, err := s.performDatabaseQuery(ctx)
    if err != nil {
        // Mark span as error
        span.SetError(err)
        return err
    }

    span.SetTag("rows.count", len(result))
    return nil
}
```

### Trace Context Propagation

The tracing context is automatically propagated through:
- HTTP requests (via middleware)
- Service layers (via context.Context)
- Database queries (if using instrumented drivers)

## Setup & Installation

### 1. Install Datadog Agent

#### Docker (Development)

```bash
docker run -d --name datadog-agent \
  -e DD_API_KEY=<YOUR_API_KEY> \
  -e DD_SITE=datadoghq.com \
  -e DD_APM_ENABLED=true \
  -e DD_APM_NON_LOCAL_TRAFFIC=true \
  -e DD_DOGSTATSD_NON_LOCAL_TRAFFIC=true \
  -p 8125:8125/udp \
  -p 8126:8126/tcp \
  -v /var/run/docker.sock:/var/run/docker.sock:ro \
  -v /proc/:/host/proc/:ro \
  -v /sys/fs/cgroup/:/host/sys/fs/cgroup:ro \
  datadog/agent:latest
```

#### Kubernetes (Production)

Use the Datadog Operator or Helm chart:

```bash
helm repo add datadog https://helm.datadoghq.com
helm install datadog-agent datadog/datadog \
  --set datadog.apiKey=<YOUR_API_KEY> \
  --set datadog.apm.enabled=true \
  --set datadog.logs.enabled=true
```

### 2. Configure Application

Update `config/app.yaml`:

```yaml
datadog:
  enabled: true
  agent_host: localhost  # or datadog-agent service name in k8s
  agent_port: "8125"
  namespace: gohexaclean
  tags:
    - env:production
    - service:gohexaclean
    - version:1.0.0
  apm_enabled: true
```

### 3. Run Application

```bash
# Development
DD_ENABLED=true DD_APM_ENABLED=true go run cmd/http/main.go

# Production (Docker)
docker run -e DD_ENABLED=true -e DD_APM_ENABLED=true your-app
```

## Viewing Metrics & Traces

### Datadog Dashboard

1. **Metrics Explorer**: https://app.datadoghq.com/metric/explorer
   - Search for metrics with prefix: `gohexaclean.*`
   - Example: `gohexaclean.http.requests.total`

2. **APM Traces**: https://app.datadoghq.com/apm/traces
   - View request traces
   - Analyze performance bottlenecks
   - Identify errors

3. **Service Map**: https://app.datadoghq.com/apm/map
   - Visualize service dependencies
   - Monitor service health

### Example Queries

**Average request duration by endpoint:**
```
avg:gohexaclean.http.request.duration{*} by {route}
```

**Error rate:**
```
sum:gohexaclean.http.requests.errors{*}.as_rate()
```

**Request throughput:**
```
sum:gohexaclean.http.requests.total{*}.as_rate()
```

## Best Practices

1. **Use Meaningful Tags**
   - Add relevant context to metrics
   - Keep cardinality in check (< 1000 unique values per tag)

2. **Span Naming**
   - Use clear, hierarchical names (e.g., `service.operation`)
   - Be consistent across the application

3. **Error Tracking**
   - Always call `span.SetError()` for errors
   - Include error details in tags

4. **Resource Cleanup**
   - Always defer `span.Finish()`
   - Ensure metrics service is closed on shutdown

5. **Performance**
   - Metrics operations are non-blocking
   - Failed metric submissions don't affect application

## Troubleshooting

### Metrics Not Appearing

1. Check Datadog agent is running:
   ```bash
   docker logs datadog-agent
   ```

2. Verify agent connectivity:
   ```bash
   telnet localhost 8125
   ```

3. Check application logs for telemetry initialization

### Traces Not Appearing

1. Verify APM is enabled in config
2. Check agent APM port (8126) is accessible
3. Review trace sampling configuration
4. Check application environment variables

### High Cardinality Warning

If you see warnings about high cardinality:
- Reduce the number of unique tag values
- Use tag aggregation
- Consider using distribution metrics instead of histograms

## Dependencies

- `github.com/DataDog/datadog-go/v5/statsd`: Datadog StatsD client
- `gopkg.in/DataDog/dd-trace-go.v1`: Datadog APM tracer

## Further Reading

- [Datadog APM Documentation](https://docs.datadoghq.com/tracing/)
- [Datadog Metrics Documentation](https://docs.datadoghq.com/metrics/)
- [Go Tracer Documentation](https://docs.datadoghq.com/tracing/setup_overview/setup/go/)
