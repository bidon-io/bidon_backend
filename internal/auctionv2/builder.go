package auctionv2

import (
	"context"
	"errors"
	"time"

	"github.com/bidon-io/bidon-backend/internal/ad"
	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/auction"
	"github.com/bidon-io/bidon-backend/internal/bidding"
	"github.com/bidon-io/bidon-backend/internal/device"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/geocoder"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
	"github.com/bidon-io/bidon-backend/internal/segment"
)

type Builder struct {
	AdUnitsMatcher               AdUnitsMatcher
	BiddingBuilder               BiddingBuilder
	BiddingAdaptersConfigBuilder BiddingAdaptersConfigBuilder
}

//go:generate go run -mod=mod github.com/matryer/moq@latest -out mocks/mocks.go -pkg mocks . AdUnitsMatcher BiddingBuilder BiddingAdaptersConfigBuilder
type AdUnitsMatcher interface {
	MatchCached(ctx context.Context, params *auction.BuildParams) ([]auction.AdUnit, error)
}

type BiddingBuilder interface {
	HoldAuction(ctx context.Context, params *bidding.BuildParams) (bidding.AuctionResult, error)
}

type BiddingAdaptersConfigBuilder interface {
	Build(ctx context.Context, appID int64, adapterKeys []adapter.Key, adUnitsMap *auction.AdUnitsMap) (adapter.ProcessedConfigsMap, error)
}

type BuildParams struct {
	AppID                int64
	AdType               ad.Type
	AdFormat             ad.Format
	DeviceType           device.Type
	Adapters             []adapter.Key
	Segment              segment.Segment
	PriceFloor           float64
	MergedAuctionRequest *schema.AuctionV2Request
	GeoData              geocoder.GeoData
	AuctionKey           string
	AdUnitIds            []int64
	AuctionConfiguration *auction.Config
}

type AuctionResult struct {
	AuctionConfiguration *auction.Config
	CPMAdUnits           *[]auction.AdUnit
	AdUnits              *[]auction.AdUnit
	BiddingAuctionResult *bidding.AuctionResult
	Stat                 *Stat
}

func (a AuctionResult) GetDuration() int64 {
	if a.Stat != nil {
		return a.Stat.DurationTS
	}

	return 0
}

type Stat struct {
	StartTS    int64
	EndTS      int64
	DurationTS int64
}

func (b *Builder) Build(ctx context.Context, params *BuildParams) (*AuctionResult, error) {
	start := time.Now()

	if params.AuctionConfiguration == nil {
		return nil, auction.ErrNoAdsFound
	}

	demandAdapters := adapter.GetCommonAdapters(params.AuctionConfiguration.Demands, params.Adapters)
	biddingAdapters := adapter.GetCommonAdapters(params.AuctionConfiguration.Bidding, params.Adapters)
	if len(demandAdapters) == 0 && len(biddingAdapters) == 0 {
		return nil, auction.ErrNoAdsFound
	}
	if len(params.AuctionConfiguration.AdUnitIDs) == 0 {
		return nil, auction.ErrNoAdsFound
	}

	adUnits, err := b.AdUnitsMatcher.MatchCached(ctx, &auction.BuildParams{
		Adapters:   params.Adapters,
		AppID:      params.AppID,
		AdType:     params.AdType,
		AdFormat:   params.AdFormat,
		DeviceType: params.DeviceType,
		AdUnitIDs:  params.AuctionConfiguration.AdUnitIDs,
	})
	if err != nil {
		return nil, err
	}

	adUnitsMap := auction.BuildAdUnitsMap(&adUnits)
	var cpmAdUnits []auction.AdUnit
	for _, adUnit := range adUnits {
		// Use auction pricefloor as BM CPM price
		if adUnit.DemandID == string(adapter.BidmachineKey) && adUnit.IsCPM() {
			adUnit.PriceFloor = &params.PriceFloor
			cpmAdUnits = append(cpmAdUnits, adUnit)
			continue
		}

		if adUnit.GetPriceFloor() > params.PriceFloor && adUnit.IsCPM() {
			cpmAdUnits = append(cpmAdUnits, adUnit)
		}
	}

	adapterConfigs, err := b.BiddingAdaptersConfigBuilder.Build(ctx, params.AppID, params.Adapters, adUnitsMap)
	if err != nil {
		return nil, err
	}

	biddingRequest := params.MergedAuctionRequest.ToBiddingRequest()
	biddingAuctionResult, err := b.BiddingBuilder.HoldAuction(ctx, &bidding.BuildParams{
		AppID:          params.AppID,
		BiddingRequest: biddingRequest,
		GeoData:        params.GeoData,
		AdapterConfigs: adapterConfigs,
		AuctionConfig:  *params.AuctionConfiguration,
		StartTS:        start.UnixMilli(),
	})
	if err != nil && !errors.Is(err, bidding.ErrNoAdaptersMatched) {
		return nil, err
	}

	if len(cpmAdUnits) == 0 && len(biddingAuctionResult.Bids) == 0 {
		return nil, auction.ErrNoAdsFound
	}
	end := time.Now()

	// Build Result
	auctionResult := AuctionResult{
		AuctionConfiguration: params.AuctionConfiguration,
		AdUnits:              &adUnits,
		CPMAdUnits:           &cpmAdUnits,
		BiddingAuctionResult: &biddingAuctionResult,
		Stat: &Stat{
			StartTS:    start.UnixMilli(),
			EndTS:      end.UnixMilli(),
			DurationTS: end.Sub(start).Microseconds(),
		},
	}

	return &auctionResult, nil
}
