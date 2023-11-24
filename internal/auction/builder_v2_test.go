package auction_test

import (
	"context"
	"errors"
	"testing"

	"github.com/bidon-io/bidon-backend/internal/ad"
	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/auction"
	auctionmocks "github.com/bidon-io/bidon-backend/internal/auction/mocks"
	"github.com/google/go-cmp/cmp"
)

func TestBuilderV2_Build(t *testing.T) {
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
	pf := 0.1
	adUnits := []auction.AdUnit{
		{DemandID: "test", PriceFloor: &pf, Label: "test", Extra: map[string]any{"placement_id": "test"}},
	}

	configFetcher := &auctionmocks.ConfigFetcherMock{
		MatchFunc: func(ctx context.Context, appID int64, adType ad.Type, segmentID int64) (*auction.Config, error) {
			return config, nil
		},
		FetchByUIDCachedFunc: func(ctx context.Context, appId int64, key string, aucUID string) *auction.Config {
			if aucUID == "1111111111111111111" {
				return config
			} else {
				return nil
			}
		},
	}
	adUnitsMatcher := &auctionmocks.AdUnitsMatcherMock{
		MatchFunc: func(ctx context.Context, params *auction.BuildParams) ([]auction.AdUnit, error) {
			return adUnits, nil
		},
	}
	builder := &auction.BuilderV2{
		ConfigFetcher:  configFetcher,
		AdUnitsMatcher: adUnitsMatcher,
	}

	testCases := []struct {
		name    string
		params  *auction.BuildParams
		want    *auction.Auction
		wantErr bool
		err     error
	}{
		{
			name:   "One round empty",
			params: &auction.BuildParams{Adapters: []adapter.Key{adapter.UnityAdsKey, adapter.BidmachineKey}},
			want: &auction.Auction{
				ConfigID:  config.ID,
				AdUnits:   adUnits,
				LineItems: []auction.LineItem{},
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
				AdUnits:   adUnits,
				LineItems: []auction.LineItem{},
				Rounds: []auction.RoundConfig{
					{ID: "ROUND_1", Demands: []adapter.Key{adapter.ApplovinKey}, Bidding: []adapter.Key{}, Timeout: 15000},
					{ID: "ROUND_3", Demands: []adapter.Key{adapter.ApplovinKey}, Bidding: []adapter.Key{}, Timeout: 15000},
					{ID: "ROUND_4", Demands: []adapter.Key{adapter.ApplovinKey}, Bidding: []adapter.Key{}, Timeout: 15000},
				},
			},
			wantErr: false,
		},
		{
			name:    "No Ads Found",
			params:  &auction.BuildParams{Adapters: []adapter.Key{}},
			want:    nil,
			wantErr: true,
			err:     auction.ErrNoAdsFound,
		},
		{
			name:   "Has auction for AuctionKey",
			params: &auction.BuildParams{AuctionKey: "GEYTCMJRGEYTCMJRGEYTCMJRGEYTCMI=", Adapters: []adapter.Key{adapter.ApplovinKey}},
			want: &auction.Auction{
				ConfigID:  config.ID,
				AdUnits:   adUnits,
				LineItems: []auction.LineItem{},
				Rounds: []auction.RoundConfig{
					{ID: "ROUND_1", Demands: []adapter.Key{adapter.ApplovinKey}, Bidding: []adapter.Key{}, Timeout: 15000},
					{ID: "ROUND_3", Demands: []adapter.Key{adapter.ApplovinKey}, Bidding: []adapter.Key{}, Timeout: 15000},
					{ID: "ROUND_4", Demands: []adapter.Key{adapter.ApplovinKey}, Bidding: []adapter.Key{}, Timeout: 15000},
				},
			},
			wantErr: false,
		},
		{
			name:    "No auction for AuctionKey",
			params:  &auction.BuildParams{AuctionKey: "GMZTGMZTGMZTGMZTGMZTGMZTGMZTGMY=", Adapters: []adapter.Key{adapter.ApplovinKey}},
			want:    nil,
			wantErr: true,
			err:     auction.InvalidAuctionKey,
		},
	}

	for _, tC := range testCases {
		got, err := builder.Build(context.Background(), tC.params)

		if tC.wantErr {
			if !errors.Is(err, tC.err) {
				t.Errorf("Expected error %v, got: %v", tC.err, err)
			}
		} else {
			if err != nil {
				t.Errorf("Error Build: %v", err)
			}

			if diff := cmp.Diff(tC.want, got); diff != "" {
				t.Errorf("builder.Build -> %+v mismatch \n(-want, +got)\n%s", tC.name, diff)
			}
		}
	}
}
