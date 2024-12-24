package grpcserver

import (
	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/auction"
	"github.com/bidon-io/bidon-backend/internal/auction/store"
	"github.com/bidon-io/bidon-backend/internal/bidding/adapters"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
	"github.com/bidon-io/bidon-backend/internal/segment"
	adcom "github.com/bidon-io/bidon-backend/pkg/proto/com/iabtechlab/adcom/v1"
	adcomctx "github.com/bidon-io/bidon-backend/pkg/proto/com/iabtechlab/adcom/v1/context"
	v3 "github.com/bidon-io/bidon-backend/pkg/proto/com/iabtechlab/openrtb/v3"
	pbctx "github.com/bidon-io/bidon-backend/pkg/proto/org/bidon/proto/v1/context"
	"google.golang.org/protobuf/proto"

	"github.com/bidon-io/bidon-backend/pkg/proto/org/bidon/proto/v1/mediation"
)

// Request helper
type RequestBuilder struct {
	req *v3.Request
}

func NewRequestBuilder() *RequestBuilder {
	return &RequestBuilder{
		req: defaultValidReq(),
	}
}

func (rb *RequestBuilder) Build() *v3.Openrtb {
	return &v3.Openrtb{
		PayloadOneof: &v3.Openrtb_Request{
			Request: rb.req,
		},
	}
}

func defaultValidReq() *v3.Request {
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
	sgmnt := DefaultSegment()
	userExt := &mediation.UserExt{
		Idfa:                        proto.String("IDFA-12345"),
		Idfv:                        proto.String("IDFV-12345"),
		Idg:                         proto.String("IDG-12345"),
		TrackingAuthorizationStatus: proto.String("authorized"),
		Segments: []*mediation.Segment{
			{
				Id:  proto.String(sgmnt.StringID()),
				Uid: proto.String(sgmnt.UID),
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
		Ip:    proto.String("8.8.8.8"),
		Ua:    proto.String("Mozilla/5.0"),
		Make:  proto.String("Apple"),
		Model: proto.String("iPhone"),
		Os:    proto.Int32(int32(adcom.OperatingSystem_IOS)),
		Osv:   proto.String("14.4"),
	}
	deviceExt := &mediation.DeviceExt{
		Id:          proto.String("session_id"),
		LaunchTs:    proto.Int64(1617187200),
		RamUsed:     proto.Int64(1024),
		RamSize:     proto.Int64(2048),
		StorageFree: proto.Int64(512),
		StorageUsed: proto.Int64(256),
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
		AuctionKey:              proto.String("1F60CVMI00400"),
		AuctionConfigurationUid: proto.String("config_uid_456"),
		Demands: map[string]*mediation.Demand{
			"demand_key": {
				Token:         proto.String("token_value"),
				Status:        proto.String("status_value"),
				TokenFinishTs: proto.Int64(1234567890),
				TokenStartTs:  proto.Int64(1234567000),
			},
		},
	}

	placement.Display = &adcom.Placement_DisplayPlacement{}
	displayPlacementExt := &mediation.DisplayPlacementExt{
		Format:      ptr(mediation.AdFormat_BANNER),
		Orientation: ptr(mediation.Orientation_PORTRAIT),
	}
	proto.SetExtension(placement.Display, mediation.E_DisplayPlacementExt, displayPlacementExt)

	proto.SetExtension(placement, mediation.E_PlacementExt, placementExt)

	placementBytes, _ := proto.Marshal(placement)
	item := &v3.Item{
		Id:   proto.String("auction_id_123"),
		Flr:  proto.Float32(0.5),
		Spec: placementBytes,
	}

	req := &v3.Request{
		Test:    proto.Bool(true),
		Tmax:    proto.Uint32(1000),
		Context: ctxBytes,
		Item:    []*v3.Item{item},
	}
	reqExt := &mediation.RequestExt{
		Adapters: map[string]*mediation.SdkAdapter{
			"applovin": {
				Version:    proto.String("0.1.0"),
				SdkVersion: proto.String("1.0.0"),
			},
		},
		AdType: ptr(mediation.AdType_AD_TYPE_BANNER),
		Ext:    proto.String(`{"mediation_mode":"bidon"}`),
	}
	proto.SetExtension(req, mediation.E_RequestExt, reqExt)

	return req
}

// Response helper
type ResponseBuilder struct {
	resp *v3.Response
}

func NewResponseBuilder() *ResponseBuilder {
	return &ResponseBuilder{
		resp: defaultValidResp(),
	}
}

func (rb *ResponseBuilder) Build() *v3.Openrtb {
	return &v3.Openrtb{
		PayloadOneof: &v3.Openrtb_Response{
			Response: rb.resp,
		},
	}
}

func (rb *ResponseBuilder) WithAdUnits(adUnits []auction.AdUnit) *ResponseBuilder {
	bids := make([]*v3.Bid, len(adUnits))
	for i, adUnit := range adUnits {
		bid, _ := adUnitToBid(&adUnit)
		bids[i] = bid
	}
	rb.resp.Seatbid[0].Bid = bids
	return rb
}

func defaultValidResp() *v3.Response {
	ac := DefaultAuctionConfig()
	sgmnt := DefaultSegment()

	adUnits := DefaultAdUnits()
	bids := make([]*v3.Bid, len(adUnits))
	for i, adUnit := range adUnits {
		bid, _ := adUnitToBid(&adUnit)
		bids[i] = bid
	}

	resp := &v3.Response{
		Id: proto.String("auction_id_123"),
		Seatbid: []*v3.SeatBid{
			{
				Bid: bids,
			},
		},
	}

	respExt := &mediation.AuctionResponseExt{
		Token:                    proto.String("{}"),
		ExternalWinNotifications: proto.Bool(ac.ExternalWinNotifications),
		Segment: &mediation.Segment{
			Id:  proto.String(sgmnt.StringID()),
			Uid: proto.String(sgmnt.UID),
		},
		AuctionConfigurationId:  proto.Int64(ac.ID),
		AuctionConfigurationUid: proto.String(ac.UID),
		AuctionTimeout:          proto.Int32(int32(ac.Timeout)),
	}
	proto.SetExtension(resp, mediation.E_AuctionResponseExt, respExt)

	return resp
}

func DefaultAuctionConfig() auction.Config {
	return auction.Config{
		ID:                       1,
		UID:                      "1701972528521547776",
		Demands:                  []adapter.Key{adapter.ApplovinKey, adapter.BidmachineKey},
		AdUnitIDs:                []int64{1, 2, 3},
		ExternalWinNotifications: true,
		Timeout:                  30000,
	}
}

func DefaultSegment() segment.Segment {
	return segment.Segment{
		ID:      1,
		UID:     "1701972528521547776",
		Filters: []segment.Filter{{Type: "country", Name: "country", Operator: "IN", Values: []string{"US", "UK"}}},
	}
}

func DefaultAdUnits() []auction.AdUnit {
	return []auction.AdUnit{
		{
			DemandID:   "meta",
			Label:      "meta",
			PriceFloor: ptr(0.8),
			UID:        "123_meta",
			BidType:    schema.RTBBidType,
			Timeout:    store.AdUnitTimeout,
			Extra: map[string]any{
				"payload":      "payload",
				"placement_id": "123",
			},
		},
		{
			DemandID: "vungle",
			Label:    "vungle",
			UID:      "123_vungle",
			BidType:  schema.RTBBidType,
			Timeout:  store.AdUnitTimeout,
			Extra: map[string]any{
				"placement_id": "123",
			},
		},
	}
}

func BuildDemandResponses(adUnits []auction.AdUnit) []adapters.DemandResponse {
	responses := make([]adapters.DemandResponse, len(adUnits))
	for i, adUnit := range adUnits {
		response := adapters.DemandResponse{
			DemandID: adapter.Key(adUnit.DemandID),
		}
		if adUnit.PriceFloor != nil {
			response.Bid = &adapters.BidDemandResponse{
				Price:    *adUnit.PriceFloor,
				Payload:  "payload",
				ID:       "123",
				ImpID:    "456",
				DemandID: adapter.Key(adUnit.DemandID),
			}
		}
		responses[i] = response
	}

	return responses
}

func ptr[T any](t T) *T {
	return &t
}
