package vungle

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
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
	"github.com/gofrs/uuid/v5"
	"github.com/prebid/openrtb/v19/adcom1"
	"github.com/prebid/openrtb/v19/openrtb2"
)

type VungleAdapter struct {
	SellerID string
	AppID    string
	TagID    string
}

var bannerFormats = map[string][2]int64{
	"BANNER":      {320, 50},
	"LEADERBOARD": {728, 90},
	"MREC":        {300, 250},
	"ADAPTIVE":    {0, 50},
	"":            {320, 50}, // Default
}

var fullscreenFormats = map[string][2]int64{
	"PHONE":  {320, 480},
	"TABLET": {768, 1024},
}

func (a *VungleAdapter) banner(br *schema.BiddingRequest) *openrtb2.Imp {
	size := bannerFormats[string(br.Imp.Format())]
	w, h := size[0], size[1]
	if !br.Imp.IsPortrait() {
		w, h = h, w
	}
	return &openrtb2.Imp{
		Instl: 0,
		Banner: &openrtb2.Banner{
			W:   &w,
			H:   &h,
			Pos: adcom1.PositionAboveFold.Ptr(),
		},
	}
}

func (a *VungleAdapter) interstitial(br *schema.BiddingRequest) *openrtb2.Imp {
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

func (a *VungleAdapter) rewarded(br *schema.BiddingRequest) *openrtb2.Imp {
	size := fullscreenFormats[string(br.Device.Type)]
	w, h := size[0], size[1]
	if !br.Imp.IsPortrait() {
		w, h = h, w
	}
	return &openrtb2.Imp{
		Video: &openrtb2.Video{
			W:     w,
			H:     h,
			MIMEs: []string{"video/mp4"},
			Ext:   json.RawMessage(`{"rewarded": 1}`),
		},
	}
}

func (a *VungleAdapter) CreateRequest(request openrtb2.BidRequest, br *schema.BiddingRequest) (openrtb2.BidRequest, error) {
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

	imp.DisplayManager = string(adapter.VungleKey)
	imp.DisplayManagerVer = br.Adapters[adapter.VungleKey].SDKVersion
	imp.Secure = &secure
	imp.BidFloor = br.Imp.GetBidFloor()
	imp.BidFloorCur = "USD"

	vungleData := make(map[string]interface{})
	vungleData["bid_token"] = br.Imp.Demands[adapter.VungleKey]["token"]

	extStructure := &map[string]interface{}{}
	_ = json.Unmarshal(imp.Ext, extStructure)
	(*extStructure)["vungle"] = vungleData
	raw, _ := json.Marshal(extStructure)
	imp.Ext = raw

	request.Imp = []openrtb2.Imp{*imp}
	request.Cur = []string{"USD"}

	request.App.Publisher.ID = a.SellerID
	request.App.ID = a.AppID

	return request, nil
}

func (a *VungleAdapter) ExecuteRequest(ctx context.Context, client *http.Client, request openrtb2.BidRequest) *adapters.DemandResponse {
	dr := &adapters.DemandResponse{
		DemandID: adapter.VungleKey,
		TagID:    a.TagID,
	}
	requestBody, err := json.Marshal(request)
	if err != nil {
		dr.Error = err
		return dr
	}
	dr.RawRequest = string(requestBody)

	url := "https://rtb.ads.vungle.com/bid/t/8ea3e9a"
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		dr.Error = err
		return dr
	}
	httpReq.Header.Add("Content-Type", "application/json")
	// It's important to set  X-OpenRTB-Version header to 2.5
	httpReq.Header.Add("X-OpenRTB-Version", "2.5")

	httpResp, err := client.Do(httpReq)
	if err != nil {
		if err == context.DeadlineExceeded {
			fmt.Println("Timeout")
			// TODO: Send Timeout Notification if bidder support, eg FB
		}
		dr.Error = err
		return dr
	}

	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		dr.Error = err
		return dr
	}
	defer httpResp.Body.Close()
	dr.RawResponse = string(respBody)
	dr.Status = httpResp.StatusCode

	return dr
}

func (a *VungleAdapter) ParseBids(dr *adapters.DemandResponse) (*adapters.DemandResponse, error) {
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

	seat := bidResponse.SeatBid[0]
	bid := seat.Bid[0]

	dr.Bid = &adapters.BidDemandResponse{
		ID:       bid.ID,
		ImpID:    bid.ImpID,
		Price:    bid.Price,
		Payload:  bid.AdM,
		DemandID: adapter.VungleKey,
		AdID:     bid.AdID,
		SeatID:   seat.Seat,
		LURL:     bid.LURL,
		NURL:     bid.NURL,
		BURL:     bid.BURL,
	}

	return dr, nil
}

// Builder builds a new instance of the Vungle adapter for the given bidder with the given config.
func Builder(cfg adapter.ProcessedConfigsMap, client *http.Client) (*adapters.Bidder, error) {
	vCfg := cfg[adapter.VungleKey]

	sellerID, ok := vCfg["seller_id"].(string)
	if !ok || sellerID == "" {
		return nil, fmt.Errorf("missing seller_id param for %s adapter", adapter.VungleKey)
	}
	appID, ok := vCfg["app_id"].(string)
	if !ok || appID == "" {
		return nil, fmt.Errorf("missing app_id param for %s adapter", adapter.MintegralKey)
	}

	adpt := &VungleAdapter{
		SellerID: sellerID,
		AppID:    appID,
		TagID:    vCfg["tag_id"].(string),
	}

	bidder := adapters.Bidder{
		Adapter: adpt,
		Client:  client,
	}

	return &bidder, nil
}
