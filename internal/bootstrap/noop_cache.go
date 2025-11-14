package bootstrap

import (
	"context"
	"fmt"
	"time"
)

// NoOpCacheService is a no-op implementation of CacheService when Redis is not available
type NoOpCacheService struct{}

func (n *NoOpCacheService) Get(ctx context.Context, key string) (string, error) {
	return "", fmt.Errorf("cache not available")
}

func (n *NoOpCacheService) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return nil // no-op
}

func (n *NoOpCacheService) Delete(ctx context.Context, key string) error {
	return nil // no-op
}

func (n *NoOpCacheService) Exists(ctx context.Context, key string) (bool, error) {
	return false, nil
}

func (n *NoOpCacheService) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	return true, nil
}
