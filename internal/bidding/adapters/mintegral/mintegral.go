package mintegral

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/bidon-io/bidon-backend/internal/ad"
	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/bidding/adapters"
	"github.com/bidon-io/bidon-backend/internal/bidding/openrtb"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
	"github.com/gofrs/uuid/v5"
	"github.com/prebid/openrtb/v19/adcom1"
	"github.com/prebid/openrtb/v19/openrtb2"
)

type MintegralAdapter struct {
	SellerID    string
	AppID       string
	TagID       string
	PlacementID string
}

var bannerFormats = map[ad.Format][2]int64{
	ad.BannerFormat:      {320, 50},
	ad.LeaderboardFormat: {728, 90},
	ad.MRECFormat:        {300, 250},
	ad.AdaptiveFormat:    {320, 50},
	ad.EmptyFormat:       {320, 50}, // Default
}

var fullscreenFormats = map[string][2]int64{
	"PHONE":  {320, 480},
	"TABLET": {320, 480},
}

func (a *MintegralAdapter) banner(br *schema.BiddingRequest) *openrtb2.Imp {
	size := bannerFormats[br.Imp.Format()]

	if br.Imp.IsAdaptive() && br.Device.IsTablet() {
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

func (a *MintegralAdapter) interstitial(br *schema.BiddingRequest) *openrtb2.Imp {
	size := fullscreenFormats[string(br.Device.Type)]
	w, h := size[0], size[1]
	if !br.Imp.IsPortrait() {
		w, h = h, w
	}
	return &openrtb2.Imp{
		Instl: 1,
		Banner: &openrtb2.Banner{
			W:   &w,
			H:   &h,
			Pos: adcom1.PositionFullScreen.Ptr(),
		},
	}
}

func (a *MintegralAdapter) rewarded(br *schema.BiddingRequest) *openrtb2.Imp {
	size := fullscreenFormats[string(br.Device.Type)]
	w, h := size[0], size[1]
	if !br.Imp.IsPortrait() {
		w, h = h, w
	}
	return &openrtb2.Imp{
		Instl: 0,
		Video: &openrtb2.Video{
			W:     w,
			H:     h,
			MIMEs: []string{"video/mp4"},
		},
		Ext: json.RawMessage(`{"is_rewarded": true}`),
	}
}

func (a *MintegralAdapter) CreateRequest(request openrtb.BidRequest, br *schema.BiddingRequest) (openrtb.BidRequest, error) {
	secure := int8(1)

	var imp *openrtb2.Imp
	switch br.Imp.Type() {
	case ad.BannerType:
		imp = a.banner(br)
	case ad.InterstitialType:
		imp = a.interstitial(br)
	case ad.RewardedType:
		imp = a.rewarded(br)
	default:
		return request, errors.New("unknown impression type")
	}

	impId, _ := uuid.NewV4()
	imp.ID = impId.String()

	if a.TagID == "" {
		return request, errors.New("TagID is empty")
	}
	imp.TagID = a.TagID

	imp.DisplayManager = string(adapter.MintegralKey)
	imp.DisplayManagerVer = br.Adapters[adapter.MintegralKey].SDKVersion
	imp.Secure = &secure
	imp.BidFloor = br.Imp.GetBidFloorForBidding()
	imp.BidFloorCur = "USD"

	request.Imp = []openrtb2.Imp{*imp}
	request.Cur = []string{"USD"}
	request.User = &openrtb.User{
		BuyerUID: br.Imp.Demands[adapter.MintegralKey]["token"].(string),
	}
	request.App.Publisher.ID = a.SellerID
	request.App.ID = a.AppID

	appExtStructure := &map[string]interface{}{}
	if br.Imp.IsPortrait() {
		(*appExtStructure)["orientation"] = 1
	} else {
		(*appExtStructure)["orientation"] = 2
	}
	raw, _ := json.Marshal(appExtStructure)
	request.App.Ext = raw

	return request, nil
}

func (a *MintegralAdapter) ExecuteRequest(ctx context.Context, client *http.Client, request openrtb.BidRequest) *adapters.DemandResponse {
	dr := &adapters.DemandResponse{
		DemandID:    adapter.MintegralKey,
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

	url := "http://hb.rayjump.com/bid"
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		dr.Error = err
		return dr
	}
	httpReq.Header.Add("Content-Type", "application/json")
	httpReq.Header.Add("openrtb", "2.5")

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

func (a *MintegralAdapter) ParseBids(dr *adapters.DemandResponse) (*adapters.DemandResponse, error) {
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
		return dr, fmt.Errorf("unauthorized request: " + strconv.Itoa(dr.Status))
	case http.StatusOK:
		break
	default:
		return dr, fmt.Errorf("unexpected status code: " + strconv.Itoa(dr.Status))
	}

	var bidResponse openrtb2.BidResponse
	err := json.Unmarshal([]byte(dr.RawResponse), &bidResponse)
	if err != nil {
		return dr, err
	}

	if bidResponse.SeatBid == nil {
		return dr, nil
	}

	seat := bidResponse.SeatBid[0]
	bid := seat.Bid[0]

	dr.Bid = &adapters.BidDemandResponse{
		ID:       bid.ID,
		ImpID:    bid.ImpID,
		Price:    bid.Price,
		Payload:  bid.AdM,
		DemandID: adapter.MintegralKey,
		AdID:     bid.AdID,
		SeatID:   seat.Seat,
		LURL:     bid.LURL,
		NURL:     bid.NURL,
		BURL:     bid.BURL,
	}

	return dr, nil
}

// Builder builds a new instance of the Mintegral adapter for the given bidder with the given config.
func Builder(cfg adapter.ProcessedConfigsMap, client *http.Client) (*adapters.Bidder, error) {
	mCfg := cfg[adapter.MintegralKey]

	sellerID, ok := mCfg["seller_id"].(string)
	if !ok || sellerID == "" {
		return nil, fmt.Errorf("missing seller_id param for %s adapter", adapter.MintegralKey)
	}
	appID, ok := mCfg["app_id"].(string)
	if !ok || appID == "" {
		return nil, fmt.Errorf("missing app_id param for %s adapter", adapter.MintegralKey)
	}
	tagID, ok := mCfg["tag_id"].(string)
	if !ok {
		tagID = ""
	}
	placementID, ok := mCfg["placement_id"].(string)
	if !ok || placementID == "" {
		placementID = ""
	}

	adpt := &MintegralAdapter{
		SellerID:    sellerID,
		AppID:       appID,
		TagID:       tagID,
		PlacementID: placementID,
	}

	bidder := adapters.Bidder{
		Adapter: adpt,
		Client:  client,
	}

	return &bidder, nil
}
