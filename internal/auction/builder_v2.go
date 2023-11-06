package auction

import (
	"context"
	"errors"

	"github.com/bidon-io/bidon-backend/internal/ad"
	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/device"
	"github.com/bidon-io/bidon-backend/internal/segment"
)

// BuilderV2 is introduced for SDK version 0.5 and above, offering enhanced ad_units structure
type BuilderV2 struct {
	ConfigMatcher  ConfigMatcher
	AdUnitsMatcher AdUnitsMatcher
}

//go:generate go run -mod=mod github.com/matryer/moq@latest -out mocks/mocks.go -pkg mocks . ConfigMatcher AdUnitsMatcher

var ErrNoAdsFound = errors.New("no ads found")

type ConfigMatcher interface {
	Match(ctx context.Context, appID int64, adType ad.Type, segmentID int64) (*Config, error)
}

type AdUnitsMatcher interface {
	Match(ctx context.Context, params *BuildParams) ([]AdUnit, error)
}

type BuildParams struct {
	AppID      int64
	AdType     ad.Type
	AdFormat   ad.Format
	DeviceType device.Type
	Adapters   []adapter.Key
	Segment    segment.Segment
	PriceFloor *float64
}

func (b *BuilderV2) Build(ctx context.Context, params *BuildParams) (*Auction, error) {
	config, err := b.ConfigMatcher.Match(ctx, params.AppID, params.AdType, params.Segment.ID)
	if err != nil {
		return nil, err
	}

	adUnits, err := b.AdUnitsMatcher.Match(ctx, params)
	if err != nil {
		return nil, err
	}

	auction := Auction{
		ConfigID:                 config.ID,
		ConfigUID:                config.UID,
		ExternalWinNotifications: config.ExternalWinNotifications,
		Rounds:                   filterRounds(config.Rounds, params.Adapters),
		AdUnits:                  adUnits,
		LineItems:                []LineItem{},
		Segment:                  Segment{ID: params.Segment.StringID(), UID: params.Segment.UID},
	}

	if len(auction.Rounds) == 0 {
		return nil, ErrNoAdsFound
	}

	return &auction, nil
}

func filterRounds(rounds []RoundConfig, sdk_adapters []adapter.Key) []RoundConfig {
	filteredRounds := []RoundConfig{}

	for _, round := range rounds {
		demands := adapter.GetCommonAdapters(round.Demands, sdk_adapters)
		bidding := adapter.GetCommonAdapters(round.Bidding, sdk_adapters)

		if len(demands) == 0 && len(bidding) == 0 {
			continue // If both demands and bidding arrays empty => remove this round from Auction
		}

		round.Demands = demands
		round.Bidding = bidding
		filteredRounds = append(filteredRounds, round)
	}

	return filteredRounds
}
