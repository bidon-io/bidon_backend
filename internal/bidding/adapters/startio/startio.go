package startio

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gofrs/uuid/v5"
	"github.com/prebid/openrtb/v19/adcom1"
	"github.com/prebid/openrtb/v19/openrtb2"

	"github.com/bidon-io/bidon-backend/internal/ad"
	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/bidding/adapters"
	"github.com/bidon-io/bidon-backend/internal/bidding/adapters/geo"
	"github.com/bidon-io/bidon-backend/internal/bidding/openrtb"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
)

// Adapter represents the Start.io bidding adapter.
type Adapter struct {
	TagID   string
	AppID   string
	Account string
}

// bannerFormats defines the supported banner formats and their dimensions.
var bannerFormats = map[ad.Format][2]int64{
	ad.BannerFormat:      {320, 50},
	ad.LeaderboardFormat: {728, 90},
	ad.MRECFormat:        {300, 250},
	ad.AdaptiveFormat:    {320, 50},
	ad.EmptyFormat:       {320, 50}, // Default
}

// banner creates a banner impression for the bid request.
func (a *Adapter) banner(auctionRequest *schema.AuctionRequest) *openrtb2.Imp {
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

// interstitial creates an interstitial impression for the bid request.
func (a *Adapter) interstitial(auctionRequest *schema.AuctionRequest) *openrtb2.Imp {
	size := adapters.FullscreenFormats[string(auctionRequest.Device.Type)]
	w, h := size[0], size[1]
	if !auctionRequest.AdObject.IsPortrait() {
		w, h = h, w
	}

	return &openrtb2.Imp{
		Instl: 1,
		Banner: &openrtb2.Banner{
			W:     &w,
			H:     &h,
			BType: []openrtb2.BannerAdType{},
			BAttr: []adcom1.CreativeAttribute{},
			Pos:   adcom1.PositionFullScreen.Ptr(),
		},
		Video: &openrtb2.Video{
			W:     w,
			H:     h,
			Pos:   adcom1.PositionFullScreen.Ptr(),
			MIMEs: []string{"video/mp4", "video/3gpp", "video/3gpp2", "video/x-m4v", "video/quicktime"},
		},
	}
}

// rewarded creates a rewarded video impression for the bid request.
func (a *Adapter) rewarded(auctionRequest *schema.AuctionRequest) *openrtb2.Imp {
	size := adapters.FullscreenFormats[string(auctionRequest.Device.Type)]
	w, h := size[0], size[1]
	if !auctionRequest.AdObject.IsPortrait() {
		w, h = h, w
	}

	skip := int8(1)

	return &openrtb2.Imp{
		Instl: 1,
		Rwdd:  1,
		Banner: &openrtb2.Banner{
			W:     &w,
			H:     &h,
			BType: []openrtb2.BannerAdType{},
			BAttr: []adcom1.CreativeAttribute{16},
			Pos:   adcom1.PositionFullScreen.Ptr(),
		},
		Video: &openrtb2.Video{
			W:         w,
			H:         h,
			BAttr:     []adcom1.CreativeAttribute{1, 2, 5, 8, 9, 14, 17},
			Pos:       adcom1.PositionFullScreen.Ptr(),
			MIMEs:     []string{"video/mp4", "video/x-m4v", "video/quicktime", "video/mpeg", "video/avi"},
			Protocols: []adcom1.MediaCreativeSubtype{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14},
			Skip:      &skip,
		},
	}
}

// CreateRequest implements the BidderInterface.CreateRequest method.
func (a *Adapter) CreateRequest(request openrtb.BidRequest, auctionRequest *schema.AuctionRequest) (openrtb.BidRequest, error) {
	if a.TagID == "" {
		return request, errors.New("startio tag ID is empty")
	}

	if a.Account == "" {
		return request, errors.New("startio account is empty")
	}

	if a.AppID == "" {
		return request, errors.New("startio app ID is empty")
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

	impID, _ := uuid.NewV4()
	imp.ID = impID.String()
	imp.TagID = a.TagID
	imp.DisplayManager = string(adapter.StartIOKey)
	if info, ok := auctionRequest.Adapters[adapter.StartIOKey]; ok {
		imp.DisplayManagerVer = info.SDKVersion
	}
	imp.Secure = &secure
	imp.BidFloor = adapters.CalculatePriceFloor(&request, auctionRequest)
	imp.BidFloorCur = "USD"

	request.Imp = []openrtb2.Imp{*imp}
	request.Cur = []string{"USD"}

	if auctionRequest.Test {
		request.Test = 1
	}

	if request.App == nil {
		request.App = &openrtb2.App{}
	}
	request.App.ID = a.AppID
	request.App.Publisher = &openrtb2.Publisher{}

	demandData, ok := auctionRequest.AdObject.Demands[adapter.StartIOKey]
	if !ok {
		return request, errors.New("startio demand data missing")
	}

	token, ok := demandData["token"].(string)
	if !ok || token == "" {
		return request, errors.New("startio token is empty")
	}

	request.User = &openrtb.User{
		BuyerUID: token,
	}

	return request, nil
}

// ExecuteRequest implements the BidderInterface.ExecuteRequest method.
func (a *Adapter) ExecuteRequest(ctx context.Context, client *http.Client, request openrtb.BidRequest) *adapters.DemandResponse {
	dr := &adapters.DemandResponse{
		DemandID:  adapter.StartIOKey,
		RequestID: request.ID,
		TagID:     a.TagID,
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		dr.Error = err
		return dr
	}

	dr.RawRequest = string(requestBody)

	alpha3 := ""
	if request.Device != nil && request.Device.Geo != nil {
		alpha3 = request.Device.Geo.Country
	}

	endpoint := endpointByRegion(alpha3)
	if endpoint == "" {
		dr.Error = errors.New("startio endpoint is empty")
		return dr
	}

	parsedURL, err := url.Parse(endpoint)
	if err != nil {
		dr.Error = fmt.Errorf("parse endpoint: %w", err)
		return dr
	}

	query := parsedURL.Query()
	query.Set("account", a.Account)
	if request.Test == 1 {
		query.Set("testAdsEnabled", "true")
	}
	parsedURL.RawQuery = query.Encode()

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, parsedURL.String(), bytes.NewBuffer(requestBody))
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

// ParseBids implements the BidderInterface.ParseBids method.
func (a *Adapter) ParseBids(dr *adapters.DemandResponse) (*adapters.DemandResponse, error) {
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
		// proceed
	default:
		return dr, fmt.Errorf("unexpected status code: %s", strconv.Itoa(dr.Status))
	}

	var bidResponse openrtb2.BidResponse
	if err := json.Unmarshal([]byte(dr.RawResponse), &bidResponse); err != nil {
		return dr, err
	}

	if len(bidResponse.SeatBid) == 0 || len(bidResponse.SeatBid[0].Bid) == 0 {
		return dr, errors.New("no seatbid or bid in response")
	}

	seat := bidResponse.SeatBid[0]
	bid := seat.Bid[0]

	dr.Bid = &adapters.BidDemandResponse{
		ID:       bid.ID,
		ImpID:    bid.ImpID,
		Price:    bid.Price,
		Payload:  bid.AdM,
		DemandID: adapter.StartIOKey,
		AdID:     bid.AdID,
		SeatID:   seat.Seat,
		LURL:     bid.LURL,
		NURL:     bid.NURL,
		BURL:     bid.BURL,
	}

	return dr, nil
}

// Builder constructs a bidder for Start.io based on processed configuration.
func Builder(cfg adapter.ProcessedConfigsMap, client *http.Client) (*adapters.Bidder, error) {
	startioCfg := cfg[adapter.StartIOKey]

	tagID, _ := startioCfg["tag_id"].(string)
	appID, _ := startioCfg["app_id"].(string)
	account, _ := startioCfg["account"].(string)

	adpt := &Adapter{
		TagID:   tagID,
		AppID:   appID,
		Account: account,
	}

	return &adapters.Bidder{
		Adapter: adpt,
		Client:  client,
	}, nil
}

func endpointByRegion(alpha3 string) string {
	switch geo.Region(alpha3) {
	case "asia":
		return "http://sin-trp-rtb.startappnetwork.com/1.3/2.5/getbid"
	case "eu":
		return "http://eu-trp-rtb.startappnetwork.com/1.3/2.5/getbid"
	default:
		return "http://trp-rtb.startappnetwork.com/1.3/2.5/getbid"
	}
}
