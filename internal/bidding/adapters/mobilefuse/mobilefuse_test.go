package mobilefuse_test

import (
	"bytes"
	"context"
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
	"github.com/bidon-io/bidon-backend/internal/bidding/adapters/mobilefuse"
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

func ptr[T any](t T) *T {
	return &t
}

func compareErrors(want, got error) bool {
	return (want == nil) == (got == nil)
}

func buildAdapter() mobilefuse.MobileFuseAdapter {
	return mobilefuse.MobileFuseAdapter{
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

func buildTestParams(adObject schema.AdObject) createRequestTestParams {
	request := buildBaseRequest()

	auctionRequest := schema.AuctionRequest{
		Adapters: schema.Adapters{
			"mobilefuse": schema.Adapter{
				Version:    "1.0.0",
				SDKVersion: "1.0.0",
			},
		},
		AdObject: schema.AdObject{
			Demands: map[adapter.Key]map[string]any{
				adapter.MobileFuseKey: {
					"token": "token",
				},
			},
			Orientation: "PORTRAIT",
		},
	}

	if adObject.Banner != nil {
		auctionRequest.AdObject.Banner = adObject.Banner
	}
	if adObject.Interstitial != nil {
		auctionRequest.AdObject.Interstitial = adObject.Interstitial
	}
	if adObject.Rewarded != nil {
		auctionRequest.AdObject.Rewarded = adObject.Rewarded
	}

	return createRequestTestParams{
		BaseBidRequest: request,
		AuctionRequest: &auctionRequest,
	}
}

func buildWantRequest(imp openrtb2.Imp) openrtb.BidRequest {
	request := openrtb.BidRequest{
		App: &openrtb2.App{Publisher: &openrtb2.Publisher{}},
		User: &openrtb.User{
			Data: []openrtb.Data{
				{
					Segment: []openrtb.Segment{
						{
							Signal: "token",
						},
					},
				},
			},
		},
		Cur: []string{"USD"},
		Imp: []openrtb2.Imp{
			{
				ID:                "1",
				DisplayManager:    "mobilefuse",
				DisplayManagerVer: "1.0.0",
				TagID:             "10182906-10192212",
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

func TestMobileFuse_CreateRequest(t *testing.T) {
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
					Ext: nil,
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
					Ext: nil,
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
						Pos: adcom1.PositionFullScreen.Ptr(),
					},
					Ext: nil,
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
					Instl: 0,
					Video: &openrtb2.Video{
						MIMEs: []string{"video/mp4"},
					},
					Ext: nil,
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

func TestMobileFuseAdapter_ExecuteRequest(t *testing.T) {
	networkAdapter := buildAdapter()
	responseBody := []byte(`{"key": "value"`)

	customClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.Method != http.MethodPost {
			t.Errorf("Expected POST request")
		}
		if req.URL.String() != "https://mfx.mobilefuse.com/openrtb?ssp=bidon" {
			t.Errorf("Expected URL: https://mfx.mobilefuse.com/openrtb?ssp=bidon")
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

	if response.DemandID != adapter.MobileFuseKey {
		t.Errorf("Expected DemandID %v, but got %v", adapter.MobileFuseKey, response.DemandID)
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

func TestMobileFuse_ParseBids(t *testing.T) {
	rawResponse := `{
	  "id": "32612272f2a17a9789aee929a9a20b94",
	  "bidid": "f5bec35899734ce7f17e4bdc56c91117",
	  "seatbid": [
		{
		  "bid": [
			{
			  "id": "9f21299b7b7b6ad7fe30226748c57abf_banner",
			  "adomain": [
				"mobilefuse.com"
			  ],
			  "impid": "1",
			  "crid": "test_js_interstitial",
			  "h": 480,
			  "price": 1.904,
			  "w": 320,
			  "lurl": "https://mfx-us-east-1.mobilefuse.com/lurl?i=60f859c11e497f9afb51877bec9c391c_0&loss=${AUCTION_LOSS}&price=${AUCTION_PRICE}",
			  "ext": {
				"mf": {
				  "format": "instl",
				  "media_type": "banner"
				},
				"signaldata": "H4sIAAAAAAAAA41Qy27CMBD8lcqqOBESm0dwBFQV6gGpKhW0t0iRYzbCkMSW16FUiH/v0l565LaPmZ3ZuTCzYxkbigkXIhWVUDxVMp1KBSCFVFKJpJQj1mfaNSzjA5mMqPa/rAAYigMWpg3gMZhgVE3I8O2AtqVqW/DUqx0x2QyJ5cIDej3P2T4Eh1ke53FTnaMOI1AYIj5obGlqqDqEgbZNHh8w+vLKOfBPZj5JqulYas5hJNNKqqoc82malqClHkquiyRni1ke/yktSBrOznhAlg0nSdJntUUsOl+TnXsN1AS/S7p3Oz5/vDx/Lj9W67fidb3dXnvOGw3/pu+b1fLlSs6aLgBlyCnM2uhjUcJenYz1t+Bs2BOirC3N8WhcgaBtu6M3xn12Ypm4/gBBRP5ctwEAAA=="
			  },
			  "exp": 3600
			}
		  ],
		  "seat": ""
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
						ID:         "9f21299b7b7b6ad7fe30226748c57abf_banner",
						ImpID:      "1",
						Price:      1.904,
						Payload:    "",
						Signaldata: "H4sIAAAAAAAAA41Qy27CMBD8lcqqOBESm0dwBFQV6gGpKhW0t0iRYzbCkMSW16FUiH/v0l565LaPmZ3ZuTCzYxkbigkXIhWVUDxVMp1KBSCFVFKJpJQj1mfaNSzjA5mMqPa/rAAYigMWpg3gMZhgVE3I8O2AtqVqW/DUqx0x2QyJ5cIDej3P2T4Eh1ke53FTnaMOI1AYIj5obGlqqDqEgbZNHh8w+vLKOfBPZj5JqulYas5hJNNKqqoc82malqClHkquiyRni1ke/yktSBrOznhAlg0nSdJntUUsOl+TnXsN1AS/S7p3Oz5/vDx/Lj9W67fidb3dXnvOGw3/pu+b1fLlSs6aLgBlyCnM2uhjUcJenYz1t+Bs2BOirC3N8WhcgaBtu6M3xn12Ypm4/gBBRP5ctwEAAA==",
						DemandID:   "mobilefuse",
						AdID:       "",
						LURL:       "https://mfx-us-east-1.mobilefuse.com/lurl?i=60f859c11e497f9afb51877bec9c391c_0&loss=${AUCTION_LOSS}&price=${AUCTION_PRICE}",
						NURL:       "",
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

func TestMobileFuse_Builder(t *testing.T) {
	client := &http.Client{}
	mobileFuseCfg := adapter.ProcessedConfigsMap{
		adapter.MobileFuseKey: map[string]any{
			"tag_id": "10182906-10192212",
		},
	}
	bidder, err := mobilefuse.Builder(mobileFuseCfg, client)
	wantAdapter := buildAdapter()
	wantBidder := &adapters.Bidder{
		Adapter: &wantAdapter,
		Client:  client,
	}
	if err != nil {
		t.Errorf("Error building adapter: %v", err)
	}
	if diff := cmp.Diff(wantBidder, bidder); diff != "" {
		t.Errorf("builder(mobileFuseCfg, client) mismatch (-want, +got):\n%s", diff)
	}
}
