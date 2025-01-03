package config

import "os"

type DemandConfig struct {
	MetaAppSecret  string
	MetaPlatformID string
}

func NewDemandConfig() *DemandConfig {
	return &DemandConfig{
		MetaAppSecret:  os.Getenv("DEMAND_META_APP_SECRET"),
		MetaPlatformID: os.Getenv("DEMAND_META_PLATFORM_ID"),
	}
}
