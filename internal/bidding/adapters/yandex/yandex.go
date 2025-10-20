package yandex

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

const yandexEndpoint = "https://mobile.yandexadexchange.net/openbidding?ssp-id=99048272"

type YandexAdapter struct {
	AdUnitID string
}

var bannerFormats = map[ad.Format][2]int64{
	ad.BannerFormat:      {320, 50},
	ad.LeaderboardFormat: {728, 90},
	ad.MRECFormat:        {300, 250},
	ad.AdaptiveFormat:    {320, 50},
	ad.EmptyFormat:       {320, 50}, // Default
}

func (a *YandexAdapter) banner(auctionRequest *schema.AuctionRequest) *openrtb2.Imp {
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

func (a *YandexAdapter) interstitial() *openrtb2.Imp {
	return &openrtb2.Imp{
		Instl: 1,
	}
}

func (a *YandexAdapter) rewarded() *openrtb2.Imp {
	return &openrtb2.Imp{
		Instl: 0,
		Video: &openrtb2.Video{
			MIMEs: []string{"video/mp4"},
		},
	}
}

func (a *YandexAdapter) CreateRequest(request openrtb.BidRequest, auctionRequest *schema.AuctionRequest) (openrtb.BidRequest, error) {
	if a.AdUnitID == "" {
		return request, errors.New("AdUnitID is empty")
	}

	secure := int8(1)

	var imp *openrtb2.Imp
	var adTypeString string
	var rwdd int8

	switch auctionRequest.AdObject.Type() {
	case ad.BannerType:
		imp = a.banner(auctionRequest)
		adTypeString = "banner"
		rwdd = 0
	case ad.InterstitialType:
		imp = a.interstitial()
		adTypeString = "interstitial"
		rwdd = 0
	case ad.RewardedType:
		imp = a.rewarded()
		adTypeString = "rewarded"
		rwdd = 1
	default:
		return request, errors.New("unknown impression type")
	}

	impId, _ := uuid.NewV4()
	imp.ID = impId.String()
	imp.TagID = a.AdUnitID
	imp.Rwdd = rwdd

	// Set imp.ext.ad_type
	impExt, err := json.Marshal(map[string]any{
		"ad_type": adTypeString,
	})
	if err != nil {
		return request, err
	}
	imp.Ext = impExt

	imp.DisplayManager = string(adapter.YandexKey)
	imp.DisplayManagerVer = auctionRequest.Adapters[adapter.YandexKey].SDKVersion
	imp.Secure = &secure
	imp.BidFloor = adapters.CalculatePriceFloor(&request, auctionRequest)
	imp.BidFloorCur = "USD"

	request.Imp = []openrtb2.Imp{*imp}
	request.Cur = []string{"USD"}

	// Extract bidder token from auction request
	demandData, ok := auctionRequest.AdObject.Demands[adapter.YandexKey]
	if !ok {
		return request, errors.New("yandex demand data is missing")
	}

	token, ok := demandData["token"].(string)
	if !ok || token == "" {
		return request, errors.New("yandex bidder token is empty")
	}

	request.User = &openrtb.User{
		Data: []openrtb.Data{
			{
				Segment: []openrtb.Segment{
					{
						Signal: token,
					},
				},
			},
		},
	}

	return request, nil
}

func (a *YandexAdapter) ExecuteRequest(ctx context.Context, client *http.Client, request openrtb.BidRequest) *adapters.DemandResponse {
	dr := &adapters.DemandResponse{
		DemandID:  adapter.YandexKey,
		RequestID: request.ID,
		TagID:     a.AdUnitID,
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		dr.Error = err
		return dr
	}
	dr.RawRequest = string(requestBody)

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, yandexEndpoint, bytes.NewBuffer(requestBody))
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

func (a *YandexAdapter) ParseBids(dr *adapters.DemandResponse) (*adapters.DemandResponse, error) {
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

	if bidResponse.SeatBid == nil || len(bidResponse.SeatBid) == 0 {
		return dr, nil
	}

	seat := bidResponse.SeatBid[0]
	if len(seat.Bid) == 0 {
		return dr, nil
	}

	bid := seat.Bid[0]

	// Extract signaldata from bid.ext
	type BidExt struct {
		SignalData string `json:"signaldata"`
	}

	var bidExt BidExt
	if bid.Ext != nil {
		if err := json.Unmarshal(bid.Ext, &bidExt); err != nil {
			return dr, err
		}
	}

	dr.Bid = &adapters.BidDemandResponse{
		ID:         bid.ID,
		ImpID:      bid.ImpID,
		Price:      bid.Price,
		Payload:    bid.AdM,
		Signaldata: bidExt.SignalData,
		DemandID:   adapter.YandexKey,
		AdID:       bid.AdID,
		SeatID:     seat.Seat,
		LURL:       bid.LURL,
		NURL:       bid.NURL,
		BURL:       bid.BURL,
	}

	return dr, nil
}

// Builder builds a new instance of the Yandex adapter for the given bidder with the given config.
func Builder(cfg adapter.ProcessedConfigsMap, client *http.Client) (*adapters.Bidder, error) {
	yandexCfg := cfg[adapter.YandexKey]

	adUnitID, ok := yandexCfg["ad_unit_id"].(string)
	if !ok {
		adUnitID = ""
	}

	adpt := &YandexAdapter{
		AdUnitID: adUnitID,
	}

	bidder := &adapters.Bidder{
		Adapter: adpt,
		Client:  client,
	}

	return bidder, nil
}
