package amazon

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/bidding/adapters"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
)

func compareErrors(x, y error) bool {
	if x == nil || y == nil {
		return x == nil && y == nil
	}
	return x.Error() == y.Error()
}

func TestAdapter_FetchBids(t *testing.T) {
	type fields struct {
		PricePointsMap PricePointsMap
	}
	type args struct {
		br *schema.BiddingRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*adapters.DemandResponse
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				PricePointsMap: PricePointsMap{
					"price_point_1": {
						Price:      1.0,
						PricePoint: "price_point_1",
					},
					"price_point_2": {
						Price:      2.0,
						PricePoint: "price_point_2",
					},
				},
			},
			args: args{
				br: &schema.BiddingRequest{
					Imp: schema.Imp{
						Demands: map[adapter.Key]map[string]interface{}{
							adapter.AmazonKey: {
								"token": `[
									{
										"slot_uuid": "slot_uuid_1",
										"price_point": "price_point_1"
									},
									{
										"slot_uuid": "slot_uuid_2",
										"price_point": "price_point_2"
									},
									{
										"slot_uuid": "slot_uuid_3",
										"price_point": "price_point_3"
									}
								]`,
							},
						},
					},
				},
			},
			want: []*adapters.DemandResponse{
				{
					DemandID: adapter.AmazonKey,
					SlotUUID: "slot_uuid_1",
					Bid: &adapters.BidDemandResponse{
						DemandID: adapter.AmazonKey,
						Price:    1.0,
					},
				},
				{
					DemandID: adapter.AmazonKey,
					SlotUUID: "slot_uuid_2",
					Bid: &adapters.BidDemandResponse{
						DemandID: adapter.AmazonKey,
						Price:    2.0,
					},
				},
				{
					DemandID: adapter.AmazonKey,
					SlotUUID: "slot_uuid_3",
					Error:    errors.New("cannot find price point"),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		adapter := Adapter{
			PricePointsMap: tt.fields.PricePointsMap,
		}
		got, err := adapter.FetchBids(tt.args.br)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. Adapter.FetchBids() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}

		// ignore id and imp_id
		for _, dr := range got {
			if dr.Bid != nil {
				dr.Bid.ID = ""
				dr.Bid.ImpID = ""
			}
		}
		if diff := cmp.Diff(tt.want, got, cmp.Comparer(compareErrors)); diff != "" {
			t.Errorf("%q. Adapter.FetchBids() mismatch (-want +got):\n%s", tt.name, diff)
		}
	}
}

func TestBuilder(t *testing.T) {
	tests := []struct {
		name        string
		config      adapter.ProcessedConfigsMap
		expected    *Adapter
		expectError bool
	}{
		{
			name: "Valid Configuration",
			config: adapter.ProcessedConfigsMap{
				adapter.AmazonKey: map[string]interface{}{
					"price_points_map": map[string]interface{}{
						"00n9g200_zzz": map[string]interface{}{
							"name":        "Interstitial",
							"price":       0.5,
							"price_point": "00n9g200_zzz",
						},
						"02bz0000_xxx": map[string]interface{}{
							"name":        "banner",
							"price":       0.7,
							"price_point": "02bz0000_xxx",
						},
					},
				},
			},
			expected: &Adapter{
				PricePointsMap: PricePointsMap{
					"00n9g200_zzz": {
						Price:      0.5,
						PricePoint: "00n9g200_zzz",
					},
					"02bz0000_xxx": {
						Price:      0.7,
						PricePoint: "02bz0000_xxx",
					},
				},
			},
			expectError: false,
		},
		{
			name: "Invalid Configuration",
			config: adapter.ProcessedConfigsMap{
				adapter.AmazonKey: map[string]interface{}{
					"price_points_map": "invalid", // Invalid value, should be a map.
				},
			},
			expected:    nil,
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			adpt, err := Builder(test.config)

			if test.expectError {
				if err == nil {
					t.Fatal("Expected an error, but got no error")
				}
			} else {
				if err != nil {
					t.Fatalf("Expected no error, but got an error: %v", err)
				}

				if adpt == nil {
					t.Fatal("Expected non-nil Adapter, but got nil")
				}

				if diff := cmp.Diff(test.expected.PricePointsMap, adpt.PricePointsMap); diff != "" {
					t.Fatalf("PricePointsMap mismatch. Got %v, expected %v", adpt.PricePointsMap, test.expected.PricePointsMap)
				}
			}
		})
	}
}
