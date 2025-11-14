package bootstrap

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gieart87/gohexaclean/internal/adapter/inbound/grpc/handler"
	"github.com/gieart87/gohexaclean/internal/adapter/outbound/datadog"
	"github.com/gieart87/gohexaclean/internal/adapter/outbound/otel"
	"github.com/gieart87/gohexaclean/internal/adapter/outbound/pgsql"
	"github.com/gieart87/gohexaclean/internal/adapter/outbound/redis"
	"github.com/gieart87/gohexaclean/internal/app"
	"github.com/gieart87/gohexaclean/internal/infra/cache"
	"github.com/gieart87/gohexaclean/internal/infra/config"
	"github.com/gieart87/gohexaclean/internal/infra/db"
	"github.com/gieart87/gohexaclean/internal/infra/logger"
	"github.com/gieart87/gohexaclean/internal/port/inbound"
	"github.com/gieart87/gohexaclean/internal/port/outbound/repository"
	"github.com/gieart87/gohexaclean/internal/port/outbound/service"
	"github.com/gieart87/gohexaclean/internal/port/outbound/telemetry"
	redisClient "github.com/redis/go-redis/v9"
)

// Container holds all application dependencies
type Container struct {
	Config *config.Config
	Logger *logger.Logger

	// Database
	DB          *sql.DB
	RedisClient *redisClient.Client

	// Repositories
	UserRepository repository.UserRepository

	// Services
	CacheService service.CacheService

	// Telemetry
	MetricsService telemetry.MetricsService
	TracingService telemetry.TracingService

	// Use Cases / Application Services
	UserService inbound.UserServicePort

	// gRPC Handlers
	UserGRPCHandler *handler.UserHandlerGRPC
}

// NewContainer creates and initializes a new dependency injection container
func NewContainer(configPath string) (*Container, error) {
	container := &Container{}

	// Load configuration
	cfg, err := config.Load(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	container.Config = cfg

	// Initialize logger
	log, err := logger.NewLogger(&cfg.Logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}
	container.Logger = log

	// Initialize database
	database, err := db.NewPostgresConnection(&cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	container.DB = database
	log.Info("Database connection established")

	// Initialize Redis
	redisConn, err := cache.NewRedisClient(&cfg.Redis)
	if err != nil {
		log.Warn("Failed to connect to Redis, cache will be disabled")
		// Continue without Redis - cache is optional
	} else {
		container.RedisClient = redisConn
		log.Info("Redis connection established")
	}

	// Initialize repositories
	container.UserRepository = pgsql.NewUserRepositoryPG(database)

	// Initialize telemetry services
	ctx := context.Background()

	// Priority: Datadog > OpenTelemetry
	if cfg.Datadog.Enabled {
		// Initialize Datadog metrics
		metricsService, err := datadog.NewMetricsServiceDatadog(
			cfg.Datadog.AgentHost+":"+cfg.Datadog.AgentPort,
			cfg.Datadog.Namespace,
			cfg.Datadog.Tags,
		)
		if err != nil {
			log.Warn("Failed to initialize Datadog metrics, continuing without metrics")
		} else {
			container.MetricsService = metricsService
			log.Info("Datadog metrics initialized")
		}

		// Initialize Datadog APM tracing
		if cfg.Datadog.APMEnabled {
			container.TracingService = datadog.NewTracingServiceDatadog(
				cfg.App.Name,
				cfg.Datadog.AgentHost,
				cfg.Datadog.AgentPort,
				cfg.App.Env,
			)
			log.Info("Datadog APM tracing initialized")
		}
	} else if cfg.Telemetry.Enabled {
		// Initialize OpenTelemetry as fallback
		log.Info("Initializing OpenTelemetry telemetry")

		// Initialize OTEL metrics
		metricsService, err := otel.NewMetricsServiceOTEL(
			ctx,
			cfg.Telemetry.ServiceName,
			cfg.Telemetry.CollectorEndpoint,
		)
		if err != nil {
			log.Warn("Failed to initialize OpenTelemetry metrics, continuing without metrics")
		} else {
			container.MetricsService = metricsService
			log.Info("OpenTelemetry metrics initialized")
		}

		// Initialize OTEL tracing
		tracingService, err := otel.NewTracingServiceOTEL(
			ctx,
			cfg.Telemetry.ServiceName,
			cfg.Telemetry.CollectorEndpoint,
		)
		if err != nil {
			log.Warn("Failed to initialize OpenTelemetry tracing, continuing without tracing")
		} else {
			container.TracingService = tracingService
			log.Info("OpenTelemetry tracing initialized")
		}
	}

	// Initialize services
	if container.RedisClient != nil {
		container.CacheService = redis.NewCacheServiceRedis(container.RedisClient)
	} else {
		// Use a no-op cache service if Redis is not available
		container.CacheService = &NoOpCacheService{}
	}

	// Initialize use cases / application services
	container.UserService = app.NewUserService(
		container.UserRepository,
		container.CacheService,
		&cfg.JWT,
	)

	// Initialize gRPC handlers
	container.UserGRPCHandler = handler.NewUserHandlerGRPC(container.UserService)

	log.Info("Container initialized successfully")

	return container, nil
}

// Close closes all resources in the container
func (c *Container) Close() error {
	if c.Logger != nil {
		c.Logger.Info("Shutting down application...")
	}

	if c.DB != nil {
		if err := db.Close(c.DB); err != nil {
			c.Logger.Error("Failed to close database connection")
		}
	}

	if c.RedisClient != nil {
		if err := cache.Close(c.RedisClient); err != nil {
			c.Logger.Error("Failed to close Redis connection")
		}
	}

	// Close telemetry services
	if c.MetricsService != nil {
		if err := c.MetricsService.Close(); err != nil {
			c.Logger.Error("Failed to close metrics service")
		}
	}

	if c.TracingService != nil {
		if err := c.TracingService.Close(); err != nil {
			c.Logger.Error("Failed to close tracing service")
		}
	}

	if c.Logger != nil {
		_ = c.Logger.Close()
	}

	return nil
}
