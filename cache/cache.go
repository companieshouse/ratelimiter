package cache

import "github.com/companieshouse/gotools/log"

// Limiter defines an interface for ratelimiter implementations
type Limiter interface {
	Limit(string, int, int) (bool, int, int, error)
}

var logger log.Glogger // Replace with generic via interface

func handleUnexpected(err error) (bool, error) {
	if err != nil {
		logger.Error("Error: [%s]", err.Error())
		return true, err
	}
	return false, nil
}