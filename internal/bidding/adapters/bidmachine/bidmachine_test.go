package bidmachine_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/prebid/openrtb/v19/adcom1"
	"github.com/prebid/openrtb/v19/openrtb2"

	"github.com/bidon-io/bidon-backend/internal/ad"
	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/bidding/adapters"
	"github.com/bidon-io/bidon-backend/internal/bidding/adapters/bidmachine"
	"github.com/bidon-io/bidon-backend/internal/bidding/openrtb"
	"github.com/bidon-io/bidon-backend/internal/device"
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

func ptr[T any](t T) *T {
	return &t
}

// compareErrors checks for error occurrence.
func compareErrors(want, got error) bool {
	return (want == nil) == (got == nil)
}

func buildAdapter() bidmachine.BidmachineAdapter {
	return bidmachine.BidmachineAdapter{
		SellerID: "1",
		Endpoint: "example.com",
	}
}

func buildBaseRequest() openrtb.BidRequest {
	return openrtb.BidRequest{
		App: &openrtb2.App{
			Publisher: &openrtb2.Publisher{ID: "1"},
		},
	}
}

func buildTestParams(imp schema.AdObject) createRequestTestParams {
	request := buildBaseRequest()

	auctionRequest := schema.AuctionRequest{
		Adapters: schema.Adapters{
			"bidmachine": schema.Adapter{
				Version:    "1.0.0",
				SDKVersion: "1.0.0",
			},
		},
		AdObject: schema.AdObject{
			Demands: map[adapter.Key]map[string]any{
				adapter.BidmachineKey: {
					"token": "token",
				},
			},
			Orientation: "PORTRAIT",
		},
		BaseRequest: schema.BaseRequest{
			App: schema.App{
				SDKVersion: "1.0.0",
			},
			Device: schema.Device{
				Type: device.PhoneType,
			},
		},
	}

	if imp.Banner != nil {
		auctionRequest.AdObject.Banner = imp.Banner
	}
	if imp.Interstitial != nil {
		auctionRequest.AdObject.Interstitial = imp.Interstitial
	}
	if imp.Rewarded != nil {
		auctionRequest.AdObject.Rewarded = imp.Rewarded
	}

	return createRequestTestParams{
		BaseBidRequest: request,
		AuctionRequest: &auctionRequest,
	}
}

func buildWantRequest(imp openrtb2.Imp) openrtb.BidRequest {
	request := openrtb.BidRequest{
		App: &openrtb2.App{
			Publisher: &openrtb2.Publisher{ID: "1"},
		},
		User: nil,
		Cur:  []string{"USD"},
		Ext:  json.RawMessage(`{"bidon_sdk_version":"1.0.0","mediation_mode":"bidon"}`),
		Imp: []openrtb2.Imp{
			{
				ID:                "1",
				DisplayManager:    "bidmachine",
				DisplayManagerVer: "1.0.0",
				TagID:             "",
				Secure:            ptr(int8(1)),
				Ext:               json.RawMessage(`{"bid_token":"token"}`),
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

func TestBidmachine_CreateRequest(t *testing.T) {
	testCases := []struct {
		name   string
		params createRequestTestParams
		want   createRequestTestOutput
	}{
		{
			name: "Banner MREC",
			params: buildTestParams(
				schema.AdObject{
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
				}),
				Err: nil,
			},
		},
		{
			name: "Banner BANNER",
			params: buildTestParams(
				schema.AdObject{
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
				}),
				Err: nil,
			},
		},
		{
			name: "Banner LEADERBOARD",
			params: buildTestParams(
				schema.AdObject{
					Banner: &schema.BannerAdObject{
						Format: ad.LeaderboardFormat,
					},
				},
			),
			want: createRequestTestOutput{
				Request: buildWantRequest(openrtb2.Imp{
					Instl: 0,
					Banner: &openrtb2.Banner{
						W:   ptr(int64(728)),
						H:   ptr(int64(90)),
						Pos: adcom1.PositionAboveFold.Ptr(),
					},
				}),
				Err: nil,
			},
		},
		{
			name: "Banner ADAPTIVE",
			params: buildTestParams(
				schema.AdObject{
					Banner: &schema.BannerAdObject{
						Format: ad.AdaptiveFormat,
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
				}),
				Err: nil,
			},
		},
		{
			name: "Banner empty format",
			params: buildTestParams(
				schema.AdObject{
					Banner: &schema.BannerAdObject{
						Format: "",
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
				}),
				Err: nil,
			},
		},
		{
			name: "Interstitial",
			params: buildTestParams(
				schema.AdObject{
					Interstitial: &schema.InterstitialAdObject{},
				},
			),
			want: createRequestTestOutput{
				Request: buildWantRequest(openrtb2.Imp{
					Instl: 1,
					Banner: &openrtb2.Banner{
						W:     ptr(int64(320)),
						H:     ptr(int64(480)),
						BType: []openrtb2.BannerAdType{},
						BAttr: []adcom1.CreativeAttribute{},
						Pos:   adcom1.PositionFullScreen.Ptr(),
					},
					Video: &openrtb2.Video{
						W:     int64(320),
						H:     int64(480),
						Pos:   adcom1.PositionFullScreen.Ptr(),
						MIMEs: []string{"video/mp4", "video/3gpp", "video/3gpp2", "video/x-m4v", "video/quicktime"},
					},
				}),
				Err: nil,
			},
		},
		{
			name: "Rewarded",
			params: buildTestParams(
				schema.AdObject{
					Rewarded: &schema.RewardedAdObject{},
				},
			),
			want: createRequestTestOutput{
				Request: buildWantRequest(openrtb2.Imp{
					Instl: 1,
					Banner: &openrtb2.Banner{
						W:     ptr(int64(320)),
						H:     ptr(int64(480)),
						BType: []openrtb2.BannerAdType{},
						BAttr: []adcom1.CreativeAttribute{16},
						Pos:   adcom1.PositionFullScreen.Ptr(),
					},
					Video: &openrtb2.Video{
						MIMEs:     []string{"video/mp4", "video/x-m4v", "video/quicktime", "video/mpeg", "video/avi"},
						Protocols: []adcom1.MediaCreativeSubtype{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14},
						W:         320,
						H:         480,
						BAttr:     []adcom1.CreativeAttribute{16},
						Pos:       adcom1.PositionFullScreen.Ptr(),
					},
					Ext: json.RawMessage(`{"bid_token":"token","rewarded":1}`),
				}),
				Err: nil,
			},
		},
	}

	adapter := buildAdapter()
	for _, tC := range testCases {
		request, err := adapter.CreateRequest(tC.params.BaseBidRequest, tC.params.AuctionRequest)
		if err == nil {
			request.Imp[0].ID = "1" // ommit random uuid
		}
		got := createRequestTestOutput{
			Request: request,
			Err:     err,
		}
		if diff := cmp.Diff(tC.want, got, cmp.Comparer(compareErrors)); diff != "" {
			t.Errorf("%s: adapter.CreateRequest(ctx, %v, %v) mismatch (-want, +got):\n%s", tC.name, tC.params.BaseBidRequest, tC.params.AuctionRequest, diff)
		}
	}
}

func TestBidmachineAdapter_ExecuteRequest(t *testing.T) {
	networkAdapter := buildAdapter()
	responseBody := []byte(`{"key": "value"`)

	customClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.Method != http.MethodPost {
			t.Errorf("Expected POST request")
		}
		if req.URL.String() != "https://api-eu.bidmachine.io/auction/prebid/bidon" {
			t.Errorf("Expected URL: https://api-eu.bidmachine.io/auction/prebid/bidon")
		}
		contentType := req.Header.Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("Expected Content-Type header: application/json")
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

	if response.DemandID != adapter.BidmachineKey {
		t.Errorf("Expected DemandID %v, but got %v", adapter.BigoAdsKey, response.DemandID)
	}
	if response.RequestID != request.ID {
		t.Errorf("Expected RequestID %v, but got %v", request.ID, response.RequestID)
	}
	if response.TagID != "" {
		t.Errorf("Expected TagID be blank, but got %v", response.TagID)
	}
	if response.PlacementID != "" {
		t.Errorf("Expected PlacementID be blank, but got %v", response.PlacementID)
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

func TestBidmachine_ParseBids(t *testing.T) {
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
						"adid": "bmad5e0471131b8a4e3c",
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
						DemandID: "bidmachine",
						AdID:     "bmad5e0471131b8a4e3c",
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
			name: "ParseBids No Conten",
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

func TestBidmachine_Builder(t *testing.T) {
	client := &http.Client{}
	bmCfg := adapter.ProcessedConfigsMap{
		adapter.BidmachineKey: map[string]any{
			"seller_id": "1",
			"endpoint":  "example.com",
		},
	}
	bidder, err := bidmachine.Builder(bmCfg, client)
	wantAdapter := buildAdapter()
	wantBidder := &adapters.Bidder{
		Adapter: &wantAdapter,
		Client:  client,
	}
	if err != nil {
		t.Errorf("Error building adapter: %v", err)
	}
	if diff := cmp.Diff(wantBidder, bidder); diff != "" {
		t.Errorf("builder(bmCfg, client) mismatch (-want, +got):\n%s", diff)
	}
}

func TestBidmachine_ExtraParams(t *testing.T) {
	testCases := []struct {
		name     string
		request  *schema.AuctionRequest
		expected map[string]any
	}{
		{
			name: "Empty SDK version",
			request: &schema.AuctionRequest{
				BaseRequest: schema.BaseRequest{
					App: schema.App{
						SDKVersion: "",
					},
				},
			},
			expected: map[string]any{
				"bidon_sdk_version": "",
				"mediation_mode":    "bidon",
			},
		},
		{
			name: "With SDK version only",
			request: &schema.AuctionRequest{
				BaseRequest: schema.BaseRequest{
					App: schema.App{
						SDKVersion: "1.2.3",
					},
				},
			},
			expected: map[string]any{
				"bidon_sdk_version": "1.2.3",
				"mediation_mode":    "bidon",
			},
		},
		{
			name: "With mediator only",
			request: &schema.AuctionRequest{
				BaseRequest: schema.BaseRequest{
					App: schema.App{
						SDKVersion: "1.0.0",
					},
					Ext: `{"mediator": "max"}`,
				},
			},
			expected: map[string]any{
				"bidon_sdk_version": "1.0.0",
				"mediation_mode":    "bidon_ca",
				"mediator":          "max",
			},
		},
		{
			name: "With nested ext data",
			request: &schema.AuctionRequest{
				BaseRequest: schema.BaseRequest{
					App: schema.App{
						SDKVersion: "2.0.0",
					},
					Ext: `{"ext": {"bidmachine": {"custom_param": "value", "another_param": 123}}}`,
				},
			},
			expected: map[string]any{
				"bidon_sdk_version": "2.0.0",
				"mediation_mode":    "bidon",
				"custom_param":      "value",
				"another_param":     float64(123), // JSON unmarshals numbers as float64
			},
		},
		{
			name: "With all parameters",
			request: &schema.AuctionRequest{
				BaseRequest: schema.BaseRequest{
					App: schema.App{
						SDKVersion: "3.0.0",
					},
					Ext: `{"mediator": "levelplay", "ext": {"bidmachine": {"test_mode": true, "placement_id": "test123"}}}`,
				},
			},
			expected: map[string]any{
				"bidon_sdk_version": "3.0.0",
				"mediation_mode":    "bidon",
				"mediator":          "levelplay",
				"test_mode":         true,
				"placement_id":      "test123",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Normalize the request to parse the Ext field
			tc.request.NormalizeValues()

			result := bidmachine.ExtraParams(tc.request)

			if diff := cmp.Diff(tc.expected, result); diff != "" {
				t.Errorf("ExtraParams() mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}
