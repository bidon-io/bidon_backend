package mobilefuse

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

type MobileFuseAdapter struct {
	TagID string
}

var bannerFormats = map[string][2]int64{
	"BANNER": {320, 50},
	"MREC":   {300, 250},
	"":       {320, 50}, // Default
}

func (a *MobileFuseAdapter) banner(br *schema.BiddingRequest) (*openrtb2.Imp, error) {
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

func (a *MobileFuseAdapter) interstitial() *openrtb2.Imp {
	return &openrtb2.Imp{
		Instl: 1,
		Banner: &openrtb2.Banner{
			Pos: adcom1.PositionFullScreen.Ptr(),
		},
	}
}

func (a *MobileFuseAdapter) rewarded() *openrtb2.Imp {
	return &openrtb2.Imp{
		Instl: 0,
		Video: &openrtb2.Video{
			MIMEs: []string{"video/mp4"},
		},
	}
}

func (a *MobileFuseAdapter) CreateRequest(request openrtb.BidRequest, br *schema.BiddingRequest) (openrtb.BidRequest, error) {
	if a.TagID == "" {
		return request, errors.New("TagID is empty")
	}

	secure := int8(1)

	var imp *openrtb2.Imp
	switch br.Imp.Type() {
	case ad.BannerType:
		bannerImp, err := a.banner(br)
		if err != nil {
			return request, err
		}
		imp = bannerImp
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

	imp.DisplayManager = string(adapter.MobileFuseKey)
	imp.DisplayManagerVer = br.Adapters[adapter.MobileFuseKey].SDKVersion
	imp.Secure = &secure
	imp.BidFloor = br.Imp.GetBidFloor()
	request.Imp = []openrtb2.Imp{*imp}
	request.Cur = []string{"USD"}

	request.User = &openrtb.User{
		Data: []openrtb.Data{
			{
				Segment: []openrtb.Segment{
					{
						Signal: br.Imp.Demands[adapter.MobileFuseKey]["token"].(string),
					},
				},
			},
		},
	}

	return request, nil
}

func (a *MobileFuseAdapter) ExecuteRequest(ctx context.Context, client *http.Client, request openrtb.BidRequest) *adapters.DemandResponse {
	dr := &adapters.DemandResponse{
		DemandID:  adapter.MobileFuseKey,
		RequestID: request.ID,
		TagID:     a.TagID,
	}
	requestBody, err := json.Marshal(request)
	if err != nil {
		dr.Error = err
		return dr
	}
	dr.RawRequest = string(requestBody)

	url := "https://mfx.mobilefuse.com/openrtb?ssp=4020"
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

func (a *MobileFuseAdapter) ParseBids(dr *adapters.DemandResponse) (*adapters.DemandResponse, error) {
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

	var extParam map[string]any
	err = json.Unmarshal(bid.Ext, &extParam)
	if err != nil {
		return dr, err
	}
	signaldata := extParam["signaldata"].(string)

	dr.Bid = &adapters.BidDemandResponse{
		ID:         bid.ID,
		ImpID:      bid.ImpID,
		Price:      bid.Price,
		Payload:    bid.AdM,
		Signaldata: signaldata,
		DemandID:   adapter.MobileFuseKey,
		AdID:       bid.AdID,
		SeatID:     seat.Seat,
		LURL:       bid.LURL,
		NURL:       bid.NURL,
		BURL:       bid.BURL,
	}

	return dr, nil
}

// Builder builds a new instance of the MobileFuse adapter for the given bidder with the given config.
func Builder(cfg adapter.ProcessedConfigsMap, client *http.Client) (*adapters.Bidder, error) {
	mobileFuseCfg := cfg[adapter.MobileFuseKey]

	adpt := &MobileFuseAdapter{
		TagID: mobileFuseCfg["tag_id"].(string),
	}

	bidder := &adapters.Bidder{
		Adapter: adpt,
		Client:  client,
	}

	return bidder, nil
}
