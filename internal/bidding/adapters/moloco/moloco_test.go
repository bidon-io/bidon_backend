package moloco_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/prebid/openrtb/v19/openrtb2"

	"github.com/bidon-io/bidon-backend/internal/ad"
	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/bidding/adapters"
	"github.com/bidon-io/bidon-backend/internal/bidding/adapters/moloco"
	"github.com/bidon-io/bidon-backend/internal/bidding/openrtb"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
)

type createRequestTestParams struct {
	BaseBidRequest openrtb.BidRequest
	AuctionRequest *schema.AuctionRequest
}

type createRequestTestOutput struct {
	Request openrtb.BidRequest
	Err     error
}

type ParseBidsTestParams struct {
	DemandsResponse adapters.DemandResponse
}

type ParseBidsTestOutput struct {
	DemandResponse adapters.DemandResponse
	Err            error
}

type TestTransport func(req *http.Request) *http.Response

func (f TestTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func NewTestClient(tr TestTransport) *http.Client {
	return &http.Client{
		Transport: tr,
	}
}

func buildAdapter() moloco.MolocoAdapter {
	return moloco.MolocoAdapter{
		TagID:  "test-tag-id",
		AppID:  "test-app-id",
		APIKey: "test-api-key",
	}
}

func buildBaseBidRequest() openrtb.BidRequest {
	return openrtb.BidRequest{
		ID: "test-request-id",
		App: &openrtb2.App{
			ID:   "test-app-id",
			Name: "Test App",
			Publisher: &openrtb2.Publisher{
				ID: "test-publisher-id",
			},
		},
		Device: &openrtb2.Device{
			UA: "test-user-agent",
		},
	}
}

func buildAuctionRequest(adType ad.Type, format ad.Format) *schema.AuctionRequest {
	adObject := schema.AdObject{
		AuctionID:  "test-auction-id",
		PriceFloor: 0.1,
		Demands: map[adapter.Key]map[string]any{
			adapter.MolocoKey: {
				"token": "test-token",
			},
		},
	}

	// Set the appropriate ad type object
	switch adType {
	case ad.BannerType:
		adObject.Banner = &schema.BannerAdObject{Format: format}
	case ad.InterstitialType:
		adObject.Interstitial = &schema.InterstitialAdObject{}
	case ad.RewardedType:
		adObject.Rewarded = &schema.RewardedAdObject{}
	}

	return &schema.AuctionRequest{
		AdObject: adObject,
		Adapters: schema.Adapters{
			adapter.MolocoKey: schema.Adapter{
				Version:    "1.0.0",
				SDKVersion: "1.0.0",
			},
		},
		BaseRequest: schema.BaseRequest{
			Device: schema.Device{
				Type: "PHONE",
			},
		},
	}
}

func TestMoloco_CreateRequest_Banner(t *testing.T) {
	adapter := buildAdapter()
	baseBidRequest := buildBaseBidRequest()
	auctionRequest := buildAuctionRequest(ad.BannerType, ad.BannerFormat)

	request, err := adapter.CreateRequest(baseBidRequest, auctionRequest)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(request.Imp) != 1 {
		t.Errorf("Expected 1 impression, got %d", len(request.Imp))
	}

	imp := request.Imp[0]
	if imp.TagID != "test-tag-id" {
		t.Errorf("Expected TagID 'test-tag-id', got '%s'", imp.TagID)
	}

	if imp.Banner == nil {
		t.Error("Expected banner impression, got nil")
	} else {
		if *imp.Banner.W != 320 || *imp.Banner.H != 50 {
			t.Errorf("Expected banner size 320x50, got %dx%d", *imp.Banner.W, *imp.Banner.H)
		}
	}

	if imp.BidFloorCur != "USD" {
		t.Errorf("Expected BidFloorCur 'USD', got '%s'", imp.BidFloorCur)
	}
}

func TestMoloco_CreateRequest_Interstitial(t *testing.T) {
	adapter := buildAdapter()
	baseBidRequest := buildBaseBidRequest()
	auctionRequest := buildAuctionRequest(ad.InterstitialType, ad.EmptyFormat)

	request, err := adapter.CreateRequest(baseBidRequest, auctionRequest)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	imp := request.Imp[0]
	if imp.Instl != 1 {
		t.Errorf("Expected interstitial flag 1, got %d", imp.Instl)
	}

	if imp.Banner == nil {
		t.Error("Expected banner impression for interstitial, got nil")
	}
}

func TestMoloco_CreateRequest_Rewarded(t *testing.T) {
	adapter := buildAdapter()
	baseBidRequest := buildBaseBidRequest()
	auctionRequest := buildAuctionRequest(ad.RewardedType, ad.EmptyFormat)

	request, err := adapter.CreateRequest(baseBidRequest, auctionRequest)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	imp := request.Imp[0]
	if imp.Video == nil {
		t.Error("Expected video impression for rewarded, got nil")
	}
}

func TestMoloco_CreateRequest_EmptyTagID(t *testing.T) {
	adapter := moloco.MolocoAdapter{
		TagID:  "",
		AppID:  "test-app-id",
		APIKey: "test-api-key",
	}
	baseBidRequest := buildBaseBidRequest()
	auctionRequest := buildAuctionRequest(ad.BannerType, ad.BannerFormat)

	_, err := adapter.CreateRequest(baseBidRequest, auctionRequest)

	if err == nil {
		t.Error("Expected error for empty TagID, got nil")
	}
}

func TestMoloco_CreateRequest_DefaultAdType(t *testing.T) {
	adapter := buildAdapter()
	baseBidRequest := buildBaseBidRequest()

	// Create an auction request with no ad type objects set - should default to banner
	auctionRequest := buildAuctionRequest(ad.BannerType, ad.BannerFormat)

	// This should succeed since it defaults to BannerType
	_, err := adapter.CreateRequest(baseBidRequest, auctionRequest)

	if err != nil {
		t.Errorf("Expected no error for default ad type, got %v", err)
	}
}

func TestMoloco_ExecuteRequest_Success(t *testing.T) {
	adapter := buildAdapter()
	baseBidRequest := buildBaseBidRequest()

	client := NewTestClient(func(req *http.Request) *http.Response {
		// Verify request headers
		if req.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type 'application/json', got '%s'", req.Header.Get("Content-Type"))
		}
		if req.Header.Get("Authorization") != "test-api-key" {
			t.Errorf("Expected Authorization 'test-api-key', got '%s'", req.Header.Get("Authorization"))
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(`{"id":"test-response"}`)),
		}
	})

	response := adapter.ExecuteRequest(context.Background(), client, baseBidRequest)

	if response.Error != nil {
		t.Errorf("Expected no error, got %v", response.Error)
	}

	if response.Status != 200 {
		t.Errorf("Expected status 200, got %d", response.Status)
	}

	if string(response.DemandID) != "moloco" {
		t.Errorf("Expected DemandID 'moloco', got '%s'", response.DemandID)
	}
}

func TestMoloco_ExecuteRequest_DefaultEndpoint(t *testing.T) {
	adapter := moloco.MolocoAdapter{
		TagID:  "test-tag-id",
		APIKey: "test-api-key",
		// Endpoint is empty, should use geographic routing with default US
	}
	baseBidRequest := buildBaseBidRequest()

	client := NewTestClient(func(req *http.Request) *http.Response {
		expectedURL := "https://sdkfnt-us.dsp-api.moloco.com/mediations/inhouse/v1"
		if req.URL.String() != expectedURL {
			t.Errorf("Expected URL '%s', got '%s'", expectedURL, req.URL.String())
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(`{}`)),
		}
	})

	adapter.ExecuteRequest(context.Background(), client, baseBidRequest)
}

func TestMoloco_ExecuteRequest_GeographicRouting(t *testing.T) {
	tests := []struct {
		name        string
		countryCode string
		expectedURL string
	}{
		{
			name:        "US region - USA",
			countryCode: "USA",
			expectedURL: "https://sdkfnt-us.dsp-api.moloco.com/mediations/inhouse/v1",
		},
		{
			name:        "Asia region - Japan",
			countryCode: "JPN",
			expectedURL: "https://sdkfnt-asia.dsp-api.moloco.com/mediations/inhouse/v1",
		},
		{
			name:        "EU region - Germany",
			countryCode: "DEU",
			expectedURL: "https://sdkfnt-eu.dsp-api.moloco.com/mediations/inhouse/v1",
		},
		{
			name:        "Unknown country - defaults to US",
			countryCode: "XXX",
			expectedURL: "https://sdkfnt-us.dsp-api.moloco.com/mediations/inhouse/v1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := moloco.MolocoAdapter{
				TagID:  "test-tag-id",
				APIKey: "test-api-key",
			}

			baseBidRequest := buildBaseBidRequest()
			// Set the country code in the request
			if baseBidRequest.Device == nil {
				baseBidRequest.Device = &openrtb2.Device{}
			}
			if baseBidRequest.Device.Geo == nil {
				baseBidRequest.Device.Geo = &openrtb2.Geo{}
			}
			baseBidRequest.Device.Geo.Country = tt.countryCode

			client := NewTestClient(func(req *http.Request) *http.Response {
				if req.URL.String() != tt.expectedURL {
					t.Errorf("Expected URL '%s', got '%s'", tt.expectedURL, req.URL.String())
				}

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{}`)),
				}
			})

			adapter.ExecuteRequest(context.Background(), client, baseBidRequest)
		})
	}
}

func TestMoloco_ParseBids_Success(t *testing.T) {
	adapter := buildAdapter()
	demandResponse := &adapters.DemandResponse{
		Status:      200,
		RawResponse: `{"id":"test-response","seatbid":[{"seat":"moloco","bid":[{"id":"test-bid","impid":"test-imp","price":1.5,"adm":"<html>test ad</html>","adid":"test-ad-id","nurl":"http://test.com/nurl","burl":"http://test.com/burl","lurl":"http://test.com/lurl","ext":{"signaldata":"test-signal"}}]}]}`,
	}

	result, err := adapter.ParseBids(demandResponse)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.Bid == nil {
		t.Error("Expected bid, got nil")
	} else {
		if result.Bid.ID != "test-bid" {
			t.Errorf("Expected bid ID 'test-bid', got '%s'", result.Bid.ID)
		}
		if result.Bid.Price != 1.5 {
			t.Errorf("Expected bid price 1.5, got %f", result.Bid.Price)
		}
		if result.Bid.Payload != "<html>test ad</html>" {
			t.Errorf("Expected payload '<html>test ad</html>', got '%s'", result.Bid.Payload)
		}
	}
}

func TestMoloco_ParseBids_NoContent(t *testing.T) {
	adapter := buildAdapter()
	demandResponse := &adapters.DemandResponse{
		Status: 204,
	}

	result, err := adapter.ParseBids(demandResponse)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.Bid != nil {
		t.Error("Expected no bid for 204 status, got bid")
	}
}

func TestMoloco_ParseBids_Unauthorized(t *testing.T) {
	adapter := buildAdapter()
	demandResponse := &adapters.DemandResponse{
		Status: 401,
	}

	_, err := adapter.ParseBids(demandResponse)

	if err == nil {
		t.Error("Expected error for 401 status, got nil")
	}
}

func TestMoloco_Builder(t *testing.T) {
	client := &http.Client{}
	molocoCfg := adapter.ProcessedConfigsMap{
		adapter.MolocoKey: map[string]any{
			"tag_id":   "test-tag-id",
			"app_id":   "test-app-id",
			"endpoint": "https://test.moloco.com/bid",
			"api_key":  "test-api-key",
		},
	}

	bidder, err := moloco.Builder(molocoCfg, client)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if bidder.Client != client {
		t.Error("Expected client to be set correctly")
	}

	// Test that adapter was created with correct configuration
	wantAdapter := moloco.MolocoAdapter{
		TagID:    "test-tag-id",
		AppID:    "test-app-id",
		APIKey:   "test-api-key",
	}
	wantBidder := &adapters.Bidder{
		Adapter: &wantAdapter,
		Client:  client,
	}

	if diff := cmp.Diff(wantBidder, bidder); diff != "" {
		t.Errorf("builder(molocoCfg, client) mismatch (-want, +got):\n%s", diff)
	}
}
