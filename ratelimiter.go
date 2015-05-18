/*
Package ratelimiter provides a redis backed ratelimiter implementation.

Based on Redis LUA examples from http://redis.io/commands/incr


Unlimited limits

There are of course times when you want a user to not be subject to ratelimiting.
You could just bypass the call to Limit, but that's needless work on your behalf.
Instead, just set the limit to -1:

    exceeded, remaining, reset, err := limiter.Limit("MyIdentity", -1, 0)
*/
package ratelimiter

import (
	"time"

	"github.com/companieshouse/ratelimiter/cache"
	"github.com/garyburd/redigo/redis"
)

// NewRateLimiter creates a new instance of rateLimiter
// If not supplied with a redis connection pool, will use in memory caching instead
func NewRateLimiter(pool *redis.Pool, logger Logger) *Limiter {

	if logger == nil {
		logger = &DefaultLogger{}
	}

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
func (lim *Limiter) Limit(identity string, limit int, window time.Duration) (rateLimitExceeded bool, remaining int, reset time.Duration, lastError error) {
	if limit == -1 {
		return false, limit, window, nil // Unlimited user
	}
	return lim.cache.Limit(identity, limit, window)
}

// QueryLimit returns the current limit defined for the given identity
func (lim *Limiter) QueryLimit(identity string, limit int, window time.Duration) (remain int, reset time.Duration, err error) {
	return lim.cache.QueryLimit(identity, limit, window)
}
