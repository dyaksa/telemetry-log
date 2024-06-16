package main_test

import (
	"errors"
	"testing"
	"time"

	"github.com/dyaksa/telemetry-log/telemetry/log"
)

type mockLogContext struct{}

func (m *mockLogContext) Any(key string, value any)           {}
func (m *mockLogContext) Bool(key string, value bool)         {}
func (m *mockLogContext) ByteString(key string, value []byte) {}
func (m *mockLogContext) String(key string, value string)     {}
func (m *mockLogContext) Float64(key string, value float64)   {}
func (m *mockLogContext) Int64(key string, value int64)       {}
func (m *mockLogContext) Uint64(key string, value uint64)     {}
func (m *mockLogContext) Time(key string, value time.Time)    {}
func (m *mockLogContext) Error(key string, value error)       {}

func TestLogContextFuncs(t *testing.T) {
	mockContext := &mockLogContext{}

	log.Any("key", "value")(mockContext)
	log.Bool("key", true)(mockContext)
	log.ByteString("key", []byte("value"))(mockContext)
	log.String("key", "value")(mockContext)
	log.Float64("key", 1.23)(mockContext)
	log.Int64("key", 123)(mockContext)
	log.Uint64("key", 123)(mockContext)
	log.Time("key", time.Now())(mockContext)
	log.Error("key", errors.New("error"))(mockContext)
}

func TestLogContextFuncsWithLoggable(t *testing.T) {
	mockContext := &mockLogContext{}
	mockLoggable := &mockLoggable{}

	log.Any("key", mockLoggable)(mockContext)
}

type mockLoggable struct{}

func (m *mockLoggable) AsLog() any {
	return "loggable"
}
