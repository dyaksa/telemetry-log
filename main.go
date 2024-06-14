package main

import (
	"errors"
	"github.com/dyaksa/telemetry-log/telemetry"
)

func main() {
	tel, err := telemetry.New(
		telemetry.WithJSONFormatter(),
	)
	if err != nil {
		panic(err)
	}

	tel.Log.WithTrace(errors.New("Internal server error")).Error("internal error")
}
