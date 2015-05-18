// Package cache provides implementations for the ratelimiter.
package cache

import (
	"time"

	"github.com/companieshouse/gotools/log"
)

// Limiter defines an interface for ratelimiter implementations
type Limiter interface {
	Limit(identity string, limit int, window time.Duration) (isLimitExceeded bool, remaining int, reset time.Duration, err error)
	QueryLimit(identity string) (remaining int, err error)
}

var logger log.Glogger // Replace with generic via interface

func handleUnexpected(err error) (bool, error) {
	if err != nil {
		logger.Error("Error: [%s]", err.Error())
		return true, err
	}
	return false, nil
}
