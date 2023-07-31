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

	"github.com/bidon-io/bidon-backend/internal/ad"
	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/bidding/adapters"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
	"github.com/gofrs/uuid/v5"
	"github.com/prebid/openrtb/v19/adcom1"
	"github.com/prebid/openrtb/v19/openrtb2"
)

type BigoAdsAdapter struct {
	SellerID string
	Endpoint string
	AppID    string
	TagID    string
}

var bannerFormats = map[string][2]int64{
	"BANNER": {320, 50},
	"MREC":   {300, 250},
	"":       {320, 50}, // Default
}

func (a *BigoAdsAdapter) banner(br *schema.BiddingRequest) (*openrtb2.Imp, error) {
	size, ok := bannerFormats[string(br.Imp.Format())]
	if !ok {
		return nil, errors.New("unknown banner format")
	}

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

func (a *BigoAdsAdapter) CreateRequest(request openrtb2.BidRequest, br *schema.BiddingRequest) (openrtb2.BidRequest, []error) {
	if a.TagID == "" {
		return request, []error{errors.New("TagID is empty")}
	}

	var errs []error
	secure := int8(1)

	var imp *openrtb2.Imp
	var impAdType int
	switch br.Imp.Type() {
	case ad.BannerType:
		bannerImp, err := a.banner(br)
		if err != nil {
			return request, []error{err}
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
		return request, []error{errors.New("unknown impression type")}
	}

	impId, _ := uuid.NewV4()
	imp.ID = impId.String()
	imp.TagID = a.TagID

	impExt, err := json.Marshal(map[string]any{
		"adtype": impAdType,
		"networkid": map[string]any{
			"appid":       a.AppID,
			"placementid": a.TagID,
		},
	})
	if err != nil {
		return request, []error{err}
	}

	imp.Ext = json.RawMessage(impExt)

	imp.DisplayManager = string(adapter.BigoAdsKey)
	imp.DisplayManagerVer = br.Adapters[adapter.BigoAdsKey].SDKVersion
	imp.Secure = &secure
	imp.BidFloor = br.Imp.GetBidFloor()
	request.Imp = []openrtb2.Imp{*imp}
	request.Cur = []string{"USD"}
	request.User = &openrtb2.User{
		BuyerUID: br.Imp.Demands[adapter.BigoAdsKey]["token"].(string),
	}
	request.App.Publisher.ID = a.SellerID
	request.App.ID = a.AppID

	return request, errs
}

func (a *BigoAdsAdapter) ExecuteRequest(ctx context.Context, client *http.Client, request openrtb2.BidRequest) *adapters.DemandResponse {
	dr := &adapters.DemandResponse{
		DemandID: adapter.BigoAdsKey,
	}
	requestBody, err := json.Marshal(request)
	if err != nil {
		dr.Error = err
		return dr
	}
	dr.RawRequest = string(requestBody)

	url := "https://api.gov-static.tech/Ad/GetUniAdS2s?id=200104"
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		dr.Error = err
		return dr
	}
	httpReq.Header.Add("Content-Type", "application/json")

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
func Builder(cfg adapter.Config, client *http.Client) (adapters.Bidder, error) {
	bigoCfg := cfg[adapter.BigoAdsKey]

	adpt := &BigoAdsAdapter{
		Endpoint: bigoCfg["endpoint"].(string),
		SellerID: bigoCfg["seller_id"].(string),
		AppID:    bigoCfg["app_id"].(string),
		TagID:    bigoCfg["tag_id"].(string),
	}

	bidder := adapters.Bidder{
		Adapter: adpt,
		Client:  client,
	}

	return bidder, nil
}
