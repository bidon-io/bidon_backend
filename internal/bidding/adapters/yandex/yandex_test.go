package yandex_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/prebid/openrtb/v19/openrtb2"

	"github.com/bidon-io/bidon-backend/internal/ad"
	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/bidding/adapters"
	"github.com/bidon-io/bidon-backend/internal/bidding/adapters/yandex"
	"github.com/bidon-io/bidon-backend/internal/bidding/openrtb"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
)

type TestTransport func(req *http.Request) *http.Response

func (f TestTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func NewTestClient(tr TestTransport) *http.Client {
	return &http.Client{
		Transport: tr,
	}
}

func ptr[T any](t T) *T {
	return &t
}

func compareErrors(want, got error) bool {
	return (want == nil) == (got == nil)
}

func buildAdapter() yandex.YandexAdapter {
	return yandex.YandexAdapter{
		AdUnitID: "demo-banner-yandex",
	}
}

func buildBaseBidRequest() openrtb.BidRequest {
	return openrtb.BidRequest{
		ID: "test-request-id",
		App: &openrtb2.App{
			Bundle: "com.example.app",
		},
	}
}

func buildAuctionRequest(adType ad.Type, format ad.Format) *schema.AuctionRequest {
	req := &schema.AuctionRequest{
		Adapters: schema.Adapters{
			adapter.YandexKey: schema.Adapter{
				Version:    "1.0.0",
				SDKVersion: "1.0.0",
			},
		},
		AdObject: schema.AdObject{
			Demands: map[adapter.Key]map[string]any{
				adapter.YandexKey: {
					"token": "test-bidder-token",
				},
			},
		},
	}

	switch adType {
	case ad.BannerType:
		req.AdObject.Banner = &schema.BannerAdObject{
			Format: format,
		}
	case ad.InterstitialType:
		req.AdObject.Interstitial = &schema.InterstitialAdObject{}
	case ad.RewardedType:
		req.AdObject.Rewarded = &schema.RewardedAdObject{}
	}

	return req
}

func TestYandex_CreateRequest_Banner(t *testing.T) {
	tests := []struct {
		name    string
		format  ad.Format
		wantW   int64
		wantH   int64
		wantExt string
	}{
		{
			name:    "Banner 320x50",
			format:  ad.BannerFormat,
			wantW:   320,
			wantH:   50,
			wantExt: `{"ad_type":"banner"}`,
		},
		{
			name:    "MREC 300x250",
			format:  ad.MRECFormat,
			wantW:   300,
			wantH:   250,
			wantExt: `{"ad_type":"banner"}`,
		},
		{
			name:    "Leaderboard 728x90",
			format:  ad.LeaderboardFormat,
			wantW:   728,
			wantH:   90,
			wantExt: `{"ad_type":"banner"}`,
		},
		{
			name:    "Adaptive",
			format:  ad.AdaptiveFormat,
			wantW:   320,
			wantH:   50,
			wantExt: `{"ad_type":"banner"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := buildAdapter()
			baseReq := buildBaseBidRequest()
			auctionReq := buildAuctionRequest(ad.BannerType, tt.format)

			got, err := adapter.CreateRequest(baseReq, auctionReq)
			if err != nil {
				t.Fatalf("CreateRequest() error = %v", err)
			}

			if len(got.Imp) != 1 {
				t.Fatalf("Expected 1 impression, got %d", len(got.Imp))
			}

			imp := got.Imp[0]

			if imp.TagID != "demo-banner-yandex" {
				t.Errorf("TagID = %v, want %v", imp.TagID, "demo-banner-yandex")
			}

			if imp.Banner == nil {
				t.Fatal("Banner is nil")
			}

			if *imp.Banner.W != tt.wantW {
				t.Errorf("Banner.W = %v, want %v", *imp.Banner.W, tt.wantW)
			}

			if *imp.Banner.H != tt.wantH {
				t.Errorf("Banner.H = %v, want %v", *imp.Banner.H, tt.wantH)
			}

			if string(imp.Ext) != tt.wantExt {
				t.Errorf("Imp.Ext = %v, want %v", string(imp.Ext), tt.wantExt)
			}

			if imp.Instl != 0 {
				t.Errorf("Instl = %v, want 0", imp.Instl)
			}

			if imp.Rwdd != 0 {
				t.Errorf("Rwdd = %v, want 0", imp.Rwdd)
			}

			// Check user token
			if got.User == nil {
				t.Fatal("User is nil")
			}

			if len(got.User.Data) != 1 {
				t.Fatalf("Expected 1 data object, got %d", len(got.User.Data))
			}

			if len(got.User.Data[0].Segment) != 1 {
				t.Fatalf("Expected 1 segment, got %d", len(got.User.Data[0].Segment))
			}

			if got.User.Data[0].Segment[0].Signal != "test-bidder-token" {
				t.Errorf("Signal = %v, want %v", got.User.Data[0].Segment[0].Signal, "test-bidder-token")
			}

			if len(got.Cur) != 1 || got.Cur[0] != "USD" {
				t.Errorf("Cur = %v, want [USD]", got.Cur)
			}
		})
	}
}

func TestYandex_CreateRequest_Interstitial(t *testing.T) {
	adapter := buildAdapter()
	baseReq := buildBaseBidRequest()
	auctionReq := buildAuctionRequest(ad.InterstitialType, ad.EmptyFormat)

	got, err := adapter.CreateRequest(baseReq, auctionReq)
	if err != nil {
		t.Fatalf("CreateRequest() error = %v", err)
	}

	if len(got.Imp) != 1 {
		t.Fatalf("Expected 1 impression, got %d", len(got.Imp))
	}

	imp := got.Imp[0]

	if imp.Instl != 1 {
		t.Errorf("Instl = %v, want 1", imp.Instl)
	}

	if imp.Rwdd != 0 {
		t.Errorf("Rwdd = %v, want 0", imp.Rwdd)
	}

	wantExt := `{"ad_type":"interstitial"}`
	if string(imp.Ext) != wantExt {
		t.Errorf("Imp.Ext = %v, want %v", string(imp.Ext), wantExt)
	}
}

func TestYandex_CreateRequest_Rewarded(t *testing.T) {
	adapter := buildAdapter()
	baseReq := buildBaseBidRequest()
	auctionReq := buildAuctionRequest(ad.RewardedType, ad.EmptyFormat)

	got, err := adapter.CreateRequest(baseReq, auctionReq)
	if err != nil {
		t.Fatalf("CreateRequest() error = %v", err)
	}

	if len(got.Imp) != 1 {
		t.Fatalf("Expected 1 impression, got %d", len(got.Imp))
	}

	imp := got.Imp[0]

	if imp.Instl != 0 {
		t.Errorf("Instl = %v, want 0", imp.Instl)
	}

	if imp.Rwdd != 1 {
		t.Errorf("Rwdd = %v, want 1", imp.Rwdd)
	}

	if imp.Video == nil {
		t.Fatal("Video is nil")
	}

	wantExt := `{"ad_type":"rewarded"}`
	if string(imp.Ext) != wantExt {
		t.Errorf("Imp.Ext = %v, want %v", string(imp.Ext), wantExt)
	}
}

func TestYandex_CreateRequest_EmptyAdUnitID(t *testing.T) {
	adapter := yandex.YandexAdapter{AdUnitID: ""}
	baseReq := buildBaseBidRequest()
	auctionReq := buildAuctionRequest(ad.BannerType, ad.BannerFormat)

	_, err := adapter.CreateRequest(baseReq, auctionReq)
	if err == nil {
		t.Error("Expected error for empty AdUnitID, got nil")
	}
}

func TestYandex_CreateRequest_MissingToken(t *testing.T) {
	adpt := buildAdapter()
	baseReq := buildBaseBidRequest()
	auctionReq := buildAuctionRequest(ad.BannerType, ad.BannerFormat)
	auctionReq.AdObject.Demands[adapter.YandexKey] = map[string]any{}

	_, err := adpt.CreateRequest(baseReq, auctionReq)
	if err == nil {
		t.Error("Expected error for missing token, got nil")
	}
}

func TestYandex_ExecuteRequest(t *testing.T) {
	adapter := buildAdapter()
	baseReq := buildBaseBidRequest()

	testClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.URL.String() != "https://mobile.yandexadexchange.net/openbidding?ssp-id=99048272" {
			t.Errorf("URL = %v, want https://mobile.yandexadexchange.net/openbidding?ssp-id=99048272", req.URL.String())
		}

		if req.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Content-Type = %v, want application/json", req.Header.Get("Content-Type"))
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(`{"id":"test"}`)),
		}
	})

	dr := adapter.ExecuteRequest(context.Background(), testClient, baseReq)

	if dr.Error != nil {
		t.Errorf("ExecuteRequest() error = %v", dr.Error)
	}

	if dr.Status != http.StatusOK {
		t.Errorf("Status = %v, want %v", dr.Status, http.StatusOK)
	}
}

func TestYandex_ParseBids_Success(t *testing.T) {
	adpt := buildAdapter()

	bidResponse := openrtb2.BidResponse{
		ID:  "test-request-id",
		Cur: "USD",
		SeatBid: []openrtb2.SeatBid{
			{
				Seat: "yandex",
				Bid: []openrtb2.Bid{
					{
						ID:    "bid-123",
						ImpID: "imp-123",
						Price: 8.2869,
						AdM:   "ad-markup",
						AdID:  "ad-123",
						NURL:  "https://yandex.ru/nurl",
						BURL:  "https://yandex.ru/burl",
						LURL:  "https://yandex.ru/lurl",
						Ext:   json.RawMessage(`{"signaldata":"{\"cache_id\":\"openbidding_123\"}"}`),
					},
				},
			},
		},
	}

	respBody, _ := json.Marshal(bidResponse)

	dr := &adapters.DemandResponse{
		DemandID:    adapter.YandexKey,
		RequestID:   "test-request-id",
		RawResponse: string(respBody),
		Status:      http.StatusOK,
	}

	result, err := adpt.ParseBids(dr)
	if err != nil {
		t.Fatalf("ParseBids() error = %v", err)
	}

	if result.Bid == nil {
		t.Fatal("Bid is nil")
	}

	if result.Bid.ID != "bid-123" {
		t.Errorf("Bid.ID = %v, want bid-123", result.Bid.ID)
	}

	if result.Bid.Price != 8.2869 {
		t.Errorf("Bid.Price = %v, want 8.2869", result.Bid.Price)
	}

	if result.Bid.Signaldata != `{"cache_id":"openbidding_123"}` {
		t.Errorf("Bid.Signaldata = %v, want {\"cache_id\":\"openbidding_123\"}", result.Bid.Signaldata)
	}

	if result.Bid.NURL != "https://yandex.ru/nurl" {
		t.Errorf("Bid.NURL = %v, want https://yandex.ru/nurl", result.Bid.NURL)
	}
}

func TestYandex_ParseBids_NoContent(t *testing.T) {
	adpt := buildAdapter()

	dr := &adapters.DemandResponse{
		DemandID:  adapter.YandexKey,
		RequestID: "test-request-id",
		Status:    http.StatusNoContent,
	}

	result, err := adpt.ParseBids(dr)
	if err != nil {
		t.Fatalf("ParseBids() error = %v", err)
	}

	if result.Bid != nil {
		t.Error("Expected nil bid for NoContent status")
	}
}

func TestYandex_ParseBids_ErrorStatuses(t *testing.T) {
	adpt := buildAdapter()

	errorStatuses := []int{
		http.StatusBadRequest,
		http.StatusUnauthorized,
		http.StatusForbidden,
		http.StatusServiceUnavailable,
	}

	for _, status := range errorStatuses {
		t.Run(http.StatusText(status), func(t *testing.T) {
			dr := &adapters.DemandResponse{
				DemandID:  adapter.YandexKey,
				RequestID: "test-request-id",
				Status:    status,
			}

			_, err := adpt.ParseBids(dr)
			if err == nil {
				t.Errorf("Expected error for status %d, got nil", status)
			}
		})
	}
}

func TestYandex_ParseBids_EmptySeatBid(t *testing.T) {
	adpt := buildAdapter()

	bidResponse := openrtb2.BidResponse{
		ID:      "test-request-id",
		Cur:     "USD",
		SeatBid: []openrtb2.SeatBid{},
	}

	respBody, _ := json.Marshal(bidResponse)

	dr := &adapters.DemandResponse{
		DemandID:    adapter.YandexKey,
		RequestID:   "test-request-id",
		RawResponse: string(respBody),
		Status:      http.StatusOK,
	}

	result, err := adpt.ParseBids(dr)
	if err != nil {
		t.Fatalf("ParseBids() error = %v", err)
	}

	if result.Bid != nil {
		t.Error("Expected nil bid for empty seatbid")
	}
}

func TestYandex_ParseBids_InvalidJSON(t *testing.T) {
	adpt := buildAdapter()

	dr := &adapters.DemandResponse{
		DemandID:    adapter.YandexKey,
		RequestID:   "test-request-id",
		RawResponse: "invalid json",
		Status:      http.StatusOK,
	}

	_, err := adpt.ParseBids(dr)
	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
}

func TestYandex_Builder(t *testing.T) {
	client := &http.Client{}
	yandexCfg := adapter.ProcessedConfigsMap{
		adapter.YandexKey: map[string]any{
			"ad_unit_id": "demo-banner-yandex",
		},
	}

	bidder, err := yandex.Builder(yandexCfg, client)
	if err != nil {
		t.Fatalf("Builder() error = %v", err)
	}

	wantAdapter := buildAdapter()
	wantBidder := &adapters.Bidder{
		Adapter: &wantAdapter,
		Client:  client,
	}

	if diff := cmp.Diff(wantBidder, bidder); diff != "" {
		t.Errorf("Builder() mismatch (-want, +got):\n%s", diff)
	}
}

func TestYandex_Builder_MissingAdUnitID(t *testing.T) {
	client := &http.Client{}
	yandexCfg := adapter.ProcessedConfigsMap{
		adapter.YandexKey: map[string]any{},
	}

	bidder, err := yandex.Builder(yandexCfg, client)
	if err != nil {
		t.Fatalf("Builder() error = %v", err)
	}

	if bidder == nil {
		t.Error("Expected bidder to be created even with missing ad_unit_id")
	}
}
