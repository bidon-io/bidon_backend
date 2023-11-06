package auction

import (
	"context"
)

// Deprecated: Builder is deprecated as of SDK version 0.5
type Builder struct {
	ConfigMatcher    ConfigMatcher
	LineItemsMatcher LineItemsMatcher
}

//go:generate go run -mod=mod github.com/matryer/moq@latest -out mocks/mocks_deprecated.go -pkg mocks . LineItemsMatcher
type LineItemsMatcher interface {
	Match(ctx context.Context, params *BuildParams) ([]LineItem, error)
}

func (b *Builder) Build(ctx context.Context, params *BuildParams) (*Auction, error) {
	config, err := b.ConfigMatcher.Match(ctx, params.AppID, params.AdType, params.Segment.ID)
	if err != nil {
		return nil, err
	}

	lineItems, err := b.LineItemsMatcher.Match(ctx, params)
	if err != nil {
		return nil, err
	}

	auction := Auction{
		ConfigID:                 config.ID,
		ConfigUID:                config.UID,
		ExternalWinNotifications: config.ExternalWinNotifications,
		Rounds:                   filterRounds(config.Rounds, params.Adapters),
		LineItems:                lineItems,
		AdUnits:                  []AdUnit{},
		Segment:                  Segment{ID: params.Segment.StringID(), UID: params.Segment.UID},
	}

	if len(auction.Rounds) == 0 {
		return nil, ErrNoAdsFound
	}

	return &auction, nil
}
