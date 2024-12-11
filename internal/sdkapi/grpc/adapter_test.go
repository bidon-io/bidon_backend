package grpcserver

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/bidon-io/bidon-backend/internal/adapter"

	adcom "github.com/bidon-io/bidon-backend/pkg/proto/com/iabtechlab/adcom/v1"
	adcomctx "github.com/bidon-io/bidon-backend/pkg/proto/com/iabtechlab/adcom/v1/context"
	pbctx "github.com/bidon-io/bidon-backend/pkg/proto/org/bidon/proto/v1/context"

	"google.golang.org/protobuf/proto"

	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
	v3 "github.com/bidon-io/bidon-backend/pkg/proto/com/iabtechlab/openrtb/v3"

	"github.com/bidon-io/bidon-backend/pkg/proto/org/bidon/proto/v1/mediation"
)

func ptr[T any](t T) *T {
	return &t
}

func TestAuctionRequestAdapter_OpenRTBToAuctionRequest(t *testing.T) {
	a := NewAuctionRequestAdapter()

	buildValidRequest := func() *v3.Openrtb {
		app := &adcomctx.DistributionChannel_App{
			Bundle: proto.String("com.example.app"),
			Ver:    proto.String("1.0"),
		}
		appExt := &mediation.AppExt{
			Key:              proto.String("app_key"),
			Framework:        proto.String("flutter"),
			FrameworkVersion: proto.String("1.22"),
			PluginVersion:    proto.String("1.0.0"),
			SdkVersion:       proto.String("2.0.0"),
			Skadn:            []string{"skadn1", "skadn2"},
		}
		proto.SetExtension(app, mediation.E_AppExt, appExt)

		user := &adcomctx.User{}
		userExt := &mediation.UserExt{
			Idfa:                        proto.String("IDFA-12345"),
			Idfv:                        proto.String("IDFV-12345"),
			Idg:                         proto.String("IDG-12345"),
			TrackingAuthorizationStatus: proto.String("authorized"),
			Segments: []*mediation.Segment{
				{
					Id:  proto.String("segment_id"),
					Uid: proto.String("segment_uid"),
				},
			},
		}
		proto.SetExtension(user, mediation.E_UserExt, userExt)

		regs := &adcomctx.Regs{
			Coppa: proto.Bool(true),
			Gdpr:  proto.Bool(true),
		}
		regsExt := &mediation.RegsExt{
			UsPrivacy: proto.String("1YNN"),
			EuPrivacy: proto.String("1"),
			Iab:       proto.String(`{"key":"value"}`),
		}
		proto.SetExtension(regs, mediation.E_RegsExt, regsExt)

		device := &adcomctx.Device{
			Ua:    proto.String("Mozilla/5.0"),
			Make:  proto.String("Apple"),
			Model: proto.String("iPhone"),
			Os:    proto.Int32(int32(adcom.OperatingSystem_IOS)),
			Osv:   proto.String("14.4"),
		}
		deviceExt := &mediation.DeviceExt{
			Id:          proto.String("session_id"),
			LaunchTs:    proto.Int32(1617187200),
			RamUsed:     proto.Int32(1024),
			RamSize:     proto.Int32(2048),
			StorageFree: proto.Int32(512),
			StorageUsed: proto.Int32(256),
			Battery:     proto.Float64(80.5),
			CpuUsage:    proto.Float64(10.6),
		}
		proto.SetExtension(device, mediation.E_DeviceExt, deviceExt)

		c := &pbctx.Context{
			DistributionChannel: &adcomctx.DistributionChannel{
				ChannelOneof: &adcomctx.DistributionChannel_App_{
					App: app,
				},
			},
			Device: device,
			User:   user,
			Regs:   regs,
		}
		ctxBytes, _ := proto.Marshal(c)

		placement := &adcom.Placement{}
		placementExt := &mediation.PlacementExt{
			AuctionId:               proto.String("auction_id_123"),
			AuctionKey:              proto.String("auction_key_789"),
			AuctionConfigurationUid: proto.String("config_uid_456"),
			Orientation:             ptr(mediation.Orientation_PORTRAIT),
			Demands: map[string]*mediation.Demand{
				"demand_key": {
					Token:         proto.String("token_value"),
					Status:        proto.String("status_value"),
					TokenFinishTs: proto.Int64(1234567890),
					TokenStartTs:  proto.Int64(1234567000),
				},
			},
		}
		proto.SetExtension(placement, mediation.E_PlacementExt, placementExt)

		placementBytes, err := proto.Marshal(placement)
		if err != nil {
			t.Fatalf("failed to marshal placement: %v", err)
		}

		item := &v3.Item{
			Id:   proto.String("auction_id_123"),
			Flr:  proto.Float32(0.5),
			Spec: placementBytes,
		}

		return &v3.Openrtb{
			PayloadOneof: &v3.Openrtb_Request{
				Request: &v3.Request{
					Test:    proto.Bool(true),
					Tmax:    proto.Uint32(1000),
					Context: ctxBytes,
					Item:    []*v3.Item{item},
				},
			},
		}
	}

	tests := []struct {
		name   string
		input  *v3.Openrtb
		want   *schema.AuctionV2Request
		errMsg string
	}{
		{
			name:  "valid request with extensions",
			input: buildValidRequest(),
			want: &schema.AuctionV2Request{
				TMax: 1000,
				Test: true,
				BaseRequest: schema.BaseRequest{
					Device: schema.Device{
						Geo:            &schema.Geo{},
						UserAgent:      "Mozilla/5.0",
						Manufacturer:   "Apple",
						Model:          "iPhone",
						OS:             "iOS",
						OSVersion:      "14.4",
						JS:             func() *int { i := 0; return &i }(),
						ConnectionType: "ConnectionType_UNKNOWN",
						Type:           "DeviceType_UNKNOWN",
					},
					Session: schema.Session{
						ID:          "session_id",
						LaunchTS:    1617187200,
						RAMUsed:     1024,
						RAMSize:     2048,
						StorageFree: 512,
						StorageUsed: 256,
						Battery:     80.5,
						CPUUsage:    func() *float64 { f := 10.6; return &f }(),
					},
					App: schema.App{
						Bundle:           "com.example.app",
						Key:              "app_key",
						Framework:        "flutter",
						Version:          "1.0",
						FrameworkVersion: "1.22",
						PluginVersion:    "1.0.0",
						SKAdN:            []string{"skadn1", "skadn2"},
						SDKVersion:       "2.0.0",
					},
					User: schema.User{
						IDFA:                        "IDFA-12345",
						TrackingAuthorizationStatus: "authorized",
						IDFV:                        "IDFV-12345",
						IDG:                         "IDG-12345",
						Consent:                     map[string]any{},
						COPPA:                       nil,
					},
					Geo: &schema.Geo{},
					Regulations: &schema.Regulations{
						COPPA:     true,
						GDPR:      true,
						USPrivacy: "1YNN",
						EUPrivacy: "1",
						IAB:       map[string]any{"key": "value"},
					},
					Ext:   "",
					Token: "",
					Segment: schema.Segment{
						ID:  "segment_id",
						UID: "segment_uid",
						Ext: "",
					},
				},
				AdObject: schema.AdObjectV2{
					AuctionID:               "auction_id_123",
					AuctionKey:              "auction_key_789",
					AuctionConfigurationID:  0,
					AuctionConfigurationUID: "config_uid_456",
					PriceFloor:              0.5,
					Orientation:             "PORTRAIT",
					Demands: map[adapter.Key]map[string]any{
						"demand_key": {
							"token":           "token_value",
							"status":          "status_value",
							"token_start_ts":  int64(1234567000),
							"token_finish_ts": int64(1234567890),
						},
					},
					Banner:       nil,
					Interstitial: nil,
					Rewarded:     nil,
				},
				AdCache: nil,
			},
		},
		{
			name:   "nil request",
			input:  &v3.Openrtb{},
			errMsg: "request is nil",
		},
		{
			name: "empty context",
			input: &v3.Openrtb{
				PayloadOneof: &v3.Openrtb_Request{
					Request: &v3.Request{
						Item: []*v3.Item{{Id: proto.String("some_id")}},
					},
				},
			},
			errMsg: "request context is empty",
		},
		{
			name: "no items",
			input: func() *v3.Openrtb {
				r := buildValidRequest()
				r.GetRequest().Item = nil
				return r
			}(),
			errMsg: "no items in request",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ar, err := a.OpenRTBToAuctionRequest(tc.input)
			if (err != nil) != (tc.errMsg != "") {
				t.Fatalf("expected error=%s, got %v", tc.errMsg, err)
			}

			if tc.errMsg != "" && err != nil {
				if msg := err.Error(); !strings.Contains(msg, tc.errMsg) {
					t.Errorf("expected error containing 'no items in request', got %q", msg)
				}
			}

			if diff := cmp.Diff(tc.want, ar, cmpopts.EquateEmpty(), cmpopts.IgnoreUnexported(schema.AuctionV2Request{}, schema.BaseRequest{})); diff != "" {
				t.Errorf("AuctionV2Request mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
