package bidding

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/bidon-io/bidon-backend/internal/bidding/adapters/amazon"

	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/auction"
	"github.com/bidon-io/bidon-backend/internal/bidding/adapters"
	"github.com/bidon-io/bidon-backend/internal/bidding/openrtb"
	"github.com/bidon-io/bidon-backend/internal/device"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/geocoder"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
	"github.com/gofrs/uuid/v5"
	"github.com/prebid/openrtb/v19/adcom1"
	"github.com/prebid/openrtb/v19/openrtb2"
	"golang.org/x/exp/maps"
)

type Builder struct {
	AdaptersBuilder     AdaptersBuilder
	NotificationHandler NotificationHandler
}

var ErrNoAdaptersMatched = errors.New("no adapters matched")

//go:generate go run -mod=mod github.com/matryer/moq@latest -out mocks/mocks.go -pkg mocks . AdaptersBuilder NotificationHandler

type AdaptersBuilder interface {
	Build(adapterKey adapter.Key, cfg adapter.ProcessedConfigsMap) (*adapters.Bidder, error)
}

type NotificationHandler interface {
	HandleBiddingRound(context.Context, *schema.Imp, AuctionResult, string, string) error
}

type BuildParams struct {
	AppID          int64
	BiddingRequest schema.BiddingRequest
	GeoData        geocoder.GeoData
	AdapterConfigs adapter.ProcessedConfigsMap
	AuctionConfig  auction.Config
	StartTS        int64
}

type AuctionResult struct {
	Bids        []adapters.DemandResponse
	RoundNumber int
}

type AmazonSlot struct {
	SlotUUID   string `json:"slot_uuid"`
	PricePoint string `json:"price_point"`
}

func (b *Builder) HoldAuction(ctx context.Context, params *BuildParams) (AuctionResult, error) {
	// get config
	// build openrtb request
	// filter adapters
	// split to bids
	// build requests and send them to adapters in parallel
	// collect results
	// build response
	emptyResponse := AuctionResult{}
	br := params.BiddingRequest
	config := params.AuctionConfig

	bidId, err := uuid.NewV4()
	if err != nil {
		return emptyResponse, fmt.Errorf("cannot generate Bid UUID: %s", err)
	}
	baseBidRequest := openrtb.BidRequest{
		ID:   bidId.String(),
		Test: *bool2int(br.Test),
		AT:   1,
		TMax: 2000,
		App: &openrtb2.App{
			Ver:    br.App.Version,
			Bundle: br.App.Bundle,
			ID:     strconv.FormatInt(params.AppID, 10),
			Publisher: &openrtb2.Publisher{
				ID: "SELLER_ID",
			},
		},
		Device: b.BuildDevice(br.Device, br.User, params.GeoData),
		Imp:    []openrtb2.Imp{},
		Regs: &openrtb2.Regs{
			COPPA: *bool2int(br.GetRegulations().COPPA),
			GDPR:  bool2int(br.GetRegulations().GDPR),
		},
	}

	var adapterKeys []adapter.Key
	roundNumber := 0
	isV2Auction := len(config.Rounds) == 0 && (len(config.Bidding) > 0 || len(config.Demands) > 0)
	if isV2Auction {
		filteredDemands := make(map[adapter.Key]map[string]any)
		for key, value := range br.Imp.Demands {
			if token, ok := value["token"]; ok && token != "" {
				filteredDemands[key] = value
			}
		}

		adapterKeys = adapter.GetCommonAdapters(config.Bidding, br.Adapters.Keys())
		adapterKeys = adapter.GetCommonAdapters(adapterKeys, maps.Keys(filteredDemands))
	} else {
		// DEPRECATED: Remove after migration to V2
		// Get adapters from request, demands from bidding request and demands from round config and merge them
		var roundConfig *auction.RoundConfig
		for idx, round := range config.Rounds {
			if round.ID == br.Imp.RoundID {
				roundConfig = &round
				roundNumber = idx
				break
			}
		}
		if roundConfig == nil {
			return emptyResponse, errors.New("round not found")
		}
		adapterKeys = adapter.GetCommonAdapters(roundConfig.Bidding, br.Adapters.Keys())
		adapterKeys = adapter.GetCommonAdapters(adapterKeys, maps.Keys(br.Imp.Demands))
	}

	if len(adapterKeys) == 0 {
		return emptyResponse, ErrNoAdaptersMatched
	}

	auctionResult := AuctionResult{
		RoundNumber: roundNumber,
		Bids:        make([]adapters.DemandResponse, 0, len(adapterKeys)),
	}

	bids := make(chan adapters.DemandResponse)
	handleError := func(adapterKey adapter.Key, err error) {
		bids <- adapters.DemandResponse{
			DemandID: adapterKey,
			Error:    err,
			StartTS:  params.StartTS,
			EndTS:    time.Now().UnixMilli(),
		}
	}
	wg := sync.WaitGroup{}

	for _, adapterKey := range adapterKeys {
		wg.Add(1)
		go b.processAdapter(ctx, adapterKey, br, baseBidRequest, params, bids, &wg, handleError)
	}

	go func() {
		wg.Wait()
		close(bids)
	}()

	for bid := range bids {
		auctionResult.Bids = append(auctionResult.Bids, bid)
	}

	b.NotificationHandler.HandleBiddingRound(ctx, &br.Imp, auctionResult, br.App.Bundle, string(br.AdType))

	return auctionResult, nil
}

func (b *Builder) processAdapter(
	ctx context.Context,
	adapterKey adapter.Key,
	br schema.BiddingRequest,
	baseBidRequest openrtb.BidRequest,
	params *BuildParams,
	bids chan adapters.DemandResponse,
	wg *sync.WaitGroup,
	handleError func(adapter.Key, error),
) {
	defer wg.Done()

	if adapterKey == adapter.AmazonKey {
		bidder, err := amazon.Builder(params.AdapterConfigs)
		if err != nil {
			handleError(adapterKey, err)
			return
		}
		demandResponses, err := bidder.FetchBids(&br)
		if err != nil {
			handleError(adapterKey, err)
			return
		}
		for _, demandResponse := range demandResponses {
			demandResponse.StartTS = params.StartTS
			demandResponse.EndTS = time.Now().UnixMilli()

			bids <- *demandResponse
		}

		return
	}

	// adapter build bid request from baseBidRequest
	// adapter send bid request
	// adapter parse bid response
	bidder, err := b.AdaptersBuilder.Build(adapterKey, params.AdapterConfigs)
	if err != nil {
		handleError(adapterKey, err)
		return
	}

	bidRequest, err := bidder.Adapter.CreateRequest(baseBidRequest, &br)
	if err != nil {
		handleError(adapterKey, err)
		return
	}

	demandResponse := bidder.Adapter.ExecuteRequest(ctx, bidder.Client, bidRequest)
	demandResponse.StartTS = params.StartTS
	demandResponse.EndTS = time.Now().UnixMilli()
	if demandResponse.Error != nil {
		bids <- *demandResponse
		return
	}

	demandResponse, err = bidder.Adapter.ParseBids(demandResponse)
	demandResponse.Error = err

	b.setTokenResponse(demandResponse, &br)

	bids <- *demandResponse
}

func (b *Builder) BuildDevice(device schema.Device, user schema.User, geo geocoder.GeoData) *openrtb2.Device {
	js := int8(0)
	if device.JS != nil {
		js = int8(*device.JS)
	}

	return &openrtb2.Device{
		IP:             geo.IPString,
		W:              int64(device.Width),
		H:              int64(device.Height),
		JS:             js,
		DeviceType:     toAdcomDeviceType(device.Type),
		ConnectionType: toAdcomConnType(device.ConnectionType),
		OS:             device.OS,
		OSV:            device.OSVersion,
		PxRatio:        device.PXRatio,
		Language:       device.Language,
		Make:           device.Manufacturer,
		HWV:            device.HardwareVersion,
		UA:             device.UserAgent,
		PPI:            int64(device.PPI),
		Model:          device.Model,
		IFA:            user.IDFA,
		Geo: &openrtb2.Geo{
			Lat:       geo.Lat,
			Lon:       geo.Lon,
			Type:      adcom1.LocationIP,
			Accuracy:  int64(geo.Accuracy),
			IPService: adcom1.LocationServiceMaxMind,
			Country:   geo.CountryCode3,
			City:      geo.CityName,
			ZIP:       geo.ZipCode,
			Region:    geo.RegionCode,
		},
	}
}

func (b *Builder) setTokenResponse(demandResponse *adapters.DemandResponse, br *schema.BiddingRequest) {
	adapterKey := demandResponse.DemandID
	demandData, ok := br.Imp.Demands[adapterKey]
	if !ok || demandData == nil {
		return
	}

	if token, ok := demandData["token"].(string); ok {
		demandResponse.Token.Value = token
	}
	if status, ok := demandData["status"].(string); ok {
		demandResponse.Token.Status = status
	}
	if tokenStartTS, ok := demandData["token_start_ts"].(float64); ok {
		demandResponse.Token.StartTS = int64(tokenStartTS)
	}
	if tokenFinishTS, ok := demandData["token_finish_ts"].(float64); ok {
		demandResponse.Token.EndTS = int64(tokenFinishTS)
	}
}

func bool2int(b bool) *int8 {
	result := int8(0)
	if b {
		result = 1
	}
	return &result
}

func toAdcomDeviceType(deviceType device.Type) adcom1.DeviceType {
	switch deviceType {
	case device.TabletType:
		return adcom1.DeviceTablet
	case device.PhoneType:
		return adcom1.DevicePhone
	default:
		return adcom1.DeviceMobile
	}
}

func toAdcomConnType(connType string) *adcom1.ConnectionType {
	ct := adcom1.ConnectionUnknown

	switch connType {
	case "ETHERNET":
		ct = adcom1.ConnectionEthernet
	case "WIFI":
		ct = adcom1.ConnectionWIFI
	case "CELLULAR":
		ct = adcom1.ConnectionCellular
	case "CELLULAR_UNKNOWN":
		ct = adcom1.ConnectionCellular
	case "CELLULAR_2_G":
		ct = adcom1.Connection2G
	case "CELLULAR_3_G":
		ct = adcom1.Connection3G
	case "CELLULAR_4_G":
		ct = adcom1.Connection4G
	case "CELLULAR_5_G":
		ct = adcom1.Connection5G
	default:
		ct = adcom1.ConnectionUnknown
	}

	return &ct
}
