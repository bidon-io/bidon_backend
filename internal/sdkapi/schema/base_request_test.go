package schema

import (
	"testing"
)

func TestBaseRequest_GetMediationMode(t *testing.T) {
	tests := []struct {
		name     string
		ext      string
		expected string
	}{
		{
			name:     "Empty ext",
			ext:      "",
			expected: "",
		},
		{
			name:     "Empty JSON",
			ext:      "{}",
			expected: "",
		},
		{
			name:     "With mediation_mode",
			ext:      `{"mediation_mode":"max"}`,
			expected: "max",
		},
		{
			name:     "With mediation_mode lp_ca",
			ext:      `{"mediation_mode":"lp_ca"}`,
			expected: "lp_ca",
		},
		{
			name:     "With mediation_mode regular",
			ext:      `{"mediation_mode":"regular"}`,
			expected: "regular",
		},
		{
			name:     "With other fields",
			ext:      `{"mediation_mode":"max","other_field":"value"}`,
			expected: "max",
		},
		{
			name:     "With mediation_mode as non-string",
			ext:      `{"mediation_mode":123}`,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &BaseRequest{
				Ext: tt.ext,
			}
			req.parseExt() // Parse the Ext field to populate extData

			result := req.GetMediationMode()
			if result != tt.expected {
				t.Errorf("GetMediationMode() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestBaseRequest_GetPrevAuctionPrice(t *testing.T) {
	tests := []struct {
		name     string
		ext      string
		expected *float64
	}{
		{
			name:     "Empty ext",
			ext:      "",
			expected: nil,
		},
		{
			name:     "Empty JSON",
			ext:      "{}",
			expected: nil,
		},
		{
			name:     "With previous_auction_price",
			ext:      `{"previous_auction_price":0.25}`,
			expected: ptrFloat64(0.25),
		},
		{
			name:     "With previous_auction_price zero",
			ext:      `{"previous_auction_price":0.0}`,
			expected: ptrFloat64(0.0),
		},
		{
			name:     "With previous_auction_price negative",
			ext:      `{"previous_auction_price":-0.1}`,
			expected: ptrFloat64(-0.1),
		},
		{
			name:     "With other fields",
			ext:      `{"previous_auction_price":0.5,"other_field":"value"}`,
			expected: ptrFloat64(0.5),
		},
		{
			name:     "With previous_auction_price as non-number",
			ext:      `{"previous_auction_price":"0.25"}`,
			expected: nil,
		},
		{
			name:     "With both mediation_mode and previous_auction_price",
			ext:      `{"mediation_mode":"max","previous_auction_price":0.75}`,
			expected: ptrFloat64(0.75),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &BaseRequest{
				Ext: tt.ext,
			}
			req.parseExt() // Parse the Ext field to populate extData

			result := req.GetPrevAuctionPrice()

			// Compare the results
			if (result == nil && tt.expected != nil) || (result != nil && tt.expected == nil) {
				t.Errorf("GetPrevAuctionPrice() = %v, want %v", result, tt.expected)
			} else if result != nil && tt.expected != nil && *result != *tt.expected {
				t.Errorf("GetPrevAuctionPrice() = %v, want %v", *result, *tt.expected)
			}
		})
	}
}

// Helper function to create a pointer to a float64
func ptrFloat64(v float64) *float64 {
	return &v
}
