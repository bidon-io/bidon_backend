package dbtest

import (
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/bidon-io/bidon-backend/internal/db"
)

var demandSourceCounter atomic.Uint64

func demandSourceDefaults(n uint64) func(*db.DemandSource) {
	return func(source *db.DemandSource) {
		if source.APIKey == "" {
			source.APIKey = fmt.Sprintf("apikey%d", n)
		}
		if source.HumanName == "" {
			source.HumanName = fmt.Sprintf("Demand Source %d", n)
		}
	}
}

func BuildDemandSource(opts ...func(*db.DemandSource)) db.DemandSource {
	var source db.DemandSource

	n := demandSourceCounter.Add(1)

	opts = append(opts, demandSourceDefaults(n))
	for _, opt := range opts {
		opt(&source)
	}

	return source
}

func CreateDemandSource(t *testing.T, tx *db.DB, opts ...func(*db.DemandSource)) db.DemandSource {
	t.Helper()

	demandSource := BuildDemandSource(opts...)
	if err := tx.Create(&demandSource).Error; err != nil {
		t.Fatalf("Failed to create demand source: %v", err)
	}

	return demandSource
}
