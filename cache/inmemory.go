package cache

import "github.com/garyburd/redigo/redis"

// InMemoryLimiter defines a redis backed rate limiter implementation
type InMemoryLimiter struct {
	Pool *redis.Pool
}

// InMemoryLimiter provides rate limiting functionality
func (rl *InMemoryLimiter) Limit(identity string, limit int, window int) (rateLimitExceeded bool, remaining int, reset int, lastError error) {

	logger.Debug("Rate limiting (in memory) for identity: [%s] Limit: [%s] Window: [%s]", identity, limit, window)

	rateLimitExceeded = false
	remaining = limit
	reset = window
	lastError = nil

	return
}
