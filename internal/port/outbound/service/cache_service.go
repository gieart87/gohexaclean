package service

import (
	"context"
	"time"
)

// CacheService defines the outbound port for caching
type CacheService interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error)
}
