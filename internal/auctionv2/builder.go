package auctionv2

import (
	"context"
	"errors"
	"math/big"
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
	ConfigFetcher                ConfigFetcher
	AdUnitsMatcher               AdUnitsMatcher
	BiddingBuilder               BiddingBuilder
	BiddingAdaptersConfigBuilder BiddingAdaptersConfigBuilder
}

//go:generate go run -mod=mod github.com/matryer/moq@latest -out mocks/mocks.go -pkg mocks . ConfigFetcher AdUnitsMatcher BiddingBuilder BiddingAdaptersConfigBuilder

type ConfigFetcher interface {
	Match(ctx context.Context, appID int64, adType ad.Type, segmentID int64, version string) (*auction.Config, error)
	FetchByUIDCached(ctx context.Context, appId int64, id, uid string) *auction.Config
}

type AdUnitsMatcher interface {
	MatchCached(ctx context.Context, params *auction.BuildParams) ([]auction.AdUnit, error)
}

type BiddingBuilder interface {
	HoldAuction(ctx context.Context, params *bidding.BuildParams) (bidding.AuctionResult, error)
}

type BiddingAdaptersConfigBuilder interface {
	Build(ctx context.Context, appID int64, adapterKeys []adapter.Key, imp schema.Imp, adUnitsMap *map[adapter.Key][]auction.AdUnit) (adapter.ProcessedConfigsMap, error)
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
}

type AuctionResult struct {
	AuctionConfiguration *auction.Config
	CPMAdUnits           *[]auction.AdUnit
	AdUnits              *[]auction.AdUnit
	BiddingAuctionResult *bidding.AuctionResult
	Stat                 *Stat
}

type Stat struct {
	StartTS    int64
	EndTS      int64
	DurationTS int64
}

func (b *Builder) Build(ctx context.Context, params *BuildParams) (*AuctionResult, error) {
	start := time.Now()

	var auctionConfig *auction.Config
	var err error

	// Fetch Auction
	if params.AuctionKey != "" {
		publicUid, success := new(big.Int).SetString(params.AuctionKey, 32)
		if !success {
			return nil, auction.InvalidAuctionKey
		}

		auctionConfig = b.ConfigFetcher.FetchByUIDCached(ctx, params.AppID, "0", publicUid.String())
		if auctionConfig == nil {
			return nil, auction.InvalidAuctionKey
		}
	} else {
		auctionConfig, err = b.ConfigFetcher.Match(ctx, params.AppID, params.AdType, params.Segment.ID, "v2")
	}

	if err != nil {
		return nil, err
	}

	demandAdapters := adapter.GetCommonAdapters(auctionConfig.Demands, params.Adapters)
	biddingAdapters := adapter.GetCommonAdapters(auctionConfig.Bidding, params.Adapters)
	if len(demandAdapters) == 0 && len(biddingAdapters) == 0 {
		return nil, auction.ErrNoAdsFound
	}
	if len(auctionConfig.AdUnitIDs) == 0 {
		return nil, auction.ErrNoAdsFound
	}

	adUnits, err := b.AdUnitsMatcher.MatchCached(ctx, &auction.BuildParams{
		Adapters:   params.Adapters,
		AppID:      params.AppID,
		AdType:     params.AdType,
		AdFormat:   params.AdFormat,
		DeviceType: params.DeviceType,
		AdUnitIDs:  auctionConfig.AdUnitIDs,
	})
	if err != nil {
		return nil, err
	}

	adUnitsMap := make(map[adapter.Key][]auction.AdUnit)
	for _, adUnit := range adUnits {
		key := adapter.Key(adUnit.DemandID)
		adUnitsMap[key] = append(adUnitsMap[key], adUnit)
	}
	var cpmAdUnits []auction.AdUnit
	for _, adUnit := range adUnits {
		if adUnit.GetPriceFloor() > params.PriceFloor && adUnit.IsCPM() {
			cpmAdUnits = append(cpmAdUnits, adUnit)
		}
	}

	// Bidding
	params.MergedAuctionRequest.AdObject.AuctionConfigurationID = auctionConfig.ID
	params.MergedAuctionRequest.AdObject.AuctionConfigurationUID = auctionConfig.UID
	imp := params.MergedAuctionRequest.AdObject.ToImp()

	adapterConfigs, err := b.BiddingAdaptersConfigBuilder.Build(ctx, params.AppID, params.Adapters, imp, &adUnitsMap)
	if err != nil {
		return nil, err
	}

	biddingRequest := params.MergedAuctionRequest.ToBiddingRequest()
	biddingAuctionResult, err := b.BiddingBuilder.HoldAuction(ctx, &bidding.BuildParams{
		AppID:          params.AppID,
		BiddingRequest: biddingRequest,
		GeoData:        params.GeoData,
		AdapterConfigs: adapterConfigs,
		AuctionConfig:  *auctionConfig,
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
		AuctionConfiguration: auctionConfig,
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
