package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

const ENV = "ENV"
const PROD = "PROD"

type Configuration struct {
	PORT      string
	LOG_LEVEL string
	LOG_FILE_PATH string
}

func MustGetEnvVars() *Configuration {
	if os.Getenv(ENV) != PROD {
		err := godotenv.Load()
		if err != nil {
			panic(fmt.Errorf("error loading .env file"))
		}
	}

	return &Configuration{
		PORT:      GetVar("PORT", true),
		LOG_LEVEL: GetVar("LOG_LEVEL", true),
		LOG_FILE_PATH: GetVar("LOG_FILE_PATH", false),
	}
}

func GetVar(key string, required bool) string {
	envVar := os.Getenv(key)

	if envVar == "" && required{
		panic(fmt.Errorf("environemnt variable %s can't be empty", key))
	}

	return envVar
}
