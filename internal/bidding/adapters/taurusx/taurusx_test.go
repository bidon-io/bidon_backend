package taurusx

import (
	"encoding/json"
	"net/http"
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

	// Verify token is in request extension
	var reqExt map[string]interface{}
	if err := json.Unmarshal(request.Ext, &reqExt); err != nil {
		t.Errorf("Failed to unmarshal request extension: %v", err)
		return
	}

	if token, ok := reqExt["token"].(string); !ok || token != "test-placement-specific-token" {
		t.Errorf("Expected token 'test-placement-specific-token' in request extension, got '%v'", reqExt["token"])
	}

	// Verify app ID is set correctly
	if request.App == nil || request.App.ID != "test-app-id" {
		t.Errorf("Expected app.id 'test-app-id', got '%v'", request.App.ID)
	}
}

func TestTaurusXAdapter_ParseBids_Success(t *testing.T) {
	taurusxAdapter := buildAdapter()

	dr := &adapters.DemandResponse{
		Status:      http.StatusOK,
		RawResponse: `{"id":"test-response","seatbid":[{"bid":[{"id":"test-bid","impid":"test-imp","price":1.5,"adm":"<html>test ad</html>","adid":"test-ad","nurl":"http://win.url","lurl":"http://loss.url"}],"seat":"test-seat"}]}`,
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

	// Should still work without token
	if len(request.Imp) != 1 {
		t.Errorf("Expected 1 impression, got %d", len(request.Imp))
	}

	// Verify request extension has no token when not provided
	var reqExt map[string]interface{}
	if err := json.Unmarshal(request.Ext, &reqExt); err != nil {
		t.Errorf("Failed to unmarshal request extension: %v", err)
		return
	}

	if _, hasToken := reqExt["token"]; hasToken {
		t.Error("Expected no token in request extension when token is not provided")
	}

	// Verify app ID is still set correctly
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
			expected: "https://sdkus.ssp.taxssp.com/ssp/v1/bidding_ad/bidon",
		},
		{
			name:     "US region - Canada",
			alpha3:   "CAN",
			expected: "https://sdkus.ssp.taxssp.com/ssp/v1/bidding_ad/bidon",
		},
		{
			name:     "EU region - Germany",
			alpha3:   "DEU",
			expected: "https://sdkeu.ssp.taxssp.com/ssp/v1/bidding_ad/bidon",
		},
		{
			name:     "EU region - United Kingdom",
			alpha3:   "GBR",
			expected: "https://sdkeu.ssp.taxssp.com/ssp/v1/bidding_ad/bidon",
		},
		{
			name:     "Asia region - Singapore",
			alpha3:   "SGP",
			expected: "https://sdksg.ssp.taxssp.com/ssp/v1/bidding_ad/bidon",
		},
		{
			name:     "Asia region - Japan",
			alpha3:   "JPN",
			expected: "https://sdksg.ssp.taxssp.com/ssp/v1/bidding_ad/bidon",
		},
		{
			name:     "Unknown country - defaults to US",
			alpha3:   "XXX",
			expected: "https://sdkus.ssp.taxssp.com/ssp/v1/bidding_ad/bidon",
		},
		{
			name:     "Empty country - defaults to US",
			alpha3:   "",
			expected: "https://sdkus.ssp.taxssp.com/ssp/v1/bidding_ad/bidon",
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
					"token": `{"different-placement-id":"some-token"}`, // Token for different placement
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

	// Should still work without token for this placement
	if len(request.Imp) != 1 {
		t.Errorf("Expected 1 impression, got %d", len(request.Imp))
	}

	// Verify request extension has no token when placement token is not found
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
