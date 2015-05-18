// Package cache provides implementations for the ratelimiter.
package cache

import "time"

// Limiter defines an interface for ratelimiter implementations
type Limiter interface {
	Limit(identity string, limit int, window time.Duration) (isLimitExceeded bool, remaining int, reset time.Duration, err error)
	QueryLimit(identity string) (remaining int, err error)
}
