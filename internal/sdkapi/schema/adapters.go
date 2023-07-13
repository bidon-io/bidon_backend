package schema

import (
	"github.com/bidon-io/bidon-backend/internal/adapter"
	"golang.org/x/exp/maps"
)

type Adapters map[adapter.Key]Adapter

func (as Adapters) Map() map[string]any {
	m := make(map[string]any, len(as))

	for k, a := range as {
		m[string(k)] = a.Map()
	}

	return m
}

func (as Adapters) Keys() []adapter.Key {
	return maps.Keys(as)
}

type Adapter struct {
	Version    string `json:"version" validate:"required"`
	SDKVersion string `json:"sdk_version" validate:"required"`
}

func (a Adapter) Map() map[string]any {
	m := map[string]any{
		"version":     a.Version,
		"sdk_version": a.SDKVersion,
	}

	return m
}
