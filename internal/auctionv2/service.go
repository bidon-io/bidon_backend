package auctionv2

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/big"
	"slices"
	"sort"
	"strconv"

	"github.com/bidon-io/bidon-backend/internal/ad"
	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/auction"
	"github.com/bidon-io/bidon-backend/internal/bidding"
	"github.com/bidon-io/bidon-backend/internal/bidding/adapters"
	"github.com/bidon-io/bidon-backend/internal/sdkapi"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/event"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/geocoder"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/store"
	"github.com/bidon-io/bidon-backend/internal/segment"
)

type Service struct {
	ConfigFetcher      ConfigFetcher
	AuctionBuilder     AuctionBuilder
	SegmentMatcher     *segment.Matcher
	AdapterKeysFetcher AdapterKeysFetcher
	EventLogger        *event.Logger
}

type Response struct {
	ConfigID                 int64            `json:"auction_configuration_id"`
	ConfigUID                string           `json:"auction_configuration_uid"`
	ExternalWinNotifications bool             `json:"external_win_notifications"`
	AdUnits                  []auction.AdUnit `json:"ad_units"`
	NoBids                   []auction.AdUnit `json:"no_bids"`
	Segment                  auction.Segment  `json:"segment"`
	Token                    string           `json:"token"`
	AuctionPriceFloor        float64          `json:"auction_pricefloor"`
	AuctionTimeout           int              `json:"auction_timeout"`
	AuctionID                string           `json:"auction_id"`
}

type ExecutionParams struct {
	Req     *schema.AuctionV2Request
	AppID   int64
	Country string
	GeoData geocoder.GeoData
	Log     func(string)
	LogErr  func(err error)
}

//go:generate go run -mod=mod github.com/matryer/moq@latest -out mocks/service_mocks.go -pkg mocks . ConfigFetcher AuctionBuilder AdapterKeysFetcher

type ConfigFetcher interface {
	Match(ctx context.Context, appID int64, adType ad.Type, segmentID int64, version string) (*auction.Config, error)
	FetchByUIDCached(ctx context.Context, appId int64, id, uid string) *auction.Config
}

type AdapterKeysFetcher interface {
	FetchEnabledAdapterKeys(ctx context.Context, appID int64, adapterKeys []adapter.Key) ([]adapter.Key, error)
}

type AuctionBuilder interface {
	Build(ctx context.Context, params *BuildParams) (*AuctionResult, error)
}

const (
	DefaultAuctionTimeout = 30000
)

var adCacheAdaptersFilter = store.NewAdCacheAdaptersFilter()

func (s *Service) Run(ctx context.Context, params *ExecutionParams) (*Response, error) {
	req := params.Req

	var auctionConfig *auction.Config
	var err error

	segmentParams := &segment.Params{
		Country: params.Country,
		Ext:     req.Segment.Ext,
		AppID:   params.AppID,
	}

	sgmnt := s.SegmentMatcher.Match(ctx, segmentParams)
	req.Segment.ID = sgmnt.StringID()
	req.Segment.UID = sgmnt.UID

	adapterKeys := adCacheAdaptersFilter.Filter(
		ad.OS(req.Device.OS),
		req.AdType,
		req.Adapters.Keys(),
		req.AdCache,
	)

	adapterKeys, err = s.AdapterKeysFetcher.FetchEnabledAdapterKeys(ctx, params.AppID, adapterKeys)
	if err != nil {
		return nil, err
	}

	if req.AdObject.AuctionKey != "" {
		publicUID, success := new(big.Int).SetString(req.AdObject.AuctionKey, 32)
		if !success {
			return nil, sdkapi.ErrInvalidAuctionKey
		}

		auctionConfig = s.ConfigFetcher.FetchByUIDCached(ctx, params.AppID, "0", publicUID.String())
		if auctionConfig == nil {
			return nil, sdkapi.ErrInvalidAuctionKey
		}
	} else {
		auctionConfig, err = s.ConfigFetcher.Match(ctx, params.AppID, req.AdType, sgmnt.ID, "v2")
	}
	if err != nil {
		return nil, sdkapi.ErrNoAdsFound
	}
	req.AdObject.AuctionConfigurationID = auctionConfig.ID
	req.AdObject.AuctionConfigurationUID = auctionConfig.UID
	req.AdObject.PriceFloor = priceFloor(req, auctionConfig)

	bp := &BuildParams{
		AppID:                params.AppID,
		AdType:               req.AdType,
		AdFormat:             req.AdObject.Format(),
		DeviceType:           req.Device.Type,
		Adapters:             adapterKeys,
		Segment:              sgmnt,
		PriceFloor:           req.AdObject.PriceFloor,
		MergedAuctionRequest: req,
		GeoData:              params.GeoData,
		AuctionKey:           req.AdObject.AuctionKey,
		AuctionConfiguration: auctionConfig,
	}

	auctionResult, err := s.AuctionBuilder.Build(ctx, bp)
	if err != nil {
		if errors.Is(err, auction.ErrNoAdsFound) {
			err = sdkapi.ErrNoAdsFound
		}

		return nil, err
	}
	params.Log(fmt.Sprintf("[AUCTION V2] auction: (%+v), err: (%s), took (%dms)", auctionResult, err, auctionResult.GetDuration()))

	adUnitsMap := auction.BuildAdUnitsMap(auctionResult.AdUnits)

	s.logEvents(req, params, auctionResult, adUnitsMap)

	return s.buildResponse(req, auctionResult, adUnitsMap)
}

var customAdapters = [...]string{"max", "level_play"}

func priceFloor(req *schema.AuctionV2Request, auctionConfig *auction.Config) float64 {
	// Default floor logic
	priceFloor := req.AdObject.PriceFloor
	for _, cacheObject := range req.AdCache {
		priceFloor = math.Max(priceFloor, cacheObject.Price)
	}
	priceFloor = math.Max(auctionConfig.PriceFloor, priceFloor)

	// Custom Adapter floor logic
	// Check if previous auction price is higher than the current price floor
	isCustomAdapter := slices.Contains(customAdapters[:], req.GetMediator())
	prevFloor := req.GetPrevAuctionPrice()
	if prevFloor != nil && isCustomAdapter {
		priceFloor = math.Max(*prevFloor, priceFloor)
	}

	return priceFloor
}

func (s *Service) buildResponse(
	req *schema.AuctionV2Request,
	auctionResult *AuctionResult,
	adUnitsMap *auction.AdUnitsMap,
) (*Response, error) {
	adObject := req.AdObject
	response := Response{
		ConfigID:                 auctionResult.AuctionConfiguration.ID,
		ConfigUID:                auctionResult.AuctionConfiguration.UID,
		Segment:                  auction.Segment{ID: req.Segment.ID, UID: req.Segment.UID},
		Token:                    "{}",
		AuctionID:                adObject.AuctionID,
		AuctionPriceFloor:        adObject.PriceFloor,
		AuctionTimeout:           auctionTimeout(auctionResult.AuctionConfiguration),
		ExternalWinNotifications: auctionResult.AuctionConfiguration.ExternalWinNotifications,
		AdUnits:                  make([]auction.AdUnit, 0),
		NoBids:                   make([]auction.AdUnit, 0),
	}

	isCOPPA := false
	if req.Regulations != nil {
		isCOPPA = req.Regulations.COPPA
	}

	// Store CPM AdUnits from AuctionConfiguration
	for _, adUnit := range *auctionResult.CPMAdUnits {
		if isCOPPA && adapter.IsDisabledForCOPPA(adapter.Key(adUnit.DemandID)) {
			continue
		}
		response.AdUnits = append(response.AdUnits, adUnit)
	}

	// Store Bids AS RTB AdUnits from BiddingAuctionResult
	for _, bidResponse := range auctionResult.BiddingAuctionResult.Bids {
		adUnit := convertBidToAdUnit(bidResponse, adUnitsMap)
		if adUnit == nil {
			continue
		}

		if bidResponse.IsBid() && bidResponse.Price() > adObject.PriceFloor {
			response.AdUnits = append(response.AdUnits, *adUnit)
		} else {
			response.NoBids = append(response.NoBids, *adUnit)
		}
	}

	// Sort AdUnits by price
	sort.Slice(response.AdUnits, func(i, j int) bool {
		return response.AdUnits[i].GetPriceFloor() > response.AdUnits[j].GetPriceFloor()
	})

	return &response, nil
}

func (s *Service) logEvents(
	req *schema.AuctionV2Request,
	params *ExecutionParams,
	auctionResult *AuctionResult,
	adUnitsMap *auction.AdUnitsMap,
) {
	auc := &auction.Auction{
		ConfigID:  auctionResult.AuctionConfiguration.ID,
		ConfigUID: auctionResult.AuctionConfiguration.UID,
	}
	auctionConfigurationUID, err := strconv.Atoi(auc.ConfigUID)
	if err != nil {
		auctionConfigurationUID = 0
	}

	events := prepareBiddingEvents(req, params, auctionResult.BiddingAuctionResult, adUnitsMap)
	aucRequestEvent := prepareAuctionRequestEvent(req, params, auc, auctionConfigurationUID)

	events = append(events, aucRequestEvent)
	for _, ev := range events {
		s.EventLogger.Log(ev, func(err error) {
			params.LogErr(fmt.Errorf("log %v event: %v", ev.EventType, err))
		})
	}
}

func convertBidToAdUnit(demandResponse adapters.DemandResponse, adUnitsMap *auction.AdUnitsMap) *auction.AdUnit {
	storeAdUnit, err := selectAdUnit(demandResponse, adUnitsMap)
	if err != nil {
		return nil
	}
	if storeAdUnit == nil {
		return nil
	}

	priceFloor := demandResponse.Price()
	ext := map[string]any{}
	if demandResponse.IsBid() {
		ext = buildDemandExt(demandResponse)
	}

	for key, value := range storeAdUnit.Extra {
		ext[key] = value
	}

	return &auction.AdUnit{
		DemandID:   string(demandResponse.DemandID),
		UID:        storeAdUnit.UID,
		Label:      storeAdUnit.Label,
		BidType:    schema.RTBBidType,
		PriceFloor: &priceFloor,
		Timeout:    storeAdUnit.Timeout,
		Extra:      ext,
	}
}

func prepareAuctionRequestEvent(
	req *schema.AuctionV2Request,
	params *ExecutionParams,
	auc *auction.Auction,
	auctionConfigurationUID int,
) *event.AdEvent {
	adRequestParams := event.AdRequestParams{
		EventType:               "auction_request",
		AdType:                  string(req.AdType),
		AdFormat:                string(req.AdObject.Format()),
		AuctionID:               req.AdObject.AuctionID,
		AuctionConfigurationID:  auc.ConfigID,
		AuctionConfigurationUID: int64(auctionConfigurationUID),
		Status:                  "",
		ImpID:                   "",
		DemandID:                "",
		AdUnitUID:               0,
		AdUnitLabel:             "",
		ECPM:                    0,
		PriceFloor:              req.AdObject.PriceFloor,
	}

	return event.NewAdEvent(&req.BaseRequest, adRequestParams, params.GeoData)
}

func prepareBiddingEvents(
	req *schema.AuctionV2Request,
	params *ExecutionParams,
	auctionResult *bidding.AuctionResult,
	adUnitsMap *auction.AdUnitsMap,
) []*event.AdEvent {
	adObject := req.AdObject
	auctionConfigurationUID, err := strconv.Atoi(adObject.AuctionConfigurationUID)
	if err != nil {
		auctionConfigurationUID = 0
	}

	events := make([]*event.AdEvent, 0, len(auctionResult.Bids))
	for _, result := range auctionResult.Bids {
		adUnit, _ := selectAdUnit(result, adUnitsMap)
		adUnitUID := int64(0)
		adUnitLabel := ""
		if adUnit != nil {
			uid, _ := strconv.ParseInt(adUnit.UID, 10, 64)
			adUnitUID = uid
			adUnitLabel = adUnit.Label
		}

		adRequestParams := event.AdRequestParams{
			EventType:               "bid_request",
			AdType:                  string(req.AdType),
			AdFormat:                string(req.AdObject.Format()),
			AuctionID:               adObject.AuctionID,
			AuctionConfigurationID:  adObject.AuctionConfigurationID,
			AuctionConfigurationUID: int64(auctionConfigurationUID),
			Status:                  fmt.Sprint(result.Status),
			ImpID:                   "",
			DemandID:                string(result.DemandID),
			AdUnitUID:               adUnitUID,
			AdUnitLabel:             adUnitLabel,
			ECPM:                    result.Price(),
			PriceFloor:              adObject.PriceFloor,
			Bidding:                 true,
			RawRequest:              result.RawRequest,
			RawResponse:             result.RawResponse,
			Error:                   result.ErrorMessage(),
			TimingMap: event.TimingMap{
				"bid":   {result.StartTS, result.EndTS},
				"token": {result.Token.StartTS, result.Token.EndTS},
			},
		}
		events = append(events, event.NewAdEvent(&req.BaseRequest, adRequestParams, params.GeoData))
		if result.IsBid() {
			adRequestParams = event.AdRequestParams{
				EventType:               "bid",
				AdType:                  string(req.AdType),
				AdFormat:                string(adObject.Format()),
				AuctionID:               adObject.AuctionID,
				AuctionConfigurationID:  adObject.AuctionConfigurationID,
				AuctionConfigurationUID: int64(auctionConfigurationUID),
				Status:                  "SUCCESS",
				ImpID:                   "",
				DemandID:                string(result.DemandID),
				AdUnitUID:               adUnitUID,
				AdUnitLabel:             adUnitLabel,
				ECPM:                    result.Bid.Price,
				PriceFloor:              adObject.PriceFloor,
				Bidding:                 true,
				TimingMap: event.TimingMap{
					"bid": {result.StartTS, result.EndTS},
				},
			}
			events = append(events, event.NewAdEvent(&req.BaseRequest, adRequestParams, params.GeoData))
		}
	}

	return events
}

func selectAdUnit(demandResponse adapters.DemandResponse, adUnitsMap *auction.AdUnitsMap) (*auction.AdUnit, error) {
	adUnits, err := adUnitsMap.All(demandResponse.DemandID, schema.RTBBidType)
	if err != nil {
		return nil, err
	}

	if demandResponse.DemandID == adapter.AmazonKey {
		for _, adUnit := range adUnits {
			if demandResponse.SlotUUID == adUnit.Extra["slot_uuid"] {
				return &adUnit, nil
			}
		}
	} else if len(adUnits) > 0 {
		adUnit := adUnits[0]
		return &adUnit, nil
	}

	return nil, fmt.Errorf("ad unit not found for demand %s", demandResponse.DemandID)
}

func buildDemandExt(demandResponse adapters.DemandResponse) map[string]any {
	switch demandResponse.DemandID {
	case adapter.AmazonKey:
		return map[string]any{}
	case adapter.MobileFuseKey:
		return map[string]any{
			"signaldata": demandResponse.Bid.Signaldata,
		}
	case adapter.VKAdsKey:
		return map[string]any{
			"bid_id": demandResponse.Bid.ID,
		}
	default:
		return map[string]any{
			"payload": demandResponse.Bid.Payload,
		}
	}
}

func auctionTimeout(conf *auction.Config) int {
	if conf.Timeout > 0 {
		return conf.Timeout
	}

	return DefaultAuctionTimeout
}
