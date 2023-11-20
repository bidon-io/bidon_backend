package config

import "go.uber.org/zap"

func NewLogger() (*zap.Logger, error) {
	if GetEnv() == ProdEnv {
		return zap.NewProduction()
	}

	return zap.NewDevelopment()
}
