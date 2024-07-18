package main

import (
	"os"
)

type Config struct {
	DBUser string
	DBName string
}

var Envs = initConfig()

func initConfig() Config {
	return Config{
		DBUser: getEnv("DB_USER", "root"),
		DBName: getEnv("DB_NAME", "go_test"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}
