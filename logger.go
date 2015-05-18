package ratelimiter

import "log"

// Logger defines the interface for the logger used by ratelimiter
type Logger interface {
	Info(format string, v ...interface{})
}

// DefaultLogger defines a wrapper around core log.Logger that satisfies the
// required interface. It is recommended that you use a more fully implemented
// logger for production use.
type DefaultLogger struct{}

// Info is an alias to log.Printf
func (l *DefaultLogger) Info(format string, v ...interface{}) {
	log.Printf(format, v)
}
