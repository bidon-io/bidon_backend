package taurusx

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/prebid/openrtb/v19/openrtb2"

	"github.com/bidon-io/bidon-backend/internal/ad"
	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/bidding/adapters"
	"github.com/bidon-io/bidon-backend/internal/bidding/openrtb"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
)

func TestTaurusXAdapter_CreateRequest(t *testing.T) {
	taurusxAdapter := buildAdapter()

	auctionRequest := &schema.AuctionRequest{
		AdObject: schema.AdObject{
			AuctionID: "test-auction-id",
			Banner: &schema.BannerAdObject{
				Format: ad.BannerFormat,
			},
			Demands: map[adapter.Key]map[string]any{
				adapter.TaurusXKey: {
					"token": `{"test-tag-id":"test-placement-specific-token"}`,
				},
			},
		},
		Adapters: schema.Adapters{
			adapter.TaurusXKey: schema.Adapter{
				Version:    "1.0.0",
				SDKVersion: "1.0.0",
			},
		},
	}

	baseRequest := openrtb.BidRequest{
		ID:  "test-request-id",
		App: &openrtb2.App{},
	}

	request, err := taurusxAdapter.CreateRequest(baseRequest, auctionRequest)
	if err != nil {
		t.Errorf("CreateRequest() error = %v", err)
		return
	}

	if len(request.Imp) != 1 {
		t.Errorf("Expected 1 impression, got %d", len(request.Imp))
	}

	imp := request.Imp[0]
	if imp.TagID != "test-tag-id" {
		t.Errorf("Expected TagID 'test-tag-id', got '%s'", imp.TagID)
	}

	if imp.BidFloorCur != "USD" {
		t.Errorf("Expected BidFloorCur 'USD', got '%s'", imp.BidFloorCur)
	}

	if len(request.Cur) != 1 || request.Cur[0] != "USD" {
		t.Errorf("Expected currency 'USD', got %v", request.Cur)
	}

	var reqExt map[string]interface{}
	if err := json.Unmarshal(request.Ext, &reqExt); err != nil {
		t.Errorf("Failed to unmarshal request extension: %v", err)
		return
	}

	if token, ok := reqExt["token"].(string); !ok || token != "test-placement-specific-token" {
		t.Errorf("Expected token 'test-placement-specific-token' in request extension, got '%v'", reqExt["token"])
	}

	if request.App == nil || request.App.ID != "test-app-id" {
		t.Errorf("Expected app.id 'test-app-id', got '%v'", request.App.ID)
	}
}

func TestTaurusXAdapter_ParseBids_Success(t *testing.T) {
	taurusxAdapter := buildAdapter()

	dr := &adapters.DemandResponse{
		Status:      http.StatusOK,
		RawResponse: `{"id":"test-response","seatbid":[{"bid":[{"id":"test-bid","impid":"test-imp","price":1.5,"adid":"test-ad","nurl":"http://win.url","lurl":"http://loss.url"}],"seat":"test-seat"}],"ext":{"payload":"<html>test ad</html>"}}`,
	}

	result, err := taurusxAdapter.ParseBids(dr)
	if err != nil {
		t.Errorf("ParseBids() error = %v", err)
		return
	}

	if result.Bid == nil {
		t.Error("Expected bid to be present")
		return
	}

	if result.Bid.Price != 1.5 {
		t.Errorf("Expected price 1.5, got %f", result.Bid.Price)
	}

	if result.Bid.Payload != "<html>test ad</html>" {
		t.Errorf("Expected payload '<html>test ad</html>', got '%s'", result.Bid.Payload)
	}

	if result.Bid.NURL != "http://win.url" {
		t.Errorf("Expected NURL 'http://win.url', got '%s'", result.Bid.NURL)
	}

	if result.Bid.LURL != "http://loss.url" {
		t.Errorf("Expected LURL 'http://loss.url', got '%s'", result.Bid.LURL)
	}
}

func TestTaurusXAdapter_ParseBids_NoContent(t *testing.T) {
	taurusxAdapter := buildAdapter()

	dr := &adapters.DemandResponse{
		Status: http.StatusNoContent,
	}

	result, err := taurusxAdapter.ParseBids(dr)
	if err != nil {
		t.Errorf("ParseBids() error = %v", err)
		return
	}

	if result.Bid != nil {
		t.Error("Expected no bid for no content response")
	}
}

func TestTaurusXAdapter_ParseBids_Error(t *testing.T) {
	taurusxAdapter := buildAdapter()

	dr := &adapters.DemandResponse{
		Status: http.StatusBadRequest,
	}

	_, err := taurusxAdapter.ParseBids(dr)
	if err == nil {
		t.Error("Expected error for bad request status")
	}
}

func TestTaurusXAdapter_ParseBids_WithoutPayload(t *testing.T) {
	taurusxAdapter := buildAdapter()

	dr := &adapters.DemandResponse{
		Status:      http.StatusOK,
		RawResponse: `{"id":"test-response","seatbid":[{"bid":[{"id":"test-bid","impid":"test-imp","price":1.5,"adid":"test-ad","nurl":"http://win.url","lurl":"http://loss.url"}],"seat":"test-seat"}]}`,
	}

	result, err := taurusxAdapter.ParseBids(dr)
	if err != nil {
		t.Errorf("ParseBids() error = %v", err)
		return
	}

	if result.Bid == nil {
		t.Error("Expected bid to be present")
		return
	}

	if result.Bid.Payload != "" {
		t.Errorf("Expected empty payload when ext.payload is not present, got '%s'", result.Bid.Payload)
	}
}

func TestTaurusXAdapter_ParseBids_InvalidExtJSON(t *testing.T) {
	taurusxAdapter := buildAdapter()

	dr := &adapters.DemandResponse{
		Status:      http.StatusOK,
		RawResponse: `{"id":"test-response","seatbid":[{"bid":[{"id":"test-bid","impid":"test-imp","price":1.5,"adid":"test-ad"}],"seat":"test-seat"}],"ext":{"invalid":"json"}}`,
	}

	// This should work fine since the ext JSON is valid, just doesn't have payload
	result, err := taurusxAdapter.ParseBids(dr)
	if err != nil {
		t.Errorf("ParseBids() error = %v", err)
		return
	}

	if result.Bid == nil {
		t.Error("Expected bid to be present")
		return
	}

	if result.Bid.Payload != "" {
		t.Errorf("Expected empty payload when ext.payload is not present, got '%s'", result.Bid.Payload)
	}
}

func TestTaurusXAdapter_ParseBids_MalformedBidResponseJSON(t *testing.T) {
	taurusxAdapter := buildAdapter()

	// Create a response with malformed JSON in the overall response
	dr := &adapters.DemandResponse{
		Status:      http.StatusOK,
		RawResponse: `{"id":"test-response","seatbid":[{"bid":[{"id":"test-bid","impid":"test-imp","price":1.5,"adid":"test-ad"}],"seat":"test-seat"}],"ext":{"malformed":json}}`,
	}

	_, err := taurusxAdapter.ParseBids(dr)
	if err == nil {
		t.Error("Expected error for malformed JSON")
		return
	}

	if !strings.Contains(err.Error(), "invalid character") {
		t.Errorf("Expected error message to contain 'invalid character', got: %v", err)
	}
}

func TestTaurusXAdapter_ParseBids_MalformedExtJSON(t *testing.T) {
	taurusxAdapter := buildAdapter()

	dr := &adapters.DemandResponse{
		Status:      http.StatusOK,
		RawResponse: `{"id":"test-response","seatbid":[{"bid":[{"id":"test-bid","impid":"test-imp","price":1.5,"adid":"test-ad"}],"seat":"test-seat"}],"ext":"{\"malformed\":json}"}`,
	}

	_, err := taurusxAdapter.ParseBids(dr)
	if err == nil {
		t.Error("Expected error for malformed ext JSON")
		return
	}

	if !strings.Contains(err.Error(), "failed to unmarshal bid response ext") {
		t.Errorf("Expected error message to contain 'failed to unmarshal bid response ext', got: %v", err)
	}
}

func TestTaurusXAdapter_ParseBids_PayloadNotString(t *testing.T) {
	taurusxAdapter := buildAdapter()

	dr := &adapters.DemandResponse{
		Status:      http.StatusOK,
		RawResponse: `{"id":"test-response","seatbid":[{"bid":[{"id":"test-bid","impid":"test-imp","price":1.5,"adid":"test-ad"}],"seat":"test-seat"}],"ext":{"payload":123}}`,
	}

	result, err := taurusxAdapter.ParseBids(dr)
	if err != nil {
		t.Errorf("ParseBids() error = %v", err)
		return
	}

	if result.Bid == nil {
		t.Error("Expected bid to be present")
		return
	}

	if result.Bid.Payload != "" {
		t.Errorf("Expected empty payload when ext.payload is not a string, got '%s'", result.Bid.Payload)
	}
}

func TestTaurusXAdapter_ParseBids_RealTaurusXResponse(t *testing.T) {
	taurusxAdapter := buildAdapter()

	dr := &adapters.DemandResponse{
		Status: http.StatusOK,
		RawResponse: `{
			"id":"ea3abdd1-191b-4b30-9b0b-f35276bd3bc9",
			"bidid":"ea3abdd1-191b-4b30-9b0b-f35276bd3bc9",
			"cur":"USD",
			"seatbid":[
				{
					"seat":"taurusx",
					"bid":[
						{
							"id":"1",
							"impid":"1",
							"price":0.11654588088605387,
							"w":320,
							"h":480,
							"lurl":"https://notice-eu.ssp.taxssp.com/v1/loss",
							"nurl":"https://notice-eu.ssp.taxssp.com/v1/win",
							"bundle":"",
							"cid":"a12933122763264",
							"adid":"",
							"crid":"2r_m13840221289921",
							"cat":[],
							"adomain":["allegro.pl"],
							"attr":[],
							"ext":{}
						}
					]
				}
			],
			"ext":{
				"payload":"SotHTdNillA6/ljMOKkJ6cVyKmNvYVamCjtWqgJtRzV9IMwyR5BoBmZujUzf7DfD"
			}
		}`,
	}

	result, err := taurusxAdapter.ParseBids(dr)
	if err != nil {
		t.Errorf("ParseBids() error = %v", err)
		return
	}

	if result.Bid == nil {
		t.Error("Expected bid to be present")
		return
	}

	if result.Bid.Price != 0.11654588088605387 {
		t.Errorf("Expected price 0.11654588088605387, got %f", result.Bid.Price)
	}

	expectedPayload := "SotHTdNillA6/ljMOKkJ6cVyKmNvYVamCjtWqgJtRzV9IMwyR5BoBmZujUzf7DfD"
	if result.Bid.Payload != expectedPayload {
		t.Errorf("Expected payload '%s', got '%s'", expectedPayload, result.Bid.Payload)
	}

	if result.Bid.ID != "1" {
		t.Errorf("Expected bid ID '1', got '%s'", result.Bid.ID)
	}

	if result.Bid.ImpID != "1" {
		t.Errorf("Expected imp ID '1', got '%s'", result.Bid.ImpID)
	}

	if result.Bid.SeatID != "taurusx" {
		t.Errorf("Expected seat ID 'taurusx', got '%s'", result.Bid.SeatID)
	}

	if result.Bid.NURL != "https://notice-eu.ssp.taxssp.com/v1/win" {
		t.Errorf("Expected NURL 'https://notice-eu.ssp.taxssp.com/v1/win', got '%s'", result.Bid.NURL)
	}

	if result.Bid.LURL != "https://notice-eu.ssp.taxssp.com/v1/loss" {
		t.Errorf("Expected LURL 'https://notice-eu.ssp.taxssp.com/v1/loss', got '%s'", result.Bid.LURL)
	}
}

func TestBuilder_Success(t *testing.T) {
	cfg := adapter.ProcessedConfigsMap{
		adapter.TaurusXKey: map[string]any{
			"app_id": "test-app-id",
			"tag_id": "test-tag-id",
		},
	}

	bidder, err := Builder(cfg, &http.Client{})
	if err != nil {
		t.Errorf("Builder() error = %v", err)
		return
	}

	if bidder == nil {
		t.Error("Expected bidder to be created")
		return
	}

	taurusxAdapter, ok := bidder.Adapter.(*TaurusXAdapter)
	if !ok {
		t.Error("Expected TaurusXAdapter")
		return
	}

	if taurusxAdapter.AppID != "test-app-id" {
		t.Errorf("Expected AppID 'test-app-key', got '%s'", taurusxAdapter.AppID)
	}

	if taurusxAdapter.TagID != "test-tag-id" {
		t.Errorf("Expected TagID 'test-tag-id', got '%s'", taurusxAdapter.TagID)
	}
}

func TestBuilder_MissingAppKey(t *testing.T) {
	cfg := adapter.ProcessedConfigsMap{
		adapter.TaurusXKey: map[string]any{
			"api_key": "test-api-key",
		},
	}

	_, err := Builder(cfg, &http.Client{})
	if err == nil {
		t.Error("Expected error for missing app_key")
	}
}

func TestBuilder_MissingAPIKey(t *testing.T) {
	cfg := adapter.ProcessedConfigsMap{
		adapter.TaurusXKey: map[string]any{
			"app_key": "test-app-key",
		},
	}

	_, err := Builder(cfg, &http.Client{})
	if err == nil {
		t.Error("Expected error for missing api_key")
	}
}

func TestTaurusXAdapter_CreateRequest_Interstitial(t *testing.T) {
	taurusxAdapter := buildAdapter()

	auctionRequest := &schema.AuctionRequest{
		AdObject: schema.AdObject{
			AuctionID:    "test-auction-id",
			Interstitial: &schema.InterstitialAdObject{},
			Demands: map[adapter.Key]map[string]any{
				adapter.TaurusXKey: {
					"token": `{"test-tag-id":"test-placement-specific-token"}`,
				},
			},
		},
		Adapters: schema.Adapters{
			adapter.TaurusXKey: schema.Adapter{
				Version:    "1.0.0",
				SDKVersion: "1.0.0",
			},
		},
	}

	baseRequest := openrtb.BidRequest{
		ID:  "test-request-id",
		App: &openrtb2.App{},
	}

	request, err := taurusxAdapter.CreateRequest(baseRequest, auctionRequest)
	if err != nil {
		t.Errorf("CreateRequest() error = %v", err)
		return
	}

	if len(request.Imp) != 1 {
		t.Errorf("Expected 1 impression, got %d", len(request.Imp))
	}

	imp := request.Imp[0]
	if imp.Instl != 1 {
		t.Errorf("Expected interstitial impression (Instl=1), got %d", imp.Instl)
	}
}

func TestTaurusXAdapter_CreateRequest_Rewarded(t *testing.T) {
	taurusxAdapter := buildAdapter()

	auctionRequest := &schema.AuctionRequest{
		AdObject: schema.AdObject{
			AuctionID: "test-auction-id",
			Rewarded:  &schema.RewardedAdObject{},
			Demands: map[adapter.Key]map[string]any{
				adapter.TaurusXKey: {
					"token": `{"test-tag-id":"test-placement-specific-token"}`,
				},
			},
		},
		Adapters: schema.Adapters{
			adapter.TaurusXKey: schema.Adapter{
				Version:    "1.0.0",
				SDKVersion: "1.0.0",
			},
		},
	}

	baseRequest := openrtb.BidRequest{
		ID:  "test-request-id",
		App: &openrtb2.App{},
	}

	request, err := taurusxAdapter.CreateRequest(baseRequest, auctionRequest)
	if err != nil {
		t.Errorf("CreateRequest() error = %v", err)
		return
	}

	if len(request.Imp) != 1 {
		t.Errorf("Expected 1 impression, got %d", len(request.Imp))
	}

	imp := request.Imp[0]
	if imp.Instl != 1 {
		t.Errorf("Expected rewarded impression (Instl=1), got %d", imp.Instl)
	}

	if imp.Video == nil {
		t.Error("Expected video object for rewarded ad")
	}
}

func TestTaurusXAdapter_CreateRequest_WithoutToken(t *testing.T) {
	taurusxAdapter := buildAdapter()

	auctionRequest := &schema.AuctionRequest{
		AdObject: schema.AdObject{
			AuctionID: "test-auction-id",
			Banner: &schema.BannerAdObject{
				Format: ad.BannerFormat,
			},
			Demands: map[adapter.Key]map[string]any{
				adapter.TaurusXKey: {},
			},
		},
		Adapters: schema.Adapters{
			adapter.TaurusXKey: schema.Adapter{
				Version:    "1.0.0",
				SDKVersion: "1.0.0",
			},
		},
	}

	baseRequest := openrtb.BidRequest{
		ID:  "test-request-id",
		App: &openrtb2.App{},
	}

	request, err := taurusxAdapter.CreateRequest(baseRequest, auctionRequest)
	if err != nil {
		t.Errorf("CreateRequest() error = %v", err)
		return
	}

	if len(request.Imp) != 1 {
		t.Errorf("Expected 1 impression, got %d", len(request.Imp))
	}

	var reqExt map[string]interface{}
	if err := json.Unmarshal(request.Ext, &reqExt); err != nil {
		t.Errorf("Failed to unmarshal request extension: %v", err)
		return
	}

	if _, hasToken := reqExt["token"]; hasToken {
		t.Error("Expected no token in request extension when token is not provided")
	}

	if request.App == nil || request.App.ID != "test-app-id" {
		t.Errorf("Expected app.id 'test-app-id', got '%v'", request.App.ID)
	}
}

func TestGetEndpoint(t *testing.T) {
	tests := []struct {
		name     string
		alpha3   string
		expected string
	}{
		{
			name:     "US region - USA",
			alpha3:   "USA",
			expected: "http://test-sdk.ssp.taxssp.com/ssp/v1/bidding_ad/testing",
		},
		{
			name:     "US region - Canada",
			alpha3:   "CAN",
			expected: "http://test-sdk.ssp.taxssp.com/ssp/v1/bidding_ad/testing",
		},
		{
			name:     "EU region - Germany",
			alpha3:   "DEU",
			expected: "http://test-sdk.ssp.taxssp.com/ssp/v1/bidding_ad/testing",
		},
		{
			name:     "EU region - United Kingdom",
			alpha3:   "GBR",
			expected: "http://test-sdk.ssp.taxssp.com/ssp/v1/bidding_ad/testing",
		},
		{
			name:     "Asia region - Singapore",
			alpha3:   "SGP",
			expected: "http://test-sdk.ssp.taxssp.com/ssp/v1/bidding_ad/testing",
		},
		{
			name:     "Asia region - Japan",
			alpha3:   "JPN",
			expected: "http://test-sdk.ssp.taxssp.com/ssp/v1/bidding_ad/testing",
		},
		{
			name:     "Unknown country - defaults to US",
			alpha3:   "XXX",
			expected: "http://test-sdk.ssp.taxssp.com/ssp/v1/bidding_ad/testing",
		},
		{
			name:     "Empty country - defaults to US",
			alpha3:   "",
			expected: "http://test-sdk.ssp.taxssp.com/ssp/v1/bidding_ad/testing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getEndpoint(tt.alpha3)
			if result != tt.expected {
				t.Errorf("getEndpoint(%s) = %s, expected %s", tt.alpha3, result, tt.expected)
			}
		})
	}
}

func TestTaurusXAdapter_ExtractPlacementToken(t *testing.T) {
	adapter := buildAdapter()

	tests := []struct {
		name          string
		tokenData     string
		placementID   string
		expectedToken string
		expectError   bool
	}{
		{
			name:          "Valid placement token extraction",
			tokenData:     `{"60001958":"rkeOZxFu2Qs4W+NrBPzPtJ0k6VePE6hXMdGSFje14LvSxAgcKdx5x9HpiWvI4h1L","12345":"another-token"}`,
			placementID:   "60001958",
			expectedToken: "rkeOZxFu2Qs4W+NrBPzPtJ0k6VePE6hXMdGSFje14LvSxAgcKdx5x9HpiWvI4h1L",
			expectError:   false,
		},
		{
			name:          "Placement ID not found",
			tokenData:     `{"60001958":"token1","12345":"token2"}`,
			placementID:   "99999",
			expectedToken: "",
			expectError:   true,
		},
		{
			name:          "Invalid JSON",
			tokenData:     `{"invalid":json}`,
			placementID:   "60001958",
			expectedToken: "",
			expectError:   true,
		},
		{
			name:          "Empty token data",
			tokenData:     "",
			placementID:   "60001958",
			expectedToken: "",
			expectError:   true,
		},
		{
			name:          "Empty placement ID",
			tokenData:     `{"60001958":"token"}`,
			placementID:   "",
			expectedToken: "",
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := adapter.extractPlacementToken(tt.tokenData, tt.placementID)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if token != tt.expectedToken {
					t.Errorf("Expected token '%s', got '%s'", tt.expectedToken, token)
				}
			}
		})
	}
}

func TestTaurusXAdapter_CreateRequest_WithPlacementTokenNotFound(t *testing.T) {
	taurusxAdapter := buildAdapter()

	auctionRequest := &schema.AuctionRequest{
		AdObject: schema.AdObject{
			AuctionID: "test-auction-id",
			Banner: &schema.BannerAdObject{
				Format: ad.BannerFormat,
			},
			Demands: map[adapter.Key]map[string]any{
				adapter.TaurusXKey: {
					"token": `{"different-placement-id":"some-token"}`,
				},
			},
		},
		Adapters: schema.Adapters{
			adapter.TaurusXKey: schema.Adapter{
				Version:    "1.0.0",
				SDKVersion: "1.0.0",
			},
		},
	}

	baseRequest := openrtb.BidRequest{
		ID:  "test-request-id",
		App: &openrtb2.App{},
	}

	request, err := taurusxAdapter.CreateRequest(baseRequest, auctionRequest)
	if err != nil {
		t.Errorf("CreateRequest() error = %v", err)
		return
	}

	if len(request.Imp) != 1 {
		t.Errorf("Expected 1 impression, got %d", len(request.Imp))
	}

	var reqExt map[string]interface{}
	if err := json.Unmarshal(request.Ext, &reqExt); err != nil {
		t.Errorf("Failed to unmarshal request extension: %v", err)
		return
	}

	if _, hasToken := reqExt["token"]; hasToken {
		t.Error("Expected no token in request extension when placement token is not found")
	}
}

func buildAdapter() *TaurusXAdapter {
	return &TaurusXAdapter{
		AppID: "test-app-id",
		TagID: "test-tag-id",
	}
}
