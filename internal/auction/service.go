package auction

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
	ConfigID                 int64    `json:"auction_configuration_id"`
	ConfigUID                string   `json:"auction_configuration_uid"`
	ExternalWinNotifications bool     `json:"external_win_notifications"`
	AdUnits                  []AdUnit `json:"ad_units"`
	NoBids                   []AdUnit `json:"no_bids"`
	Segment                  Segment  `json:"segment"`
	Token                    string   `json:"token"`
	AuctionPriceFloor        float64  `json:"auction_pricefloor"`
	AuctionTimeout           int      `json:"auction_timeout"`
	AuctionID                string   `json:"auction_id"`
}

type ExecutionParams struct {
	Req     *schema.AuctionRequest
	AppID   int64
	Country string
	GeoData geocoder.GeoData
	Log     func(string)
	LogErr  func(err error)
}

//go:generate go run -mod=mod github.com/matryer/moq@latest -out mocks/service_mocks.go -pkg mocks . ConfigFetcher AuctionBuilder AdapterKeysFetcher

type ConfigFetcher interface {
	Match(ctx context.Context, appID int64, adType ad.Type, segmentID int64, version string) (*Config, error)
	FetchByUIDCached(ctx context.Context, appID int64, id, uid string) *Config
}

type AdapterKeysFetcher interface {
	FetchEnabledAdapterKeys(ctx context.Context, appID int64, adapterKeys []adapter.Key) ([]adapter.Key, error)
}

type AuctionBuilder interface { //nolint:revive
	Build(ctx context.Context, params *BuildParams) (*Result, error)
}

const (
	DefaultAuctionTimeout = 30000
)

var adCacheAdaptersFilter = store.NewAdCacheAdaptersFilter()

func (s *Service) Run(ctx context.Context, params *ExecutionParams) (*Response, error) {
	req := params.Req

	var auctionConfig *Config
	var auctionResult *Result
	var adUnitsMap *AdUnitsMap
	var err error

	// Ensure events are always logged, even on errors
	defer func() {
		s.logEvents(req, params, auctionConfig, auctionResult, adUnitsMap, err)
	}()

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
			err = sdkapi.ErrInvalidAuctionKey
			return nil, err
		}

		auctionConfig = s.ConfigFetcher.FetchByUIDCached(ctx, params.AppID, "0", publicUID.String())
		if auctionConfig == nil {
			err = sdkapi.ErrInvalidAuctionKey
			return nil, err
		}
	} else {
		auctionConfig, err = s.ConfigFetcher.Match(ctx, params.AppID, req.AdType, sgmnt.ID, "v2")
	}
	if err != nil {
		err = sdkapi.ErrNoAdsFound
		return nil, err
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
		AuctionRequest:       req,
		GeoData:              params.GeoData,
		AuctionKey:           req.AdObject.AuctionKey,
		AuctionConfiguration: auctionConfig,
	}

	auctionResult, err = s.AuctionBuilder.Build(ctx, bp)
	if err != nil {
		if errors.Is(err, ErrNoAdsFound) {
			err = sdkapi.ErrNoAdsFound
		}
		return nil, err
	}
	params.Log(fmt.Sprintf("[AUCTION V2] auction: (%+v), err: (%s), took (%dms)", auctionResult, err, auctionResult.GetDuration()))

	adUnitsMap = buildAdUnitsMap(auctionResult.AdUnits)

	return s.buildResponse(req, auctionResult, adUnitsMap)
}

var customAdapters = [...]string{"max", "level_play"}

// Auction keys that should have disabled price floors when using custom adapters
var disabledFloorAuctionKeys = [...]string{
	"1LOQ1LROG0000", // Inter
	"1LOQ2BFG00000", // Banner
	"1LOQ2KLES0400", // Rewarded
	"1LPGM6QBC0000", // Banner
	"1LPGMABLK0000", // Interstitial
	"1LRBRTLCS0400", // Interstitial
	"1LVHHK37O0400", // Banner
	"1LVHHK51K0400", // Interstitial
	"1LVHEVUJO0400", // Interstitial
	"1LQP712M00000", // Banner
	"1LQP75UU00400", // Interstitial
	"1LQP7A4UG0400", // Rewarded
}

func priceFloor(req *schema.AuctionRequest, auctionConfig *Config) float64 {
	// Check if price floor should be disabled for this auction key with custom adapter
	isCustomAdapter := slices.Contains(customAdapters[:], req.GetMediator())
	if isCustomAdapter && slices.Contains(disabledFloorAuctionKeys[:], req.AdObject.AuctionKey) {
		return 0 // Disable price floor for specified auction keys with custom adapters
	}

	// Default floor logic
	priceFloor := req.AdObject.PriceFloor
	for _, cacheObject := range req.AdCache {
		priceFloor = math.Max(priceFloor, cacheObject.Price)
	}
	priceFloor = math.Max(auctionConfig.PriceFloor, priceFloor)

	// Custom Adapter floor logic
	// Check if previous auction price is higher than the current price floor
	prevFloor := req.GetPrevAuctionPrice()
	if prevFloor != nil && isCustomAdapter {
		priceFloor = math.Max(*prevFloor, priceFloor)
	}

	return priceFloor
}

func (s *Service) buildResponse(
	req *schema.AuctionRequest,
	auctionResult *Result,
	adUnitsMap *AdUnitsMap,
) (*Response, error) {
	adObject := req.AdObject
	response := Response{
		ConfigID:                 auctionResult.AuctionConfiguration.ID,
		ConfigUID:                auctionResult.AuctionConfiguration.UID,
		Segment:                  Segment{ID: req.Segment.ID, UID: req.Segment.UID},
		Token:                    "{}",
		AuctionID:                adObject.AuctionID,
		AuctionPriceFloor:        adObject.PriceFloor,
		AuctionTimeout:           auctionTimeout(auctionResult.AuctionConfiguration),
		ExternalWinNotifications: auctionResult.AuctionConfiguration.ExternalWinNotifications,
		AdUnits:                  make([]AdUnit, 0),
		NoBids:                   make([]AdUnit, 0),
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
		if adUnit.DemandID == string(adapter.BidmachineKey) && req.GetMediator() != "" {
			adUnit.Extra["custom_parameters"] = map[string]any{
				"mediator": req.GetMediator(),
			}
		}

		response.AdUnits = append(response.AdUnits, adUnit)
	}

	// Store Bids AS RTB AdUnits from BiddingAuctionResult
	for _, bidResponse := range auctionResult.BiddingAuctionResult.Bids {
		adUnit := convertBidToAdUnit(req, bidResponse, adUnitsMap)
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
	req *schema.AuctionRequest,
	params *ExecutionParams,
	auctionConfig *Config,
	auctionResult *Result,
	adUnitsMap *AdUnitsMap,
	auctionErr error,
) {
	// Prepare auction info from available data
	var auc *Auction
	var auctionConfigurationUID int

	if auctionResult != nil && auctionResult.AuctionConfiguration != nil {
		// Use auction result configuration (success case)
		auc = &Auction{
			ConfigID:  auctionResult.AuctionConfiguration.ID,
			ConfigUID: auctionResult.AuctionConfiguration.UID,
		}
	} else if auctionConfig != nil {
		// Use provided configuration (error case with config available)
		auc = &Auction{
			ConfigID:  auctionConfig.ID,
			ConfigUID: auctionConfig.UID,
		}
	} else {
		// Create minimal auction info for early errors
		auc = &Auction{
			ConfigID:  0,
			ConfigUID: "0",
		}
	}

	uid, err := strconv.Atoi(auc.ConfigUID)
	if err != nil {
		auctionConfigurationUID = 0
	} else {
		auctionConfigurationUID = uid
	}

	// Prepare events
	var events []*event.AdEvent

	// Add bidding events if available
	if auctionResult != nil && auctionResult.BiddingAuctionResult != nil && adUnitsMap != nil {
		events = prepareBiddingEvents(req, params, auctionResult.BiddingAuctionResult, adUnitsMap)
	}

	// Add auction request event
	aucRequestEvent := prepareAuctionRequestEvent(req, params, auc, auctionConfigurationUID, auctionErr)
	events = append(events, aucRequestEvent)

	// Log all events
	for _, ev := range events {
		s.EventLogger.Log(ev, func(err error) {
			params.LogErr(fmt.Errorf("log %v event: %v", ev.EventType, err))
		})
	}
}

func convertBidToAdUnit(req *schema.AuctionRequest, demandResponse adapters.DemandResponse, adUnitsMap *AdUnitsMap) *AdUnit {
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
		ext = buildDemandExt(req, demandResponse)
	}

	for key, value := range storeAdUnit.Extra {
		ext[key] = value
	}

	return &AdUnit{
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
	req *schema.AuctionRequest,
	params *ExecutionParams,
	auc *Auction,
	auctionConfigurationUID int,
	auctionErr error,
) *event.AdEvent {
	status := event.SuccessAdRequestStatus
	errorMsg := ""

	if auctionErr != nil {
		status = event.ErrorAdRequestStatus
		errorMsg = auctionErr.Error()
	}

	adRequestParams := event.AdRequestParams{
		EventType:               "auction_request",
		AdType:                  string(req.AdType),
		AdFormat:                string(req.AdObject.Format()),
		AuctionID:               req.AdObject.AuctionID,
		AuctionConfigurationID:  auc.ConfigID,
		AuctionConfigurationUID: int64(auctionConfigurationUID),
		Status:                  status,
		ImpID:                   "",
		DemandID:                "",
		AdUnitUID:               0,
		AdUnitLabel:             "",
		ECPM:                    0,
		PriceFloor:              req.AdObject.PriceFloor,
		Error:                   errorMsg,
	}

	return event.NewAdEvent(&req.BaseRequest, adRequestParams, params.GeoData)
}

func prepareBiddingEvents(
	req *schema.AuctionRequest,
	params *ExecutionParams,
	auctionResult *bidding.AuctionResult,
	adUnitsMap *AdUnitsMap,
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

func selectAdUnit(demandResponse adapters.DemandResponse, adUnitsMap *AdUnitsMap) (*AdUnit, error) {
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

func buildDemandExt(req *schema.AuctionRequest, demandResponse adapters.DemandResponse) map[string]any {
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
	case adapter.BidmachineKey:
		extra := map[string]any{
			"payload": demandResponse.Bid.Payload,
		}
		if req.GetMediator() != "" {
			extra["custom_parameters"] = map[string]any{
				"mediator": req.GetMediator(),
			}
		}
		return extra
	default:
		return map[string]any{
			"payload": demandResponse.Bid.Payload,
		}
	}
}

func auctionTimeout(conf *Config) int {
	if conf.Timeout > 0 {
		return conf.Timeout
	}

	return DefaultAuctionTimeout
}
