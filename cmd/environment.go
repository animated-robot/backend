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
}

func MustGetEnvVars() *Configuration {
	if os.Getenv(ENV) != PROD {
		err := godotenv.Load()
		if err != nil {
			panic(fmt.Errorf("error loading .env file"))
		}
	}

	return &Configuration{
		PORT:      GetVar("PORT"),
		LOG_LEVEL: GetVar("LOG_LEVEL"),
	}
}

func GetVar(key string) string {
	envVar := os.Getenv(key)

	if envVar == "" {
		panic(fmt.Errorf("environemnt variable %s can't be empty", key))
	}

	return envVar
}
