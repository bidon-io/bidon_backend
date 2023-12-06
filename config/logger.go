package config

import "go.uber.org/zap"

func NewZapLogger() (*zap.Logger, error) {
	if GetEnv() == ProdEnv {
		return zap.NewProduction()
	}

	return zap.NewDevelopment()
}
