package moloco

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

// MolocoAdapter represents the Moloco bidding adapter
type MolocoAdapter struct { //nolint:revive
	TagID  string
	AppID  string
	APIKey string
}

// bannerFormats defines the supported banner formats and their dimensions
var bannerFormats = map[ad.Format][2]int64{
	ad.BannerFormat:      {320, 50},
	ad.LeaderboardFormat: {728, 90},
	ad.MRECFormat:        {300, 250},
	ad.AdaptiveFormat:    {320, 50},
	ad.EmptyFormat:       {320, 50}, // Default
}

// banner creates a banner impression for the bid request
func (a *MolocoAdapter) banner(auctionRequest *schema.AuctionRequest) *openrtb2.Imp {
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

// interstitial creates an interstitial impression for the bid request
func (a *MolocoAdapter) interstitial(auctionRequest *schema.AuctionRequest) *openrtb2.Imp {
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

// rewarded creates a rewarded video impression for the bid request
func (a *MolocoAdapter) rewarded(auctionRequest *schema.AuctionRequest) *openrtb2.Imp {
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

// CreateRequest implements the BidderInterface.CreateRequest method
func (a *MolocoAdapter) CreateRequest(request openrtb.BidRequest, auctionRequest *schema.AuctionRequest) (openrtb.BidRequest, error) {
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

	if a.TagID == "" {
		return request, errors.New("moloco AdUnitID is empty")
	}
	imp.TagID = a.TagID

	imp.DisplayManager = string(adapter.MolocoKey)
	imp.DisplayManagerVer = auctionRequest.Adapters[adapter.MolocoKey].SDKVersion
	imp.Secure = &secure
	imp.BidFloor = adapters.CalculatePriceFloor(&request, auctionRequest)
	imp.BidFloorCur = "USD"

	request.Imp = []openrtb2.Imp{*imp}
	request.Cur = []string{"USD"}

	// Set app ID if configured
	if a.AppID != "" && request.App != nil {
		request.App.ID = a.AppID
	}

	// Add user data if token is available
	request.User = &openrtb.User{
		BuyerUID: auctionRequest.AdObject.Demands[adapter.MolocoKey]["token"].(string),
	}

	return request, nil
}

// ExecuteRequest implements the BidderInterface.ExecuteRequest method
func (a *MolocoAdapter) ExecuteRequest(ctx context.Context, client *http.Client, request openrtb.BidRequest) *adapters.DemandResponse {
	dr := &adapters.DemandResponse{
		DemandID:  adapter.MolocoKey,
		RequestID: request.ID,
		TagID:     a.TagID,
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		dr.Error = err
		return dr
	}
	dr.RawRequest = string(requestBody)

	// Get country code for geographic routing
	alpha3 := ""
	if request.Device != nil && request.Device.Geo != nil {
		alpha3 = request.Device.Geo.Country
	}

	// Use geographic endpoint selection or configured endpoint
	url := getEndpoint(alpha3)
	if url == "" {
		dr.Error = errors.New("moloco endpoint is empty")
		return dr
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(requestBody))
	if err != nil {
		dr.Error = err
		return dr
	}
	httpReq.Header.Add("Content-Type", "application/json")
	//httpReq.Header.Add("Accept-Encoding", "gzip")

	// Add Authorization header with API key
	if a.APIKey == "" {
		dr.Error = errors.New("moloco API key is empty")
		return dr
	}
	httpReq.Header.Add("Authorization", a.APIKey)

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

// ParseBids implements the BidderInterface.ParseBids method
func (a *MolocoAdapter) ParseBids(dr *adapters.DemandResponse) (*adapters.DemandResponse, error) {
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
		DemandID: adapter.MolocoKey,
		AdID:     bid.AdID,
		SeatID:   seat.Seat,
		LURL:     bid.LURL,
		NURL:     bid.NURL,
		BURL:     bid.BURL,
	}

	return dr, nil
}

// Builder builds a new instance of the Moloco adapter for the given bidder with the given config.
func Builder(cfg adapter.ProcessedConfigsMap, client *http.Client) (*adapters.Bidder, error) {
	molocoCfg := cfg[adapter.MolocoKey]

	tagID, ok := molocoCfg["tag_id"].(string)
	if !ok {
		tagID = ""
	}

	appID, ok := molocoCfg["app_id"].(string)
	if !ok {
		appID = ""
	}

	apiKey, ok := molocoCfg["api_key"].(string)
	if !ok {
		apiKey = ""
	}

	adpt := &MolocoAdapter{
		TagID:  tagID,
		AppID:  appID,
		APIKey: apiKey,
	}

	bidder := &adapters.Bidder{
		Adapter: adpt,
		Client:  client,
	}

	return bidder, nil
}

// alpha3ToRegionMapping maps country codes to Moloco regions based on country_mapping.json
var alpha3ToRegionMapping = map[string]string{
	// US region countries
	"ABW": "us", "AIA": "us", "ARG": "us", "ATG": "us", "BES": "us", "BHS": "us", "BLM": "us", "BLZ": "us",
	"BOL": "us", "BRA": "us", "BRB": "us", "CAN": "us", "CHL": "us", "COL": "us", "CRI": "us", "CUB": "us",
	"CUW": "us", "CYM": "us", "DMA": "us", "DOM": "us", "ECU": "us", "GLP": "us", "GRD": "us", "GRL": "us",
	"GTM": "us", "GUF": "us", "GUY": "us", "HND": "us", "HTI": "us", "JAM": "us", "KNA": "us", "LCA": "us",
	"MAF": "us", "MEX": "us", "MSR": "us", "MTQ": "us", "NIC": "us", "PAN": "us", "PER": "us", "PRI": "us",
	"PRY": "us", "SLV": "us", "SUR": "us", "SXM": "us", "TCA": "us", "TST": "us", "TTO": "us", "UMI": "us",
	"URY": "us", "USA": "us", "VCT": "us", "VEN": "us", "VGB": "us", "VIR": "us",

	// Asia region countries
	"AFG": "asia", "ARE": "asia", "ARM": "asia", "ASM": "asia", "ATA": "asia", "ATF": "asia", "AUS": "asia",
	"BGD": "asia", "BHR": "asia", "BRN": "asia", "BTN": "asia", "CCK": "asia", "CHN": "asia", "COK": "asia",
	"COM": "asia", "CXR": "asia", "FJI": "asia", "FSM": "asia", "GUM": "asia", "HKG": "asia", "HMD": "asia",
	"IDN": "asia", "IND": "asia", "IOT": "asia", "IRN": "asia", "IRQ": "asia", "ISR": "asia", "JPN": "asia",
	"KAZ": "asia", "KHM": "asia", "KIR": "asia", "KOR": "asia", "KWT": "asia", "LAO": "asia", "LBN": "asia",
	"LKA": "asia", "MAC": "asia", "MDV": "asia", "MHL": "asia", "MMR": "asia", "MNG": "asia", "MNP": "asia",
	"MYS": "asia", "MYT": "asia", "NCL": "asia", "NFK": "asia", "NIU": "asia", "NPL": "asia", "NRU": "asia",
	"NZL": "asia", "OMN": "asia", "PAK": "asia", "PCN": "asia", "PHL": "asia", "PLW": "asia", "PNG": "asia",
	"PRK": "asia", "PYF": "asia", "QAT": "asia", "SAU": "asia", "SGP": "asia", "SLB": "asia", "SSG": "asia",
	"SYC": "asia", "THA": "asia", "TJK": "asia", "TKL": "asia", "TKM": "asia", "TLS": "asia", "TON": "asia",
	"TUV": "asia", "TWN": "asia", "UZB": "asia", "VNM": "asia", "VUT": "asia", "WLF": "asia", "WSM": "asia",
	"YEM": "asia",

	// EU region countries
	"AGO": "eu", "ALA": "eu", "ALB": "eu", "AND": "eu", "AUT": "eu", "AZE": "eu", "BDI": "eu", "BEL": "eu",
	"BEN": "eu", "BFA": "eu", "BGR": "eu", "BIH": "eu", "BLR": "eu", "BMU": "eu", "BVT": "eu", "BWA": "eu",
	"CAF": "eu", "CHE": "eu", "CIV": "eu", "CMR": "eu", "COD": "eu", "COG": "eu", "CPV": "eu", "CYP": "eu",
	"CZE": "eu", "DEU": "eu", "DJI": "eu", "DNK": "eu", "DZA": "eu", "EGY": "eu", "ERI": "eu", "ESH": "eu",
	"ESP": "eu", "EST": "eu", "ETH": "eu", "FIN": "eu", "FRA": "eu", "FRO": "eu", "GAB": "eu", "GBR": "eu",
	"GEO": "eu", "GGY": "eu", "GHA": "eu", "GIB": "eu", "GIN": "eu", "GMB": "eu", "GNB": "eu", "GNQ": "eu",
	"GRC": "eu", "HRV": "eu", "HUN": "eu", "IMN": "eu", "IRL": "eu", "ISL": "eu", "ITA": "eu", "JEY": "eu",
	"JOR": "eu", "KEN": "eu", "KGZ": "eu", "LBR": "eu", "LBY": "eu", "LIE": "eu", "LSO": "eu", "LTU": "eu",
	"LUX": "eu", "LVA": "eu", "MAR": "eu", "MCO": "eu", "MDA": "eu", "MDG": "eu", "MKD": "eu", "MLI": "eu",
	"MLT": "eu", "MNE": "eu", "MOZ": "eu", "MRT": "eu", "MUS": "eu", "MWI": "eu", "NAM": "eu", "NER": "eu",
	"NGA": "eu", "NLD": "eu", "NOR": "eu", "POL": "eu", "PRT": "eu", "PSE": "eu", "REU": "eu", "ROU": "eu",
	"RUS": "eu", "RWA": "eu", "SDN": "eu", "SEN": "eu", "SGS": "eu", "SHN": "eu", "SJM": "eu", "SLE": "eu",
	"SMR": "eu", "SOM": "eu", "SRB": "eu", "SSD": "eu", "STP": "eu", "SVK": "eu", "SVN": "eu", "SWE": "eu",
	"SWZ": "eu", "SYR": "eu", "TCD": "eu", "TGO": "eu", "TUN": "eu", "TUR": "eu", "TZA": "eu", "UGA": "eu",
	"UKR": "eu", "VAT": "eu", "ZAF": "eu", "ZMB": "eu", "ZWE": "eu",
}

const defaultRegion = "us"

// getEndpoint returns the appropriate Moloco endpoint based on country code and fallback endpoint
func getEndpoint(alpha3 string) string {
	// Determine region based on country code
	region := defaultRegion
	if mappedRegion, ok := alpha3ToRegionMapping[alpha3]; ok {
		region = mappedRegion
	}

	return "https://sdkfnt-" + region + ".dsp-api.moloco.com/mediations/inhouse/v1"
}
