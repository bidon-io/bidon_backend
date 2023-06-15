package auction

import (
	"context"
	"errors"

	"github.com/bidon-io/bidon-backend/internal/ad"
	"github.com/bidon-io/bidon-backend/internal/device"
	"golang.org/x/exp/slices"
)

type Builder struct {
	ConfigMatcher    ConfigMatcher
	LineItemsMatcher LineItemsMatcher
}

//go:generate go run -mod=mod github.com/matryer/moq@latest -out mocks_test.go . ConfigMatcher LineItemsMatcher

var ErrNoAdsFound = errors.New("no ads found")

type ConfigMatcher interface {
	Match(ctx context.Context, appID int64, adType ad.Type) (*Config, error)
}

type LineItemsMatcher interface {
	Match(ctx context.Context, params *BuildParams) ([]LineItem, error)
}

type BuildParams struct {
	AppID      int64
	AdType     ad.Type
	AdFormat   ad.Format
	DeviceType device.Type
	Adapters   []string
}

func (b *Builder) Build(ctx context.Context, params *BuildParams) (*Auction, error) {
	config, err := b.ConfigMatcher.Match(ctx, params.AppID, params.AdType)
	if err != nil {
		return nil, err
	}

	lineItems, err := b.LineItemsMatcher.Match(ctx, params)
	if err != nil {
		return nil, err
	}

	auction := Auction{
		ConfigID:  config.ID,
		Rounds:    filterRounds(config.Rounds, params.Adapters),
		LineItems: lineItems,
	}

	return &auction, nil
}

func filterRounds(rounds []RoundConfig, adapters []string) []RoundConfig {
	filteredRounds := []RoundConfig{}

	for _, round := range rounds {
		filteredDemands := []string{}
		for _, demand := range round.Demands {
			if slices.Contains(adapters, demand) {
				filteredDemands = append(filteredDemands, demand)
			}
		}
		if len(filteredDemands) == 0 {
			continue
		}

		round.Demands = filteredDemands
		filteredRounds = append(filteredRounds, round)
	}

	return filteredRounds
}
