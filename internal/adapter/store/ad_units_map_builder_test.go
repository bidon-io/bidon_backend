package store_test

import (
	"context"
	"github.com/bidon-io/bidon-backend/internal/auction"
	"github.com/bidon-io/bidon-backend/internal/auction/mocks"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
	"testing"

	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/adapter/store"
	"github.com/google/go-cmp/cmp"
)

func TestAdUnitsMapBuilder_Build(t *testing.T) {
	adUnitsMatcher := mocks.AdUnitsMatcherMock{
		MatchFunc: func(ctx context.Context, params *auction.BuildParams) ([]auction.AdUnit, error) {
			return []auction.AdUnit{
				{DemandID: "applovin"},
				{DemandID: "bidmachine"},
				{DemandID: "amazon"},
			}, nil
		},
	}

	testCases := []struct {
		name        string
		adapterKeys []adapter.Key
		want        store.AdUnitsMap
	}{
		{
			name:        "",
			adapterKeys: adapter.Keys,
			want: store.AdUnitsMap{
				adapter.ApplovinKey: {
					{DemandID: "applovin"},
				},
				adapter.BidmachineKey: {
					{DemandID: "bidmachine"},
				},
				adapter.AmazonKey: {
					{DemandID: "amazon"},
				},
			},
		},
	}

	builder := store.AdUnitsMapBuilder{
		AdUnitsMatcher: &adUnitsMatcher,
	}

	for _, tC := range testCases {
		got, err := builder.Build(context.Background(), 1, tC.adapterKeys, schema.Imp{})
		if err != nil {
			t.Fatalf("failed to fetch app demand profiles: %v", err)
		}

		if diff := cmp.Diff(tC.want, got); diff != "" {
			t.Errorf("fetcher.Fetch -> %v mismatch (-want +got):\n%s", tC.name, diff)
		}
	}
}
