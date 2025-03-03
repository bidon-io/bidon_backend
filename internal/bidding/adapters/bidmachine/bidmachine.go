package bidmachine

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

type BidmachineAdapter struct {
	SellerID string
	Endpoint string
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
	"TABLET": {768, 1024},
}

func (a *BidmachineAdapter) banner(br *schema.BiddingRequest) *openrtb2.Imp {
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

func (a *BidmachineAdapter) interstitial(br *schema.BiddingRequest) *openrtb2.Imp {
	size := fullscreenFormats[string(br.Device.Type)]
	w, h := size[0], size[1]
	if !br.Imp.IsPortrait() {
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

func (a *BidmachineAdapter) rewarded(br *schema.BiddingRequest) *openrtb2.Imp {
	size := fullscreenFormats[string(br.Device.Type)]
	w, h := size[0], size[1]
	if !br.Imp.IsPortrait() {
		w, h = h, w
	}
	return &openrtb2.Imp{
		Instl: 1,
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
			BAttr:     []adcom1.CreativeAttribute{16},
			Pos:       adcom1.PositionFullScreen.Ptr(),
			MIMEs:     []string{"video/mp4", "video/x-m4v", "video/quicktime", "video/mpeg", "video/avi"},
			Protocols: []adcom1.MediaCreativeSubtype{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14},
		},
		Ext: json.RawMessage(`{"rewarded": 1}`),
	}
}

func (a *BidmachineAdapter) CreateRequest(request openrtb.BidRequest, br *schema.BiddingRequest) (openrtb.BidRequest, error) {
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
	imp.DisplayManager = string(adapter.BidmachineKey)
	imp.DisplayManagerVer = br.Adapters[adapter.BidmachineKey].SDKVersion
	imp.Secure = &secure
	imp.BidFloor = adapters.CalculatePriceFloor(&request, br)

	request.App.Publisher.ID = a.SellerID

	extStructure := &map[string]interface{}{}
	_ = json.Unmarshal(imp.Ext, extStructure)

	(*extStructure)["bid_token"] = br.Imp.Demands[adapter.BidmachineKey]["token"]

	raw, _ := json.Marshal(extStructure)

	imp.Ext = raw

	request.Imp = []openrtb2.Imp{*imp}
	request.Cur = []string{"USD"}

	return request, nil
}

func (a *BidmachineAdapter) ExecuteRequest(ctx context.Context, client *http.Client, request openrtb.BidRequest) *adapters.DemandResponse {
	dr := &adapters.DemandResponse{
		DemandID:  adapter.BidmachineKey,
		RequestID: request.ID,
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
	url := getEndpoint(alpha3)
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

func (a *BidmachineAdapter) ParseBids(dr *adapters.DemandResponse) (*adapters.DemandResponse, error) {
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
		DemandID: adapter.BidmachineKey,
		AdID:     bid.AdID,
		SeatID:   seat.Seat,
		LURL:     bid.LURL,
		NURL:     bid.NURL,
		BURL:     bid.BURL,
	}

	return dr, nil
}

// Builder builds a new instance of the Bidmachine adapter for the given bidder with the given config.
func Builder(cfg adapter.ProcessedConfigsMap, client *http.Client) (*adapters.Bidder, error) {
	bmCfg := cfg[adapter.BidmachineKey]
	endpoint, ok := bmCfg["endpoint"].(string)
	if !ok || endpoint == "" {
		return nil, fmt.Errorf("missing endpoint param for BM adapter")
	}
	sellerID, ok := bmCfg["seller_id"].(string)
	if !ok || sellerID == "" {
		return nil, fmt.Errorf("missing seller_id param for %s adapter", adapter.BidmachineKey)
	}

	adpt := &BidmachineAdapter{
		Endpoint: endpoint,
		SellerID: sellerID,
	}

	bidder := &adapters.Bidder{
		Adapter: adpt,
		Client:  client,
	}

	return bidder, nil
}

var alpha3ToDcMapping = map[string]string{
	"AFG": "apac",
	"AUS": "apac",
	"BHR": "apac",
	"BGD": "apac",
	"BTN": "apac",
	"BRN": "apac",
	"KHM": "apac",
	"CHN": "apac",
	"TLS": "apac",
	"FJI": "apac",
	"HKG": "apac",
	"IND": "apac",
	"IDN": "apac",
	"IRN": "apac",
	"IRQ": "apac",
	"ISR": "apac",
	"JPN": "apac",
	"JOR": "apac",
	"KAZ": "apac",
	"KGZ": "apac",
	"KWT": "apac",
	"LAO": "apac",
	"LBN": "apac",
	"MYS": "apac",
	"MDV": "apac",
	"MNG": "apac",
	"MMR": "apac",
	"NPL": "apac",
	"NZL": "apac",
	"PRK": "apac",
	"OMN": "apac",
	"PAK": "apac",
	"PNG": "apac",
	"PHL": "apac",
	"QAT": "apac",
	"SAU": "apac",
	"SGP": "apac",
	"KOR": "apac",
	"LKA": "apac",
	"SYR": "apac",
	"TWN": "apac",
	"TJK": "apac",
	"THA": "apac",
	"TKM": "apac",
	"ARE": "apac",
	"UZB": "apac",
	"VNM": "apac",
	"YEM": "apac",
	"AGL": "us",
	"ARG": "us",
	"BHS": "us",
	"BRB": "us",
	"BLZ": "us",
	"BMU": "us",
	"BOL": "us",
	"BRA": "us",
	"CAN": "us",
	"CYM": "us",
	"CHL": "us",
	"COL": "us",
	"CRI": "us",
	"CUB": "us",
	"DMA": "us",
	"DOM": "us",
	"ECU": "us",
	"SLV": "us",
	"GRD": "us",
	"GTM": "us",
	"GUY": "us",
	"HTI": "us",
	"HND": "us",
	"JAM": "us",
	"MEX": "us",
	"NIC": "us",
	"PAN": "us",
	"PRY": "us",
	"PER": "us",
	"PRI": "us",
	"KNA": "us",
	"LCA": "us",
	"VCT": "us",
	"SUR": "us",
	"TTO": "us",
	"URY": "us",
	"USA": "us",
	"VEN": "us",
}

const defaultDc = "eu"

func getEndpoint(alpha3 string) string {
	dc := defaultDc
	if rewrittenDc, ok := alpha3ToDcMapping[alpha3]; ok {
		dc = rewrittenDc
	}

	return "https://api-" + dc + ".bidmachine.io/auction/prebid/bidon"
}
