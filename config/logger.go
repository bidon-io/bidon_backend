package config

import "go.uber.org/zap"

func NewLogger() (*zap.Logger, error) {
	if Env == ProdEnv {
		return zap.NewProduction()
	}

	return zap.NewDevelopment()
}
