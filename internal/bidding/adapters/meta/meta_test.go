package meta_test

import (
	"encoding/json"
	"testing"

	"github.com/bidon-io/bidon-backend/internal/ad"
	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/bidding/adapters/meta"
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

func ptr[T any](t T) *T {
	return &t
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
			"meta": schema.Adapter{
				Version:    "1.0.0",
				SDKVersion: "1.0.0",
			},
		},
		Imp: schema.Imp{
			Demands: map[adapter.Key]map[string]any{
				adapter.MetaKey: {
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
			ID:        "",
			Publisher: &openrtb2.Publisher{ID: "10182906"},
			Ext:       json.RawMessage(`{"orientation":1}`),
		},
		User: &openrtb.User{BuyerUID: "token"},
		Cur:  []string{"USD"},
		Imp: []openrtb2.Imp{
			{
				ID:                "1",
				DisplayManager:    "meta",
				DisplayManagerVer: "1.0.0",
				TagID:             "10182906",
				BidFloorCur:       "USD",
				Secure:            ptr(int8(1)),
			},
		},
		Ext: json.RawMessage(`{"authentication_id":"b613679a0814d9ec772f95d778c35fc5ff1697c493715653c6c712144292c5ad","platformid":"687579938617452"}`),
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

func TestMeta_CreateRequestTest(t *testing.T) {
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
					Ext: nil,
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
					Ext: nil,
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
					Ext: nil,
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
						MIMEs: nil,
						W:     int64(320),
						H:     int64(480),
						Ext:   json.RawMessage(`{"videotype": "rewarded"}`),
					},
					Ext: nil,
				}),
				Err: nil,
			},
		},
	}

	adapter := &meta.MetaAdapter{
		AppID: "10182906",
		TagID: "10182906",
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
