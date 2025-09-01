package inmobi_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/prebid/openrtb/v19/adcom1"
	"github.com/prebid/openrtb/v19/openrtb2"

	"github.com/bidon-io/bidon-backend/internal/ad"
	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/bidding/adapters"
	"github.com/bidon-io/bidon-backend/internal/bidding/adapters/inmobi"
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

func buildAdapter() inmobi.InMobiAdapter {
	return inmobi.InMobiAdapter{
		AppID:       "test-app-id",
		PlacementID: "1621323861540",
	}
}

func buildBaseRequest() openrtb.BidRequest {
	return openrtb.BidRequest{
		App: &openrtb2.App{
			Publisher: &openrtb2.Publisher{},
		},
		Device: &openrtb2.Device{
			Geo: &openrtb2.Geo{
				Country: "USA", // Default to USA for existing tests (ISO-3166-1-alpha-3)
			},
		},
	}
}

func buildTestParams(adObject schema.AdObject) createRequestTestParams {
	request := buildBaseRequest()

	auctionRequest := schema.AuctionRequest{
		Adapters: schema.Adapters{
			"inmobi": schema.Adapter{
				Version:    "1.0.0",
				SDKVersion: "1.0.0",
			},
		},
		AdObject: schema.AdObject{
			Demands: map[adapter.Key]map[string]any{
				adapter.InmobiKey: {
					"token": "test-token",
				},
			},
			Orientation: "PORTRAIT",
		},
	}

	auctionRequest.AdObject.AuctionID = adObject.AuctionID
	auctionRequest.AdObject.AuctionConfigurationID = adObject.AuctionConfigurationID
	auctionRequest.AdObject.AuctionConfigurationUID = adObject.AuctionConfigurationUID
	auctionRequest.AdObject.PriceFloor = adObject.PriceFloor

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

func TestInMobiAdapter_CreateRequest(t *testing.T) {
	tests := []struct {
		name   string
		params createRequestTestParams
		want   createRequestTestOutput
	}{
		{
			name: "Banner request",
			params: buildTestParams(schema.AdObject{
				AuctionID:               "test-auction-id",
				AuctionConfigurationID:  1,
				AuctionConfigurationUID: "test-config-uid",
				PriceFloor:              0.1,
				Banner: &schema.BannerAdObject{
					Format: ad.BannerFormat,
				},
			}),
			want: createRequestTestOutput{
				Request: openrtb.BidRequest{
					App: &openrtb2.App{
						ID:        "test-app-id",
						Publisher: &openrtb2.Publisher{},
					},
					Device: &openrtb2.Device{
						Geo: &openrtb2.Geo{
							Country: "USA",
						},
					},
					Imp: []openrtb2.Imp{
						{
							TagID:             "1621323861540",
							DisplayManager:    "inmobi",
							DisplayManagerVer: "1.0.0",
							Instl:             0,
							Banner: &openrtb2.Banner{
								W:   ptr(int64(320)),
								H:   ptr(int64(50)),
								Pos: adcom1.PositionAboveFold.Ptr(),
								API: []adcom1.APIFramework{3, 5}, // MRAID 1.0, MRAID 2.0
							},
							Secure:      ptr(int8(1)),
							BidFloor:    0.100001,
							BidFloorCur: "USD",
						},
					},
					Cur: []string{"USD"},
					User: &openrtb.User{
						BuyerUID: "test-token",
					},
				},
				Err: nil,
			},
		},
		{
			name: "Missing placement ID",
			params: func() createRequestTestParams {
				params := buildTestParams(schema.AdObject{
					AuctionID:               "test-auction-id",
					AuctionConfigurationID:  1,
					AuctionConfigurationUID: "test-config-uid",
					PriceFloor:              0.1,
					Banner: &schema.BannerAdObject{
						Format: ad.BannerFormat,
					},
				})
				return params
			}(),
			want: createRequestTestOutput{
				Request: openrtb.BidRequest{},
				Err:     errors.New("PlacementID is empty"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := buildAdapter()
			if tt.name == "Missing placement ID" {
				adapter.PlacementID = ""
			}

			got, err := adapter.CreateRequest(tt.params.BaseBidRequest, tt.params.AuctionRequest)

			if !compareErrors(tt.want.Err, err) {
				t.Errorf("CreateRequest() error = %v, wantErr %v", err, tt.want.Err)
				return
			}

			if err == nil {
				// Clear dynamic fields for comparison
				if len(got.Imp) > 0 {
					got.Imp[0].ID = ""
				}

				if diff := cmp.Diff(tt.want.Request, got); diff != "" {
					t.Errorf("CreateRequest() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestInMobiAdapter_Builder(t *testing.T) {
	tests := []struct {
		name    string
		cfg     adapter.ProcessedConfigsMap
		wantErr bool
	}{
		{
			name: "Valid configuration",
			cfg: adapter.ProcessedConfigsMap{
				adapter.InmobiKey: map[string]any{
					"app_id":       "test-app-id",
					"placement_id": "test-placement-id",
				},
			},
			wantErr: false,
		},
		{
			name: "Missing app_id",
			cfg: adapter.ProcessedConfigsMap{
				adapter.InmobiKey: map[string]any{
					"placement_id": "test-placement-id",
				},
			},
			wantErr: true,
		},
		{
			name: "Missing placement_id",
			cfg: adapter.ProcessedConfigsMap{
				adapter.InmobiKey: map[string]any{
					"app_id": "test-app-id",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := inmobi.Builder(tt.cfg, &http.Client{})
			if (err != nil) != tt.wantErr {
				t.Errorf("Builder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
