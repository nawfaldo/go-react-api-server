package config

import (
	"github.com/joho/godotenv"
)

type Config struct {
	DBUser    string
	DBName    string
	DBPwd     string
	SesSecret string
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()

	return Config{
		DBUser:    "tasgembo_nawfaldo",
		DBPwd:     "Ariyanto88$",
		DBName:    "tasgembo_dc_nawfaldo",
		SesSecret: "510147",
	}
}
