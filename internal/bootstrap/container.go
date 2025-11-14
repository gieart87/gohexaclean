package bootstrap

import (
	"database/sql"
	"fmt"

	"github.com/gieart87/gohexaclean/internal/adapter/inbound/grpc/handler"
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

	if c.Logger != nil {
		_ = c.Logger.Close()
	}

	return nil
}
