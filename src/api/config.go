package api

import "os"

// Config defines the shape of the configurable parameters for the API.
type Config struct {
	goEnv        string
	port         string
	PostgresUser string
	PostgresHost string
}

// NewConfig constructs a config struct, reading values from ENV s.t it's easy
// to swap these values out for different environments.
func NewConfig() *Config {
	return &Config{
		goEnv:        os.Getenv("GO_ENV"),
		port:         os.Getenv("PORT"),
		PostgresUser: os.Getenv("POSTGRES_USER"),
		PostgresHost: os.Getenv("POSTGRES_HOST"),
	}
}
