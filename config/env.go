package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUser    string
	DBName    string
	SesSecret string
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()

	return Config{
		DBUser:    os.Getenv("DB_USER"),
		DBName:    os.Getenv("DB_NAME"),
		SesSecret: os.Getenv("SESSION_SECRET"),
	}
}
