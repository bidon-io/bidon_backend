package schema

import (
	"strconv"
	"strings"

	"github.com/bidon-io/bidon-backend/internal/ad"
)

type ShowRequest struct {
	BaseRequest
	AdType ad.Type `param:"ad_type"`
	Bid    *Bid    `json:"bid" validate:"required_without=Show"`
	Show   *Bid    `json:"show" validate:"required_without=Bid"`
}

func (r *ShowRequest) Map() map[string]any {
	m := r.BaseRequest.Map()

	m["ad_type"] = r.AdType

	if r.Bid != nil {
		m["bid"] = r.Bid.Map()
	}
	if r.Show != nil {
		m["show"] = r.Show.Map()
	}

	return m
}

func (b *ShowRequest) NormalizeValues() {
	b.BaseRequest.NormalizeValues()

	if b.Bid != nil {
		// Some SDK versions can send lower case bid_type
		b.Bid.BidType = BidType(strings.ToUpper(b.Bid.BidType.String()))
	}
}

func (r *ShowRequest) GetAuctionConfigurationParams() (string, string) {
	return strconv.FormatInt(int64(r.Bid.AuctionConfigurationID), 10), r.Bid.AuctionConfigurationUID
}

func (r *ShowRequest) SetAuctionConfigurationParams(id int64, uid string) {
	r.Bid.AuctionConfigurationID = id
	r.Bid.AuctionConfigurationUID = uid
}
