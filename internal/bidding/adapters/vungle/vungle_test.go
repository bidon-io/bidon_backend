package vungle_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/bidon-io/bidon-backend/internal/bidding/adapters"

	"github.com/bidon-io/bidon-backend/internal/ad"
	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/bidding/adapters/vungle"
	"github.com/bidon-io/bidon-backend/internal/bidding/openrtb"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
	"github.com/google/go-cmp/cmp"
	"github.com/prebid/openrtb/v19/adcom1"
	"github.com/prebid/openrtb/v19/openrtb2"
)

type createRequestTestParams struct {
	BaseBidRequest openrtb.BidRequest
	Br             *schema.BiddingRequest
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

func ptr[T any](t T) *T {
	return &t
}

func compareErrors(want, got error) bool {
	return (want == nil) == (got == nil)
}

func buildAdapter() vungle.VungleAdapter {
	return vungle.VungleAdapter{
		SellerID: "1",
		AppID:    "10182906",
		TagID:    "10182906-10192212",
	}
}

func buildBaseRequest() openrtb.BidRequest {
	return openrtb.BidRequest{
		App: &openrtb2.App{
			Publisher: &openrtb2.Publisher{},
			Ext:       json.RawMessage(`{"orientation":1}`),
		},
	}
}

func buildTestParams(imp schema.Imp) createRequestTestParams {
	request := buildBaseRequest()

	br := schema.BiddingRequest{
		Adapters: schema.Adapters{
			"vungle": schema.Adapter{
				Version:    "1.0.0",
				SDKVersion: "1.0.0",
			},
		},
		Imp: schema.Imp{
			Demands: map[adapter.Key]map[string]any{
				adapter.VungleKey: {
					"token": "token",
				},
			},
			Orientation: "PORTRAIT",
		},
		BaseRequest: schema.BaseRequest{
			Device: schema.Device{
				Type: "PHONE",
			},
		},
	}

	if imp.Banner != nil {
		br.Imp.Banner = imp.Banner
	}
	if imp.Interstitial != nil {
		br.Imp.Interstitial = imp.Interstitial
	}
	if imp.Rewarded != nil {
		br.Imp.Rewarded = imp.Rewarded
	}

	return createRequestTestParams{
		BaseBidRequest: request,
		Br:             &br,
	}
}

func buildWantRequest(imp openrtb2.Imp) openrtb.BidRequest {
	request := openrtb.BidRequest{
		App: &openrtb2.App{
			ID:        "10182906",
			Publisher: &openrtb2.Publisher{ID: "1"},
			Ext:       json.RawMessage(`{"orientation":1}`),
		},
		User: nil,
		Cur:  []string{"USD"},
		Imp: []openrtb2.Imp{
			{
				ID:                "1",
				DisplayManager:    "vungle",
				DisplayManagerVer: "1.0.0",
				TagID:             "10182906",
				BidFloorCur:       "USD",
				Secure:            ptr(int8(1)),
				BidFloor:          schema.MinBidFloor,
			},
		},
	}
	if imp.Banner != nil {
		request.Imp[0].Banner = imp.Banner
	}
	if imp.Video != nil {
		request.Imp[0].Video = imp.Video
	}
	if imp.Instl != 0 {
		request.Imp[0].Instl = imp.Instl
	}
	if imp.Ext != nil {
		request.Imp[0].Ext = imp.Ext
	}

	return request
}

func TestVungle_CreateRequestTest(t *testing.T) {
	testCases := []struct {
		name   string
		params createRequestTestParams
		want   createRequestTestOutput
	}{
		{
			name: "Banner MREC",
			params: buildTestParams(
				schema.Imp{
					Banner: &schema.BannerAdObject{
						Format: ad.MRECFormat,
					},
				},
			),
			want: createRequestTestOutput{
				Request: buildWantRequest(openrtb2.Imp{
					Instl: 0,
					Banner: &openrtb2.Banner{
						W:   ptr(int64(300)),
						H:   ptr(int64(250)),
						Pos: adcom1.PositionAboveFold.Ptr(),
					},
					Ext: json.RawMessage(`{"vungle":{"bid_token":"token"}}`),
				}),
				Err: nil,
			},
		},
		{
			name: "Banner BANNER",
			params: buildTestParams(
				schema.Imp{
					Banner: &schema.BannerAdObject{
						Format: ad.BannerFormat,
					},
				},
			),
			want: createRequestTestOutput{
				Request: buildWantRequest(openrtb2.Imp{
					Instl: 0,
					Banner: &openrtb2.Banner{
						W:   ptr(int64(320)),
						H:   ptr(int64(50)),
						Pos: adcom1.PositionAboveFold.Ptr(),
					},
					Ext: json.RawMessage(`{"vungle":{"bid_token":"token"}}`),
				}),
				Err: nil,
			},
		},
		{
			name: "Interstitial",
			params: buildTestParams(
				schema.Imp{
					Interstitial: &schema.InterstitialAdObject{},
				},
			),
			want: createRequestTestOutput{
				Request: buildWantRequest(openrtb2.Imp{
					Instl: 1,
					Banner: &openrtb2.Banner{
						W:   ptr(int64(320)),
						H:   ptr(int64(480)),
						Pos: adcom1.PositionFullScreen.Ptr(),
					},
					Ext: json.RawMessage(`{"vungle":{"bid_token":"token"}}`),
				}),
				Err: nil,
			},
		},
		{
			name: "Rewarded",
			params: buildTestParams(
				schema.Imp{
					Rewarded: &schema.RewardedAdObject{},
				},
			),
			want: createRequestTestOutput{
				Request: buildWantRequest(openrtb2.Imp{
					Instl: 0,
					Video: &openrtb2.Video{
						MIMEs: []string{"video/mp4"},
						W:     int64(320),
						H:     int64(480),
						Ext:   json.RawMessage(`{"rewarded": 1}`),
					},
					Ext: json.RawMessage(`{"vungle":{"bid_token":"token"}}`),
				}),
				Err: nil,
			},
		},
	}

	adapter := &vungle.VungleAdapter{
		SellerID: "1",
		AppID:    "10182906",
		TagID:    "10182906",
	}

	for _, tC := range testCases {
		request, err := adapter.CreateRequest(tC.params.BaseBidRequest, tC.params.Br)
		if err == nil {
			request.Imp[0].ID = "1" // ommit random uuid
		}
		got := createRequestTestOutput{
			Request: request,
			Err:     err,
		}
		if diff := cmp.Diff(tC.want, got, cmp.Comparer(func(x, y error) bool {
			return x == y
		})); diff != "" {
			t.Errorf("%s: adapter.CreateRequest(ctx, %v, %v) mismatch (-want, +got):\n%s", tC.name, tC.params.BaseBidRequest, tC.params.Br, diff)
		}
	}
}

func TestVungleAdapter_ExecuteRequest(t *testing.T) {
	networkAdapter := buildAdapter()
	responseBody := []byte(`{"key": "value"`)

	customClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.Method != "POST" {
			t.Errorf("Expected POST request")
		}
		if req.URL.String() != "https://rtb.ads.vungle.com/bid/t/8ea3e9a" {
			t.Errorf("Expected URL: https://rtb.ads.vungle.com/bid/t/8ea3e9a")
		}
		contentType := req.Header.Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("Expected Content-Type header: application/json")
		}
		openrtbHeader := req.Header.Get("X-OpenRTB-Version")
		if openrtbHeader != "2.5" {
			t.Errorf("Expected X-OpenRTB-Version header: 2.5")
		}
		return &http.Response{
			Status:        http.StatusText(http.StatusOK),
			StatusCode:    http.StatusOK,
			Body:          io.NopCloser(bytes.NewBuffer(responseBody)),
			ContentLength: int64(len(responseBody)),
		}
	})
	request := openrtb.BidRequest{
		ID: "test-request-id",
	}

	response := networkAdapter.ExecuteRequest(context.Background(), customClient, request)

	if response.DemandID != adapter.VungleKey {
		t.Errorf("Expected DemandID %v, but got %v", adapter.VungleKey, response.DemandID)
	}
	if response.RequestID != request.ID {
		t.Errorf("Expected RequestID %v, but got %v", request.ID, response.RequestID)
	}
	if response.TagID != networkAdapter.TagID {
		t.Errorf("Expected TagID %v, but got %v", networkAdapter.TagID, response.TagID)
	}
	if response.Error != nil {
		t.Errorf("Expected no error, but got an error: %v", response.Error)
	}
	if response.RawResponse != string(responseBody) {
		t.Errorf("Expected client response body as RawResponse but got: %v", response.RawResponse)
	}
	if response.Status != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, response.Status)
	}
}

func TestVungle_ParseBids(t *testing.T) {
	rawResponse := `{
		"id": "47611e59-e05b-4e1e-9074-5a65eb4501e4",
		"seatbid": [
			{
				"bid": [
					{
						"id": "0",
						"impid": "6579ca7b-7e2c-48b6-8915-46efa6530fb5",
						"price": 1.5,
						"nurl": "https://api.gov-static.tech/Ad/AdxEvent?sid=0&sslot=10182906-10163778&adtype=4",
						"lurl": "https://api.gov-static.tech/Ad/AdxEvent?sid=0&sslot=10182906-10163778",
						"adm": "0692d0a0efdbd5bd470dafea742cef6a1f6b840c5c83240e165bc33a038b3d5487e25a52",
						"adid": "Vunglead5e0471131b8a4e3c",
						"crid": "e2d42134881d5b45134f3cf77989dec7"
					}
				]
			}
		]
	}`
	adapter := buildAdapter()

	testCases := []struct {
		name   string
		params ParseBidsTestParams
		want   ParseBidsTestOutput
	}{
		{
			name: "ParseBids Success",
			params: ParseBidsTestParams{
				DemandsResponse: adapters.DemandResponse{
					Status:      200,
					RawResponse: rawResponse,
				},
			},
			want: ParseBidsTestOutput{
				DemandResponse: adapters.DemandResponse{
					Status:      200,
					RawResponse: rawResponse,
					Bid: &adapters.BidDemandResponse{
						ID:       "0",
						ImpID:    "6579ca7b-7e2c-48b6-8915-46efa6530fb5",
						Price:    1.5,
						Payload:  "0692d0a0efdbd5bd470dafea742cef6a1f6b840c5c83240e165bc33a038b3d5487e25a52",
						DemandID: "vungle",
						AdID:     "Vunglead5e0471131b8a4e3c",
						LURL:     "https://api.gov-static.tech/Ad/AdxEvent?sid=0&sslot=10182906-10163778",
						NURL:     "https://api.gov-static.tech/Ad/AdxEvent?sid=0&sslot=10182906-10163778&adtype=4",
					},
				},
				Err: nil,
			},
		},
		{
			name: "ParseBids Bad Request",
			params: ParseBidsTestParams{
				DemandsResponse: adapters.DemandResponse{
					Status:      400,
					RawResponse: rawResponse,
				},
			},
			want: ParseBidsTestOutput{
				DemandResponse: adapters.DemandResponse{
					Status:      400,
					RawResponse: rawResponse,
				},
				Err: errors.New("unauthorized request: 400"),
			},
		},
		{
			name: "ParseBids No Content",
			params: ParseBidsTestParams{
				DemandsResponse: adapters.DemandResponse{
					Status:      204,
					RawResponse: rawResponse,
				},
			},
			want: ParseBidsTestOutput{
				DemandResponse: adapters.DemandResponse{
					Status:      204,
					RawResponse: rawResponse,
				},
				Err: nil,
			},
		},
	}
	for _, tC := range testCases {
		response, err := adapter.ParseBids(&tC.params.DemandsResponse)
		got := ParseBidsTestOutput{
			DemandResponse: *response,
			Err:            err,
		}
		if diff := cmp.Diff(tC.want, got, cmp.Comparer(compareErrors)); diff != "" {
			t.Errorf("%s: adapter.ParseBids(ctx, %v) mismatch (-want, +got):\n%s", tC.name, tC.params.DemandsResponse, diff)
		}
	}
}

func TestVungle_Builder(t *testing.T) {
	client := &http.Client{}
	bigoCfg := adapter.ProcessedConfigsMap{
		adapter.VungleKey: map[string]any{
			"seller_id": "1",
			"app_id":    "10182906",
			"tag_id":    "10182906-10192212",
		},
	}
	bidder, err := vungle.Builder(bigoCfg, client)
	wantAdapter := buildAdapter()
	wantBidder := &adapters.Bidder{
		Adapter: &wantAdapter,
		Client:  client,
	}
	if err != nil {
		t.Errorf("Error building adapter: %v", err)
	}
	if diff := cmp.Diff(wantBidder, bidder); diff != "" {
		t.Errorf("builder(bigoCfg, client) mismatch (-want, +got):\n%s", diff)
	}
}
