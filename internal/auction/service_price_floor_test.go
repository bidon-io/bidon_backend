package auction

import (
	"testing"

	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
)

func TestPriceFloor(t *testing.T) {
	tests := []struct {
		name          string
		req           *schema.AuctionRequest
		auctionConfig *Config
		expected      float64
	}{
		{
			name: "Custom adapter (max) with previous auction price",
			req: &schema.AuctionRequest{
				AdObject: schema.AdObject{
					PriceFloor: 0.01,
				},
				BaseRequest: schema.BaseRequest{
					Ext: `{"mediator":"max","previous_auction_price":0.25}`,
				},
				AdCache: []schema.AdCacheObject{
					{Price: 0.02},
				},
			},
			auctionConfig: &Config{
				PriceFloor: 0.05,
			},
			expected: 0.25, // Should use previous auction price since it's higher than calculated floor
		},
		{
			name: "Custom adapter (max) with previous auction price lower than calculated floor",
			req: &schema.AuctionRequest{
				AdObject: schema.AdObject{
					PriceFloor: 0.01,
				},
				BaseRequest: schema.BaseRequest{
					Ext: `{"mediator":"max","previous_auction_price":0.03}`,
				},
				AdCache: []schema.AdCacheObject{
					{Price: 0.02},
				},
			},
			auctionConfig: &Config{
				PriceFloor: 0.05,
			},
			expected: 0.05, // Should use calculated floor since it's higher than previous auction price
		},
		{
			name: "Custom adapter (max) without previous auction price",
			req: &schema.AuctionRequest{
				AdObject: schema.AdObject{
					PriceFloor: 0.01,
				},
				BaseRequest: schema.BaseRequest{
					Ext: `{"mediator":"max"}`,
				},
				AdCache: []schema.AdCacheObject{
					{Price: 0.02},
				},
			},
			auctionConfig: &Config{
				PriceFloor: 0.05,
			},
			expected: 0.05, // Should use max of request, cache, and config
		},
		{
			name: "Non-custom adapter with previous auction price",
			req: &schema.AuctionRequest{
				AdObject: schema.AdObject{
					PriceFloor: 0.01,
				},
				BaseRequest: schema.BaseRequest{
					Ext: `{""mediator":"regular","previous_auction_price":0.25}`,
				},
				AdCache: []schema.AdCacheObject{
					{Price: 0.02},
				},
			},
			auctionConfig: &Config{
				PriceFloor: 0.05,
			},
			expected: 0.05, // Should use max of request, cache, and config (ignore previous auction price)
		},
		{
			name: "Multiple ad cache objects",
			req: &schema.AuctionRequest{
				AdObject: schema.AdObject{
					PriceFloor: 0.01,
				},
				AdCache: []schema.AdCacheObject{
					{Price: 0.02},
					{Price: 0.07},
					{Price: 0.03},
				},
			},
			auctionConfig: &Config{
				PriceFloor: 0.05,
			},
			expected: 0.07, // Should use max of request, highest cache, and config
		},
		{
			name: "Request price floor higher than config",
			req: &schema.AuctionRequest{
				AdObject: schema.AdObject{
					PriceFloor: 0.1,
				},
				AdCache: []schema.AdCacheObject{
					{Price: 0.02},
				},
			},
			auctionConfig: &Config{
				PriceFloor: 0.05,
			},
			expected: 0.1, // Should use request price floor
		},
		{
			name: "Empty ad cache",
			req: &schema.AuctionRequest{
				AdObject: schema.AdObject{
					PriceFloor: 0.01,
				},
				AdCache: []schema.AdCacheObject{},
			},
			auctionConfig: &Config{
				PriceFloor: 0.05,
			},
			expected: 0.05, // Should use config price floor
		},
		{
			name: "Custom adapter (level_play) with previous auction price",
			req: &schema.AuctionRequest{
				AdObject: schema.AdObject{
					PriceFloor: 0.01,
				},
				BaseRequest: schema.BaseRequest{
					Ext: `{"mediator":"level_play","previous_auction_price":0.25}`,
				},
				AdCache: []schema.AdCacheObject{
					{Price: 0.02},
				},
			},
			auctionConfig: &Config{
				PriceFloor: 0.05,
			},
			expected: 0.25, // Should use previous auction price
		},
		{
			name: "Custom adapter (level_play) with previous auction price of zero",
			req: &schema.AuctionRequest{
				AdObject: schema.AdObject{
					PriceFloor: 0.01,
				},
				BaseRequest: schema.BaseRequest{
					Ext: `{"mediator":"level_play","previous_auction_price":0.0}`,
				},
				AdCache: []schema.AdCacheObject{
					{Price: 0.02},
				},
			},
			auctionConfig: &Config{
				PriceFloor: 0.05,
			},
			expected: 0.05, // Should use calculated floor since previous auction price is lower
		},
		{
			name: "Disabled floor - Custom adapter (max) with disabled auction key",
			req: &schema.AuctionRequest{
				AdObject: schema.AdObject{
					AuctionKey: "1LOQ1LROG0000", // Inter - should disable floor
					PriceFloor: 0.01,
				},
				BaseRequest: schema.BaseRequest{
					Ext: `{"mediator":"max"}`,
				},
				AdCache: []schema.AdCacheObject{
					{Price: 0.02},
				},
			},
			auctionConfig: &Config{
				PriceFloor: 0.05,
			},
			expected: 0.0, // Should return 0 (disabled floor) for specified auction key with custom adapter
		},
		{
			name: "Disabled floor - Custom adapter (level_play) with disabled auction key",
			req: &schema.AuctionRequest{
				AdObject: schema.AdObject{
					AuctionKey: "1LOQ2BFG00000", // Banner - should disable floor
					PriceFloor: 0.01,
				},
				BaseRequest: schema.BaseRequest{
					Ext: `{"mediator":"level_play"}`,
				},
				AdCache: []schema.AdCacheObject{
					{Price: 0.02},
				},
			},
			auctionConfig: &Config{
				PriceFloor: 0.05,
			},
			expected: 0.0, // Should return 0 (disabled floor) for specified auction key with custom adapter
		},
		{
			name: "Disabled floor - Custom adapter (max) with rewarded auction key",
			req: &schema.AuctionRequest{
				AdObject: schema.AdObject{
					AuctionKey: "1LOQ2KLES0400", // Rewarded - should disable floor
					PriceFloor: 0.01,
				},
				BaseRequest: schema.BaseRequest{
					Ext: `{"mediator":"max"}`,
				},
				AdCache: []schema.AdCacheObject{
					{Price: 0.02},
				},
			},
			auctionConfig: &Config{
				PriceFloor: 0.05,
			},
			expected: 0.0, // Should return 0 (disabled floor) for specified auction key with custom adapter
		},
		{
			name: "No disabled floor - Non-custom adapter with disabled auction key",
			req: &schema.AuctionRequest{
				AdObject: schema.AdObject{
					AuctionKey: "1LOQ1LROG0000", // Inter key but non-custom adapter
					PriceFloor: 0.01,
				},
				BaseRequest: schema.BaseRequest{
					Ext: `{"mediator":"regular"}`,
				},
				AdCache: []schema.AdCacheObject{
					{Price: 0.02},
				},
			},
			auctionConfig: &Config{
				PriceFloor: 0.05,
			},
			expected: 0.05, // Should use normal floor logic since adapter is not custom
		},
		{
			name: "No disabled floor - Custom adapter with non-disabled auction key",
			req: &schema.AuctionRequest{
				AdObject: schema.AdObject{
					AuctionKey: "SOMEOTHERKEY", // Different auction key
					PriceFloor: 0.01,
				},
				BaseRequest: schema.BaseRequest{
					Ext: `{"mediator":"max"}`,
				},
				AdCache: []schema.AdCacheObject{
					{Price: 0.02},
				},
			},
			auctionConfig: &Config{
				PriceFloor: 0.05,
			},
			expected: 0.05, // Should use normal floor logic since auction key is not in disabled list
		},
		{
			name: "No disabled floor - Custom adapter with empty auction key",
			req: &schema.AuctionRequest{
				AdObject: schema.AdObject{
					AuctionKey: "", // Empty auction key
					PriceFloor: 0.01,
				},
				BaseRequest: schema.BaseRequest{
					Ext: `{"mediator":"max"}`,
				},
				AdCache: []schema.AdCacheObject{
					{Price: 0.02},
				},
			},
			auctionConfig: &Config{
				PriceFloor: 0.05,
			},
			expected: 0.05, // Should use normal floor logic since auction key is empty
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
