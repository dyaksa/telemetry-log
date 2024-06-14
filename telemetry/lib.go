package telemetry

import (
	"fmt"
	"github.com/dyaksa/telemetry-log/telemetry/log"
	"github.com/dyaksa/telemetry-log/telemetry/mongo"
	"time"
)

type OptFunc func(*Lib) error

type Lib struct {
	Level string `env:"TELEMETRY_LOG_LEVEL" envDefault:"debug" json:"level"`

	Host     string `env:"TELEMETRY_HOST" json:"host"`
	Port     string `env:"TELEMETRY_PORT" json:"port"`
	Username string `env:"TELEMETRY_USERNAME" json:"username"`
	Password string `env:"TELEMETRY_PASSWORD" json:"password"`

	MC  *mongo.Mongo
	Log log.Logger

	logOpt []log.OptFunc
}

func WithJSONFormatter() OptFunc {
	return func(li *Lib) (err error) {
		li.logOpt = append(li.logOpt, log.JSONFormatter())
		return
	}
}

func New(opts ...OptFunc) (li *Lib, err error) {
	li = &Lib{}

	if err = LoadEnv(li); err != nil {
		return nil, fmt.Errorf("fail to load env: %w", err)
	}

	for _, opt := range opts {
		if err := opt(li); err != nil {
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

func (li *Lib) initEnv() (err error) {
	if err = LoadEnv(li); err != nil {
		return fmt.Errorf("fail to load env: %w", err)
	}
	return
}

func (li *Lib) initCMD() (err error) {
	mongoHook := &MongoHook{
		Client:  li.MC,
		Timeout: 5 * time.Second,
	}

	li.logOpt = append(li.logOpt, log.WithLogLevel(li.Level))
	li.logOpt = append(li.logOpt, log.WithHook(mongoHook))

	l, err := log.New(li.logOpt...)

	if err != nil {
		return fmt.Errorf("fail to create log: %w", err)
	}

	li.Log = l
	return
}

func (li *Lib) initConnection() (err error) {
	opts := []mongo.OptFunc{}
	opts = append(opts, mongo.WithConnection(li.Host, li.Port, li.Username, li.Password))
	li.MC, err = mongo.New(opts...)

	if err != nil {
		return fmt.Errorf("fail to create mongo connection: %w", err)
	}

	return
}
