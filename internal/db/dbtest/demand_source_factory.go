package dbtest

import (
	"fmt"

	"github.com/bidon-io/bidon-backend/internal/db"
)

type DemandSourceFactory struct {
	APIKey    func(int) string
	HumanName func(int) string
}

func (f DemandSourceFactory) Build(i int) db.DemandSource {
	demandSource := db.DemandSource{}

	if f.APIKey == nil {
		demandSource.APIKey = fmt.Sprintf("apikey%d", i)
	} else {
		demandSource.APIKey = f.APIKey(i)
	}

	if f.HumanName == nil {
		demandSource.HumanName = fmt.Sprintf("Demand Source %d", i)
	} else {
		demandSource.HumanName = f.HumanName(i)
	}

	return demandSource
}
