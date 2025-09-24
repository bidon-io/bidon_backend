package taurusx

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofrs/uuid/v5"
	"github.com/prebid/openrtb/v19/adcom1"
	"github.com/prebid/openrtb/v19/openrtb2"

	"github.com/bidon-io/bidon-backend/internal/ad"
	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/bidding/adapters"
	"github.com/bidon-io/bidon-backend/internal/bidding/openrtb"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
)

type TaurusXAdapter struct {
	AppID string
	TagID string
}

var bannerFormats = map[ad.Format][2]int64{
	ad.BannerFormat:   {320, 50},
	ad.MRECFormat:     {300, 250},
	ad.AdaptiveFormat: {320, 50},
	ad.EmptyFormat:    {320, 50}, // Default
}

func (a *TaurusXAdapter) banner(auctionRequest *schema.AuctionRequest) *openrtb2.Imp {
	size, ok := bannerFormats[auctionRequest.AdObject.Format()]
	if !ok {
		size = bannerFormats[ad.EmptyFormat] // Use default
	}

	// Handle adaptive format for tablets
	if auctionRequest.AdObject.IsAdaptive() && auctionRequest.Device.IsTablet() {
		size = [2]int64{728, 90} // Leaderboard format
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

func (a *TaurusXAdapter) interstitial(auctionRequest *schema.AuctionRequest) *openrtb2.Imp {
	size := adapters.FullscreenFormats[string(auctionRequest.Device.Type)]
	w, h := size[0], size[1]
	return &openrtb2.Imp{
		Instl: 1,
		Banner: &openrtb2.Banner{
			W:   &w,
			H:   &h,
			Pos: adcom1.PositionFullScreen.Ptr(),
		},
	}
}

func (a *TaurusXAdapter) rewarded(auctionRequest *schema.AuctionRequest) *openrtb2.Imp {
	size := adapters.FullscreenFormats[string(auctionRequest.Device.Type)]
	w, h := size[0], size[1]
	return &openrtb2.Imp{
		Instl: 1,
		Video: &openrtb2.Video{
			W:         w,
			H:         h,
			Pos:       adcom1.PositionFullScreen.Ptr(),
			MIMEs:     []string{"video/mp4", "video/x-m4v", "video/quicktime", "video/mpeg", "video/avi"},
			Protocols: []adcom1.MediaCreativeSubtype{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14},
		},
	}
}

func (a *TaurusXAdapter) CreateRequest(request openrtb.BidRequest, auctionRequest *schema.AuctionRequest) (openrtb.BidRequest, error) {
	if a.TagID == "" {
		return request, errors.New("TagID is empty")
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

	imp.DisplayManager = string(adapter.TaurusXKey)
	imp.DisplayManagerVer = auctionRequest.Adapters[adapter.TaurusXKey].SDKVersion
	imp.Secure = &secure
	imp.BidFloor = adapters.CalculatePriceFloor(&request, auctionRequest)
	imp.BidFloorCur = "USD"

	request.Imp = []openrtb2.Imp{*imp}
	request.Cur = []string{"USD"}

	// Set app ID
	if request.App != nil {
		request.App.ID = a.AppID
	}

	reqExt := make(map[string]interface{})
	if request.Ext != nil {
		_ = json.Unmarshal(request.Ext, &reqExt)
	}

	if demandData, ok := auctionRequest.AdObject.Demands[adapter.TaurusXKey]; ok {
		if tokenData, ok := demandData["token"].(string); ok && tokenData != "" {
			// Parse the token JSON to extract placement-specific token
			placementToken, err := a.extractPlacementToken(tokenData, a.TagID)
			if err == nil && placementToken != "" {
				reqExt["token"] = placementToken
			}
		}
	}

	extBytes, _ := json.Marshal(reqExt)
	request.Ext = extBytes

	return request, nil
}

func (a *TaurusXAdapter) extractPlacementToken(tokenData, placementID string) (string, error) {
	if tokenData == "" || placementID == "" {
		return "", errors.New("empty token data or placement ID")
	}

	var tokenMap map[string]string
	err := json.Unmarshal([]byte(tokenData), &tokenMap)
	if err != nil {
		return "", fmt.Errorf("failed to parse token JSON: %v", err)
	}

	if placementToken, exists := tokenMap[placementID]; exists {
		return placementToken, nil
	}

	return "", fmt.Errorf("no token found for placement ID: %s", placementID)
}

// logCurlCommand logs the HTTP request as a readable CURL command for debugging purposes
func logCurlCommand(req *http.Request, body []byte) {
	var curlCmd strings.Builder
	curlCmd.WriteString("curl -X ")
	curlCmd.WriteString(req.Method)

	// Add headers
	for name, values := range req.Header {
		for _, value := range values {
			curlCmd.WriteString(" -H '")
			curlCmd.WriteString(name)
			curlCmd.WriteString(": ")
			curlCmd.WriteString(value)
			curlCmd.WriteString("'")
		}
	}

	// Add body if present
	if len(body) > 0 {
		curlCmd.WriteString(" -d '")
		// Escape single quotes in JSON body
		bodyStr := strings.ReplaceAll(string(body), "'", "'\"'\"'")
		curlCmd.WriteString(bodyStr)
		curlCmd.WriteString("'")
	}

	// Add URL
	curlCmd.WriteString(" '")
	curlCmd.WriteString(req.URL.String())
	curlCmd.WriteString("'")

	log.Printf("[TaurusX Debug] CURL Command: %s", curlCmd.String())
}

func (a *TaurusXAdapter) ExecuteRequest(ctx context.Context, client *http.Client, request openrtb.BidRequest) *adapters.DemandResponse {
	dr := &adapters.DemandResponse{
		DemandID:  adapter.TaurusXKey,
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

	// Use geographic endpoint selection
	url := getEndpoint(alpha3)
	if url == "" {
		dr.Error = errors.New("taurusx endpoint is empty")
		return dr
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(requestBody))
	if err != nil {
		dr.Error = err
		return dr
	}
	httpReq.Header.Add("Content-Type", "application/json")
	httpReq.Header.Add("X-OpenRTB-Version", "2.5")

	// Log CURL command for debugging
	logCurlCommand(httpReq, requestBody)

	httpResp, err := client.Do(httpReq)
	if err != nil {
		dr.Error = err
		return dr
	}
	defer httpResp.Body.Close()

	dr.Status = httpResp.StatusCode
	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		dr.Error = err
		return dr
	}
	dr.RawResponse = string(respBody)

	parsedDr, err := a.ParseBids(dr)
	if err != nil {
		dr.Error = err
		return dr
	}
	return parsedDr
}

func (a *TaurusXAdapter) ParseBids(dr *adapters.DemandResponse) (*adapters.DemandResponse, error) {
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
		return dr, nil
	}

	seat := bidResponse.SeatBid[0]
	bid := seat.Bid[0]

	dr.Bid = &adapters.BidDemandResponse{
		ID:       bid.ID,
		ImpID:    bid.ImpID,
		Price:    bid.Price,
		Payload:  bid.AdM,
		DemandID: adapter.TaurusXKey,
		AdID:     bid.AdID,
		SeatID:   seat.Seat,
		LURL:     bid.LURL,
		NURL:     bid.NURL,
		BURL:     bid.BURL,
	}

	return dr, nil
}

// Builder builds a new instance of the TaurusX adapter for the given bidder with the given config.
func Builder(cfg adapter.ProcessedConfigsMap, client *http.Client) (*adapters.Bidder, error) {
	tCfg := cfg[adapter.TaurusXKey]

	appID, ok := tCfg["app_id"].(string)
	if !ok || appID == "" {
		return nil, fmt.Errorf("missing app_id param for %s adapter", adapter.TaurusXKey)
	}
	tagID, ok := tCfg["tag_id"].(string)
	if !ok {
		tagID = ""
	}

	adpt := &TaurusXAdapter{
		AppID: appID,
		TagID: tagID,
	}

	bidder := adapters.Bidder{
		Adapter: adpt,
		Client:  client,
	}

	return &bidder, nil
}

// alpha3ToRegionMapping maps country codes to TaurusX regions
var alpha3ToRegionMapping = map[string]string{
	// US region countries (Americas)
	"ABW": "us", "AIA": "us", "ARG": "us", "ATG": "us", "BES": "us", "BHS": "us", "BLM": "us", "BLZ": "us",
	"BOL": "us", "BRA": "us", "BRB": "us", "CAN": "us", "CHL": "us", "COL": "us", "CRI": "us", "CUB": "us",
	"CUW": "us", "CYM": "us", "DMA": "us", "DOM": "us", "ECU": "us", "GLP": "us", "GRD": "us", "GRL": "us",
	"GTM": "us", "GUF": "us", "GUY": "us", "HND": "us", "HTI": "us", "JAM": "us", "KNA": "us", "LCA": "us",
	"MAF": "us", "MEX": "us", "MSR": "us", "MTQ": "us", "NIC": "us", "PAN": "us", "PER": "us", "PRI": "us",
	"PRY": "us", "SLV": "us", "SUR": "us", "SXM": "us", "TCA": "us", "TST": "us", "TTO": "us", "UMI": "us",
	"URY": "us", "USA": "us", "VCT": "us", "VEN": "us", "VGB": "us", "VIR": "us",

	// Asia region countries
	"AFG": "sg", "ARE": "sg", "ARM": "sg", "ASM": "sg", "ATA": "sg", "ATF": "sg", "AUS": "sg",
	"BGD": "sg", "BHR": "sg", "BRN": "sg", "BTN": "sg", "CCK": "sg", "CHN": "sg", "COK": "sg",
	"COM": "sg", "CXR": "sg", "FJI": "sg", "FSM": "sg", "GUM": "sg", "HKG": "sg", "HMD": "sg",
	"IDN": "sg", "IND": "sg", "IOT": "sg", "IRN": "sg", "IRQ": "sg", "ISR": "sg", "JPN": "sg",
	"KAZ": "sg", "KHM": "sg", "KIR": "sg", "KOR": "sg", "KWT": "sg", "LAO": "sg", "LBN": "sg",
	"LKA": "sg", "MAC": "sg", "MDV": "sg", "MHL": "sg", "MMR": "sg", "MNG": "sg", "MNP": "sg",
	"MYS": "sg", "MYT": "sg", "NCL": "sg", "NFK": "sg", "NIU": "sg", "NPL": "sg", "NRU": "sg",
	"NZL": "sg", "OMN": "sg", "PAK": "sg", "PCN": "sg", "PHL": "sg", "PLW": "sg", "PNG": "sg",
	"PRK": "sg", "PYF": "sg", "QAT": "sg", "SAU": "sg", "SGP": "sg", "SLB": "sg", "SSG": "sg",
	"SYC": "sg", "THA": "sg", "TJK": "sg", "TKL": "sg", "TKM": "sg", "TLS": "sg", "TON": "sg",
	"TUV": "sg", "TWN": "sg", "UZB": "sg", "VNM": "sg", "VUT": "sg", "WLF": "sg", "WSM": "sg",
	"YEM": "sg",

	// EU region countries (Europe, Africa, Middle East)
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

// getEndpoint returns the appropriate TaurusX endpoint based on country code
func getEndpoint(alpha3 string) string {
	// Determine region based on country code
	region := defaultRegion
	if mappedRegion, ok := alpha3ToRegionMapping[alpha3]; ok {
		region = mappedRegion
	}

	// Return the appropriate regional endpoint
	switch region {
	case "eu":
		return "https://sdkeu.ssp.taxssp.com/ssp/v1/bidding_ad/bidon"
	case "sg":
		return "https://sdksg.ssp.taxssp.com/ssp/v1/bidding_ad/bidon"
	case "us":
		fallthrough
	default:
		return "https://sdkus.ssp.taxssp.com/ssp/v1/bidding_ad/bidon"
	}
}
