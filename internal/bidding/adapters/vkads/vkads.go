package vkads

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

type VKAdsAdapter struct {
	TagID string
	AppID string
}

var bannerFormats = map[ad.Format][2]int64{
	ad.BannerFormat:      {320, 50},
	ad.LeaderboardFormat: {728, 90},
	ad.MRECFormat:        {300, 250},
	ad.AdaptiveFormat:    {320, 50},
	ad.EmptyFormat:       {320, 50}, // Default
}

const (
	rewardedWidth  int64 = 1920
	rewardedHeight int64 = 1080
)

func (a *VKAdsAdapter) banner(auctionRequest *schema.AuctionRequest) *openrtb2.Imp {
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
		},
	}
}

func (a *VKAdsAdapter) interstitial() *openrtb2.Imp {
	return &openrtb2.Imp{
		Instl: 1,
		Banner: &openrtb2.Banner{
			Pos: adcom1.PositionFullScreen.Ptr(),
		},
	}
}

func (a *VKAdsAdapter) rewarded() *openrtb2.Imp {
	w, h := rewardedWidth, rewardedHeight
	return &openrtb2.Imp{
		Banner: &openrtb2.Banner{
			W: &w,
			H: &h,
		},
	}
}

func (a *VKAdsAdapter) CreateRequest(request openrtb.BidRequest, auctionRequest *schema.AuctionRequest) (openrtb.BidRequest, error) {
	if a.AppID == "" {
		return request, errors.New("AppID is empty")
	}
	if a.TagID == "" {
		return request, errors.New("TagID is empty")
	}
	token, ok := auctionRequest.AdObject.Demands[adapter.VKAdsKey]["token"].(string)
	if !ok || token == "" {
		return request, errors.New("token is empty")
	}

	var imp *openrtb2.Imp
	switch auctionRequest.AdObject.Type() {
	case ad.BannerType:
		imp = a.banner(auctionRequest)
	case ad.InterstitialType:
		imp = a.interstitial()
	case ad.RewardedType:
		imp = a.rewarded()
	default:
		return request, errors.New("unknown impression type")
	}

	impId, _ := uuid.NewV4()
	imp.ID = impId.String()
	imp.TagID = a.TagID
	imp.BidFloor = adapters.CalculatePriceFloor(&request, auctionRequest)
	imp.BidFloorCur = "USD"

	request.Imp = []openrtb2.Imp{*imp}
	request.Cur = []string{"USD"}

	request.User = &openrtb.User{
		ID:  auctionRequest.User.IDG,
		Ext: json.RawMessage(fmt.Sprintf(`{"buyeruid": "%s"}`, token)),
	}

	request.App.ID = a.AppID
	request.Ext = json.RawMessage(`{"pid":111}`)

	return request, nil
}

func (a *VKAdsAdapter) ExecuteRequest(ctx context.Context, client *http.Client, request openrtb.BidRequest) *adapters.DemandResponse {
	dr := &adapters.DemandResponse{
		DemandID:  adapter.VKAdsKey,
		RequestID: request.ID,
		TagID:     a.TagID,
	}
	requestBody, err := json.Marshal(request)
	if err != nil {
		dr.Error = err
		return dr
	}
	dr.RawRequest = string(requestBody)

	url := "https://ad.mail.ru/api/bid"
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

func (a *VKAdsAdapter) ParseBids(dr *adapters.DemandResponse) (*adapters.DemandResponse, error) {
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

	var bidResponse openrtb.BidResponse
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
		DemandID: adapter.VKAdsKey,
		AdID:     bid.AdID,
		SeatID:   seat.Seat,
		LURL:     bid.LURL,
		NURL:     bid.NURL,
		BURL:     bid.BURL,
	}

	return dr, nil
}

// Builder builds a new instance of the VKAds adapter for the given bidder with the given config.
func Builder(cfg adapter.ProcessedConfigsMap, client *http.Client) (*adapters.Bidder, error) {
	vkCfg := cfg[adapter.VKAdsKey]

	appID, ok := vkCfg["app_id"].(string)
	if !ok || appID == "" {
		appID = ""
	}

	tagID, ok := vkCfg["tag_id"].(string)
	if !ok {
		tagID = ""
	}

	adpt := &VKAdsAdapter{
		AppID: appID,
		TagID: tagID,
	}

	bidder := &adapters.Bidder{
		Adapter: adpt,
		Client:  client,
	}

	return bidder, nil
}
