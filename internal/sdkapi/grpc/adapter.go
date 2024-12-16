package grpcserver

import (
	"encoding/json"
	"fmt"

	"github.com/bidon-io/bidon-backend/internal/auction"
	"github.com/bidon-io/bidon-backend/internal/auctionv2"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/bidon-io/bidon-backend/internal/ad"
	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/device"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
	adcom "github.com/bidon-io/bidon-backend/pkg/proto/com/iabtechlab/adcom/v1"
	v3 "github.com/bidon-io/bidon-backend/pkg/proto/com/iabtechlab/openrtb/v3"
	pbctx "github.com/bidon-io/bidon-backend/pkg/proto/org/bidon/proto/v1/context"
	"github.com/bidon-io/bidon-backend/pkg/proto/org/bidon/proto/v1/mediation"
	"google.golang.org/protobuf/proto"
)

type AuctionAdapter struct{}

func NewAuctionAdapter() *AuctionAdapter {
	return &AuctionAdapter{}
}

// OpenRTBToAuctionRequest converts a v3.Openrtb request into an AuctionV2Request.
func (*AuctionAdapter) OpenRTBToAuctionRequest(o *v3.Openrtb) (*schema.AuctionV2Request, error) {
	req := o.GetRequest()
	if req == nil {
		return nil, fmt.Errorf("OpenRTBToAuctionRequest: request is nil")
	}

	ar := &schema.AuctionV2Request{}
	err := parseAuctionRequest(ar, req)
	if err != nil {
		return nil, fmt.Errorf("OpenRTBToAuctionRequest: %w", err)
	}

	return ar, nil
}

// AuctionResponseToOpenRTB converts an AuctionResponse into a v3.Openrtb response.
func (*AuctionAdapter) AuctionResponseToOpenRTB(r *auctionv2.Response) (*v3.Openrtb, error) {
	bids := make([]*v3.Bid, 0, len(r.AdUnits)+len(r.NoBids))
	adUnits := append(r.AdUnits, r.NoBids...)
	for _, a := range adUnits {
		bid, err := adUnitToBid(&a)
		if err != nil {
			return nil, fmt.Errorf("AuctionResponseToOpenRTB: %w", err)
		}
		bids = append(bids, bid)
	}

	resp := &v3.Response{
		Id: proto.String(r.AuctionID),
		Seatbid: []*v3.SeatBid{
			{Bid: bids},
		},
	}

	respExt := &mediation.AuctionResponseExt{
		Token:                    proto.String(r.Token),
		ExternalWinNotifications: proto.Bool(r.ExternalWinNotifications),
		Segment: &mediation.Segment{
			Id:  proto.String(r.Segment.ID),
			Uid: proto.String(r.Segment.UID),
		},
		AuctionConfigurationId:  proto.Int64(r.ConfigID),
		AuctionConfigurationUid: proto.String(r.ConfigUID),
		AuctionTimeout:          proto.Int32(int32(r.AuctionTimeout)),
	}
	proto.SetExtension(resp, mediation.E_AuctionResponseExt, respExt)

	return &v3.Openrtb{
		PayloadOneof: &v3.Openrtb_Response{
			Response: resp,
		},
	}, nil
}

func parseAuctionRequest(ar *schema.AuctionV2Request, req *v3.Request) error {
	ext, err := getMediationExtension[*mediation.RequestExt](req, mediation.E_RequestExt)
	if err != nil {
		return fmt.Errorf("parseAuctionRequest: %w", err)
	}

	br, err := parseBaseRequest(req)
	if err != nil {
		return err
	}
	ar.BaseRequest = br

	adObject, err := parseAdObject(req)
	if err != nil {
		return err
	}
	ar.AdObject = adObject

	ar.Test = req.GetTest()
	ar.TMax = int64(req.GetTmax())
	ar.Ext = ext.GetExt()
	ar.Adapters = parseAdapters(ext)
	ar.AdType = parseAdType(ext)

	return nil
}

func parseAdapters(ext *mediation.RequestExt) schema.Adapters {
	mAdapters := ext.GetAdapters()
	adapters := make(schema.Adapters, len(mAdapters))
	for key, ma := range mAdapters {
		adapters[adapter.Key(key)] = schema.Adapter{
			Version:    ma.GetVersion(),
			SDKVersion: ma.GetSdkVersion(),
		}
	}

	return adapters
}

func parseAdType(ext *mediation.RequestExt) ad.Type {
	mAdType := ext.GetAdType()
	switch mAdType {
	case mediation.AdType_AD_TYPE_BANNER:
		return ad.BannerType
	case mediation.AdType_AD_TYPE_INTERSTITIAL:
		return ad.InterstitialType
	case mediation.AdType_AD_TYPE_REWARDED:
		return ad.RewardedType
	default:
		return ad.UnknownType
	}
}

func parseBaseRequest(req *v3.Request) (schema.BaseRequest, error) {
	reqCtx := req.GetContext()
	if reqCtx == nil {
		return schema.BaseRequest{}, fmt.Errorf("parseBaseRequest: request context is empty")
	}

	c := &pbctx.Context{}
	if err := proto.Unmarshal(reqCtx, c); err != nil {
		return schema.BaseRequest{}, fmt.Errorf("parseBaseRequest: failed to unmarshal context: %w", err)
	}

	device, err := parseDevice(c)
	if err != nil {
		return schema.BaseRequest{}, err
	}

	session, err := parseSession(c)
	if err != nil {
		return schema.BaseRequest{}, err
	}

	app, err := parseApp(c)
	if err != nil {
		return schema.BaseRequest{}, err
	}

	user, err := parseUser(c)
	if err != nil {
		return schema.BaseRequest{}, err
	}

	segment, err := parseSegment(c)
	if err != nil {
		return schema.BaseRequest{}, err
	}

	regs, err := parseRegs(c)
	if err != nil {
		return schema.BaseRequest{}, err
	}

	return schema.BaseRequest{
		Device:      device,
		Geo:         device.Geo,
		Session:     session,
		App:         app,
		User:        user,
		Segment:     segment,
		Regulations: regs,
	}, nil
}

func parseAdObject(r *v3.Request) (schema.AdObjectV2, error) {
	items := r.GetItem()
	if len(items) == 0 {
		return schema.AdObjectV2{}, fmt.Errorf("parseAdObject: no items in request")
	}
	i := items[0]
	placementBytes := i.GetSpec()
	if placementBytes == nil {
		return schema.AdObjectV2{}, fmt.Errorf("parseAdObject: placement is empty")
	}

	var placement adcom.Placement
	if err := proto.Unmarshal(placementBytes, &placement); err != nil {
		return schema.AdObjectV2{}, fmt.Errorf("parseAdObject: failed to unmarshal placement: %w", err)
	}

	mi, err := getMediationExtension[*mediation.PlacementExt](&placement, mediation.E_PlacementExt)
	if err != nil {
		return schema.AdObjectV2{}, fmt.Errorf("parseAdObject: %w", err)
	}

	var banner *schema.BannerAdObject
	if b := mi.GetBanner(); b != nil {
		banner = &schema.BannerAdObject{
			Format: ad.Format(b.GetFormat().String()),
		}
	}

	var interstitial *schema.InterstitialAdObject
	if inter := mi.GetInterstitial(); inter != "" {
		interstitial = &schema.InterstitialAdObject{}
	}

	var rewarded *schema.RewardedAdObject
	if rew := mi.GetRewarded(); rew != "" {
		rewarded = &schema.RewardedAdObject{}
	}

	demands := make(map[adapter.Key]map[string]any)
	mdemands := mi.GetDemands()
	for k, v := range mdemands {
		demands[adapter.Key(k)] = map[string]any{
			"token":           v.GetToken(),
			"status":          v.GetStatus(),
			"token_finish_ts": v.GetTokenFinishTs(),
			"token_start_ts":  v.GetTokenStartTs(),
		}
	}

	return schema.AdObjectV2{
		AuctionID:               i.GetId(),
		AuctionConfigurationUID: mi.GetAuctionConfigurationUid(),
		Orientation:             mi.GetOrientation().String(),
		PriceFloor:              float64(i.GetFlr()),
		AuctionKey:              mi.GetAuctionKey(),
		Demands:                 demands,
		Banner:                  banner,
		Interstitial:            interstitial,
		Rewarded:                rewarded,
	}, nil
}

func parseApp(c *pbctx.Context) (schema.App, error) {
	a := c.DistributionChannel.GetApp()
	if a == nil {
		return schema.App{}, fmt.Errorf("parseApp: app is empty in context")
	}

	ma, err := getMediationExtension[*mediation.AppExt](a, mediation.E_AppExt)
	if err != nil {
		return schema.App{}, fmt.Errorf("parseApp: %w", err)
	}

	return schema.App{
		Bundle:           a.GetBundle(),
		Key:              ma.GetKey(),
		Framework:        ma.GetFramework(),
		Version:          a.GetVer(),
		FrameworkVersion: ma.GetFrameworkVersion(),
		PluginVersion:    ma.GetPluginVersion(),
		SKAdN:            ma.GetSkadn(),
		SDKVersion:       ma.GetSdkVersion(),
	}, nil
}

func parseUser(c *pbctx.Context) (schema.User, error) {
	u := c.GetUser()
	if u == nil {
		return schema.User{}, fmt.Errorf("parseUser: user is empty in context")
	}

	mu, err := getMediationExtension[*mediation.UserExt](u, mediation.E_UserExt)
	if err != nil {
		return schema.User{}, fmt.Errorf("parseUser: %w", err)
	}

	return schema.User{
		IDFA:                        mu.GetIdfa(),
		TrackingAuthorizationStatus: mu.GetTrackingAuthorizationStatus(),
		IDFV:                        mu.GetIdfv(),
		IDG:                         mu.GetIdg(),
		Consent:                     parseJson(u.GetConsent()),
	}, nil
}

func parseSegment(c *pbctx.Context) (schema.Segment, error) {
	u := c.GetUser()
	if u == nil {
		return schema.Segment{}, fmt.Errorf("parseSegment: user is empty in context")
	}

	mu, err := getMediationExtension[*mediation.UserExt](u, mediation.E_UserExt)
	if err != nil {
		return schema.Segment{}, fmt.Errorf("parseSegment: %w", err)
	}

	segments := mu.GetSegments()
	if len(segments) == 0 {
		return schema.Segment{}, fmt.Errorf("parseSegment: segments is empty")
	}

	ms := segments[0]
	return schema.Segment{
		ID:  ms.GetId(),
		UID: ms.GetUid(),
		Ext: ms.GetExt(),
	}, nil
}

func parseSession(c *pbctx.Context) (schema.Session, error) {
	d := c.GetDevice()
	if d == nil {
		return schema.Session{}, fmt.Errorf("parseSession: device is empty in context")
	}

	sess, err := getMediationExtension[*mediation.DeviceExt](d, mediation.E_DeviceExt)
	if err != nil {
		return schema.Session{}, fmt.Errorf("parseSession: %w", err)
	}

	return schema.Session{
		ID:                        sess.GetId(),
		LaunchTS:                  int(sess.GetLaunchTs()),
		LaunchMonotonicTS:         int(sess.GetLaunchMonotonicTs()),
		StartTS:                   int(sess.GetStartTs()),
		StartMonotonicTS:          int(sess.GetStartMonotonicTs()),
		TS:                        int(sess.GetTs()),
		MonotonicTS:               int(sess.GetMonotonicTs()),
		MemoryWarningsTS:          sliceToInt(sess.GetMemoryWarningsTs()),
		MemoryWarningsMonotonicTS: sliceToInt(sess.GetMemoryWarningsMonotonicTs()),
		RAMUsed:                   int(sess.GetRamUsed()),
		RAMSize:                   int(sess.GetRamSize()),
		StorageFree:               int(sess.GetStorageFree()),
		StorageUsed:               int(sess.GetStorageUsed()),
		Battery:                   float64(sess.GetBattery()),
		CPUUsage:                  proto.Float64(sess.GetCpuUsage()),
	}, nil
}

func parseDevice(c *pbctx.Context) (schema.Device, error) {
	d := c.GetDevice()
	if d == nil {
		return schema.Device{}, fmt.Errorf("parseDevice: device is empty in context")
	}

	g := d.GetGeo()
	return schema.Device{
		Geo: &schema.Geo{
			Lat:       float64(g.GetLat()),
			Lon:       float64(g.GetLon()),
			Accuracy:  float64(g.GetAccur()),
			LastFix:   int(g.GetLastfix()),
			Country:   g.GetCountry(),
			City:      g.GetCity(),
			ZIP:       g.GetZip(),
			UTCOffset: int(g.GetUtcoffset()),
		},
		UserAgent:       d.GetUa(),
		Manufacturer:    d.GetMake(),
		Model:           d.GetModel(),
		OS:              parseOs(adcom.OperatingSystem(d.GetOs())),
		OSVersion:       d.GetOsv(),
		HardwareVersion: d.GetHwv(),
		Height:          int(d.GetH()),
		Width:           int(d.GetW()),
		PPI:             int(d.GetPpi()),
		PXRatio:         float64(d.GetPxratio()),
		JS:              parseJS(d.GetJs()),
		Language:        d.GetLang(),
		IP:              d.GetIp(),
		Carrier:         d.GetCarrier(),
		MCCMNC:          d.GetMccmnc(),
		ConnectionType:  adcom.ConnectionType(d.GetContype()).String(),
		Type:            device.Type(adcom.DeviceType(d.GetType()).String()),
	}, nil
}

func parseRegs(c *pbctx.Context) (*schema.Regulations, error) {
	r := c.GetRegs()
	if r == nil {
		return &schema.Regulations{}, fmt.Errorf("parseRegs: regs is empty in context")
	}

	mr, err := getMediationExtension[*mediation.RegsExt](r, mediation.E_RegsExt)
	if err != nil {
		return &schema.Regulations{}, fmt.Errorf("parseRegs: %w", err)
	}

	return &schema.Regulations{
		COPPA:     r.GetCoppa(),
		GDPR:      r.GetGdpr(),
		USPrivacy: mr.GetUsPrivacy(),
		EUPrivacy: mr.GetEuPrivacy(),
		IAB:       parseJson(mr.GetIab()),
	}, nil
}

func adUnitToBid(a *auction.AdUnit) (*v3.Bid, error) {
	var price float32
	if a.PriceFloor != nil {
		price = float32(*a.PriceFloor)
	}
	bid := &v3.Bid{
		Item:  proto.String(a.UID),
		Price: proto.Float32(price),
		Cid:   proto.String(a.DemandID),
	}

	ext := make(map[string]string, len(a.Extra))
	for k, v := range a.Extra {
		ext[k] = fmt.Sprintf("%v", v)
	}
	bidExt := &mediation.BidExt{
		Label:   proto.String(a.Label),
		BidType: proto.String(a.BidType.String()),
		Ext:     ext,
	}
	proto.SetExtension(bid, mediation.E_BidExt, bidExt)

	return bid, nil
}

func getMediationExtension[T proto.Message](m proto.Message, ext protoreflect.ExtensionType) (T, error) {
	e := proto.GetExtension(m, ext)
	if e == nil {
		return *new(T), fmt.Errorf("getMediationExtension: extension %q not found", ext)
	}
	casted, ok := e.(T)
	if !ok {
		return *new(T), fmt.Errorf("getMediationExtension: extension %q invalid type", ext)
	}
	return casted, nil
}

var osMap = map[adcom.OperatingSystem]string{
	adcom.OperatingSystem_ANDROID: string(ad.AndroidOS),
	adcom.OperatingSystem_IOS:     string(ad.IOSOS),
}

func parseOs(os adcom.OperatingSystem) string {
	osStr, ok := osMap[os]
	if !ok {
		return string(ad.UnknownOS)
	}
	return osStr
}

func parseJS(js bool) *int {
	v := 0
	if js {
		v = 1
	}
	return &v
}

func parseJson(str string) map[string]any {
	if str == "" {
		return map[string]any{}
	}

	m := make(map[string]any)
	err := json.Unmarshal([]byte(str), &m)
	if err != nil {
		return map[string]any{}
	}
	return m
}

func sliceToInt(in []int64) []int {
	out := make([]int, len(in))
	for i, v := range in {
		out[i] = int(v)
	}
	return out
}
