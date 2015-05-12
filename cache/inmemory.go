package cache

import "time"

// InMemoryLimiter defines a redis backed rate limiter implementation
type InMemoryLimiter struct {
}

// Limit provides rate limiting functionality
func (rl *InMemoryLimiter) Limit(identity string, limit int, window time.Duration) (rateLimitExceeded bool, remaining int, reset time.Duration, lastError error) {

	logger.Debug("Rate limiting (in memory) for identity: [%s] Limit: [%d] Window: [%d]", identity, limit, window)

	rateLimitExceeded = false
	remaining = limit
	reset = window
	lastError = nil

	return
}

// QueryLimit provides a stubbed fake QueryLimit for inmemory operation
func (rl *InMemoryLimiter) QueryLimit(identity string) (remain int, err error) {
	return 0, nil
}
