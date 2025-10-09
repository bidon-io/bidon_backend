package startio_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/prebid/openrtb/v19/openrtb2"

	"github.com/bidon-io/bidon-backend/internal/ad"
	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/bidding/adapters"
	"github.com/bidon-io/bidon-backend/internal/bidding/adapters/startio"
	"github.com/bidon-io/bidon-backend/internal/bidding/openrtb"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
)

type transportFunc func(req *http.Request) *http.Response

func (f transportFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func newTestClient(tr transportFunc) *http.Client {
	return &http.Client{Transport: tr}
}

func buildAdapter() startio.Adapter {
	return startio.Adapter{
		TagID:   "test-tag-id",
		AppID:   "test-app-id",
		Account: "trp",
	}
}

func buildBidRequest() openrtb.BidRequest {
	return openrtb.BidRequest{
		ID: "test-request-id",
		App: &openrtb2.App{ID: ""},
		Device: &openrtb2.Device{
			UA: "test-user-agent",
			Geo: &openrtb2.Geo{Country: "USA"},
		},
	}
}

func buildAuctionRequest(adType ad.Type, format ad.Format) *schema.AuctionRequest {
	adObject := schema.AdObject{
		AuctionID:  "auction",
		PriceFloor: 0.1,
		Demands: map[adapter.Key]map[string]any{
			adapter.StartIOKey: {
				"token": "test-token",
			},
		},
	}

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
			adapter.StartIOKey: {SDKVersion: "1.2.3"},
		},
		BaseRequest: schema.BaseRequest{
			Device: schema.Device{Type: "PHONE"},
		},
	}
}

func TestAdapter_CreateRequest_Banner(t *testing.T) {
	adapterInstance := buildAdapter()
	base := buildBidRequest()
	auction := buildAuctionRequest(ad.BannerType, ad.BannerFormat)
	auction.Test = true

	request, err := adapterInstance.CreateRequest(base, auction)
	if err != nil {
		t.Fatalf("CreateRequest returned error: %v", err)
	}

	if len(request.Imp) != 1 {
		t.Fatalf("expected 1 imp, got %d", len(request.Imp))
	}

	if request.User == nil || request.User.BuyerUID != "test-token" {
		t.Fatalf("expected BuyerUID 'test-token', got %+v", request.User)
	}

	if request.App == nil || request.App.ID != "test-app-id" {
		t.Fatalf("expected App.ID to be 'test-app-id', got %+v", request.App)
	}

	if request.Test != 1 {
		t.Fatalf("expected request.Test=1, got %d", request.Test)
	}

	imp := request.Imp[0]
	if imp.TagID != "test-tag-id" {
		t.Fatalf("expected TagID 'test-tag-id', got %s", imp.TagID)
	}

	if imp.DisplayManager != string(adapter.StartIOKey) {
		t.Fatalf("unexpected DisplayManager: %s", imp.DisplayManager)
	}
}

func TestAdapter_CreateRequest_Errors(t *testing.T) {
	base := buildBidRequest()
	auction := buildAuctionRequest(ad.BannerType, ad.BannerFormat)

	cases := []struct {
		name    string
		adapter startio.Adapter
		wantErr string
	}{
		{name: "missing tag", adapter: startio.Adapter{AppID: "app", Account: "acc"}, wantErr: "tag ID"},
		{name: "missing account", adapter: startio.Adapter{TagID: "tag", AppID: "app"}, wantErr: "account"},
		{name: "missing app", adapter: startio.Adapter{TagID: "tag", Account: "acc"}, wantErr: "app ID"},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			_, err := tc.adapter.CreateRequest(base, auction)
			if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("expected error containing %q, got %v", tc.wantErr, err)
			}
		})
	}

	t.Run("missing token", func(t *testing.T) {
		adapterInstance := buildAdapter()
		auction := buildAuctionRequest(ad.BannerType, ad.BannerFormat)
		auction.AdObject.Demands[adapter.StartIOKey]["token"] = ""

		_, err := adapterInstance.CreateRequest(base, auction)
		if err == nil || !strings.Contains(err.Error(), "token") {
			t.Fatalf("expected token error, got %v", err)
		}
	})
}

func TestAdapter_ExecuteRequest(t *testing.T) {
	adapterInstance := buildAdapter()
	base := buildBidRequest()
	auction := buildAuctionRequest(ad.InterstitialType, ad.EmptyFormat)

	request, err := adapterInstance.CreateRequest(base, auction)
	if err != nil {
		t.Fatalf("CreateRequest error: %v", err)
	}

	client := newTestClient(func(req *http.Request) *http.Response {
		if req.URL.Scheme != "http" {
			t.Fatalf("unexpected scheme: %s", req.URL.Scheme)
		}
		if req.URL.Host != "trp-rtb.startappnetwork.com" {
			t.Fatalf("unexpected host: %s", req.URL.Host)
		}
		if req.URL.Query().Get("account") != "trp" {
			t.Fatalf("expected account query param")
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(`{"id":"resp"}`)),
			Header:     make(http.Header),
		}
	})

	dr := adapterInstance.ExecuteRequest(context.Background(), client, request)
	if dr.Error != nil {
		t.Fatalf("ExecuteRequest error: %v", dr.Error)
	}

	if dr.Status != http.StatusOK {
		t.Fatalf("unexpected status: %d", dr.Status)
	}
}

func TestAdapter_ParseBids(t *testing.T) {
	adapterInstance := buildAdapter()
	resp := adapters.DemandResponse{
		DemandID: adapter.StartIOKey,
		Status:   http.StatusOK,
		RawResponse: `{"id":"response","seatbid":[{"seat":"startio","bid":[{"id":"bid-1","impid":"imp-1","price":1.23,"adm":"<html></html>","adid":"creative","nurl":"http://nurl","lurl":"http://lurl","burl":"http://burl"}]}]}`,
	}

	parsed, err := adapterInstance.ParseBids(&resp)
	if err != nil {
		t.Fatalf("ParseBids error: %v", err)
	}

	if parsed.Bid == nil {
		t.Fatal("expected bid to be set")
	}

	if parsed.Bid.Price != 1.23 {
		t.Fatalf("unexpected price: %f", parsed.Bid.Price)
	}
}

func TestBuilder(t *testing.T) {
	cfg := adapter.ProcessedConfigsMap{
		adapter.StartIOKey: {
			"tag_id":  "tag",
			"app_id":  "app",
			"account": "acc",
		},
	}

	client := newTestClient(func(req *http.Request) *http.Response {
		return &http.Response{StatusCode: http.StatusNoContent, Body: io.NopCloser(bytes.NewBuffer(nil))}
	})

	bidder, err := startio.Builder(cfg, client)
	if err != nil {
		t.Fatalf("Builder returned error: %v", err)
	}

	adpt, ok := bidder.Adapter.(*startio.Adapter)
	if !ok {
		t.Fatalf("unexpected adapter type %T", bidder.Adapter)
	}

	want := startio.Adapter{TagID: "tag", AppID: "app", Account: "acc"}
	if diff := cmp.Diff(want, *adpt); diff != "" {
		t.Fatalf("builder mismatch (-want +got):\n%s", diff)
	}
}
