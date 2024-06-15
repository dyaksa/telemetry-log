// Package log provides an interface and functions for logging.
package log

import "time"

// Logger is an interface that defines methods for logging at different levels.
type Logger interface {
	Debug(message string, fn ...LogContextFunc) // Debug logs a debug message.
	Info(message string, fn ...LogContextFunc)  // Info logs an informational message.
	Warn(message string, fn ...LogContextFunc)  // Warn logs a warning message.
	Error(message string, fn ...LogContextFunc) // Error logs an error message.
	Fatal(message string, fn ...LogContextFunc) // Fatal logs a fatal error message.

	WithCtx(LogContextFunc) Logger // WithCtx returns a new Logger with the specified context.
	WithTrace(err error) Logger    // WithTrace returns a new Logger with the specified error trace.
}

// LogContextFunc is a function that modifies a LogContext.
type LogContextFunc func(LogContext)

// LogContext is an interface that defines methods for setting context values of different types.
type LogContext interface {
	Any(key string, value any)           // Any sets a context value of any type.
	Bool(key string, value bool)         // Bool sets a context value of type bool.
	ByteString(key string, value []byte) // ByteString sets a context value of type []byte.
	String(key string, value string)     // String sets a context value of type string.
	Float64(key string, value float64)   // Float64 sets a context value of type float64.
	Int64(key string, value int64)       // Int64 sets a context value of type int64.
	Uint64(key string, value uint64)     // Uint64 sets a context value of type uint64.
	Time(key string, value time.Time)    // Time sets a context value of type time.Time.
	Error(key string, value error)       // Error sets a context value of type error.
}

// Loggable is an interface that defines a method for converting an object to a loggable format.
type Loggable interface {
	AsLog() any // AsLog converts the object to a loggable format.
}

// Any is a function that returns a LogContextFunc which sets a context value of any type.
func Any(key string, value any) LogContextFunc {
	if l, ok := value.(Loggable); ok {
		value = l.AsLog()
	}

	return func(lc LogContext) {
		lc.Any(key, value)
	}
}

// Bool is a function that returns a LogContextFunc which sets a context value of type bool.
func Bool(key string, value bool) LogContextFunc {
	return func(lc LogContext) {
		lc.Bool(key, value)
	}
}

// ByteString is a function that returns a LogContextFunc which sets a context value of type []byte.
func ByteString(key string, value []byte) LogContextFunc {
	return func(lc LogContext) {
		lc.ByteString(key, value)
	}
}

// String is a function that returns a LogContextFunc which sets a context value of type string.
func String(key string, value string) LogContextFunc {
	return func(lc LogContext) {
		lc.String(key, value)
	}
}

// Float64 is a function that returns a LogContextFunc which sets a context value of type float64.
func Float64(key string, value float64) LogContextFunc {
	return func(lc LogContext) {
		lc.Float64(key, value)
	}
}

// Int64 is a function that returns a LogContextFunc which sets a context value of type int64.
func Int64(key string, value int64) LogContextFunc {
	return func(lc LogContext) {
		lc.Int64(key, value)
	}
}

// Uint64 is a function that returns a LogContextFunc which sets a context value of type uint64.
func Uint64(key string, value uint64) LogContextFunc {
	return func(lc LogContext) {
		lc.Uint64(key, value)
	}
}

// Time is a function that returns a LogContextFunc which sets a context value of type time.Time.
func Time(key string, value time.Time) LogContextFunc {
	return func(lc LogContext) {
		lc.Time(key, value)
	}
}

// Error is a function that returns a LogContextFunc which sets a context value of type error.
func Error(key string, value error) LogContextFunc {
	return func(lc LogContext) {
		lc.Error(key, value)
	}
}
