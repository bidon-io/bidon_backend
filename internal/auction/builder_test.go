package auction_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bidon-io/bidon-backend/internal/ad"
	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/auction"
	auctionmocks "github.com/bidon-io/bidon-backend/internal/auction/mocks"
)

func TestBuilder_Build(t *testing.T) {
	config := &auction.Config{
		ID: 1,
		Rounds: []auction.RoundConfig{
			{
				ID:      "ROUND_1",
				Demands: []adapter.Key{adapter.ApplovinKey, adapter.BidmachineKey},
				Timeout: 15000,
			},
			{
				ID:      "ROUND_2",
				Demands: []adapter.Key{adapter.UnityAdsKey},
				Bidding: []adapter.Key{adapter.BidmachineKey},
				Timeout: 15000,
			},
			{
				ID:      "ROUND_3",
				Demands: []adapter.Key{adapter.ApplovinKey},
				Timeout: 15000,
			},
			{
				ID:      "ROUND_4",
				Demands: []adapter.Key{adapter.UnityAdsKey, adapter.ApplovinKey},
				Timeout: 15000,
			},
			{
				ID:      "ROUND_5",
				Bidding: []adapter.Key{adapter.BidmachineKey},
				Timeout: 15000,
			},
		},
	}
	lineItems := []auction.LineItem{
		{ID: "test", PriceFloor: 0.1, AdUnitID: "test_id"},
	}

	configFetcher := &auctionmocks.ConfigFetcherMock{
		MatchFunc: func(_ context.Context, _ int64, _ ad.Type, _ int64, _ string) (*auction.Config, error) {
			return config, nil
		},
	}
	lineItemsMatcher := &auctionmocks.LineItemsMatcherMock{
		MatchCachedFunc: func(_ context.Context, _ *auction.BuildParams) ([]auction.LineItem, error) {
			return lineItems, nil
		},
	}
	builder := &auction.Builder{
		ConfigFetcher:    configFetcher,
		LineItemsMatcher: lineItemsMatcher,
	}

	testCases := []struct {
		name   string
		params *auction.BuildParams
		want   *auction.Auction
	}{
		{
			name:   "One round empty",
			params: &auction.BuildParams{Adapters: []adapter.Key{adapter.UnityAdsKey, adapter.BidmachineKey}},
			want: &auction.Auction{
				ConfigID:  config.ID,
				LineItems: lineItems,
				AdUnits:   []auction.AdUnit{},
				Rounds: []auction.RoundConfig{
					{ID: "ROUND_1", Demands: []adapter.Key{adapter.BidmachineKey}, Bidding: []adapter.Key{}, Timeout: 15000},
					{ID: "ROUND_2", Demands: []adapter.Key{adapter.UnityAdsKey}, Bidding: []adapter.Key{adapter.BidmachineKey}, Timeout: 15000},
					{ID: "ROUND_4", Demands: []adapter.Key{adapter.UnityAdsKey}, Bidding: []adapter.Key{}, Timeout: 15000},
					{ID: "ROUND_5", Demands: []adapter.Key{}, Bidding: []adapter.Key{adapter.BidmachineKey}, Timeout: 15000},
				},
			},
		},
		{
			name:   "Single adapter available",
			params: &auction.BuildParams{Adapters: []adapter.Key{adapter.ApplovinKey}},
			want: &auction.Auction{
				ConfigID:  config.ID,
				LineItems: lineItems,
				AdUnits:   []auction.AdUnit{},
				Rounds: []auction.RoundConfig{
					{ID: "ROUND_1", Demands: []adapter.Key{adapter.ApplovinKey}, Bidding: []adapter.Key{}, Timeout: 15000},
					{ID: "ROUND_3", Demands: []adapter.Key{adapter.ApplovinKey}, Bidding: []adapter.Key{}, Timeout: 15000},
					{ID: "ROUND_4", Demands: []adapter.Key{adapter.ApplovinKey}, Bidding: []adapter.Key{}, Timeout: 15000},
				},
			},
		},
		{
			name:   "Empty Response",
			params: &auction.BuildParams{Adapters: []adapter.Key{}},
			want:   nil,
		},
	}

	for _, tC := range testCases {
		got, _ := builder.Build(context.Background(), tC.params)

		if diff := cmp.Diff(tC.want, got); diff != "" {
			t.Errorf("builder.Build -> %+v mismatch \n(-want, +got)\n%s", tC.name, diff)
		}
	}
}
