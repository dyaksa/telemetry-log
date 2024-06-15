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

type OptFunc func(cmd *CMD) error

type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

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

func WithLevel(l Level) OptFunc {
	return func(z *CMD) (err error) {
		if l < LevelDebug || l > LevelFatal {
			return fmt.Errorf("invalid level: %d", l)
		}

		z.lvl = l
		return
	}
}

func JSONFormatter() OptFunc {
	return func(z *CMD) (err error) {
		z.lg.SetFormatter(&logrus.JSONFormatter{})
		return
	}
}

func WithHook(hook logrus.Hook) OptFunc {
	return func(z *CMD) (err error) {
		z.lg.AddHook(hook)
		return
	}
}

type CMD struct {
	lg       *logrus.Logger
	lvl      Level
	ctxFunc  []log.LogContextFunc
	errTrace *err.ErrorTracer
}

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

func (l *CMD) logWithFields(fn ...log.LogContextFunc) (entry *logrus.Entry) {
	ctx := newLoggerContext(append(l.ctxFunc, fn...)...)
	mergedFields := mergeFields(ctx.fields)
	entry = l.lg.WithFields(mergedFields)
	return
}

func mergeFields(fields []logrus.Fields) logrus.Fields {
	merged := logrus.Fields{}
	for _, fieldSet := range fields {
		for key, value := range fieldSet {
			merged[key] = value
		}
	}
	return merged
}

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

func (l CMD) Debug(message string, fn ...log.LogContextFunc) {
	if l.lvl > LevelDebug {
		return
	}

	fn = append(fn, addTraceInfo())

	l.logWithFields(fn...).Debug(message)
}

func (l CMD) Info(message string, fn ...log.LogContextFunc) {
	if l.lvl > LevelInfo {
		return
	}

	fn = append(fn, addTraceInfo())

	l.logWithFields(fn...).Info(message)
}

func (l CMD) Warn(message string, fn ...log.LogContextFunc) {
	if l.lvl > LevelWarn {
		return
	}

	fn = append(fn, addTraceInfo())

	l.logWithFields(fn...).Warn(message)
}

func (l CMD) Error(message string, fn ...log.LogContextFunc) {
	if l.lvl > LevelError {
		return
	}

	fn = append(fn, addTraceInfo())

	l.logWithFields(fn...).Error(message)
}

func (l CMD) Fatal(message string, fn ...log.LogContextFunc) {
	if l.lvl > LevelFatal {
		return
	}

	fn = append(fn, addTraceInfo())

	l.logWithFields(fn...).Fatal(message)
}

func (l CMD) WithCtx(fn log.LogContextFunc) log.Logger {
	newLogger := l
	newLogger.ctxFunc = append(newLogger.ctxFunc, fn)
	return &newLogger
}

func (l CMD) WithTrace(err error) log.Logger {
	newLogger := l
	if err != nil {
		err = l.errTrace.Err()
		errors.As(err, &l.errTrace)
	}
	newLogger.ctxFunc = append(newLogger.ctxFunc, log.Any("trace", l.errTrace.Print()))
	return &newLogger
}

type loggerContext struct {
	fields []logrus.Fields
}

func newLoggerContext(fn ...log.LogContextFunc) *loggerContext {
	lc := &loggerContext{fields: make([]logrus.Fields, 0, len(fn))}
	for _, fn := range fn {
		fn(lc)
	}
	return lc
}

func (lc *loggerContext) Any(key string, value any) {
	lc.fields = append(lc.fields, logrus.Fields{key: value})
}

func (lc *loggerContext) Bool(key string, value bool) {
	lc.fields = append(lc.fields, logrus.Fields{key: value})
}

func (lc *loggerContext) ByteString(key string, value []byte) {
	lc.fields = append(lc.fields, logrus.Fields{key: value})
}

func (lc *loggerContext) String(key string, value string) {
	lc.fields = append(lc.fields, logrus.Fields{key: value})
}

func (lc *loggerContext) Float64(key string, value float64) {
	lc.fields = append(lc.fields, logrus.Fields{key: value})
}

func (lc *loggerContext) Int64(key string, value int64) {
	lc.fields = append(lc.fields, logrus.Fields{key: value})
}

func (lc *loggerContext) Uint64(key string, value uint64) {
	lc.fields = append(lc.fields, logrus.Fields{key: value})
}

func (lc *loggerContext) Time(key string, value time.Time) {
	lc.fields = append(lc.fields, logrus.Fields{key: value})
}

func (lc *loggerContext) Error(key string, value error) {
	lc.fields = append(lc.fields, logrus.Fields{key: value})
}
