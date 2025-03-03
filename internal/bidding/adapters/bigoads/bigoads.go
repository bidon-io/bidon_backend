package bigoads

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

type BigoAdsAdapter struct {
	SellerID    string
	AppID       string
	TagID       string
	PlacementID string
}

var bannerFormats = map[ad.Format][2]int64{
	ad.BannerFormat:   {320, 50},
	ad.MRECFormat:     {300, 250},
	ad.AdaptiveFormat: {320, 50},
	ad.EmptyFormat:    {320, 50}, // Default
}

func (a *BigoAdsAdapter) banner(br *schema.BiddingRequest) (*openrtb2.Imp, error) {
	size, ok := bannerFormats[br.Imp.Format()]
	if !ok || br.Imp.IsAdaptive() && br.Device.IsTablet() { // Does not support leaderboard format
		return nil, fmt.Errorf("unknown banner format: %s", br.Imp.Format())
	}

	w, h := size[0], size[1]

	return &openrtb2.Imp{
		Instl: 0,
		Banner: &openrtb2.Banner{
			W:   &w,
			H:   &h,
			Pos: adcom1.PositionAboveFold.Ptr(),
		},
	}, nil
}

func (a *BigoAdsAdapter) interstitial() *openrtb2.Imp {
	return &openrtb2.Imp{
		Instl: 1,
		Banner: &openrtb2.Banner{
			Pos: adcom1.PositionFullScreen.Ptr(),
		},
	}
}

func (a *BigoAdsAdapter) rewarded() *openrtb2.Imp {
	return &openrtb2.Imp{
		Instl: 0,
		Video: &openrtb2.Video{
			MIMEs: []string{"video/mp4"},
		},
	}
}

func (a *BigoAdsAdapter) CreateRequest(request openrtb.BidRequest, br *schema.BiddingRequest) (openrtb.BidRequest, error) {
	if a.TagID == "" {
		return request, errors.New("TagID is empty")
	}

	secure := int8(1)

	var imp *openrtb2.Imp
	var impAdType int
	switch br.Imp.Type() {
	case ad.BannerType:
		bannerImp, err := a.banner(br)
		if err != nil {
			return request, err
		}
		imp = bannerImp
		impAdType = 2
	case ad.InterstitialType:
		imp = a.interstitial()
		impAdType = 3
	case ad.RewardedType:
		imp = a.rewarded()
		impAdType = 4
	default:
		return request, errors.New("unknown impression type")
	}

	impId, _ := uuid.NewV4()
	imp.ID = impId.String()
	imp.TagID = a.TagID

	impExt, err := json.Marshal(map[string]any{
		"adtype": impAdType,
		"networkid": map[string]any{
			"appid":       a.AppID,
			"placementid": a.PlacementID,
		},
	})
	if err != nil {
		return request, err
	}

	imp.Ext = impExt

	imp.DisplayManager = string(adapter.BigoAdsKey)
	imp.DisplayManagerVer = br.Adapters[adapter.BigoAdsKey].SDKVersion
	imp.Secure = &secure
	imp.BidFloor = adapters.CalculatePriceFloor(&request, br)
	request.Imp = []openrtb2.Imp{*imp}
	request.Cur = []string{"USD"}
	request.User = &openrtb.User{
		BuyerUID: br.Imp.Demands[adapter.BigoAdsKey]["token"].(string),
	}
	request.App.Publisher.ID = a.SellerID
	request.App.ID = a.AppID

	return request, nil
}

func (a *BigoAdsAdapter) ExecuteRequest(ctx context.Context, client *http.Client, request openrtb.BidRequest) *adapters.DemandResponse {
	dr := &adapters.DemandResponse{
		DemandID:    adapter.BigoAdsKey,
		RequestID:   request.ID,
		TagID:       a.TagID,
		PlacementID: a.PlacementID,
	}
	requestBody, err := json.Marshal(request)
	if err != nil {
		dr.Error = err
		return dr
	}
	dr.RawRequest = string(requestBody)

	url := "https://api.gov-static.tech/Ad/GetUniAdS2s?id=200104"
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(requestBody))
	if err != nil {
		dr.Error = err
		return dr
	}
	httpReq.Header.Add("Content-Type", "application/json")

	httpResp, err := client.Do(httpReq)
	if err != nil {
		dr.Error = err
		return dr
	}
	defer httpResp.Body.Close()

	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		dr.Error = err
		return dr
	}

	dr.RawResponse = string(respBody)
	dr.Status = httpResp.StatusCode

	return dr
}

func (a *BigoAdsAdapter) ParseBids(dr *adapters.DemandResponse) (*adapters.DemandResponse, error) {
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

	seat := bidResponse.SeatBid[0]
	bid := seat.Bid[0]

	dr.Bid = &adapters.BidDemandResponse{
		ID:       bid.ID,
		ImpID:    bid.ImpID,
		Price:    bid.Price,
		Payload:  bid.AdM,
		DemandID: adapter.BigoAdsKey,
		AdID:     bid.AdID,
		SeatID:   seat.Seat,
		LURL:     bid.LURL,
		NURL:     bid.NURL,
		BURL:     bid.BURL,
	}

	return dr, nil
}

// Builder builds a new instance of the BigoAds adapter for the given bidder with the given config.
func Builder(cfg adapter.ProcessedConfigsMap, client *http.Client) (*adapters.Bidder, error) {
	bigoCfg := cfg[adapter.BigoAdsKey]

	sellerID, ok := bigoCfg["seller_id"].(string)
	if !ok || sellerID == "" {
		return nil, fmt.Errorf("missing seller_id param for %s adapter", adapter.BigoAdsKey)
	}
	appID, ok := bigoCfg["app_id"].(string)
	if !ok || appID == "" {
		return nil, fmt.Errorf("missing app_id param for %s adapter", adapter.BigoAdsKey)
	}
	tagID, ok := bigoCfg["tag_id"].(string)
	if !ok || tagID == "" {
		return nil, fmt.Errorf("missing tag_id param for %s adapter", adapter.BigoAdsKey)
	}
	placementID, ok := bigoCfg["placement_id"].(string)
	if !ok {
		placementID = ""
	}

	adpt := &BigoAdsAdapter{
		SellerID:    sellerID,
		AppID:       appID,
		TagID:       tagID,
		PlacementID: placementID,
	}

	bidder := &adapters.Bidder{
		Adapter: adpt,
		Client:  client,
	}

	return bidder, nil
}
