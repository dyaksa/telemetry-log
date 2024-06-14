package log

import "time"

type Logger interface {
	Debug(message string, fn ...LogContextFunc)
	Info(message string, fn ...LogContextFunc)
	Warn(message string, fn ...LogContextFunc)
	Error(message string, fn ...LogContextFunc)
	Fatal(message string, fn ...LogContextFunc)

	WithCtx(LogContextFunc) Logger
	WithTrace(err error) Logger
}

type LogContextFunc func(LogContext)

type LogContext interface {
	Any(key string, value any)
	Bool(key string, value bool)
	ByteString(key string, value []byte)
	String(key string, value string)
	Float64(key string, value float64)
	Int64(key string, value int64)
	Uint64(key string, value uint64)
	Time(key string, value time.Time)
	Error(key string, value error)
}

type Loggable interface {
	AsLog() any
}

func Any(key string, value any) LogContextFunc {
	if l, ok := value.(Loggable); ok {
		value = l.AsLog()
	}

	return func(lc LogContext) {
		lc.Any(key, value)
	}
}

func Bool(key string, value bool) LogContextFunc {
	return func(lc LogContext) {
		lc.Bool(key, value)
	}
}

func ByteString(key string, value []byte) LogContextFunc {
	return func(lc LogContext) {
		lc.ByteString(key, value)
	}
}

func String(key string, value string) LogContextFunc {
	return func(lc LogContext) {
		lc.String(key, value)
	}
}

func Float64(key string, value float64) LogContextFunc {
	return func(lc LogContext) {
		lc.Float64(key, value)
	}
}

func Int64(key string, value int64) LogContextFunc {
	return func(lc LogContext) {
		lc.Int64(key, value)
	}
}

func Uint64(key string, value uint64) LogContextFunc {
	return func(lc LogContext) {
		lc.Uint64(key, value)
	}
}

func Time(key string, value time.Time) LogContextFunc {
	return func(lc LogContext) {
		lc.Time(key, value)
	}
}

func Error(key string, value error) LogContextFunc {
	return func(lc LogContext) {
		lc.Error(key, value)
	}
}
