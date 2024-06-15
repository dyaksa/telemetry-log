// Package telemetry provides functionality for telemetry logging.
package telemetry

import (
	"fmt"
	"github.com/dyaksa/telemetry-log/cmd"
	"github.com/dyaksa/telemetry-log/telemetry/log"
	"github.com/dyaksa/telemetry-log/telemetry/mongo"
	"time"
)

// OptFunc is a type that defines a function that modifies a Lib instance.
type OptFunc func(*Lib) error

// Lib is a struct that holds the necessary information for telemetry logging.
type Lib struct {
	Level string `env:"TELEMETRY_LOG_LEVEL" envDefault:"debug" json:"level"`

	Host     string `env:"TELEMETRY_HOST" envDefault:"127.0.0.1" json:"host"`
	Port     string `env:"TELEMETRY_PORT" envDefault:"27017" json:"port"`
	Username string `env:"TELEMETRY_USERNAME" envDefault:"username" json:"username"`
	Password string `env:"TELEMETRY_PASSWORD" envDefault:"password" json:"password"`
	withHook bool

	mc  *mongo.Mongo
	Log log.Logger

	logOpt    []cmd.OptFunc
	mongoOpts []mongo.OptFunc
}

// WithJSONFormatter is a function that returns an OptFunc which sets the JSON formatter for a Lib instance.
func WithJSONFormatter() OptFunc {
	return func(li *Lib) (err error) {
		li.logOpt = append(li.logOpt, cmd.JSONFormatter())
		return
	}
}

// WithHook is a function that returns an OptFunc which sets the hook for a Lib instance.
func WithHook() OptFunc {
	return func(li *Lib) (err error) {
		li.withHook = true
		return
	}
}

// New is a function that creates a new Lib instance.
// It applies the provided options to the Lib instance and then attempts to initialize the environment and command.
func New(opts ...OptFunc) (li *Lib, err error) {
	li = &Lib{withHook: false}

	if err = LoadEnv(li); err != nil {
		return nil, fmt.Errorf("fail to load env: %w", err)
	}

	for _, opt := range opts {
		if err = opt(li); err != nil {
			return nil, fmt.Errorf("fail to apply options: %w", err)
		}
	}

	if err = li.initConnection(); err != nil {
		return nil, fmt.Errorf("fail to init connection: %w", err)
	}

	if err = li.initCMD(); err != nil {
		return nil, fmt.Errorf("fail to init cmd: %w", err)
	}

	return li, nil
}

// initEnv is a method that loads the environment variables for a Lib instance.
func (li *Lib) initEnv() (err error) {
	if err = LoadEnv(li); err != nil {
		return fmt.Errorf("fail to load env: %w", err)
	}
	return
}

// initCMD is a method that initializes the command for a Lib instance.
func (li *Lib) initCMD() (err error) {
	mongoHook := &MongoHook{
		Client:   li.mc,
		Timeout:  5 * time.Second,
		WithHook: li.withHook,
	}

	li.logOpt = append(li.logOpt, cmd.WithLogLevel(li.Level))
	li.logOpt = append(li.logOpt, cmd.WithHook(mongoHook))

	li.Log, err = cmd.New(li.logOpt...)

	if err != nil {
		return fmt.Errorf("fail to create log: %w", err)
	}

	return
}

// initConnection is a method that initializes the connection for a Lib instance.
func (li *Lib) initConnection() (err error) {
	li.mongoOpts = append(li.mongoOpts, mongo.WithConnection(li.Host, li.Port, li.Username, li.Password))
	li.mc, err = mongo.New(li.mongoOpts...)

	if err != nil {
		return fmt.Errorf("fail to create mongo connection: %w", err)
	}

	return
}
