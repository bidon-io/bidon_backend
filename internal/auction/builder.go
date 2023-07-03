package auction

import (
	"context"
	"errors"
	"strconv"

	"github.com/bidon-io/bidon-backend/internal/ad"
	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/device"
)

type Builder struct {
	ConfigMatcher    ConfigMatcher
	LineItemsMatcher LineItemsMatcher
}

//go:generate go run -mod=mod github.com/matryer/moq@latest -out mocks_test.go . ConfigMatcher LineItemsMatcher

var ErrNoAdsFound = errors.New("no ads found")

type ConfigMatcher interface {
	Match(ctx context.Context, appID int64, adType ad.Type, segmentID int64) (*Config, error)
}

type LineItemsMatcher interface {
	Match(ctx context.Context, params *BuildParams) ([]LineItem, error)
}

type BuildParams struct {
	AppID      int64
	AdType     ad.Type
	AdFormat   ad.Format
	DeviceType device.Type
	Adapters   []adapter.Key
	SegmentID  int64
}

func (b *Builder) Build(ctx context.Context, params *BuildParams) (*Auction, error) {
	config, err := b.ConfigMatcher.Match(ctx, params.AppID, params.AdType, params.SegmentID)
	if err != nil {
		return nil, err
	}

	lineItems, err := b.LineItemsMatcher.Match(ctx, params)
	if err != nil {
		return nil, err
	}

	var segmentID string
	if params.SegmentID != 0 {
		segmentID = strconv.Itoa(int(params.SegmentID))
	} else {
		segmentID = ""
	}

	auction := Auction{
		ConfigID:                 config.ID,
		ExternalWinNotifications: config.ExternalWinNotifications,
		Rounds:                   filterRounds(config.Rounds, params.Adapters),
		LineItems:                lineItems,
		Segment:                  Segment{ID: segmentID},
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
