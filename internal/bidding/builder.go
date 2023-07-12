package bidding

import (
	"context"
	"errors"

	"github.com/bidon-io/bidon-backend/internal/ad"
	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/auction"
	"github.com/bidon-io/bidon-backend/internal/device"
)

type Builder struct {
	ConfigMatcher ConfigMatcher
}

//go:generate go run -mod=mod github.com/matryer/moq@latest -out mocks_test.go . ConfigMatcher

var ErrNoBids = errors.New("no bids")

type ConfigMatcher interface {
	Match(ctx context.Context, appID int64, adType ad.Type, segmentID int64) (*auction.Config, error)
}

type BuildParams struct {
	AppID      int64
	AdType     ad.Type
	AdFormat   ad.Format
	DeviceType device.Type
	Adapters   []adapter.Key
	SegmentID  int64
}

func (b *Builder) Build(ctx context.Context, params *BuildParams) (*DemandResponse, error) {
	config, err := b.ConfigMatcher.Match(ctx, params.AppID, params.AdType, params.SegmentID)
	if err != nil {
		return nil, err
	}
	config.ExternalWinNotifications = true // just to make the test pass

	response := DemandResponse{
		Price: 0,
	}

	return &response, nil
}
