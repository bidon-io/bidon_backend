// Package config provides configuration for different parts of Bidon services, that is shared between them
package config

import "os"

var Env = getEnv()

const (
	ProdEnv = "production"
	DevEnv  = "development"
)

func getEnv() string {
	if os.Getenv("ENVIRONMENT") == ProdEnv {
		return ProdEnv
	}

	return DevEnv
}
