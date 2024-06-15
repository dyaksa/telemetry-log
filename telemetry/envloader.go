// Package telemetry provides functionality for telemetry logging.
package telemetry

import (
	"errors"
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

// LoadEnv is a function that loads environment variables from a .env file into a given struct.
// It uses the "env" tag in the struct to match the environment variables.
// If the .env file cannot be loaded or the environment variables cannot be parsed, it returns an error.
func LoadEnv(v interface{}) (err error) {
	// Load environment variables from .env file
	err = godotenv.Load()

	// Define options for parsing environment variables
	opts := env.Options{TagName: "env"}

	// Parse environment variables into the given struct
	if err = env.ParseWithOptions(v, opts); err != nil {
		// If an error occurs, join it with a new error message
		err = errors.Join(err, errors.New("fail to load env"))
	}

	// Return any error that occurred
	return
}
