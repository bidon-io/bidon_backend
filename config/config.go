// Package config provides configuration for different parts of Bidon services, that is shared between them
package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

const (
	ProdEnv    = "production"
	DevEnv     = "development"
	TestEnv    = "test"
	UnknownEnv = ""
)

func GetEnv() string {
	switch env := os.Getenv("ENVIRONMENT"); env {
	case ProdEnv:
		return ProdEnv
	case TestEnv:
		return TestEnv
	case DevEnv:
		return DevEnv
	default:
		return UnknownEnv
	}
}

func Debug() bool {
	return GetEnv() != ProdEnv
}

func LoadEnvFile() {
	var err error
	if GetEnv() == TestEnv {
		err = godotenv.Load(".env.test")
	} else {
		err = godotenv.Load()
	}
	if err != nil {
		log.Printf("Did not load .env file: %v", err)
	}
}
