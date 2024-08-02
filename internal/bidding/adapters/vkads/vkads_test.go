package vkads_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/bidon-io/bidon-backend/internal/ad"
	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/bidding/adapters"
	"github.com/bidon-io/bidon-backend/internal/bidding/adapters/vkads"
	"github.com/bidon-io/bidon-backend/internal/bidding/openrtb"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
	"github.com/google/go-cmp/cmp"
	"github.com/prebid/openrtb/v19/adcom1"
	"github.com/prebid/openrtb/v19/openrtb2"
	"io/ioutil"
	"net/http"
	"testing"
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

func compareErrors(want, got error) bool {
	return (want == nil) == (got == nil)
}

func ptr[T any](t T) *T {
	return &t
}

func buildAdapter() vkads.VKAdsAdapter {
	return vkads.VKAdsAdapter{
		AppID: "10182906",
		TagID: "10182906-10192212",
	}
}

func buildBaseRequest() openrtb.BidRequest {
	return openrtb.BidRequest{
		App: &openrtb2.App{
			Publisher: &openrtb2.Publisher{},
		},
	}
}

func buildTestParams(imp schema.Imp) createRequestTestParams {
	request := buildBaseRequest()

	br := schema.BiddingRequest{
		Adapters: schema.Adapters{
			"vkads": schema.Adapter{
				Version:    "1.0.0",
				SDKVersion: "1.0.0",
			},
		},
		Imp: schema.Imp{
			Demands: map[adapter.Key]map[string]any{
				adapter.VKAdsKey: {
					"token": "token",
				},
			},
			Orientation: "PORTRAIT",
			BidFloor:    ptr(schema.MinBidFloor),
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
		App:  &openrtb2.App{ID: "10182906", Publisher: &openrtb2.Publisher{}},
		User: &openrtb.User{Ext: json.RawMessage(`{"buyeruid": "token"}`)},
		Cur:  []string{"USD"},
		Imp:  []openrtb2.Imp{imp},
		Ext:  json.RawMessage(`{"pid":111}`),
	}

	request.Imp[0].TagID = "10182906-10192212"
	request.Imp[0].ID = "1" // ommit random uuid

	return request
}

func TestVKAds_CreateRequest(t *testing.T) {
	testCases := []struct {
		name   string
		params createRequestTestParams
		want   createRequestTestOutput
	}{
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
				}),
				Err: nil,
			},
		},
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
				}),
				Err: nil,
			},
		},
		{
			name: "Banner LEADERBOARD",
			params: buildTestParams(
				schema.Imp{
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
						Pos: adcom1.PositionFullScreen.Ptr(),
					},
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
					Banner: &openrtb2.Banner{
						W: ptr(int64(1920)),
						H: ptr(int64(1080)),
					},
					BidFloor:    schema.MinBidFloor,
					BidFloorCur: "USD",
				}),
				Err: nil,
			},
		},
	}

	adapter := buildAdapter()
	for _, tC := range testCases {
		request, err := adapter.CreateRequest(tC.params.BaseBidRequest, tC.params.Br)
		if err == nil {
			request.Imp[0].ID = "1" // ommit random uuid
		}
		got := createRequestTestOutput{
			Request: request,
			Err:     err,
		}
		if diff := cmp.Diff(tC.want, got, cmp.Comparer(compareErrors)); diff != "" {
			t.Errorf("%s: adapter.CreateRequest(ctx, %v, %v) mismatch (-want, +got):\n%s", tC.name, tC.params.BaseBidRequest, tC.params.Br, diff)
		}
	}
}

func TestAdapter_ExecuteRequest(t *testing.T) {
	a := buildAdapter()

	testCases := []struct {
		name                string
		responseBody        []byte
		expectedURL         string
		expectedDemandID    adapter.Key
		expectedRequestID   string
		expectedTagID       string
		expectedStatus      int
		expectedRawResponse string
	}{
		{
			name:                "Valid Request",
			responseBody:        []byte(`{"key": "value"}`),
			expectedURL:         "https://ad.mail.ru/api/bid",
			expectedDemandID:    adapter.VKAdsKey,
			expectedRequestID:   "test-request-id",
			expectedTagID:       a.TagID,
			expectedStatus:      http.StatusOK,
			expectedRawResponse: `{"key": "value"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			customClient := NewTestClient(func(req *http.Request) *http.Response {
				if req.Method != http.MethodPost {
					t.Errorf("Expected POST request, got %s", req.Method)
				}
				if req.URL.String() != tc.expectedURL {
					t.Errorf("Expected URL %s, got %s", tc.expectedURL, req.URL.String())
				}
				if contentType := req.Header.Get("Content-Type"); contentType != "application/json" {
					t.Errorf("Expected Content-Type application/json, got %s", contentType)
				}
				return &http.Response{
					Status:        http.StatusText(http.StatusOK),
					StatusCode:    http.StatusOK,
					Body:          ioutil.NopCloser(bytes.NewBuffer(tc.responseBody)),
					ContentLength: int64(len(tc.responseBody)),
				}
			})

			request := openrtb.BidRequest{
				ID: tc.expectedRequestID,
			}

			response := a.ExecuteRequest(context.Background(), customClient, request)

			if response.DemandID != tc.expectedDemandID {
				t.Errorf("Expected DemandID %v, but got %v", tc.expectedDemandID, response.DemandID)
			}
			if response.RequestID != tc.expectedRequestID {
				t.Errorf("Expected RequestID %v, but got %v", tc.expectedRequestID, response.RequestID)
			}
			if response.TagID != tc.expectedTagID {
				t.Errorf("Expected TagID %v, but got %v", tc.expectedTagID, response.TagID)
			}
			if response.RawResponse != tc.expectedRawResponse {
				t.Errorf("Expected RawResponse %v, but got %v", tc.expectedRawResponse, response.RawResponse)
			}
			if response.Status != tc.expectedStatus {
				t.Errorf("Expected status %d, but got %d", tc.expectedStatus, response.Status)
			}
		})
	}
}

func TestVKAds_ParseBids(t *testing.T) {
	rawResponse := `{
	  "cur": "USD",
	  "bidid": "fd2a9b71-7b6e-4beb-84d2-2d47c8bf0c9b",
	  "seatbid": [
		{
		  "bid": [
			{
			  "price": 1.5,
			  "id": "2:1::669e29a737559894",
			  "adm": "adm",
			  "impid": "7703af66-0ec1-475f-b5a8-eda9d65c44e6",
			  "nurl": "https://rs.mail.ru/pixel/Q.gif?price=37.27&currency=RUB",
			  "cid": "96295833",
			  "lurl": "https://rs.mail.ru",
			  "adomain": [
				"ya.ru"
			  ],
			  "adid": "162456424",
			  "crid": "1:162456424",
			  "language": "ru"
			}
		  ]
		}
	  ],
	  "id": "fd2a9b71-7b6e-4beb-84d2-2d47c8bf0c9b"
	}`
	a := buildAdapter()

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
						ID:       "2:1::669e29a737559894",
						ImpID:    "7703af66-0ec1-475f-b5a8-eda9d65c44e6",
						Price:    1.5,
						DemandID: "vkads",
						AdID:     "162456424",
						LURL:     "https://rs.mail.ru",
						NURL:     "https://rs.mail.ru/pixel/Q.gif?price=37.27&currency=RUB",
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
		response, err := a.ParseBids(&tC.params.DemandsResponse)
		got := ParseBidsTestOutput{
			DemandResponse: *response,
			Err:            err,
		}
		if diff := cmp.Diff(tC.want, got, cmp.Comparer(compareErrors)); diff != "" {
			t.Errorf("%s: a.ParseBids(ctx, %v) mismatch (-want, +got):\n%s", tC.name, tC.params.DemandsResponse, diff)
		}
	}
}

func TestVKAds_Builder(t *testing.T) {
	client := &http.Client{}

	testCases := []struct {
		name           string
		cfg            adapter.ProcessedConfigsMap
		expectedBidder *adapters.Bidder
		expectedErr    error
	}{
		{
			name: "Valid Configuration",
			cfg: adapter.ProcessedConfigsMap{
				adapter.VKAdsKey: map[string]any{
					"app_id": "10182906",
					"tag_id": "10182906-10192212",
				},
			},
			expectedBidder: &adapters.Bidder{
				Adapter: ptr(buildAdapter()),
				Client:  client,
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			bidder, err := vkads.Builder(tc.cfg, client)

			if !errors.Is(err, tc.expectedErr) {
				t.Fatalf("Expected error %v, but got %v", tc.expectedErr, err)
			}

			if diff := cmp.Diff(tc.expectedBidder, bidder); diff != "" {
				t.Errorf("builder(cfg, client) mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}
