package cache

import "errors"

// Cache infrastructure errors
var (
	ErrCacheConnection  = errors.New("cache connection failed")
	ErrCacheTimeout     = errors.New("cache operation timeout")
	ErrCacheKeyNotFound = errors.New("key not found in cache")
	ErrCacheMarshal     = errors.New("failed to marshal cache data")
	ErrCacheUnmarshal   = errors.New("failed to unmarshal cache data")
	ErrCacheExpired     = errors.New("cache entry expired")
)
