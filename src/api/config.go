package api

import "os"

type Config struct {
	// TODO: reconcile public/private
	goEnv string
	port string
	PostgresUser string
	PostgresHost string
}

func NewConfig() *Config {
	// TODO: log error if one of these values is null
	return &Config{
		goEnv: os.Getenv("GO_ENV"),
		port: os.Getenv("PORT"),
		PostgresUser: os.Getenv("POSTGRES_USER"),
		PostgresHost: os.Getenv("POSTGRES_HOST"),
	}
}