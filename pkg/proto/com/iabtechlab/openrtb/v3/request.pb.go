//*
// The request object contains minimal high level attributes (e.g., its ID, test
// mode, auction type, maximum auction time, buyer restrictions, etc.) and
// subordinate objects that cover the source of the request and the actual offer
// of sale. The latter includes the item(s) being offered and any applicable
// deals.
//
// There are two points in this model that interface to Layer-4 domain objects:
// the Request object and the Item object. Domain objects included under Request
// would include those that provide context for the overall offer. These would
// include objects that describe the site or app, the device, the user, and
// others. Domain objects included in an Item object would specify details about
// the item being offered (e.g., the impression opportunity) and specifications
// and restrictions on the media that can be associated with acceptable bids.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        (unknown)
// source: com/iabtechlab/openrtb/v3/request.proto

package openrtbv3

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// *
// The Request object contains a globally unique bid request ID. This id
// attribute is required as is an Item array with at least one object (i.e., at
// least one item for sale). Other attributes establish rules and restrictions
// that apply to all items being offered. This object also interfaces to Layer-4
// domain objects for context such as the user, device, site or app, etc.
type Request struct {
	state           protoimpl.MessageState
	sizeCache       protoimpl.SizeCache
	unknownFields   protoimpl.UnknownFields
	extensionFields protoimpl.ExtensionFields

	// Unique ID of the bid request; provided by the exchange.
	Id *string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	// Indicator of test mode in which auctions are not billable, where 0 =
	// live mode, 1 = test mode.
	Test *bool `protobuf:"varint,2,opt,name=test" json:"test,omitempty"`
	// Maximum time in milliseconds the exchange allows for bids to be received
	// including Internet latency to avoid timeout. This value supersedes any a
	// priori guidance from the exchange. If an exchange acts as an intermediary,
	// it should decrease the outbound tmax value from what it received to account
	// for its latency and the additional internet hop.
	Tmax *uint32 `protobuf:"varint,3,opt,name=tmax" json:"tmax,omitempty"`
	// Auction type, where 1 = First Price, 2 = Second Price Plus. Values greater
	// than 500 can be used for exchange-specific auction types. Common values
	// are defined in the enumeration com.iabtechlab.openrtb.v3.AuctionType.
	At *int32 `protobuf:"varint,4,opt,name=at" json:"at,omitempty"`
	// Array of accepted currencies for bids on this bid request using ISO-4217
	// alpha codes. Recommended if the exchange accepts multiple currencies. If
	// omitted, the single currency of “USD” is assumed.
	Cur []string `protobuf:"bytes,5,rep,name=cur" json:"cur,omitempty"`
	// Restriction list of buyer seats for bidding on this item. Knowledge of
	// buyer’s customers and their seat IDs must be coordinated between parties a
	// priori. Omission implies no restrictions.
	Seat []string `protobuf:"bytes,6,rep,name=seat" json:"seat,omitempty"`
	// Flag that determines the restriction interpretation of the seat array,
	// where 0 = blocklist, 1 = allowlist.
	Wseat *bool `protobuf:"varint,7,opt,name=wseat" json:"wseat,omitempty"`
	// Allows bidder to retrieve data set on its behalf in the exchange’s cookie
	// (refer to cdata in Object: Response) if supported by the exchange. The
	// string must be in base85 cookie-safe characters.
	Cdata *string `protobuf:"bytes,8,opt,name=cdata" json:"cdata,omitempty"`
	// A Source object that provides data about the inventory source and which
	// entity makes the final decision. Refer to Object: Source.
	Source *Source `protobuf:"bytes,9,opt,name=source" json:"source,omitempty"`
	// Array of Item objects (at least one) that constitute the set of goods being
	// offered for sale. Refer to Object: Item.
	Item []*Item `protobuf:"bytes,10,rep,name=item" json:"item,omitempty"`
	// Flag to indicate if the Exchange can verify that the items offered
	// represent all of the items available in context (e.g., all impressions on a
	// web page, all video spots such as pre/mid/post roll) to support
	// road-blocking, where 0 = no, 1 = yes.
	Package *bool `protobuf:"varint,11,opt,name=package" json:"package,omitempty"`
	// Layer-4 domain object structure that provides context for the items being
	// offered conforming to the specification and version referenced in
	// openrtb.domainspec and openrtb.domainver.
	// For AdCOM v1.x, the objects allowed here all of which are optional are one
	// of the DistributionChannel subtypes (i.e., Site, App, or Dooh), User,
	// Device, Regs, Restrictions, and any objects subordinate to these as
	// specified by AdCOM.
	Context []byte `protobuf:"bytes,12,opt,name=context" json:"context,omitempty"`
}

func (x *Request) Reset() {
	*x = Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_com_iabtechlab_openrtb_v3_request_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Request) ProtoMessage() {}

func (x *Request) ProtoReflect() protoreflect.Message {
	mi := &file_com_iabtechlab_openrtb_v3_request_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Request.ProtoReflect.Descriptor instead.
func (*Request) Descriptor() ([]byte, []int) {
	return file_com_iabtechlab_openrtb_v3_request_proto_rawDescGZIP(), []int{0}
}

func (x *Request) GetId() string {
	if x != nil && x.Id != nil {
		return *x.Id
	}
	return ""
}

func (x *Request) GetTest() bool {
	if x != nil && x.Test != nil {
		return *x.Test
	}
	return false
}

func (x *Request) GetTmax() uint32 {
	if x != nil && x.Tmax != nil {
		return *x.Tmax
	}
	return 0
}

func (x *Request) GetAt() int32 {
	if x != nil && x.At != nil {
		return *x.At
	}
	return 0
}

func (x *Request) GetCur() []string {
	if x != nil {
		return x.Cur
	}
	return nil
}

func (x *Request) GetSeat() []string {
	if x != nil {
		return x.Seat
	}
	return nil
}

func (x *Request) GetWseat() bool {
	if x != nil && x.Wseat != nil {
		return *x.Wseat
	}
	return false
}

func (x *Request) GetCdata() string {
	if x != nil && x.Cdata != nil {
		return *x.Cdata
	}
	return ""
}

func (x *Request) GetSource() *Source {
	if x != nil {
		return x.Source
	}
	return nil
}

func (x *Request) GetItem() []*Item {
	if x != nil {
		return x.Item
	}
	return nil
}

func (x *Request) GetPackage() bool {
	if x != nil && x.Package != nil {
		return *x.Package
	}
	return false
}

func (x *Request) GetContext() []byte {
	if x != nil {
		return x.Context
	}
	return nil
}

// *
// This object carries data about the source of the transaction including the
// unique ID of the transaction itself, source authentication information, and
// the chain of custody.
//
// NOTE: Attributes ds, dsmap, cert, and digest support digitally signed bid
// requests as defined by the Ads.cert: Signed Bid Requests specification. As
// the Ads.cert specification is still in its BETA state, these attributes
// should be considered to be in a similar state.
type Source struct {
	state           protoimpl.MessageState
	sizeCache       protoimpl.SizeCache
	unknownFields   protoimpl.UnknownFields
	extensionFields protoimpl.ExtensionFields

	// Transaction ID that must be common across all participants throughout the
	// entire supply chain of this transaction. This also applies across all
	// participating exchanges in a header bidding or similar publisher-centric
	// broadcast scenario.
	Tid *string `protobuf:"bytes,1,opt,name=tid" json:"tid,omitempty"`
	// Timestamp when the request originated at the beginning of the supply chain
	// in Unix format (i.e., milliseconds since the epoch). This value must be
	// held as immutable throughout subsequent intermediaries.
	Ts *uint64 `protobuf:"varint,2,opt,name=ts" json:"ts,omitempty"`
	// Digital signature used to authenticate the origin of this request computed
	// by the publisher or its trusted agent from a digest string composed of a
	// set of immutable attributes found in the bid request. Refer to Section
	// “Inventory Authentication” for more details.
	Ds *string `protobuf:"bytes,3,opt,name=ds" json:"ds,omitempty"`
	// An ordered list of identifiers that indicates the attributes used to create
	// the digest. This map provides the essential instructions for recreating the
	// digest from the bid request, which is a necessary step in validating the
	// digital signature in the ds attribute. Refer to Section “Inventory
	// Authentication” for more details.
	Dsmap *string `protobuf:"bytes,4,opt,name=dsmap" json:"dsmap,omitempty"`
	// File name of the certificate (i.e., the public key) used to generate the
	// digital signature in the ds attribute. Refer to Section “Inventory
	// Authentication” for more details.
	Cert *string `protobuf:"bytes,5,opt,name=cert" json:"cert,omitempty"`
	// The full digest string that was signed to produce the digital signature.
	// Refer to Section “Inventory Authentication” for more details.
	// NOTE: This is only intended for debugging purposes as needed. It is not
	// intended for normal Production traffic due to the bandwidth impact.
	Digest *string `protobuf:"bytes,6,opt,name=digest" json:"digest,omitempty"`
	// Payment ID chain string containing embedded syntax described in the TAG
	// Payment ID Protocol.
	// NOTE: Authentication features in this Source object combined with the
	// “ads.txt” specification may lead to the deprecation of this attribute.
	Pchain *string `protobuf:"bytes,7,opt,name=pchain" json:"pchain,omitempty"`
}

func (x *Source) Reset() {
	*x = Source{}
	if protoimpl.UnsafeEnabled {
		mi := &file_com_iabtechlab_openrtb_v3_request_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Source) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Source) ProtoMessage() {}

func (x *Source) ProtoReflect() protoreflect.Message {
	mi := &file_com_iabtechlab_openrtb_v3_request_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Source.ProtoReflect.Descriptor instead.
func (*Source) Descriptor() ([]byte, []int) {
	return file_com_iabtechlab_openrtb_v3_request_proto_rawDescGZIP(), []int{1}
}

func (x *Source) GetTid() string {
	if x != nil && x.Tid != nil {
		return *x.Tid
	}
	return ""
}

func (x *Source) GetTs() uint64 {
	if x != nil && x.Ts != nil {
		return *x.Ts
	}
	return 0
}

func (x *Source) GetDs() string {
	if x != nil && x.Ds != nil {
		return *x.Ds
	}
	return ""
}

func (x *Source) GetDsmap() string {
	if x != nil && x.Dsmap != nil {
		return *x.Dsmap
	}
	return ""
}

func (x *Source) GetCert() string {
	if x != nil && x.Cert != nil {
		return *x.Cert
	}
	return ""
}

func (x *Source) GetDigest() string {
	if x != nil && x.Digest != nil {
		return *x.Digest
	}
	return ""
}

func (x *Source) GetPchain() string {
	if x != nil && x.Pchain != nil {
		return *x.Pchain
	}
	return ""
}

// *
// This object represents a unit of goods being offered for sale either on the
// open market or in relation to a private marketplace deal. The id attribute is
// required since there may be multiple items being offered in the same bid
// request and bids must reference the specific item of interest. This object
// interfaces to Layer-4 domain objects for deeper specification of the item
// being offered (e.g., an impression).
type Item struct {
	state           protoimpl.MessageState
	sizeCache       protoimpl.SizeCache
	unknownFields   protoimpl.UnknownFields
	extensionFields protoimpl.ExtensionFields

	// A unique identifier for this item within the context of the offer
	// (typically starts with “1” and increments).
	Id *string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	// Types that are assignable to QtyOneof:
	//
	//	*Item_Qty
	//	*Item_Qtyflt
	QtyOneof isItem_QtyOneof `protobuf_oneof:"qty_oneof"`
	// If multiple items are offered in the same bid request, the sequence number
	// allows for the coordinated delivery.
	Seq *uint32 `protobuf:"varint,4,opt,name=seq" json:"seq,omitempty"`
	// Minimum bid price for this item expressed in CPM.
	Flr *float32 `protobuf:"fixed32,5,opt,name=flr" json:"flr,omitempty"`
	// Currency of the flr attribute specified using ISO-4217 alpha codes.
	Flrcur *string `protobuf:"bytes,6,opt,name=flrcur" json:"flrcur,omitempty"`
	// Advisory as to the number of seconds that may elapse between auction and
	// fulfilment.
	Exp *uint64 `protobuf:"varint,7,opt,name=exp" json:"exp,omitempty"`
	// Timestamp when the item is expected to be fulfilled (e.g. when a DOOH
	// impression will be displayed) in Unix format (i.e., milliseconds since the
	// epoch).
	Dt *uint64 `protobuf:"varint,8,opt,name=dt" json:"dt,omitempty"`
	// Item (e.g., an Ad object) delivery method required, where 0 = either
	// method, 1 = the item must be sent as part of the transaction (e.g., by
	// value in the bid itself, fetched by URL included in the bid), and 2 = an
	// item previously uploaded to the exchange must be referenced by its ID. Note
	// that if an exchange does not supported prior upload, then the default of 0
	// is effectively the same as 1 since there can be no items to reference.
	// Common values are defined in the enumeration
	// com.iabtechlab.openrtb.v3.ItemDeliveryMethod.
	Dlvy *int32 `protobuf:"varint,9,opt,name=dlvy" json:"dlvy,omitempty"`
	// An array of Metric objects. Refer to Object: Metric.
	Metric []*Metric `protobuf:"bytes,10,rep,name=metric" json:"metric,omitempty"`
	// Array of Deal objects that convey special terms applicable to this item.
	// Refer to Object: Deal.
	Deal []*Deal `protobuf:"bytes,11,rep,name=deal" json:"deal,omitempty"`
	// Indicator of auction eligibility to seats named in Deal objects, where 0 =
	// all bids are accepted, 1 = bids are restricted to the deals specified and
	// the terms thereof.
	Private *bool `protobuf:"varint,12,opt,name=private" json:"private,omitempty"`
	// Layer-4 domain object structure that provides specifies the item being
	// offered conforming to the specification and version referenced in
	// openrtb.domainspec and openrtb.domainver.
	// For AdCOM v1.x, the objects allowed here are Placement and any objects
	// subordinate to these as specified by AdCOM.
	Spec []byte `protobuf:"bytes,13,opt,name=spec" json:"spec,omitempty"`
}

func (x *Item) Reset() {
	*x = Item{}
	if protoimpl.UnsafeEnabled {
		mi := &file_com_iabtechlab_openrtb_v3_request_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Item) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Item) ProtoMessage() {}

func (x *Item) ProtoReflect() protoreflect.Message {
	mi := &file_com_iabtechlab_openrtb_v3_request_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Item.ProtoReflect.Descriptor instead.
func (*Item) Descriptor() ([]byte, []int) {
	return file_com_iabtechlab_openrtb_v3_request_proto_rawDescGZIP(), []int{2}
}

func (x *Item) GetId() string {
	if x != nil && x.Id != nil {
		return *x.Id
	}
	return ""
}

func (m *Item) GetQtyOneof() isItem_QtyOneof {
	if m != nil {
		return m.QtyOneof
	}
	return nil
}

func (x *Item) GetQty() uint32 {
	if x, ok := x.GetQtyOneof().(*Item_Qty); ok {
		return x.Qty
	}
	return 0
}

func (x *Item) GetQtyflt() float32 {
	if x, ok := x.GetQtyOneof().(*Item_Qtyflt); ok {
		return x.Qtyflt
	}
	return 0
}

func (x *Item) GetSeq() uint32 {
	if x != nil && x.Seq != nil {
		return *x.Seq
	}
	return 0
}

func (x *Item) GetFlr() float32 {
	if x != nil && x.Flr != nil {
		return *x.Flr
	}
	return 0
}

func (x *Item) GetFlrcur() string {
	if x != nil && x.Flrcur != nil {
		return *x.Flrcur
	}
	return ""
}

func (x *Item) GetExp() uint64 {
	if x != nil && x.Exp != nil {
		return *x.Exp
	}
	return 0
}

func (x *Item) GetDt() uint64 {
	if x != nil && x.Dt != nil {
		return *x.Dt
	}
	return 0
}

func (x *Item) GetDlvy() int32 {
	if x != nil && x.Dlvy != nil {
		return *x.Dlvy
	}
	return 0
}

func (x *Item) GetMetric() []*Metric {
	if x != nil {
		return x.Metric
	}
	return nil
}

func (x *Item) GetDeal() []*Deal {
	if x != nil {
		return x.Deal
	}
	return nil
}

func (x *Item) GetPrivate() bool {
	if x != nil && x.Private != nil {
		return *x.Private
	}
	return false
}

func (x *Item) GetSpec() []byte {
	if x != nil {
		return x.Spec
	}
	return nil
}

type isItem_QtyOneof interface {
	isItem_QtyOneof()
}

type Item_Qty struct {
	// The quantity of billable events which will be deemed to have occured if
	// this item is purchased. In most cases, this represents impressions. For
	// example, a single display of an ad on a DOOH placement may count as multiple
	// impressions on the basis of expected viewership. In such a case, qty would
	// be greater than 1.
	Qty uint32 `protobuf:"varint,2,opt,name=qty,oneof"`
}

type Item_Qtyflt struct {
	// The quantity of billable events which will be deemed to have occured if this
	// item is purchased. This version of the fields exists for cases where the
	// quantity is not expressed as a whole number. For example, a DOOH opportunity
	// may be considered to be 14.2 impressions.
	Qtyflt float32 `protobuf:"fixed32,3,opt,name=qtyflt,oneof"`
}

func (*Item_Qty) isItem_QtyOneof() {}

func (*Item_Qtyflt) isItem_QtyOneof() {}

// *
// This object constitutes a specific deal that was struck a priori between
// a seller and a buyer. Its presence indicates that this item is available
// under the terms of that deal.
type Deal struct {
	state           protoimpl.MessageState
	sizeCache       protoimpl.SizeCache
	unknownFields   protoimpl.UnknownFields
	extensionFields protoimpl.ExtensionFields

	// A unique identifier for the deal.
	Id *string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	// Minimum deal price for this item expressed in CPM.
	Flr *float32 `protobuf:"fixed32,2,opt,name=flr" json:"flr,omitempty"`
	// Currency of the flr attribute specified using ISO-4217 alpha codes.
	Flrcur *string `protobuf:"bytes,3,opt,name=flrcur" json:"flrcur,omitempty"`
	// Optional override of the overall auction type of the request, where 1 =
	// First Price, 2 = Second Price Plus, 3 = the value passed in flr is the
	// agreed upon deal price. Additional auction types can be defined by the
	// exchange using 500+ values. Common values are defined in the enumeration
	// com.iabtechlab.openrtb.v3.AuctionType.
	At *int32 `protobuf:"varint,4,opt,name=at" json:"at,omitempty"`
	// Allowlist of buyer seats allowed to bid on this deal. IDs of seats and the
	// buyer’s customers to which they refer must be coordinated between bidders
	// and the exchange a priori. Omission implies no restrictions.
	Wseat []string `protobuf:"bytes,5,rep,name=wseat" json:"wseat,omitempty"`
	// Array of advertiser domains (e.g., advertiser.com) allowed to bid on this
	// deal. Omission implies no restrictions.
	Wadomain []string `protobuf:"bytes,6,rep,name=wadomain" json:"wadomain,omitempty"`
}

func (x *Deal) Reset() {
	*x = Deal{}
	if protoimpl.UnsafeEnabled {
		mi := &file_com_iabtechlab_openrtb_v3_request_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Deal) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Deal) ProtoMessage() {}

func (x *Deal) ProtoReflect() protoreflect.Message {
	mi := &file_com_iabtechlab_openrtb_v3_request_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Deal.ProtoReflect.Descriptor instead.
func (*Deal) Descriptor() ([]byte, []int) {
	return file_com_iabtechlab_openrtb_v3_request_proto_rawDescGZIP(), []int{3}
}

func (x *Deal) GetId() string {
	if x != nil && x.Id != nil {
		return *x.Id
	}
	return ""
}

func (x *Deal) GetFlr() float32 {
	if x != nil && x.Flr != nil {
		return *x.Flr
	}
	return 0
}

func (x *Deal) GetFlrcur() string {
	if x != nil && x.Flrcur != nil {
		return *x.Flrcur
	}
	return ""
}

func (x *Deal) GetAt() int32 {
	if x != nil && x.At != nil {
		return *x.At
	}
	return 0
}

func (x *Deal) GetWseat() []string {
	if x != nil {
		return x.Wseat
	}
	return nil
}

func (x *Deal) GetWadomain() []string {
	if x != nil {
		return x.Wadomain
	}
	return nil
}

// *
// This object is associated with an item as an array of metrics. These metrics
// can offer insight to assist with decisioning such as average recent
// viewability, click-through rate, etc. Each metric is identified by its type,
// reports the value of the metric, and optionally identifies the source or
// vendor measuring the value.
type Metric struct {
	state           protoimpl.MessageState
	sizeCache       protoimpl.SizeCache
	unknownFields   protoimpl.UnknownFields
	extensionFields protoimpl.ExtensionFields

	// Type of metric being presented using exchange curated string names which
	// should be published to bidders a priori.
	Type *string `protobuf:"bytes,1,opt,name=type" json:"type,omitempty"`
	// Number representing the value of the metric. Probabilities must be in the
	// range 0.0 – 1.0.
	Value *float32 `protobuf:"fixed32,2,opt,name=value" json:"value,omitempty"`
	// Source of the value using exchange curated string names which should be
	// published to bidders a priori. If the exchange itself is the source versus
	// a third party, “EXCHANGE” is recommended.
	Vendor *string `protobuf:"bytes,3,opt,name=vendor" json:"vendor,omitempty"`
}

func (x *Metric) Reset() {
	*x = Metric{}
	if protoimpl.UnsafeEnabled {
		mi := &file_com_iabtechlab_openrtb_v3_request_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Metric) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Metric) ProtoMessage() {}

func (x *Metric) ProtoReflect() protoreflect.Message {
	mi := &file_com_iabtechlab_openrtb_v3_request_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Metric.ProtoReflect.Descriptor instead.
func (*Metric) Descriptor() ([]byte, []int) {
	return file_com_iabtechlab_openrtb_v3_request_proto_rawDescGZIP(), []int{4}
}

func (x *Metric) GetType() string {
	if x != nil && x.Type != nil {
		return *x.Type
	}
	return ""
}

func (x *Metric) GetValue() float32 {
	if x != nil && x.Value != nil {
		return *x.Value
	}
	return 0
}

func (x *Metric) GetVendor() string {
	if x != nil && x.Vendor != nil {
		return *x.Vendor
	}
	return ""
}

var File_com_iabtechlab_openrtb_v3_request_proto protoreflect.FileDescriptor

var file_com_iabtechlab_openrtb_v3_request_proto_rawDesc = []byte{
	0x0a, 0x27, 0x63, 0x6f, 0x6d, 0x2f, 0x69, 0x61, 0x62, 0x74, 0x65, 0x63, 0x68, 0x6c, 0x61, 0x62,
	0x2f, 0x6f, 0x70, 0x65, 0x6e, 0x72, 0x74, 0x62, 0x2f, 0x76, 0x33, 0x2f, 0x72, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x19, 0x63, 0x6f, 0x6d, 0x2e, 0x69,
	0x61, 0x62, 0x74, 0x65, 0x63, 0x68, 0x6c, 0x61, 0x62, 0x2e, 0x6f, 0x70, 0x65, 0x6e, 0x72, 0x74,
	0x62, 0x2e, 0x76, 0x33, 0x1a, 0x25, 0x63, 0x6f, 0x6d, 0x2f, 0x69, 0x61, 0x62, 0x74, 0x65, 0x63,
	0x68, 0x6c, 0x61, 0x62, 0x2f, 0x6f, 0x70, 0x65, 0x6e, 0x72, 0x74, 0x62, 0x2f, 0x76, 0x33, 0x2f,
	0x65, 0x6e, 0x75, 0x6d, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xce, 0x02, 0x0a, 0x07,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x65, 0x73, 0x74, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x04, 0x74, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x74,
	0x6d, 0x61, 0x78, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x04, 0x74, 0x6d, 0x61, 0x78, 0x12,
	0x0e, 0x0a, 0x02, 0x61, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x61, 0x74, 0x12,
	0x10, 0x0a, 0x03, 0x63, 0x75, 0x72, 0x18, 0x05, 0x20, 0x03, 0x28, 0x09, 0x52, 0x03, 0x63, 0x75,
	0x72, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x65, 0x61, 0x74, 0x18, 0x06, 0x20, 0x03, 0x28, 0x09, 0x52,
	0x04, 0x73, 0x65, 0x61, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x77, 0x73, 0x65, 0x61, 0x74, 0x18, 0x07,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x77, 0x73, 0x65, 0x61, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x63,
	0x64, 0x61, 0x74, 0x61, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x63, 0x64, 0x61, 0x74,
	0x61, 0x12, 0x39, 0x0a, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x21, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x69, 0x61, 0x62, 0x74, 0x65, 0x63, 0x68, 0x6c,
	0x61, 0x62, 0x2e, 0x6f, 0x70, 0x65, 0x6e, 0x72, 0x74, 0x62, 0x2e, 0x76, 0x33, 0x2e, 0x53, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x52, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x33, 0x0a, 0x04,
	0x69, 0x74, 0x65, 0x6d, 0x18, 0x0a, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x63, 0x6f, 0x6d,
	0x2e, 0x69, 0x61, 0x62, 0x74, 0x65, 0x63, 0x68, 0x6c, 0x61, 0x62, 0x2e, 0x6f, 0x70, 0x65, 0x6e,
	0x72, 0x74, 0x62, 0x2e, 0x76, 0x33, 0x2e, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x04, 0x69, 0x74, 0x65,
	0x6d, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x61, 0x63, 0x6b, 0x61, 0x67, 0x65, 0x18, 0x0b, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x07, 0x70, 0x61, 0x63, 0x6b, 0x61, 0x67, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x63,
	0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x63, 0x6f,
	0x6e, 0x74, 0x65, 0x78, 0x74, 0x2a, 0x05, 0x08, 0x64, 0x10, 0x90, 0x4e, 0x22, 0x9b, 0x01, 0x0a,
	0x06, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x74, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x74, 0x69, 0x64, 0x12, 0x0e, 0x0a, 0x02, 0x74, 0x73, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x74, 0x73, 0x12, 0x0e, 0x0a, 0x02, 0x64, 0x73, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x64, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x64, 0x73, 0x6d,
	0x61, 0x70, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x64, 0x73, 0x6d, 0x61, 0x70, 0x12,
	0x12, 0x0a, 0x04, 0x63, 0x65, 0x72, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x63,
	0x65, 0x72, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x64, 0x69, 0x67, 0x65, 0x73, 0x74, 0x18, 0x06, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x64, 0x69, 0x67, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x70,
	0x63, 0x68, 0x61, 0x69, 0x6e, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x70, 0x63, 0x68,
	0x61, 0x69, 0x6e, 0x2a, 0x05, 0x08, 0x64, 0x10, 0x90, 0x4e, 0x22, 0xe8, 0x02, 0x0a, 0x04, 0x49,
	0x74, 0x65, 0x6d, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x03, 0x71, 0x74, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d,
	0x48, 0x00, 0x52, 0x03, 0x71, 0x74, 0x79, 0x12, 0x18, 0x0a, 0x06, 0x71, 0x74, 0x79, 0x66, 0x6c,
	0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x02, 0x48, 0x00, 0x52, 0x06, 0x71, 0x74, 0x79, 0x66, 0x6c,
	0x74, 0x12, 0x10, 0x0a, 0x03, 0x73, 0x65, 0x71, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x03,
	0x73, 0x65, 0x71, 0x12, 0x10, 0x0a, 0x03, 0x66, 0x6c, 0x72, 0x18, 0x05, 0x20, 0x01, 0x28, 0x02,
	0x52, 0x03, 0x66, 0x6c, 0x72, 0x12, 0x16, 0x0a, 0x06, 0x66, 0x6c, 0x72, 0x63, 0x75, 0x72, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x66, 0x6c, 0x72, 0x63, 0x75, 0x72, 0x12, 0x10, 0x0a,
	0x03, 0x65, 0x78, 0x70, 0x18, 0x07, 0x20, 0x01, 0x28, 0x04, 0x52, 0x03, 0x65, 0x78, 0x70, 0x12,
	0x0e, 0x0a, 0x02, 0x64, 0x74, 0x18, 0x08, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x64, 0x74, 0x12,
	0x12, 0x0a, 0x04, 0x64, 0x6c, 0x76, 0x79, 0x18, 0x09, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x64,
	0x6c, 0x76, 0x79, 0x12, 0x39, 0x0a, 0x06, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x18, 0x0a, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x69, 0x61, 0x62, 0x74, 0x65, 0x63,
	0x68, 0x6c, 0x61, 0x62, 0x2e, 0x6f, 0x70, 0x65, 0x6e, 0x72, 0x74, 0x62, 0x2e, 0x76, 0x33, 0x2e,
	0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x06, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x12, 0x33,
	0x0a, 0x04, 0x64, 0x65, 0x61, 0x6c, 0x18, 0x0b, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x63,
	0x6f, 0x6d, 0x2e, 0x69, 0x61, 0x62, 0x74, 0x65, 0x63, 0x68, 0x6c, 0x61, 0x62, 0x2e, 0x6f, 0x70,
	0x65, 0x6e, 0x72, 0x74, 0x62, 0x2e, 0x76, 0x33, 0x2e, 0x44, 0x65, 0x61, 0x6c, 0x52, 0x04, 0x64,
	0x65, 0x61, 0x6c, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x18, 0x0c,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x70, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x12, 0x12, 0x0a,
	0x04, 0x73, 0x70, 0x65, 0x63, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x73, 0x70, 0x65,
	0x63, 0x2a, 0x05, 0x08, 0x64, 0x10, 0x90, 0x4e, 0x42, 0x0b, 0x0a, 0x09, 0x71, 0x74, 0x79, 0x5f,
	0x6f, 0x6e, 0x65, 0x6f, 0x66, 0x22, 0x89, 0x01, 0x0a, 0x04, 0x44, 0x65, 0x61, 0x6c, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x10,
	0x0a, 0x03, 0x66, 0x6c, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x02, 0x52, 0x03, 0x66, 0x6c, 0x72,
	0x12, 0x16, 0x0a, 0x06, 0x66, 0x6c, 0x72, 0x63, 0x75, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x66, 0x6c, 0x72, 0x63, 0x75, 0x72, 0x12, 0x0e, 0x0a, 0x02, 0x61, 0x74, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x61, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x77, 0x73, 0x65, 0x61,
	0x74, 0x18, 0x05, 0x20, 0x03, 0x28, 0x09, 0x52, 0x05, 0x77, 0x73, 0x65, 0x61, 0x74, 0x12, 0x1a,
	0x0a, 0x08, 0x77, 0x61, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x18, 0x06, 0x20, 0x03, 0x28, 0x09,
	0x52, 0x08, 0x77, 0x61, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x2a, 0x05, 0x08, 0x64, 0x10, 0x90,
	0x4e, 0x22, 0x51, 0x0a, 0x06, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x12, 0x12, 0x0a, 0x04, 0x74,
	0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12,
	0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x02, 0x52, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x76, 0x65, 0x6e, 0x64, 0x6f, 0x72, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x76, 0x65, 0x6e, 0x64, 0x6f, 0x72, 0x2a, 0x05, 0x08,
	0x64, 0x10, 0x90, 0x4e, 0x42, 0x85, 0x02, 0x0a, 0x1d, 0x63, 0x6f, 0x6d, 0x2e, 0x63, 0x6f, 0x6d,
	0x2e, 0x69, 0x61, 0x62, 0x74, 0x65, 0x63, 0x68, 0x6c, 0x61, 0x62, 0x2e, 0x6f, 0x70, 0x65, 0x6e,
	0x72, 0x74, 0x62, 0x2e, 0x76, 0x33, 0x42, 0x0c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x50,
	0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x4f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x62, 0x69, 0x64, 0x6f, 0x6e, 0x2d, 0x69, 0x6f, 0x2f, 0x62, 0x69, 0x64, 0x6f,
	0x6e, 0x2d, 0x62, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x6d, 0x2f, 0x69, 0x61, 0x62, 0x74, 0x65, 0x63, 0x68, 0x6c,
	0x61, 0x62, 0x2f, 0x6f, 0x70, 0x65, 0x6e, 0x72, 0x74, 0x62, 0x2f, 0x76, 0x33, 0x3b, 0x6f, 0x70,
	0x65, 0x6e, 0x72, 0x74, 0x62, 0x76, 0x33, 0xa2, 0x02, 0x03, 0x43, 0x49, 0x4f, 0xaa, 0x02, 0x19,
	0x43, 0x6f, 0x6d, 0x2e, 0x49, 0x61, 0x62, 0x74, 0x65, 0x63, 0x68, 0x6c, 0x61, 0x62, 0x2e, 0x4f,
	0x70, 0x65, 0x6e, 0x72, 0x74, 0x62, 0x2e, 0x56, 0x33, 0xca, 0x02, 0x19, 0x43, 0x6f, 0x6d, 0x5c,
	0x49, 0x61, 0x62, 0x74, 0x65, 0x63, 0x68, 0x6c, 0x61, 0x62, 0x5c, 0x4f, 0x70, 0x65, 0x6e, 0x72,
	0x74, 0x62, 0x5c, 0x56, 0x33, 0xe2, 0x02, 0x25, 0x43, 0x6f, 0x6d, 0x5c, 0x49, 0x61, 0x62, 0x74,
	0x65, 0x63, 0x68, 0x6c, 0x61, 0x62, 0x5c, 0x4f, 0x70, 0x65, 0x6e, 0x72, 0x74, 0x62, 0x5c, 0x56,
	0x33, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x1c,
	0x43, 0x6f, 0x6d, 0x3a, 0x3a, 0x49, 0x61, 0x62, 0x74, 0x65, 0x63, 0x68, 0x6c, 0x61, 0x62, 0x3a,
	0x3a, 0x4f, 0x70, 0x65, 0x6e, 0x72, 0x74, 0x62, 0x3a, 0x3a, 0x56, 0x33,
}

var (
	file_com_iabtechlab_openrtb_v3_request_proto_rawDescOnce sync.Once
	file_com_iabtechlab_openrtb_v3_request_proto_rawDescData = file_com_iabtechlab_openrtb_v3_request_proto_rawDesc
)

func file_com_iabtechlab_openrtb_v3_request_proto_rawDescGZIP() []byte {
	file_com_iabtechlab_openrtb_v3_request_proto_rawDescOnce.Do(func() {
		file_com_iabtechlab_openrtb_v3_request_proto_rawDescData = protoimpl.X.CompressGZIP(file_com_iabtechlab_openrtb_v3_request_proto_rawDescData)
	})
	return file_com_iabtechlab_openrtb_v3_request_proto_rawDescData
}

var file_com_iabtechlab_openrtb_v3_request_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_com_iabtechlab_openrtb_v3_request_proto_goTypes = []any{
	(*Request)(nil), // 0: com.iabtechlab.openrtb.v3.Request
	(*Source)(nil),  // 1: com.iabtechlab.openrtb.v3.Source
	(*Item)(nil),    // 2: com.iabtechlab.openrtb.v3.Item
	(*Deal)(nil),    // 3: com.iabtechlab.openrtb.v3.Deal
	(*Metric)(nil),  // 4: com.iabtechlab.openrtb.v3.Metric
}
var file_com_iabtechlab_openrtb_v3_request_proto_depIdxs = []int32{
	1, // 0: com.iabtechlab.openrtb.v3.Request.source:type_name -> com.iabtechlab.openrtb.v3.Source
	2, // 1: com.iabtechlab.openrtb.v3.Request.item:type_name -> com.iabtechlab.openrtb.v3.Item
	4, // 2: com.iabtechlab.openrtb.v3.Item.metric:type_name -> com.iabtechlab.openrtb.v3.Metric
	3, // 3: com.iabtechlab.openrtb.v3.Item.deal:type_name -> com.iabtechlab.openrtb.v3.Deal
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_com_iabtechlab_openrtb_v3_request_proto_init() }
func file_com_iabtechlab_openrtb_v3_request_proto_init() {
	if File_com_iabtechlab_openrtb_v3_request_proto != nil {
		return
	}
	file_com_iabtechlab_openrtb_v3_enums_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_com_iabtechlab_openrtb_v3_request_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*Request); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			case 3:
				return &v.extensionFields
			default:
				return nil
			}
		}
		file_com_iabtechlab_openrtb_v3_request_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*Source); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			case 3:
				return &v.extensionFields
			default:
				return nil
			}
		}
		file_com_iabtechlab_openrtb_v3_request_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*Item); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			case 3:
				return &v.extensionFields
			default:
				return nil
			}
		}
		file_com_iabtechlab_openrtb_v3_request_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*Deal); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			case 3:
				return &v.extensionFields
			default:
				return nil
			}
		}
		file_com_iabtechlab_openrtb_v3_request_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*Metric); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			case 3:
				return &v.extensionFields
			default:
				return nil
			}
		}
	}
	file_com_iabtechlab_openrtb_v3_request_proto_msgTypes[2].OneofWrappers = []any{
		(*Item_Qty)(nil),
		(*Item_Qtyflt)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_com_iabtechlab_openrtb_v3_request_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_com_iabtechlab_openrtb_v3_request_proto_goTypes,
		DependencyIndexes: file_com_iabtechlab_openrtb_v3_request_proto_depIdxs,
		MessageInfos:      file_com_iabtechlab_openrtb_v3_request_proto_msgTypes,
	}.Build()
	File_com_iabtechlab_openrtb_v3_request_proto = out.File
	file_com_iabtechlab_openrtb_v3_request_proto_rawDesc = nil
	file_com_iabtechlab_openrtb_v3_request_proto_goTypes = nil
	file_com_iabtechlab_openrtb_v3_request_proto_depIdxs = nil
}
