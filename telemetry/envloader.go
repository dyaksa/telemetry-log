package telemetry

import (
	"errors"
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

func LoadEnv(v interface{}) (err error) {
	err = godotenv.Load()
	opts := env.Options{TagName: "env"}
	if err = env.ParseWithOptions(v, opts); err != nil {
		err = errors.Join(err, errors.New("fail to load env"))
	}
	return
}
