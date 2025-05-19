package auctionv2

import (
	"testing"

	"github.com/bidon-io/bidon-backend/internal/auction"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
)

func TestPriceFloor(t *testing.T) {
	tests := []struct {
		name          string
		req           *schema.AuctionV2Request
		auctionConfig *auction.Config
		expected      float64
	}{
		{
			name: "Custom adapter (max) with previous auction price",
			req: &schema.AuctionV2Request{
				AdObject: schema.AdObjectV2{
					PriceFloor: 0.01,
				},
				BaseRequest: schema.BaseRequest{
					Ext: `{"mediator":"max","previous_auction_price":0.25}`,
				},
				AdCache: []schema.AdCacheObject{
					{Price: 0.02},
				},
			},
			auctionConfig: &auction.Config{
				PriceFloor: 0.05,
			},
			expected: 0.25, // Should use previous auction price since it's higher than calculated floor
		},
		{
			name: "Custom adapter (max) with previous auction price lower than calculated floor",
			req: &schema.AuctionV2Request{
				AdObject: schema.AdObjectV2{
					PriceFloor: 0.01,
				},
				BaseRequest: schema.BaseRequest{
					Ext: `{"mediator":"max","previous_auction_price":0.03}`,
				},
				AdCache: []schema.AdCacheObject{
					{Price: 0.02},
				},
			},
			auctionConfig: &auction.Config{
				PriceFloor: 0.05,
			},
			expected: 0.05, // Should use calculated floor since it's higher than previous auction price
		},
		{
			name: "Custom adapter (max) without previous auction price",
			req: &schema.AuctionV2Request{
				AdObject: schema.AdObjectV2{
					PriceFloor: 0.01,
				},
				BaseRequest: schema.BaseRequest{
					Ext: `{"mediator":"max"}`,
				},
				AdCache: []schema.AdCacheObject{
					{Price: 0.02},
				},
			},
			auctionConfig: &auction.Config{
				PriceFloor: 0.05,
			},
			expected: 0.05, // Should use max of request, cache, and config
		},
		{
			name: "Non-custom adapter with previous auction price",
			req: &schema.AuctionV2Request{
				AdObject: schema.AdObjectV2{
					PriceFloor: 0.01,
				},
				BaseRequest: schema.BaseRequest{
					Ext: `{""mediator":"regular","previous_auction_price":0.25}`,
				},
				AdCache: []schema.AdCacheObject{
					{Price: 0.02},
				},
			},
			auctionConfig: &auction.Config{
				PriceFloor: 0.05,
			},
			expected: 0.05, // Should use max of request, cache, and config (ignore previous auction price)
		},
		{
			name: "Multiple ad cache objects",
			req: &schema.AuctionV2Request{
				AdObject: schema.AdObjectV2{
					PriceFloor: 0.01,
				},
				AdCache: []schema.AdCacheObject{
					{Price: 0.02},
					{Price: 0.07},
					{Price: 0.03},
				},
			},
			auctionConfig: &auction.Config{
				PriceFloor: 0.05,
			},
			expected: 0.07, // Should use max of request, highest cache, and config
		},
		{
			name: "Request price floor higher than config",
			req: &schema.AuctionV2Request{
				AdObject: schema.AdObjectV2{
					PriceFloor: 0.1,
				},
				AdCache: []schema.AdCacheObject{
					{Price: 0.02},
				},
			},
			auctionConfig: &auction.Config{
				PriceFloor: 0.05,
			},
			expected: 0.1, // Should use request price floor
		},
		{
			name: "Empty ad cache",
			req: &schema.AuctionV2Request{
				AdObject: schema.AdObjectV2{
					PriceFloor: 0.01,
				},
				AdCache: []schema.AdCacheObject{},
			},
			auctionConfig: &auction.Config{
				PriceFloor: 0.05,
			},
			expected: 0.05, // Should use config price floor
		},
		{
			name: "Custom adapter (level_play) with previous auction price",
			req: &schema.AuctionV2Request{
				AdObject: schema.AdObjectV2{
					PriceFloor: 0.01,
				},
				BaseRequest: schema.BaseRequest{
					Ext: `{"mediator":"level_play","previous_auction_price":0.25}`,
				},
				AdCache: []schema.AdCacheObject{
					{Price: 0.02},
				},
			},
			auctionConfig: &auction.Config{
				PriceFloor: 0.05,
			},
			expected: 0.25, // Should use previous auction price
		},
		{
			name: "Custom adapter (level_play) with previous auction price of zero",
			req: &schema.AuctionV2Request{
				AdObject: schema.AdObjectV2{
					PriceFloor: 0.01,
				},
				BaseRequest: schema.BaseRequest{
					Ext: `{"mediator":"level_play","previous_auction_price":0.0}`,
				},
				AdCache: []schema.AdCacheObject{
					{Price: 0.02},
				},
			},
			auctionConfig: &auction.Config{
				PriceFloor: 0.05,
			},
			expected: 0.05, // Should use calculated floor since previous auction price is lower
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Need to parse the Ext field to populate extData
			tt.req.NormalizeValues()

			result := priceFloor(tt.req, tt.auctionConfig)
			if result != tt.expected {
				t.Errorf("priceFloor() = %v, want %v", result, tt.expected)
			}
		})
	}
}
