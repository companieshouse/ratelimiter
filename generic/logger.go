package generic

import "log"

// Logger defines the interface for the logger used by ratelimiter
type Logger interface {
	Debug(format string, v ...interface{})
	Info(format string, v ...interface{})
	Error(format string, v ...interface{})
}

// DefaultLogger defines a wrapper around core log.Logger that satisfies the
// required interface. It is recommended that you use a more fully implemented
// logger for production use.
type DefaultLogger struct{}

// Debug is an alias to log.Printf
func (l *DefaultLogger) Debug(format string, v ...interface{}) {
	log.Printf(format, v)
}

// Info is an alias to log.Printf
func (l *DefaultLogger) Info(format string, v ...interface{}) {
	log.Printf(format, v)
}

// Error is an alias to log.Printf (error doesn't count as fatal)
func (l *DefaultLogger) Error(format string, v ...interface{}) {
	log.Printf(format, v)
}
