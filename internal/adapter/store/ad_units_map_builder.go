package store

import (
	"context"
	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/auction"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
)

type AdUnitsMapBuilder struct {
	AdUnitsMatcher AdUnitsMatcher
}

type AdUnitsMap map[adapter.Key][]auction.AdUnit

type AdUnitsMatcher interface {
	Match(ctx context.Context, params *auction.BuildParams) ([]auction.AdUnit, error)
}

func (b *AdUnitsMapBuilder) Build(ctx context.Context, appID int64, adapterKeys []adapter.Key, imp schema.Imp) (AdUnitsMap, error) {
	adUnits, err := b.AdUnitsMatcher.Match(ctx, &auction.BuildParams{
		Adapters: adapterKeys,
		AppID:    appID,
		AdType:   imp.Type(),
		AdFormat: imp.Format(),
	})
	if err != nil {
		return nil, err
	}

	adUnitsMap := make(map[adapter.Key][]auction.AdUnit)
	for _, adUnit := range adUnits {
		key := adapter.Key(adUnit.DemandID)
		adUnitsMap[key] = append(adUnitsMap[key], adUnit)
	}

	return adUnitsMap, nil
}
