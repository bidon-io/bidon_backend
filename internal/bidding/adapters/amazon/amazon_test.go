package amazon

import (
	"errors"
	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/bidding/adapters"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
	"github.com/google/go-cmp/cmp"
	"testing"
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
										"slotUuid": "slot_uuid_1",
										"pricePoint": "price_point_1"
									},
									{
										"slotUuid": "slot_uuid_2",
										"pricePoint": "price_point_2"
									},
									{
										"slotUuid": "slot_uuid_3",
										"pricePoint": "price_point_3"
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
						Price: 1.0,
					},
				},
				{
					DemandID: adapter.AmazonKey,
					SlotUUID: "slot_uuid_2",
					Bid: &adapters.BidDemandResponse{
						Price: 2.0,
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
		for _, dr := range got {
			if dr.Bid != nil {
				dr.Bid.ImpID = "" // ignore imp_id
			}
		}
		if diff := cmp.Diff(tt.want, got, cmp.Comparer(compareErrors)); diff != "" {
			t.Errorf("%q. Adapter.FetchBids() mismatch (-want +got):\n%s", tt.name, diff)
		}
	}
}
