package adapters

import (
	"github.com/bidon-io/bidon-backend/internal/bidding/openrtb"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
)

func CalculatePriceFloor(rtbRequest *openrtb.BidRequest, incomingRequest *schema.AuctionRequest) float64 {
	if rtbRequest == nil || incomingRequest == nil {
		return 0
	}

	if len(rtbRequest.Imp) == 1 {
		return rtbRequest.Imp[0].BidFloor
	} else {
		return incomingRequest.AdObject.GetBidFloorForBidding()
	}
}

var FullscreenFormats = map[string][2]int64{
	"PHONE":  {320, 480},
	"TABLET": {768, 1024},
}
