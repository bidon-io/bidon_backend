package schema

import (
	"golang.org/x/exp/maps"

	"github.com/bidon-io/bidon-backend/internal/adapter"
)

type Adapters map[adapter.Key]Adapter

func (as Adapters) Keys() []adapter.Key {
	return maps.Keys(as)
}

type Adapter struct {
	Version    string `json:"version" validate:"required"`
	SDKVersion string `json:"sdk_version" validate:"required"`
}
