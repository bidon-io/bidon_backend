//*
// The response object contains minimal high level attributes (e.g., reference
// to the request ID, bid currency, etc.) and an array of seat bids, each of
// which is a set of bids on behalf of a buyer seat.
//
// The individual bid references the item in the request to which it pertains
// and buying information such as the price, a deal ID if applicable, and
// notification URLs. The media related to a bid is conveyed via Layer-4 domain
// objects (i.e., ad creative, markup) included in each bid.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        (unknown)
// source: com/iabtechlab/openrtb/v3/response.proto

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
// This object is the bid response object under the Openrtb root. Its id
// attribute is a reflection of the bid request ID. The bidid attribute is an
// optional response tracking ID for bidders. If specified, it will be available
// for use in substitution macros placed in markup and notification URLs. At
// least one Seatbid object is required, which contains at least one Bid for an
// item. Other attributes are optional.
//
// To express a “no-bid”, the most compact option is simply to return an empty
// response with HTTP 204. However, if the bidder wishes to convey a reason for
// not bidding, a Response object can be returned with just a reason code in the
// nbr attribute.
type Response struct {
	state           protoimpl.MessageState
	sizeCache       protoimpl.SizeCache
	unknownFields   protoimpl.UnknownFields
	extensionFields protoimpl.ExtensionFields

	// ID of the bid request to which this is a response; must match the
	// request.id attribute.
	Id *string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	// Bidder generated response ID to assist with logging/tracking.
	Bidid *string `protobuf:"bytes,2,opt,name=bidid" json:"bidid,omitempty"`
	// Reason for not bidding if applicable (see List: No-Bid Reason Codes). Note
	// that while many exchanges prefer a simple HTTP 204 response to indicate a
	// no-bid, responses indicating a reason code can be useful in debugging
	// scenarios. Common values are defined in the enumeration
	// com.iabtechlab.openrtb.v3.NoBidReason.
	Nbr *int32 `protobuf:"varint,3,opt,name=nbr" json:"nbr,omitempty"`
	// Bid currency using ISO-4217 alpha codes.
	Cur *string `protobuf:"bytes,4,opt,name=cur" json:"cur,omitempty"`
	// Allows bidder to set data in the exchange’s cookie, which can be retrieved
	// on bid requests (refer to cdata in Object: Request) if supported by the
	// exchange. The string must be in base85 cookie-safe characters.
	Cdata *string `protobuf:"bytes,5,opt,name=cdata" json:"cdata,omitempty"`
	// Array of Seatbid objects; 1+ required if a bid is to be made. Refer to
	// Object: Seatbid.
	Seatbid []*SeatBid `protobuf:"bytes,6,rep,name=seatbid" json:"seatbid,omitempty"`
}

func (x *Response) Reset() {
	*x = Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_com_iabtechlab_openrtb_v3_response_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Response) ProtoMessage() {}

func (x *Response) ProtoReflect() protoreflect.Message {
	mi := &file_com_iabtechlab_openrtb_v3_response_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Response.ProtoReflect.Descriptor instead.
func (*Response) Descriptor() ([]byte, []int) {
	return file_com_iabtechlab_openrtb_v3_response_proto_rawDescGZIP(), []int{0}
}

func (x *Response) GetId() string {
	if x != nil && x.Id != nil {
		return *x.Id
	}
	return ""
}

func (x *Response) GetBidid() string {
	if x != nil && x.Bidid != nil {
		return *x.Bidid
	}
	return ""
}

func (x *Response) GetNbr() int32 {
	if x != nil && x.Nbr != nil {
		return *x.Nbr
	}
	return 0
}

func (x *Response) GetCur() string {
	if x != nil && x.Cur != nil {
		return *x.Cur
	}
	return ""
}

func (x *Response) GetCdata() string {
	if x != nil && x.Cdata != nil {
		return *x.Cdata
	}
	return ""
}

func (x *Response) GetSeatbid() []*SeatBid {
	if x != nil {
		return x.Seatbid
	}
	return nil
}

// *
// A bid response can contain multiple Seatbid objects, each on behalf of a
// different buyer seat and each containing one or more individual bids. If
// multiple items are presented in the request offer, the package attribute can
// be used to specify if a seat is willing to accept any impressions that it can
// win (default) or if it is interested in winning any only if it can win them
// all as a group.
type SeatBid struct {
	state           protoimpl.MessageState
	sizeCache       protoimpl.SizeCache
	unknownFields   protoimpl.UnknownFields
	extensionFields protoimpl.ExtensionFields

	// ID of the buyer seat on whose behalf this bid is made.
	Seat *string `protobuf:"bytes,1,opt,name=seat" json:"seat,omitempty"`
	// For offers with multiple items, this flag Indicates if the bidder is
	// willing to accept wins on a subset of bids or requires the full group as a
	// package, where 0 = individual wins accepted; 1 = package win or loss only.
	Package *bool `protobuf:"varint,2,opt,name=package" json:"package,omitempty"`
	// Array of 1+ Bid objects each related to an item. Multiple bids can relate
	// to the same item. Refer to Object: Bid.
	Bid []*Bid `protobuf:"bytes,3,rep,name=bid" json:"bid,omitempty"`
}

func (x *SeatBid) Reset() {
	*x = SeatBid{}
	if protoimpl.UnsafeEnabled {
		mi := &file_com_iabtechlab_openrtb_v3_response_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SeatBid) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SeatBid) ProtoMessage() {}

func (x *SeatBid) ProtoReflect() protoreflect.Message {
	mi := &file_com_iabtechlab_openrtb_v3_response_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SeatBid.ProtoReflect.Descriptor instead.
func (*SeatBid) Descriptor() ([]byte, []int) {
	return file_com_iabtechlab_openrtb_v3_response_proto_rawDescGZIP(), []int{1}
}

func (x *SeatBid) GetSeat() string {
	if x != nil && x.Seat != nil {
		return *x.Seat
	}
	return ""
}

func (x *SeatBid) GetPackage() bool {
	if x != nil && x.Package != nil {
		return *x.Package
	}
	return false
}

func (x *SeatBid) GetBid() []*Bid {
	if x != nil {
		return x.Bid
	}
	return nil
}

// *
// A Seatbid object contains one or more Bid objects, each of which relates to a
// specific item in the bid request offer via the “item” attribute and
// constitutes an offer to buy that item for a given price.
type Bid struct {
	state           protoimpl.MessageState
	sizeCache       protoimpl.SizeCache
	unknownFields   protoimpl.UnknownFields
	extensionFields protoimpl.ExtensionFields

	// Bidder generated bid ID to assist with logging/tracking.
	Id *string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	// ID of the item object in the related bid request; specifically item.id.
	Item *string `protobuf:"bytes,2,opt,name=item" json:"item,omitempty"`
	// Bid price expressed as CPM although the actual transaction is for a unit
	// item only. Note that while the type indicates float, integer math is highly
	// recommended when handling currencies (e.g., BigDecimal in Java).
	Price *float32 `protobuf:"fixed32,3,opt,name=price" json:"price,omitempty"`
	// Reference to a deal from the bid request if this bid pertains to a private
	// marketplace deal; specifically deal.id.
	Deal *string `protobuf:"bytes,4,opt,name=deal" json:"deal,omitempty"`
	// Campaign ID or other similar grouping of brand-related ads. Typically used
	// to increase the efficiency of audit processes.
	Cid *string `protobuf:"bytes,5,opt,name=cid" json:"cid,omitempty"`
	// Tactic ID to enable buyers to label bids for reporting to the exchange the
	// tactic through which their bid was submitted. The specific usage and
	// meaning of the tactic ID should be communicated between buyer and exchanges
	// a priori.
	Tactic *string `protobuf:"bytes,6,opt,name=tactic" json:"tactic,omitempty"`
	// Pending notice URL called by the exchange when a bid has been declared the
	// winner within the scope of an OpenRTB compliant supply chain (i.e., there
	// may still be non-compliant decisioning such as header bidding).
	// Substitution macros may be included.
	Purl *string `protobuf:"bytes,7,opt,name=purl" json:"purl,omitempty"`
	// Billing notice URL called by the exchange when a winning bid becomes
	// billable based on exchange-specific business policy (e.g., markup
	// rendered). Substitution macros may be included.
	Burl *string `protobuf:"bytes,8,opt,name=burl" json:"burl,omitempty"`
	// Loss notice URL called by the exchange when a bid is known to have been
	// lost. Substitution macros may be included. Exchange-specific policy may
	// preclude support for loss notices or the disclosure of winning clearing
	// prices resulting in ${OPENRTB_PRICE} macros being removed (i.e., replaced
	// with a zero-length string).
	Lurl *string `protobuf:"bytes,9,opt,name=lurl" json:"lurl,omitempty"`
	// Advisory as to the number of seconds the buyer is willing to wait between
	// auction and fulfilment.
	Exp *uint64 `protobuf:"varint,10,opt,name=exp" json:"exp,omitempty"`
	// ID to enable media to be specified by reference if previously uploaded to
	// the exchange rather than including it by value in the domain objects.
	Mid *string `protobuf:"bytes,11,opt,name=mid" json:"mid,omitempty"`
	// Array of Macro objects that enable bid specific values to be substituted
	// into markup; especially useful for previously uploaded media referenced via
	// the mid attribute. Refer to Object: Macro.
	Macro []*Macro `protobuf:"bytes,12,rep,name=macro" json:"macro,omitempty"`
	// Layer-4 domain object structure that specifies the media to be presented if
	// the bid is won conforming to the specification and version referenced in
	// openrtb.domainspec and openrtb.domainver. For AdCOM v1.x, the objects
	// allowed here are “Ad” and any objects subordinate thereto as specified by
	// AdCOM.
	Media []byte `protobuf:"bytes,13,opt,name=media" json:"media,omitempty"`
}

func (x *Bid) Reset() {
	*x = Bid{}
	if protoimpl.UnsafeEnabled {
		mi := &file_com_iabtechlab_openrtb_v3_response_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Bid) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Bid) ProtoMessage() {}

func (x *Bid) ProtoReflect() protoreflect.Message {
	mi := &file_com_iabtechlab_openrtb_v3_response_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Bid.ProtoReflect.Descriptor instead.
func (*Bid) Descriptor() ([]byte, []int) {
	return file_com_iabtechlab_openrtb_v3_response_proto_rawDescGZIP(), []int{2}
}

func (x *Bid) GetId() string {
	if x != nil && x.Id != nil {
		return *x.Id
	}
	return ""
}

func (x *Bid) GetItem() string {
	if x != nil && x.Item != nil {
		return *x.Item
	}
	return ""
}

func (x *Bid) GetPrice() float32 {
	if x != nil && x.Price != nil {
		return *x.Price
	}
	return 0
}

func (x *Bid) GetDeal() string {
	if x != nil && x.Deal != nil {
		return *x.Deal
	}
	return ""
}

func (x *Bid) GetCid() string {
	if x != nil && x.Cid != nil {
		return *x.Cid
	}
	return ""
}

func (x *Bid) GetTactic() string {
	if x != nil && x.Tactic != nil {
		return *x.Tactic
	}
	return ""
}

func (x *Bid) GetPurl() string {
	if x != nil && x.Purl != nil {
		return *x.Purl
	}
	return ""
}

func (x *Bid) GetBurl() string {
	if x != nil && x.Burl != nil {
		return *x.Burl
	}
	return ""
}

func (x *Bid) GetLurl() string {
	if x != nil && x.Lurl != nil {
		return *x.Lurl
	}
	return ""
}

func (x *Bid) GetExp() uint64 {
	if x != nil && x.Exp != nil {
		return *x.Exp
	}
	return 0
}

func (x *Bid) GetMid() string {
	if x != nil && x.Mid != nil {
		return *x.Mid
	}
	return ""
}

func (x *Bid) GetMacro() []*Macro {
	if x != nil {
		return x.Macro
	}
	return nil
}

func (x *Bid) GetMedia() []byte {
	if x != nil {
		return x.Media
	}
	return nil
}

// *
// This object constitutes a buyer defined key/value pair used to inject dynamic
// values into media markup. While they apply to any media markup irrespective
// of how it is conveyed, the principle use case is for media that was uploaded
// to the exchange prior to the transaction (e.g., pre-registered for creative
// quality review) and referenced in bid. The full form of the macro to be
// substituted at runtime is ${CUSTOM_KEY}, where “KEY” is the name supplied in
// the key attribute. This ensures no conflict with standard OpenRTB macros.
type Macro struct {
	state           protoimpl.MessageState
	sizeCache       protoimpl.SizeCache
	unknownFields   protoimpl.UnknownFields
	extensionFields protoimpl.ExtensionFields

	// Name of a buyer specific macro.
	Key *string `protobuf:"bytes,1,opt,name=key" json:"key,omitempty"`
	// Value to substitute for each instance of the macro found in markup.
	Value *string `protobuf:"bytes,2,opt,name=value" json:"value,omitempty"`
}

func (x *Macro) Reset() {
	*x = Macro{}
	if protoimpl.UnsafeEnabled {
		mi := &file_com_iabtechlab_openrtb_v3_response_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Macro) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Macro) ProtoMessage() {}

func (x *Macro) ProtoReflect() protoreflect.Message {
	mi := &file_com_iabtechlab_openrtb_v3_response_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Macro.ProtoReflect.Descriptor instead.
func (*Macro) Descriptor() ([]byte, []int) {
	return file_com_iabtechlab_openrtb_v3_response_proto_rawDescGZIP(), []int{3}
}

func (x *Macro) GetKey() string {
	if x != nil && x.Key != nil {
		return *x.Key
	}
	return ""
}

func (x *Macro) GetValue() string {
	if x != nil && x.Value != nil {
		return *x.Value
	}
	return ""
}

var File_com_iabtechlab_openrtb_v3_response_proto protoreflect.FileDescriptor

var file_com_iabtechlab_openrtb_v3_response_proto_rawDesc = []byte{
	0x0a, 0x28, 0x63, 0x6f, 0x6d, 0x2f, 0x69, 0x61, 0x62, 0x74, 0x65, 0x63, 0x68, 0x6c, 0x61, 0x62,
	0x2f, 0x6f, 0x70, 0x65, 0x6e, 0x72, 0x74, 0x62, 0x2f, 0x76, 0x33, 0x2f, 0x72, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x19, 0x63, 0x6f, 0x6d, 0x2e,
	0x69, 0x61, 0x62, 0x74, 0x65, 0x63, 0x68, 0x6c, 0x61, 0x62, 0x2e, 0x6f, 0x70, 0x65, 0x6e, 0x72,
	0x74, 0x62, 0x2e, 0x76, 0x33, 0x1a, 0x25, 0x63, 0x6f, 0x6d, 0x2f, 0x69, 0x61, 0x62, 0x74, 0x65,
	0x63, 0x68, 0x6c, 0x61, 0x62, 0x2f, 0x6f, 0x70, 0x65, 0x6e, 0x72, 0x74, 0x62, 0x2f, 0x76, 0x33,
	0x2f, 0x65, 0x6e, 0x75, 0x6d, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xaf, 0x01, 0x0a,
	0x08, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x62, 0x69, 0x64,
	0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x62, 0x69, 0x64, 0x69, 0x64, 0x12,
	0x10, 0x0a, 0x03, 0x6e, 0x62, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x6e, 0x62,
	0x72, 0x12, 0x10, 0x0a, 0x03, 0x63, 0x75, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x63, 0x75, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x64, 0x61, 0x74, 0x61, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x63, 0x64, 0x61, 0x74, 0x61, 0x12, 0x3c, 0x0a, 0x07, 0x73, 0x65, 0x61,
	0x74, 0x62, 0x69, 0x64, 0x18, 0x06, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x63, 0x6f, 0x6d,
	0x2e, 0x69, 0x61, 0x62, 0x74, 0x65, 0x63, 0x68, 0x6c, 0x61, 0x62, 0x2e, 0x6f, 0x70, 0x65, 0x6e,
	0x72, 0x74, 0x62, 0x2e, 0x76, 0x33, 0x2e, 0x53, 0x65, 0x61, 0x74, 0x42, 0x69, 0x64, 0x52, 0x07,
	0x73, 0x65, 0x61, 0x74, 0x62, 0x69, 0x64, 0x2a, 0x05, 0x08, 0x64, 0x10, 0x90, 0x4e, 0x22, 0x70,
	0x0a, 0x07, 0x53, 0x65, 0x61, 0x74, 0x42, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x65, 0x61,
	0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x73, 0x65, 0x61, 0x74, 0x12, 0x18, 0x0a,
	0x07, 0x70, 0x61, 0x63, 0x6b, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07,
	0x70, 0x61, 0x63, 0x6b, 0x61, 0x67, 0x65, 0x12, 0x30, 0x0a, 0x03, 0x62, 0x69, 0x64, 0x18, 0x03,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x69, 0x61, 0x62, 0x74, 0x65,
	0x63, 0x68, 0x6c, 0x61, 0x62, 0x2e, 0x6f, 0x70, 0x65, 0x6e, 0x72, 0x74, 0x62, 0x2e, 0x76, 0x33,
	0x2e, 0x42, 0x69, 0x64, 0x52, 0x03, 0x62, 0x69, 0x64, 0x2a, 0x05, 0x08, 0x64, 0x10, 0x90, 0x4e,
	0x22, 0xb2, 0x02, 0x0a, 0x03, 0x42, 0x69, 0x64, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x69, 0x74, 0x65, 0x6d,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x69, 0x74, 0x65, 0x6d, 0x12, 0x14, 0x0a, 0x05,
	0x70, 0x72, 0x69, 0x63, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x02, 0x52, 0x05, 0x70, 0x72, 0x69,
	0x63, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x65, 0x61, 0x6c, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x64, 0x65, 0x61, 0x6c, 0x12, 0x10, 0x0a, 0x03, 0x63, 0x69, 0x64, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x03, 0x63, 0x69, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x74, 0x61, 0x63, 0x74,
	0x69, 0x63, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x74, 0x61, 0x63, 0x74, 0x69, 0x63,
	0x12, 0x12, 0x0a, 0x04, 0x70, 0x75, 0x72, 0x6c, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x70, 0x75, 0x72, 0x6c, 0x12, 0x12, 0x0a, 0x04, 0x62, 0x75, 0x72, 0x6c, 0x18, 0x08, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x62, 0x75, 0x72, 0x6c, 0x12, 0x12, 0x0a, 0x04, 0x6c, 0x75, 0x72, 0x6c,
	0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6c, 0x75, 0x72, 0x6c, 0x12, 0x10, 0x0a, 0x03,
	0x65, 0x78, 0x70, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x04, 0x52, 0x03, 0x65, 0x78, 0x70, 0x12, 0x10,
	0x0a, 0x03, 0x6d, 0x69, 0x64, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6d, 0x69, 0x64,
	0x12, 0x36, 0x0a, 0x05, 0x6d, 0x61, 0x63, 0x72, 0x6f, 0x18, 0x0c, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x20, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x69, 0x61, 0x62, 0x74, 0x65, 0x63, 0x68, 0x6c, 0x61, 0x62,
	0x2e, 0x6f, 0x70, 0x65, 0x6e, 0x72, 0x74, 0x62, 0x2e, 0x76, 0x33, 0x2e, 0x4d, 0x61, 0x63, 0x72,
	0x6f, 0x52, 0x05, 0x6d, 0x61, 0x63, 0x72, 0x6f, 0x12, 0x14, 0x0a, 0x05, 0x6d, 0x65, 0x64, 0x69,
	0x61, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x2a, 0x05,
	0x08, 0x64, 0x10, 0x90, 0x4e, 0x22, 0x36, 0x0a, 0x05, 0x4d, 0x61, 0x63, 0x72, 0x6f, 0x12, 0x10,
	0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79,
	0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x2a, 0x05, 0x08, 0x64, 0x10, 0x90, 0x4e, 0x42, 0x86, 0x02,
	0x0a, 0x1d, 0x63, 0x6f, 0x6d, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x69, 0x61, 0x62, 0x74, 0x65, 0x63,
	0x68, 0x6c, 0x61, 0x62, 0x2e, 0x6f, 0x70, 0x65, 0x6e, 0x72, 0x74, 0x62, 0x2e, 0x76, 0x33, 0x42,
	0x0d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01,
	0x5a, 0x4f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x62, 0x69, 0x64,
	0x6f, 0x6e, 0x2d, 0x69, 0x6f, 0x2f, 0x62, 0x69, 0x64, 0x6f, 0x6e, 0x2d, 0x62, 0x61, 0x63, 0x6b,
	0x65, 0x6e, 0x64, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f,
	0x6d, 0x2f, 0x69, 0x61, 0x62, 0x74, 0x65, 0x63, 0x68, 0x6c, 0x61, 0x62, 0x2f, 0x6f, 0x70, 0x65,
	0x6e, 0x72, 0x74, 0x62, 0x2f, 0x76, 0x33, 0x3b, 0x6f, 0x70, 0x65, 0x6e, 0x72, 0x74, 0x62, 0x76,
	0x33, 0xa2, 0x02, 0x03, 0x43, 0x49, 0x4f, 0xaa, 0x02, 0x19, 0x43, 0x6f, 0x6d, 0x2e, 0x49, 0x61,
	0x62, 0x74, 0x65, 0x63, 0x68, 0x6c, 0x61, 0x62, 0x2e, 0x4f, 0x70, 0x65, 0x6e, 0x72, 0x74, 0x62,
	0x2e, 0x56, 0x33, 0xca, 0x02, 0x19, 0x43, 0x6f, 0x6d, 0x5c, 0x49, 0x61, 0x62, 0x74, 0x65, 0x63,
	0x68, 0x6c, 0x61, 0x62, 0x5c, 0x4f, 0x70, 0x65, 0x6e, 0x72, 0x74, 0x62, 0x5c, 0x56, 0x33, 0xe2,
	0x02, 0x25, 0x43, 0x6f, 0x6d, 0x5c, 0x49, 0x61, 0x62, 0x74, 0x65, 0x63, 0x68, 0x6c, 0x61, 0x62,
	0x5c, 0x4f, 0x70, 0x65, 0x6e, 0x72, 0x74, 0x62, 0x5c, 0x56, 0x33, 0x5c, 0x47, 0x50, 0x42, 0x4d,
	0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x1c, 0x43, 0x6f, 0x6d, 0x3a, 0x3a, 0x49,
	0x61, 0x62, 0x74, 0x65, 0x63, 0x68, 0x6c, 0x61, 0x62, 0x3a, 0x3a, 0x4f, 0x70, 0x65, 0x6e, 0x72,
	0x74, 0x62, 0x3a, 0x3a, 0x56, 0x33,
}

var (
	file_com_iabtechlab_openrtb_v3_response_proto_rawDescOnce sync.Once
	file_com_iabtechlab_openrtb_v3_response_proto_rawDescData = file_com_iabtechlab_openrtb_v3_response_proto_rawDesc
)

func file_com_iabtechlab_openrtb_v3_response_proto_rawDescGZIP() []byte {
	file_com_iabtechlab_openrtb_v3_response_proto_rawDescOnce.Do(func() {
		file_com_iabtechlab_openrtb_v3_response_proto_rawDescData = protoimpl.X.CompressGZIP(file_com_iabtechlab_openrtb_v3_response_proto_rawDescData)
	})
	return file_com_iabtechlab_openrtb_v3_response_proto_rawDescData
}

var file_com_iabtechlab_openrtb_v3_response_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_com_iabtechlab_openrtb_v3_response_proto_goTypes = []any{
	(*Response)(nil), // 0: com.iabtechlab.openrtb.v3.Response
	(*SeatBid)(nil),  // 1: com.iabtechlab.openrtb.v3.SeatBid
	(*Bid)(nil),      // 2: com.iabtechlab.openrtb.v3.Bid
	(*Macro)(nil),    // 3: com.iabtechlab.openrtb.v3.Macro
}
var file_com_iabtechlab_openrtb_v3_response_proto_depIdxs = []int32{
	1, // 0: com.iabtechlab.openrtb.v3.Response.seatbid:type_name -> com.iabtechlab.openrtb.v3.SeatBid
	2, // 1: com.iabtechlab.openrtb.v3.SeatBid.bid:type_name -> com.iabtechlab.openrtb.v3.Bid
	3, // 2: com.iabtechlab.openrtb.v3.Bid.macro:type_name -> com.iabtechlab.openrtb.v3.Macro
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_com_iabtechlab_openrtb_v3_response_proto_init() }
func file_com_iabtechlab_openrtb_v3_response_proto_init() {
	if File_com_iabtechlab_openrtb_v3_response_proto != nil {
		return
	}
	file_com_iabtechlab_openrtb_v3_enums_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_com_iabtechlab_openrtb_v3_response_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*Response); i {
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
		file_com_iabtechlab_openrtb_v3_response_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*SeatBid); i {
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
		file_com_iabtechlab_openrtb_v3_response_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*Bid); i {
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
		file_com_iabtechlab_openrtb_v3_response_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*Macro); i {
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
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_com_iabtechlab_openrtb_v3_response_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_com_iabtechlab_openrtb_v3_response_proto_goTypes,
		DependencyIndexes: file_com_iabtechlab_openrtb_v3_response_proto_depIdxs,
		MessageInfos:      file_com_iabtechlab_openrtb_v3_response_proto_msgTypes,
	}.Build()
	File_com_iabtechlab_openrtb_v3_response_proto = out.File
	file_com_iabtechlab_openrtb_v3_response_proto_rawDesc = nil
	file_com_iabtechlab_openrtb_v3_response_proto_goTypes = nil
	file_com_iabtechlab_openrtb_v3_response_proto_depIdxs = nil
}
