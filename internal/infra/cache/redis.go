package cache

import (
	"context"
	"fmt"

	"github.com/gieart87/gohexaclean/internal/infra/config"
	"github.com/redis/go-redis/v9"
)

// NewRedisClient creates a new Redis client
func NewRedisClient(cfg *config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.GetRedisAddr(),
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
	})

	// Test connection
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	return client, nil
}

// Close closes the Redis client
func Close(client *redis.Client) error {
	if client != nil {
		return client.Close()
	}
	return nil
}
