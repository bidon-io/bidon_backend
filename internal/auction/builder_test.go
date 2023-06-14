package auction_test

import (
	"context"
	"testing"

	"github.com/bidon-io/bidon-backend/internal/ad"
	"github.com/bidon-io/bidon-backend/internal/auction"
	"github.com/google/go-cmp/cmp"
)

func TestBuilder_Build(t *testing.T) {
	config := &auction.Config{
		ID: 1,
		Rounds: []auction.RoundConfig{
			{ID: "ROUND_1", Demands: []string{"applovin", "bidmachine"}, Timeout: 15000},
			{ID: "ROUND_2", Demands: []string{"unityads", "bidmachine"}, Timeout: 15000},
			{ID: "ROUND_3", Demands: []string{"applovin"}, Timeout: 15000},
			{ID: "ROUND_4", Demands: []string{"unityads", "applovin"}, Timeout: 15000},
		},
	}
	lineItems := []auction.LineItem{
		{ID: "test", PriceFloor: 0.1, AdUnitID: "test_id"},
	}

	configFetcher := &auction.ConfigMatcherMock{
		MatchFunc: func(ctx context.Context, appID int64, adType ad.Type) (*auction.Config, error) {
			return config, nil
		},
	}
	lineItemsMatcher := &auction.LineItemsMatcherMock{
		MatchFunc: func(ctx context.Context, params *auction.BuildParams) ([]auction.LineItem, error) {
			return lineItems, nil
		},
	}
	builder := &auction.Builder{
		ConfigMatcher:    configFetcher,
		LineItemsMatcher: lineItemsMatcher,
	}

	testCases := []struct {
		params *auction.BuildParams
		want   *auction.Auction
	}{
		{
			params: &auction.BuildParams{Adapters: []string{"unityads", "bidmachine"}},
			want: &auction.Auction{
				ConfigID:  config.ID,
				LineItems: lineItems,
				Rounds: []auction.RoundConfig{
					{ID: "ROUND_1", Demands: []string{"bidmachine"}, Timeout: 15000},
					{ID: "ROUND_2", Demands: []string{"unityads", "bidmachine"}, Timeout: 15000},
					{ID: "ROUND_4", Demands: []string{"unityads"}, Timeout: 15000},
				},
			},
		},
		{
			params: &auction.BuildParams{Adapters: []string{"applovin"}},
			want: &auction.Auction{
				ConfigID:  config.ID,
				LineItems: lineItems,
				Rounds: []auction.RoundConfig{
					{ID: "ROUND_1", Demands: []string{"applovin"}, Timeout: 15000},
					{ID: "ROUND_3", Demands: []string{"applovin"}, Timeout: 15000},
					{ID: "ROUND_4", Demands: []string{"applovin"}, Timeout: 15000},
				},
			},
		},
		{
			params: &auction.BuildParams{Adapters: []string{}},
			want: &auction.Auction{
				ConfigID:  config.ID,
				LineItems: lineItems,
				Rounds:    []auction.RoundConfig{},
			},
		},
	}

	for _, tC := range testCases {
		got, _ := builder.Build(context.Background(), tC.params)

		if diff := cmp.Diff(tC.want, got); diff != "" {
			t.Errorf("builder.Build(ctx, %+v) mismatch (-want, +got)\n%s", tC.params, diff)
		}
	}
}
