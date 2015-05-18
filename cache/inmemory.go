package cache

import (
	"time"

	"github.com/companieshouse/ratelimiter/generic"
)

// InMemoryLimiter defines a stubbed in-memory rate limiter implementation.
//
// This implementation does not provide real rate limiting. Instead it is intended
// as a stubbed interface for mocking/testing when you don't have (or want to use)
// a real redis instance.
type InMemoryLimiter struct {
	Logger generic.Logger
}

// Limit provides stubbed rate limiting functionality. For the in-memory implementation
// this will always return rateLimitExceeded=false
func (rl *InMemoryLimiter) Limit(identity string, limit int, window time.Duration) (rateLimitExceeded bool, remaining int, reset time.Duration, lastError error) {

	rl.Logger.Debug("Rate limiting (in memory) for identity: [%s] Limit: [%d] Window: [%d]", identity, limit, window)

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
