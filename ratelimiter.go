package ratelimiter

import (
	"github.com/companieshouse/gotools/log"
	"github.com/companieshouse/ratelimiter/cache"
	"github.com/garyburd/redigo/redis"
)

var logger log.Glogger // Replace with generic via interface?

// NewRateLimiter creates a new instance of rateLimiter
// If not supplied with a redis connection pool, will use in memory caching instead
func NewRateLimiter(pool *redis.Pool) *Limiter {
	if pool != nil {
		logger.Info("Creating rate limiter with redis cache")
		return &Limiter{
			cache: &cache.RedisLimiter{Pool: pool},
		}
	}
	logger.Info("Creating rate limiter with in-memory cache")
	return &Limiter{
		cache: &cache.InMemoryLimiter{},
	}
}

// Limiter implements a default rate limiter wrapper
type Limiter struct {
	cache cache.Limiter
}

// Limit implements a rate limiter
func (lim *Limiter) Limit(identity string, limit int, window int) (rateLimitExceeded bool, remaining int, reset int, lastError error) {
	if limit == -1 {
		return false, limit, window, nil // Unlimited user
	}
	return lim.cache.Limit(identity, limit, window)
}
