# Telemetry Quick Start Guide

Quick guide to get started with telemetry in GoHexaClean.

## Choose Your Backend

### Option 1: Datadog (Recommended for Datadog users)

**1. Update config/app.yaml:**
```yaml
datadog:
  enabled: true
  agent_host: localhost
  agent_port: "8125"
  namespace: gohexaclean
  tags:
    - env:production
    - service:gohexaclean
  apm_enabled: true
```

**2. Start Datadog Agent:**
```bash
docker run -d --name datadog-agent \
  -e DD_API_KEY=<YOUR_API_KEY> \
  -e DD_SITE=datadoghq.com \
  -e DD_APM_ENABLED=true \
  -p 8125:8125/udp \
  -p 8126:8126/tcp \
  datadog/agent:latest
```

**3. Run your application:**
```bash
go run cmd/http/main.go
```

**4. View metrics:**
- Dashboard: https://app.datadoghq.com
- Metrics will appear as: `gohexaclean.*`

---

### Option 2: OpenTelemetry (Vendor-agnostic)

**1. Update config/app.yaml:**
```yaml
telemetry:
  enabled: true
  service_name: gohexaclean
  collector_endpoint: localhost:4317
```

**2. Start OTEL Collector:**
```bash
# Using Jaeger all-in-one
docker run -d --name jaeger \
  -p 4317:4317 \
  -p 16686:16686 \
  jaegertracing/all-in-one:latest
```

**3. Run your application:**
```bash
go run cmd/http/main.go
```

**4. View traces:**
- Jaeger UI: http://localhost:16686

---

## Configuration via Environment Variables

### Datadog
```bash
export DD_ENABLED=true
export DD_AGENT_HOST=localhost
export DD_AGENT_PORT=8125
export DD_APM_ENABLED=true
go run cmd/http/main.go
```

### OpenTelemetry
```bash
# Datadog must be disabled for OTEL to activate
export DD_ENABLED=false
# No specific OTEL env vars needed, uses YAML config
go run cmd/http/main.go
```

---

## What Gets Collected?

### Automatic HTTP Metrics
- Request count by endpoint
- Request duration
- Status codes (2xx, 4xx, 5xx)
- Error rates

### Automatic HTTP Traces
- Full request trace with timing
- Method, URL, route, status code
- Error information

---

## Development vs Production

**Development:**
```yaml
datadog:
  enabled: false  # or use free Jaeger

telemetry:
  enabled: true
  collector_endpoint: localhost:4317
```

**Production:**
```yaml
datadog:
  enabled: true
  agent_host: datadog-agent  # k8s service name
  namespace: gohexaclean
  tags:
    - env:production
    - version:1.0.0
  apm_enabled: true
```

---

## Troubleshooting

**Metrics/traces not appearing:**
1. Check logs for "Telemetry middleware enabled"
2. Verify agent/collector is running
3. Check network connectivity
4. Confirm configuration is loaded

**Which backend is active:**
Check application logs on startup:
- "Datadog metrics initialized" → Using Datadog
- "OpenTelemetry metrics initialized" → Using OTEL

---

## Next Steps

- Read full documentation: [TELEMETRY.md](./TELEMETRY.md)
- Add custom metrics to your services
- Create dashboards in your backend
- Set up alerts
