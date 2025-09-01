package inmobi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gofrs/uuid/v5"
	"github.com/prebid/openrtb/v19/adcom1"
	"github.com/prebid/openrtb/v19/openrtb2"

	"github.com/bidon-io/bidon-backend/internal/ad"
	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/bidding/adapters"
	"github.com/bidon-io/bidon-backend/internal/bidding/openrtb"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
)

type InMobiAdapter struct {
	AppID       string
	PlacementID string
}

var bannerFormats = map[ad.Format][2]int64{
	ad.BannerFormat:      {320, 50},
	ad.LeaderboardFormat: {728, 90},
	ad.MRECFormat:        {300, 250},
	ad.AdaptiveFormat:    {320, 50},
	ad.EmptyFormat:       {320, 50}, // Default
}

func (a *InMobiAdapter) banner(auctionRequest *schema.AuctionRequest) *openrtb2.Imp {
	size := bannerFormats[auctionRequest.AdObject.Format()]

	if auctionRequest.AdObject.IsAdaptive() && auctionRequest.Device.IsTablet() {
		size = bannerFormats[ad.LeaderboardFormat]
	}

	w, h := size[0], size[1]

	return &openrtb2.Imp{
		Instl: 0,
		Banner: &openrtb2.Banner{
			W:   &w,
			H:   &h,
			Pos: adcom1.PositionAboveFold.Ptr(),
			API: []adcom1.APIFramework{3, 5}, // MRAID 1.0, MRAID 2.0
		},
	}
}

func (a *InMobiAdapter) interstitial(auctionRequest *schema.AuctionRequest) *openrtb2.Imp {
	size := adapters.FullscreenFormats[string(auctionRequest.Device.Type)]
	w, h := size[0], size[1]
	if !auctionRequest.AdObject.IsPortrait() {
		w, h = h, w
	}
	return &openrtb2.Imp{
		Instl: 1,
		Banner: &openrtb2.Banner{
			W:   &w,
			H:   &h,
			Pos: adcom1.PositionFullScreen.Ptr(),
			API: []adcom1.APIFramework{3, 5}, // MRAID 1.0, MRAID 2.0
		},
	}
}

func (a *InMobiAdapter) rewarded(auctionRequest *schema.AuctionRequest) *openrtb2.Imp {
	size := adapters.FullscreenFormats[string(auctionRequest.Device.Type)]
	w, h := size[0], size[1]
	if !auctionRequest.AdObject.IsPortrait() {
		w, h = h, w
	}
	return &openrtb2.Imp{
		Instl: 0,
		Video: &openrtb2.Video{
			W:           w,
			H:           h,
			MIMEs:       []string{"video/mp4"},
			MinDuration: 0,
			MaxDuration: 6000,
			Protocols:   []adcom1.MediaCreativeSubtype{2, 3, 5, 6}, // VAST 2.0, VAST 3.0, VAST 4.0 Wrapper, VAST 4.1 Wrapper
			StartDelay:  adcom1.StartDelay(0).Ptr(),                // Pre-roll
			API:         []adcom1.APIFramework{1, 2, 3, 5, 6, 7},   // VPAID 1.0, VPAID 2.0, MRAID 1.0, MRAID 2.0, MRAID 3.0, OMID 1.0
			Pos:         adcom1.PositionAboveFold.Ptr(),
		},
		Ext: json.RawMessage(`{"is_rewarded": true}`),
	}
}

func (a *InMobiAdapter) CreateRequest(request openrtb.BidRequest, auctionRequest *schema.AuctionRequest) (openrtb.BidRequest, error) {
	if a.PlacementID == "" {
		return request, errors.New("PlacementID is empty")
	}

	secure := int8(1)
	var imp *openrtb2.Imp

	switch auctionRequest.AdObject.Type() {
	case ad.BannerType:
		imp = a.banner(auctionRequest)
	case ad.InterstitialType:
		imp = a.interstitial(auctionRequest)
	case ad.RewardedType:
		imp = a.rewarded(auctionRequest)
	default:
		return request, errors.New("unknown impression type")
	}

	impId, _ := uuid.NewV4()
	imp.ID = impId.String()
	imp.TagID = a.PlacementID
	imp.DisplayManager = string(adapter.InmobiKey)
	imp.DisplayManagerVer = auctionRequest.Adapters[adapter.InmobiKey].SDKVersion
	imp.Secure = &secure
	imp.BidFloor = adapters.CalculatePriceFloor(&request, auctionRequest)
	imp.BidFloorCur = "USD"

	request.Imp = []openrtb2.Imp{*imp}
	request.Cur = []string{"USD"}

	// Set user with bidding token
	if token, exists := auctionRequest.AdObject.Demands[adapter.InmobiKey]["token"]; exists {
		request.User = &openrtb.User{
			BuyerUID: token.(string),
		}
	}

	// Set app ID
	request.App.ID = a.AppID

	return request, nil
}

func (a *InMobiAdapter) ExecuteRequest(ctx context.Context, client *http.Client, request openrtb.BidRequest) *adapters.DemandResponse {
	dr := &adapters.DemandResponse{
		DemandID:    adapter.InmobiKey,
		RequestID:   request.ID,
		PlacementID: a.PlacementID,
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		dr.Error = err
		return dr
	}
	dr.RawRequest = string(requestBody)

	url := "https://api.w.inmobi.com/ortb/imsdk"
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(requestBody))
	if err != nil {
		dr.Error = err
		return dr
	}

	httpReq.Header.Add("Content-Type", "application/json")
	httpReq.Header.Add("X-OpenRTB-Version", "2.5")

	httpResp, err := client.Do(httpReq)
	if err != nil {
		dr.Error = err
		return dr
	}
	defer httpResp.Body.Close()

	dr.Status = httpResp.StatusCode
	responseBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		dr.Error = err
		return dr
	}
	dr.RawResponse = string(responseBody)

	return dr
}

func (a *InMobiAdapter) ParseBids(dr *adapters.DemandResponse) (*adapters.DemandResponse, error) {
	switch dr.Status {
	case http.StatusNoContent:
		return dr, nil
	case http.StatusServiceUnavailable:
		fallthrough
	case http.StatusBadRequest:
		fallthrough
	case http.StatusUnauthorized:
		fallthrough
	case http.StatusForbidden:
		return dr, fmt.Errorf("unauthorized request: %s", strconv.Itoa(dr.Status))
	case http.StatusOK:
		break
	default:
		return dr, fmt.Errorf("unexpected status code: %s", strconv.Itoa(dr.Status))
	}

	var bidResponse openrtb2.BidResponse
	err := json.Unmarshal([]byte(dr.RawResponse), &bidResponse)
	if err != nil {
		return dr, err
	}

	if len(bidResponse.SeatBid) == 0 {
		return dr, nil
	}

	seat := bidResponse.SeatBid[0]
	if len(seat.Bid) == 0 {
		return dr, nil
	}

	bid := seat.Bid[0]

	dr.Bid = &adapters.BidDemandResponse{
		ID:       bid.ID,
		ImpID:    bid.ImpID,
		Price:    bid.Price,
		Payload:  bid.AdM,
		DemandID: adapter.InmobiKey,
		AdID:     bid.AdID,
		SeatID:   seat.Seat,
		LURL:     bid.LURL,
		NURL:     bid.NURL,
		BURL:     bid.BURL,
	}

	return dr, nil
}

// Builder builds a new instance of the InMobi adapter for the given bidder with the given config.
func Builder(cfg adapter.ProcessedConfigsMap, client *http.Client) (*adapters.Bidder, error) {
	inmobiCfg := cfg[adapter.InmobiKey]

	appID, ok := inmobiCfg["app_id"].(string)
	if !ok || appID == "" {
		return nil, fmt.Errorf("missing app_id param for %s adapter", adapter.InmobiKey)
	}

	placementID, ok := inmobiCfg["placement_id"].(string)
	if !ok || placementID == "" {
		return nil, fmt.Errorf("missing placement_id param for %s adapter", adapter.InmobiKey)
	}

	adpt := &InMobiAdapter{
		AppID:       appID,
		PlacementID: placementID,
	}

	bidder := adapters.Bidder{
		Adapter: adpt,
		Client:  client,
	}

	return &bidder, nil
}
