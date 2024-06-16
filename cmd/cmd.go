// Package cmd provides functionality for command line operations.
package cmd

import (
	"errors"
	"fmt"
	"github.com/dyaksa/telemetry-log/err"
	"github.com/dyaksa/telemetry-log/telemetry/log"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// OptFunc is a type that defines a function that modifies a CMD instance.
type OptFunc func(cmd *CMD) error

// Level is a type that defines the logging level.
type Level int

// These constants represent the different logging levels.
const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

// WithLogLevel is a function that returns an OptFunc which sets the logging level of a CMD instance.
func WithLogLevel(s string) OptFunc {
	var l Level = -1
	switch strings.ToLower(s) {
	case "debug":
		l = 0
	case "info":
		l = 1
	case "warn":
		l = 2
	case "error":
		l = 3
	case "fatal":
		l = 4
	}
	return WithLevel(l)
}

// WithLevel is a function that returns an OptFunc which sets the logging level of a CMD instance.
func WithLevel(l Level) OptFunc {
	return func(z *CMD) (err error) {
		if l < LevelDebug || l > LevelFatal {
			return fmt.Errorf("invalid level: %d", l)
		}

		z.lvl = l
		return
	}
}

// JSONFormatter is a function that returns an OptFunc which sets the JSON formatter for a CMD instance.
func JSONFormatter() OptFunc {
	return func(l *CMD) (err error) {
		l.lg.SetFormatter(&logrus.JSONFormatter{})
		return
	}
}

// WithHook is a function that returns an OptFunc which sets the hook for a CMD instance.
func WithHook(hook logrus.Hook) OptFunc {
	return func(l *CMD) (err error) {
		l.lg.AddHook(hook)
		return
	}
}

// CMD is a struct that holds the necessary information for command line operations.
type CMD struct {
	lg       *logrus.Logger
	lvl      Level
	ctxFunc  []log.LogContextFunc
	errTrace *err.ErrorTracer
}

// New is a function that creates a new CMD instance.
// It applies the provided options to the CMD instance.
func New(opts ...OptFunc) (l log.Logger, err error) {
	logr := logrus.New()
	lg := &CMD{
		lg:  logr,
		lvl: LevelInfo,
	}

	for _, opt := range opts {
		if err = opt(lg); err != nil {
			return
		}
	}

	l = lg
	return
}

// logWithFields is a method that logs an entry with the specified context.
func (l *CMD) logWithFields(fn ...log.LogContextFunc) (entry *logrus.Entry) {
	ctx := newLoggerContext(append(l.ctxFunc, fn...)...)
	mergedFields := mergeFields(ctx.fields)
	entry = l.lg.WithFields(mergedFields)
	return
}

// mergeFields is a function that merges multiple sets of logrus fields into one.
func mergeFields(fields []logrus.Fields) logrus.Fields {
	merged := logrus.Fields{}
	for _, fieldSet := range fields {
		for key, value := range fieldSet {
			merged[key] = value
		}
	}
	return merged
}

// addTraceInfo is a function that returns a LogContextFunc which adds trace information to the context.
func addTraceInfo() log.LogContextFunc {
	return func(ctx log.LogContext) {
		if pc, file, line, ok := runtime.Caller(4); ok {
			function := runtime.FuncForPC(pc).Name()
			ctx.Any("file", path.Base(file))
			ctx.Any("line", line)
			ctx.Any("func", function)
		}
	}
}

// Debug is a method that logs a debug message.
func (l CMD) Debug(message string, fn ...log.LogContextFunc) {
	if l.lvl > LevelDebug {
		return
	}

	fn = append(fn, addTraceInfo())

	l.logWithFields(fn...).Debug(message)
}

// Info is a method that logs an informational message.
func (l CMD) Info(message string, fn ...log.LogContextFunc) {
	if l.lvl > LevelInfo {
		return
	}

	fn = append(fn, addTraceInfo())

	l.logWithFields(fn...).Info(message)
}

// Warn is a method that logs a warning message.
func (l CMD) Warn(message string, fn ...log.LogContextFunc) {
	if l.lvl > LevelWarn {
		return
	}

	fn = append(fn, addTraceInfo())

	l.logWithFields(fn...).Warn(message)
}

// Error is a method that logs an error message.
func (l CMD) Error(message string, fn ...log.LogContextFunc) {
	if l.lvl > LevelError {
		return
	}

	fn = append(fn, addTraceInfo())

	l.logWithFields(fn...).Error(message)
}

// Fatal is a method that logs a fatal error message.
func (l CMD) Fatal(message string, fn ...log.LogContextFunc) {
	if l.lvl > LevelFatal {
		return
	}

	fn = append(fn, addTraceInfo())

	l.logWithFields(fn...).Fatal(message)
}

// WithCtx is a method that returns a new Logger with the specified context.
func (l CMD) WithCtx(fn log.LogContextFunc) log.Logger {
	newLogger := l
	newLogger.ctxFunc = append(newLogger.ctxFunc, fn)
	return &newLogger
}

// WithTrace is a method that returns a new Logger with the specified error trace.
func (l CMD) WithTrace(err error) log.Logger {
	newLogger := l
	if err != nil {
		err = l.errTrace.Err()
		errors.As(err, &l.errTrace)
	}
	newLogger.ctxFunc = append(newLogger.ctxFunc, log.Any("trace", l.errTrace.Print()))
	return &newLogger
}

func (l CMD) WithFields(fields map[string]interface{}) log.Logger {
	newLogger := l
	for key, field := range fields {
		newLogger.ctxFunc = append(newLogger.ctxFunc, log.Any(key, field))
	}
	return &newLogger
}

// loggerContext is a struct that holds the necessary information for a logger context.
type loggerContext struct {
	fields []logrus.Fields
}

// newLoggerContext is a function that creates a new logger context with the specified context functions.
func newLoggerContext(fn ...log.LogContextFunc) *loggerContext {
	lc := &loggerContext{fields: make([]logrus.Fields, 0, len(fn))}
	for _, fn := range fn {
		fn(lc)
	}
	return lc
}

// Any is a method that sets a context value of any type.
func (lc *loggerContext) Any(key string, value any) {
	lc.fields = append(lc.fields, logrus.Fields{key: value})
}

// Bool is a method that sets a context value of type bool.
func (lc *loggerContext) Bool(key string, value bool) {
	lc.fields = append(lc.fields, logrus.Fields{key: value})
}

// ByteString is a method that sets a context value of type []byte.
func (lc *loggerContext) ByteString(key string, value []byte) {
	lc.fields = append(lc.fields, logrus.Fields{key: value})
}

// String is a method that sets a context value of type string.
func (lc *loggerContext) String(key string, value string) {
	lc.fields = append(lc.fields, logrus.Fields{key: value})
}

// Float64 is a method that sets a context value of type float64.
func (lc *loggerContext) Float64(key string, value float64) {
	lc.fields = append(lc.fields, logrus.Fields{key: value})
}

// Int64 is a method that sets a context value of type int64.
func (lc *loggerContext) Int64(key string, value int64) {
	lc.fields = append(lc.fields, logrus.Fields{key: value})
}

// Uint64 is a method that sets a context value of type uint64.
func (lc *loggerContext) Uint64(key string, value uint64) {
	lc.fields = append(lc.fields, logrus.Fields{key: value})
}

// Time is a method that sets a context value of type time.Time.
func (lc *loggerContext) Time(key string, value time.Time) {
	lc.fields = append(lc.fields, logrus.Fields{key: value})
}

// Error is a method that sets a context value of type error.
func (lc *loggerContext) Error(key string, value error) {
	lc.fields = append(lc.fields, logrus.Fields{key: value})
}
