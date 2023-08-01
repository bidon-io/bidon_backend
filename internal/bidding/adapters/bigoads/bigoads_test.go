package bigoads_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/bidon-io/bidon-backend/internal/ad"
	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/bidding/adapters"
	"github.com/bidon-io/bidon-backend/internal/bidding/adapters/bigoads"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
	"github.com/google/go-cmp/cmp"
	"github.com/prebid/openrtb/v19/adcom1"
	"github.com/prebid/openrtb/v19/openrtb2"
)

type createRequestTestParams struct {
	BaseBidRequest openrtb2.BidRequest
	Br             *schema.BiddingRequest
}

type createRequestTestOutput struct {
	Request openrtb2.BidRequest
	Err     []error
}

type ParseBidsTestParams struct {
	DemandsResponse adapters.DemandResponse
}

type ParseBidsTestOutput struct {
	DemandResponse adapters.DemandResponse
	Err            error
}

func ptr[T any](t T) *T {
	return &t
}

// compareErrors checks for error occurrence.
func compareErrors(want, got error) bool {
	return (want == nil) == (got == nil)
}

func buildAdapter() bigoads.BigoAdsAdapter {
	return bigoads.BigoAdsAdapter{
		SellerID:    "1",
		AppID:       "10182906",
		TagID:       "10182906-10192212",
		PlacementID: "10182906-10192212",
	}
}

func buildBaseRequest() openrtb2.BidRequest {
	return openrtb2.BidRequest{
		App: &openrtb2.App{
			Publisher: &openrtb2.Publisher{},
		},
	}
}

func buildTestParams(imp schema.Imp) createRequestTestParams {
	request := buildBaseRequest()

	br := schema.BiddingRequest{
		Adapters: schema.Adapters{
			"bigoads": schema.Adapter{
				Version:    "1.0.0",
				SDKVersion: "1.0.0",
			},
		},
		Imp: schema.Imp{
			Demands: map[adapter.Key]map[string]any{
				adapter.BigoAdsKey: {
					"token": "token",
				},
			},
			Orientation: "PORTRAIT",
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

func buildWantRequest(imp openrtb2.Imp) openrtb2.BidRequest {
	request := openrtb2.BidRequest{
		App:  &openrtb2.App{ID: "10182906", Publisher: &openrtb2.Publisher{ID: "1"}},
		User: &openrtb2.User{BuyerUID: "token"},
		Cur:  []string{"USD"},
		Imp: []openrtb2.Imp{
			{
				ID:                "1",
				DisplayManager:    "bigoads",
				DisplayManagerVer: "1.0.0",
				TagID:             "10182906-10192212",
				Secure:            ptr(int8(1)),
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

func TestBigoAds_CreateRequest(t *testing.T) {
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
					Ext: json.RawMessage(`{"adtype":2,"networkid":{"appid":"10182906","placementid":"10182906-10192212"}}`),
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
					Ext: json.RawMessage(`{"adtype":2,"networkid":{"appid":"10182906","placementid":"10182906-10192212"}}`),
				}),
				Err: nil,
			},
		},
		{
			name: "Banner unsupported format",
			params: buildTestParams(
				schema.Imp{
					Banner: &schema.BannerAdObject{
						Format: ad.LeaderboardFormat,
					},
				},
			),
			want: createRequestTestOutput{
				Request: buildBaseRequest(),
				Err:     []error{errors.New("unknown banner format")},
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
					Ext: json.RawMessage(`{"adtype":3,"networkid":{"appid":"10182906","placementid":"10182906-10192212"}}`),
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
					},
					Ext: json.RawMessage(`{"adtype":4,"networkid":{"appid":"10182906","placementid":"10182906-10192212"}}`),
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

func TestBigoAds_ParseBids(t *testing.T) {
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
						"adid": "Bigoad5e0471131b8a4e3c",
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
						DemandID: "bigoads",
						AdID:     "Bigoad5e0471131b8a4e3c",
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

func TestBigoAds_Builder(t *testing.T) {
	client := &http.Client{}
	bigoCfg := adapter.ProcessedConfigsMap{
		adapter.BigoAdsKey: map[string]string{
			"seller_id":    "1",
			"app_id":       "10182906",
			"tag_id":       "10182906-10192212",
			"placement_id": "10182906-10192212",
		},
	}
	bidder, err := bigoads.Builder(bigoCfg, client)
	wantAdapter := buildAdapter()
	wantBidder := adapters.Bidder{
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
