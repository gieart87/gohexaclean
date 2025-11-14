package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gieart87/gohexaclean/internal/port/outbound/service"
	"github.com/redis/go-redis/v9"
)

// CacheServiceRedis implements CacheService interface for Redis
type CacheServiceRedis struct {
	client *redis.Client
}

// NewCacheServiceRedis creates a new Redis cache service
func NewCacheServiceRedis(client *redis.Client) service.CacheService {
	return &CacheServiceRedis{client: client}
}

// Get retrieves a value from cache
func (s *CacheServiceRedis) Get(ctx context.Context, key string) (string, error) {
	val, err := s.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("key not found")
	}
	if err != nil {
		return "", fmt.Errorf("failed to get cache: %w", err)
	}
	return val, nil
}

// Set sets a value in cache
func (s *CacheServiceRedis) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	var val string
	switch v := value.(type) {
	case string:
		val = v
	default:
		b, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal value: %w", err)
		}
		val = string(b)
	}

	err := s.client.Set(ctx, key, val, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set cache: %w", err)
	}
	return nil
}

// Delete deletes a value from cache
func (s *CacheServiceRedis) Delete(ctx context.Context, key string) error {
	err := s.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete cache: %w", err)
	}
	return nil
}

// Exists checks if a key exists in cache
func (s *CacheServiceRedis) Exists(ctx context.Context, key string) (bool, error) {
	val, err := s.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check existence: %w", err)
	}
	return val > 0, nil
}

// SetNX sets a value only if it doesn't exist
func (s *CacheServiceRedis) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	var val string
	switch v := value.(type) {
	case string:
		val = v
	default:
		b, err := json.Marshal(value)
		if err != nil {
			return false, fmt.Errorf("failed to marshal value: %w", err)
		}
		val = string(b)
	}

	result, err := s.client.SetNX(ctx, key, val, expiration).Result()
	if err != nil {
		return false, fmt.Errorf("failed to set cache: %w", err)
	}
	return result, nil
}
