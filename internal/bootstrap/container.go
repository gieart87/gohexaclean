package bootstrap

import (
	"context"
	"fmt"

	"github.com/gieart87/gohexaclean/internal/adapter/inbound/consumer"
	"github.com/gieart87/gohexaclean/internal/adapter/inbound/grpc/handler"
	"github.com/gieart87/gohexaclean/internal/adapter/outbound/datadog"
	"github.com/gieart87/gohexaclean/internal/adapter/outbound/event"
	"github.com/gieart87/gohexaclean/internal/adapter/outbound/otel"
	"github.com/gieart87/gohexaclean/internal/adapter/outbound/pgsql"
	"github.com/gieart87/gohexaclean/internal/adapter/outbound/redis"
	"github.com/gieart87/gohexaclean/internal/app"
	brokerFactory "github.com/gieart87/gohexaclean/internal/infra/broker"
	"github.com/gieart87/gohexaclean/internal/infra/cache"
	"github.com/gieart87/gohexaclean/internal/infra/config"
	"github.com/gieart87/gohexaclean/internal/infra/db"
	"github.com/gieart87/gohexaclean/internal/infra/logger"
	asynqInfra "github.com/gieart87/gohexaclean/internal/infrastructure/asynq"
	"github.com/gieart87/gohexaclean/internal/port/inbound"
	"github.com/gieart87/gohexaclean/internal/port/outbound/broker"
	"github.com/gieart87/gohexaclean/internal/port/outbound/repository"
	"github.com/gieart87/gohexaclean/internal/port/outbound/service"
	"github.com/gieart87/gohexaclean/internal/port/outbound/telemetry"
	"github.com/hibiken/asynq"
	redisClient "github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// Container holds all application dependencies
type Container struct {
	Config *config.Config
	Logger *logger.Logger

	// Database
	DB          *gorm.DB
	RedisClient *redisClient.Client

	// Repositories
	UserRepository repository.UserRepository

	// Services
	CacheService service.CacheService

	// Message Broker
	MessageBroker   broker.MessageBroker
	EventPublisher  *event.UserEventPublisher
	EventConsumer   *consumer.UserEventConsumer

	// Background Jobs
	TaskClient *asynq.Client

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

	// Initialize database with GORM
	database, err := db.NewGormConnection(&cfg.Database)
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

	// Initialize Asynq task client for background jobs
	if container.RedisClient != nil {
		redisAddr := fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port)
		container.TaskClient = asynqInfra.NewClient(redisAddr)
		log.Info("Asynq task client initialized")
	} else {
		log.Warn("Redis not available, background jobs will be disabled")
	}

	// Initialize message broker
	if cfg.Broker.Enabled {
		messageBroker, err := brokerFactory.NewMessageBroker(&cfg.Broker)
		if err != nil {
			log.Warn("Failed to create message broker, events will be disabled: " + err.Error())
		} else {
			if err := messageBroker.Connect(ctx); err != nil {
				log.Warn("Failed to connect to message broker, events will be disabled: " + err.Error())
			} else {
				container.MessageBroker = messageBroker
				log.Info("Message broker connected successfully")

				// Initialize event publisher
				container.EventPublisher = event.NewUserEventPublisher(messageBroker)

				// Initialize event consumer
				container.EventConsumer = consumer.NewUserEventConsumer(messageBroker)
				if err := container.EventConsumer.Start(ctx); err != nil {
					log.Warn("Failed to start event consumer: " + err.Error())
				} else {
					log.Info("Event consumer started successfully")
				}
			}
		}
	} else {
		log.Info("Message broker is disabled")
	}

	// Initialize use cases / application services
	container.UserService = app.NewUserService(
		container.UserRepository,
		container.CacheService,
		&cfg.JWT,
		container.EventPublisher,
		container.TaskClient,
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

	// Close Asynq task client
	if c.TaskClient != nil {
		if err := c.TaskClient.Close(); err != nil {
			c.Logger.Error("Failed to close Asynq task client")
		}
	}

	// Close message broker
	if c.EventConsumer != nil {
		if err := c.EventConsumer.Stop(); err != nil {
			c.Logger.Error("Failed to stop event consumer")
		}
	}

	if c.MessageBroker != nil {
		if err := c.MessageBroker.Close(); err != nil {
			c.Logger.Error("Failed to close message broker")
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
