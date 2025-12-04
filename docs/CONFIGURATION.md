# Configuration Guide

This guide explains how to configure the GoHexaClean application for different environments.

## Configuration Methods

GoHexaClean supports multiple configuration methods with the following priority order:

1. **Environment Variables** (highest priority)
2. **`.env` file**
3. **YAML configuration files** (`config/app.yaml`, `config/app.dev.yaml`, `config/app.prod.yaml`)
4. **Default values** (lowest priority)

## Environment Variables (.env)

The application automatically loads environment variables from the `.env` file in the project root directory.

### Creating Your .env File

1. Copy the example file (if available):
   ```bash
   cp .env.example .env
   ```

2. Or create a new `.env` file with the following variables:

```bash
# Application
APP_NAME=gohexaclean
APP_ENV=development
APP_DEBUG=true

# Server Ports
HTTP_PORT=8080
GRPC_PORT=50051

# Database PostgreSQL
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=gohexaclean
DB_SSL_MODE=disable
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=5
DB_MAX_LIFETIME=5m

# Redis Cache
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
REDIS_POOL_SIZE=10

# Logger
LOG_LEVEL=debug
LOG_FORMAT=json
LOG_OUTPUT=stdout

# JWT Authentication
JWT_SECRET=your-secret-key-change-this-in-production
JWT_EXPIRED=24h

# CORS
CORS_ALLOW_ORIGINS=*
CORS_ALLOW_METHODS=GET,POST,PUT,DELETE,PATCH
CORS_ALLOW_HEADERS=Origin,Content-Type,Accept,Authorization

# Rate Limiting
RATE_LIMIT_ENABLED=true
RATE_LIMIT_MAX=100
RATE_LIMIT_WINDOW=1m

# Telemetry
OTEL_ENABLED=true
OTEL_SERVICE_NAME=gohexaclean
OTEL_COLLECTOR_ENDPOINT=localhost:4317

# Metrics
METRICS_ENABLED=true
METRICS_PORT=9090

# Message Broker (RabbitMQ)
BROKER_ENABLED=false
BROKER_TYPE=rabbitmq
RABBITMQ_URL=amqp://guest:guest@localhost:5672/

# Background Jobs (Asynq)
REDIS_ADDR=localhost:6379
```

## Configuration Variables Reference

### Application Settings

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `APP_NAME` | Application name | `gohexaclean` | Yes |
| `APP_ENV` | Environment (development/staging/production) | `development` | Yes |
| `APP_DEBUG` | Enable debug mode | `true` | No |

### Server Settings

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `HTTP_PORT` | HTTP server port | `8080` | Yes |
| `GRPC_PORT` | gRPC server port | `50051` | Yes |

### Database Settings

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `DB_HOST` | PostgreSQL host | `localhost` | Yes |
| `DB_PORT` | PostgreSQL port | `5432` | Yes |
| `DB_USER` | Database user | `postgres` | Yes |
| `DB_PASSWORD` | Database password | `postgres` | Yes |
| `DB_NAME` | Database name | `gohexaclean` | Yes |
| `DB_SSL_MODE` | SSL mode (disable/require/verify-ca/verify-full) | `disable` | No |
| `DB_MAX_OPEN_CONNS` | Maximum open connections | `25` | No |
| `DB_MAX_IDLE_CONNS` | Maximum idle connections | `5` | No |
| `DB_MAX_LIFETIME` | Connection max lifetime | `5m` | No |

### Redis Settings

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `REDIS_HOST` | Redis host | `localhost` | No |
| `REDIS_PORT` | Redis port | `6379` | No |
| `REDIS_PASSWORD` | Redis password | (empty) | No |
| `REDIS_DB` | Redis database number | `0` | No |
| `REDIS_POOL_SIZE` | Connection pool size | `10` | No |
| `REDIS_ADDR` | Redis address for Asynq | `localhost:6379` | No |

### JWT Settings

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `JWT_SECRET` | Secret key for JWT signing | - | Yes |
| `JWT_EXPIRED` | Token expiration duration | `24h` | Yes |

**⚠️ IMPORTANT:** Always use a strong, unique `JWT_SECRET` in production!

### Logger Settings

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `LOG_LEVEL` | Log level (debug/info/warn/error) | `debug` | No |
| `LOG_FORMAT` | Log format (json/text) | `json` | No |
| `LOG_OUTPUT` | Log output (stdout/file) | `stdout` | No |

### CORS Settings

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `CORS_ALLOW_ORIGINS` | Allowed origins (* or comma-separated URLs) | `*` | No |
| `CORS_ALLOW_METHODS` | Allowed HTTP methods | `GET,POST,PUT,DELETE,PATCH` | No |
| `CORS_ALLOW_HEADERS` | Allowed headers | `Origin,Content-Type,Accept,Authorization` | No |

### Rate Limiting

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `RATE_LIMIT_ENABLED` | Enable rate limiting | `true` | No |
| `RATE_LIMIT_MAX` | Maximum requests | `100` | No |
| `RATE_LIMIT_WINDOW` | Time window | `1m` | No |

### Telemetry (OpenTelemetry)

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `OTEL_ENABLED` | Enable OpenTelemetry | `true` | No |
| `OTEL_SERVICE_NAME` | Service name for tracing | `gohexaclean` | No |
| `OTEL_COLLECTOR_ENDPOINT` | OTEL collector endpoint | `localhost:4317` | No |

### Metrics

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `METRICS_ENABLED` | Enable metrics endpoint | `true` | No |
| `METRICS_PORT` | Metrics server port | `9090` | No |

### Message Broker (RabbitMQ)

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `BROKER_ENABLED` | Enable message broker | `false` | No |
| `BROKER_TYPE` | Broker type (rabbitmq/kafka/nats) | `rabbitmq` | No |
| `RABBITMQ_URL` | RabbitMQ connection URL | `amqp://guest:guest@localhost:5672/` | No |

## Environment-Specific Configuration

### Development Environment

For local development, use the default `.env` file:

```bash
APP_ENV=development
APP_DEBUG=true
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=postgres
LOG_LEVEL=debug
```

### Docker Environment

When running with Docker, the `.env` file is automatically loaded by:
- `docker-compose.yml` - for container configuration
- `Makefile` - for database migrations and seeds

The Dockerfile and docker-compose automatically use these variables.

### Production Environment

For production, create a `.env.production` file or use environment variables directly:

```bash
APP_ENV=production
APP_DEBUG=false
DB_HOST=production-db-host
DB_USER=prod_user
DB_PASSWORD=strong-password-here
JWT_SECRET=very-strong-secret-key-here
LOG_LEVEL=info
CORS_ALLOW_ORIGINS=https://yourdomain.com
RATE_LIMIT_MAX=50
```

**Security Best Practices:**
1. Never commit `.env` files to version control
2. Add `.env` to `.gitignore`
3. Use strong passwords and secrets
4. Rotate credentials regularly
5. Use different credentials for each environment

## YAML Configuration Files

In addition to `.env`, you can use YAML configuration files located in the `config/` directory:

- `config/app.yaml` - Base configuration
- `config/app.dev.yaml` - Development overrides
- `config/app.prod.yaml` - Production overrides

### Example config/app.yaml:

```yaml
app:
  name: ${APP_NAME}
  env: ${APP_ENV}
  debug: ${APP_DEBUG}

http:
  port: ${HTTP_PORT}

grpc:
  port: ${GRPC_PORT}

database:
  host: ${DB_HOST}
  port: ${DB_PORT}
  user: ${DB_USER}
  password: ${DB_PASSWORD}
  name: ${DB_NAME}
  ssl_mode: ${DB_SSL_MODE}
  max_open_conns: ${DB_MAX_OPEN_CONNS}
  max_idle_conns: ${DB_MAX_IDLE_CONNS}
  max_lifetime: ${DB_MAX_LIFETIME}

redis:
  host: ${REDIS_HOST}
  port: ${REDIS_PORT}
  password: ${REDIS_PASSWORD}
  db: ${REDIS_DB}
  pool_size: ${REDIS_POOL_SIZE}

jwt:
  secret: ${JWT_SECRET}
  expired: ${JWT_EXPIRED}

logger:
  level: ${LOG_LEVEL}
  format: ${LOG_FORMAT}
  output: ${LOG_OUTPUT}
```

## Makefile Integration

The `Makefile` automatically loads the `.env` file and uses these variables:

### Database Migrations

```bash
# All commands use DB_* variables from .env
make migrate-up      # Run migrations
make migrate-down    # Rollback migrations
make migrate-status  # Check migration status
make migrate-reset   # Reset database
```

### Database Seeding

```bash
# Uses DB_* variables from .env
make seed
```

### Docker Operations

```bash
# docker-compose automatically loads .env
make docker-up       # Start containers with .env config
make docker-down     # Stop containers
make docker-logs     # View logs
```

## Verifying Configuration

### Check Loaded Environment Variables

```bash
# In Makefile
make version

# Or manually
echo $DB_HOST
echo $DB_USER
```

### Test Database Connection

```bash
# Using psql with .env values
PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "SELECT version();"
```

### Check Docker Compose Config

```bash
# See what docker-compose will use
docker-compose config
```

## Troubleshooting

### Variables Not Loading

**Problem:** Environment variables from `.env` not being used

**Solution:**
1. Ensure `.env` file exists in project root
2. Check file permissions: `ls -la .env`
3. Verify no syntax errors in `.env`
4. Restart Docker containers: `make docker-down && make docker-up`

### Database Connection Failed

**Problem:** Cannot connect to PostgreSQL

**Solution:**
1. Verify credentials in `.env`:
   ```bash
   cat .env | grep DB_
   ```
2. Check PostgreSQL is running:
   ```bash
   docker ps | grep postgres
   ```
3. Test connection manually:
   ```bash
   PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -U $DB_USER -d postgres
   ```

### Migration Failed

**Problem:** `make migrate-up` fails

**Solution:**
1. Check database credentials in `.env`
2. Ensure PostgreSQL container is healthy
3. Verify database exists:
   ```bash
   docker exec gohexaclean-postgres psql -U postgres -c "\l"
   ```

## Best Practices

### 1. Use Different .env Files

```bash
.env              # Development (committed as .env.example)
.env.local        # Local overrides (gitignored)
.env.production   # Production (gitignored)
.env.staging      # Staging (gitignored)
```

### 2. Template File

Maintain a `.env.example` with all required variables (without sensitive values):

```bash
# .env.example
APP_NAME=gohexaclean
APP_ENV=development
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_username
DB_PASSWORD=your_password
DB_NAME=gohexaclean
JWT_SECRET=change_this_in_production
```

### 3. Validation

Add validation for required variables in your application startup.

### 4. Documentation

Keep this configuration guide updated when adding new variables.

## Migration from Hardcoded Values

If migrating from hardcoded values:

1. Create `.env` file with current values
2. Update code to read from environment
3. Test locally
4. Update deployment scripts
5. Remove hardcoded values

## Additional Resources

- [12-Factor App Configuration](https://12factor.net/config)
- [Docker Compose Environment Variables](https://docs.docker.com/compose/environment-variables/)
- [PostgreSQL Environment Variables](https://www.postgresql.org/docs/current/libpq-envars.html)

---

**Note:** Always keep sensitive configuration like passwords and API keys out of version control. Use `.env` files for local development and proper secret management tools for production.
