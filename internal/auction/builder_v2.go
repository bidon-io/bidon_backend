package auction

import (
	"context"
	"errors"
	"math/big"

	"github.com/bidon-io/bidon-backend/internal/ad"
	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/device"
	"github.com/bidon-io/bidon-backend/internal/segment"
)

// BuilderV2 is introduced for SDK version 0.5 and above, offering enhanced ad_units structure
type BuilderV2 struct {
	ConfigFetcher  ConfigFetcher
	AdUnitsMatcher AdUnitsMatcher
}

//go:generate go run -mod=mod github.com/matryer/moq@latest -out mocks/mocks.go -pkg mocks . ConfigFetcher AdUnitsMatcher

var ErrNoAdsFound = errors.New("no ads found")
var ErrInvalidAuctionKey = errors.New("invalid auction_key")

type ConfigFetcher interface {
	Match(ctx context.Context, appID int64, adType ad.Type, segmentID int64, version string) (*Config, error)
	FetchByUIDCached(ctx context.Context, appID int64, id, uid string) *Config
}

type AdUnitsMatcher interface {
	MatchCached(ctx context.Context, params *BuildParams) ([]AdUnit, error)
}

type BuildParams struct {
	AppID      int64
	AdType     ad.Type
	AdFormat   ad.Format
	DeviceType device.Type
	Adapters   []adapter.Key
	Segment    segment.Segment
	PriceFloor *float64
	AuctionKey string
	AdUnitIDs  []int64
}

func (b *BuilderV2) Build(ctx context.Context, params *BuildParams) (*Auction, error) {
	var config *Config
	var err error
	if params.AuctionKey != "" {
		publicUID, success := new(big.Int).SetString(params.AuctionKey, 32)
		if !success {
			return nil, ErrInvalidAuctionKey
		}

		config = b.ConfigFetcher.FetchByUIDCached(ctx, params.AppID, "0", publicUID.String())
		if config == nil {
			return nil, ErrInvalidAuctionKey
		}
	} else {
		config, err = b.ConfigFetcher.Match(ctx, params.AppID, params.AdType, params.Segment.ID, "v1")
	}

	if err != nil {
		return nil, err
	}

	adUnits, err := b.AdUnitsMatcher.MatchCached(ctx, params)
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

func filterRounds(rounds []RoundConfig, sdkAdapters []adapter.Key) []RoundConfig {
	filteredRounds := []RoundConfig{}

	for _, round := range rounds {
		demands := adapter.GetCommonAdapters(round.Demands, sdkAdapters)
		bidding := adapter.GetCommonAdapters(round.Bidding, sdkAdapters)

		if len(demands) == 0 && len(bidding) == 0 {
			continue // If both demands and bidding arrays empty => remove this round from Auction
		}

		round.Demands = demands
		round.Bidding = bidding
		filteredRounds = append(filteredRounds, round)
	}

	return filteredRounds
}
