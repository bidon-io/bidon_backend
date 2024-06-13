package auctionv2_test

import (
	"context"
	"errors"
	"github.com/google/go-cmp/cmp"
	"testing"

	"github.com/bidon-io/bidon-backend/internal/ad"
	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/auction"
	"github.com/bidon-io/bidon-backend/internal/auctionv2"
	"github.com/bidon-io/bidon-backend/internal/auctionv2/mocks"
	"github.com/bidon-io/bidon-backend/internal/bidding"
	"github.com/bidon-io/bidon-backend/internal/bidding/adapters"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
)

type BuilderMocks struct {
	ConfigFetcher                *mocks.ConfigFetcherMock
	AdUnitsMatcher               *mocks.AdUnitsMatcherMock
	BiddingBuilder               *mocks.BiddingBuilderMock
	BiddingAdaptersConfigBuilder *mocks.BiddingAdaptersConfigBuilderMock
}

type BuilderOption func(*BuilderMocks)

func WithConfigFetcher(cf *mocks.ConfigFetcherMock) BuilderOption {
	return func(b *BuilderMocks) {
		b.ConfigFetcher = cf
	}
}

func WithAdUnitsMatcher(au *mocks.AdUnitsMatcherMock) BuilderOption {
	return func(b *BuilderMocks) {
		b.AdUnitsMatcher = au
	}
}

func WithBiddingBuilder(bb *mocks.BiddingBuilderMock) BuilderOption {
	return func(b *BuilderMocks) {
		b.BiddingBuilder = bb
	}
}

func testHelperDefaultAuctionBuilderMocks() *BuilderMocks {
	auctionConfig := &auction.Config{
		ID:      1,
		Demands: []adapter.Key{adapter.GAMKey, adapter.DTExchangeKey},
		Bidding: []adapter.Key{adapter.BidmachineKey, adapter.AmazonKey, adapter.MetaKey},
		Timeout: 15000,
	}
	adUnits := []auction.AdUnit{
		{
			DemandID:   "gam",
			Label:      "gam",
			PriceFloor: ptr(0.1),
			UID:        "123_gam",
			BidType:    schema.CPMBidType,
			Extra: map[string]any{
				"placement_id": "123",
			},
		},
		{
			DemandID:   "dtexchange",
			Label:      "dtexchange",
			PriceFloor: ptr(0.01),
			UID:        "123_dtexchange",
			BidType:    schema.CPMBidType,
			Extra: map[string]any{
				"placement_id": "123",
			},
		},
	}
	configFetcher := &mocks.ConfigFetcherMock{
		MatchFunc: func(ctx context.Context, appID int64, adType ad.Type, segmentID int64, version string) (*auction.Config, error) {
			return auctionConfig, nil
		},
		FetchByUIDCachedFunc: func(ctx context.Context, appId int64, key string, aucUID string) *auction.Config {
			if aucUID == "1688565055735595008" {
				return auctionConfig
			} else {
				return nil
			}
		},
	}
	adUnitsMatcher := &mocks.AdUnitsMatcherMock{
		MatchCachedFunc: func(ctx context.Context, params *auction.BuildParams) ([]auction.AdUnit, error) {
			return adUnits, nil
		},
	}
	biddingAdaptersConfigBuilder := &mocks.BiddingAdaptersConfigBuilderMock{
		BuildFunc: func(ctx context.Context, appID int64, adapterKeys []adapter.Key, imp schema.Imp, adUnitsMap *map[adapter.Key][]auction.AdUnit) (adapter.ProcessedConfigsMap, error) {
			return adapter.ProcessedConfigsMap{}, nil
		},
	}
	biddingBuilder := &mocks.BiddingBuilderMock{
		HoldAuctionFunc: func(ctx context.Context, params *bidding.BuildParams) (bidding.AuctionResult, error) {
			return bidding.AuctionResult{
				RoundNumber: 0,
				Bids:        []adapters.DemandResponse{},
			}, nil
		},
	}

	return &BuilderMocks{
		ConfigFetcher:                configFetcher,
		AdUnitsMatcher:               adUnitsMatcher,
		BiddingBuilder:               biddingBuilder,
		BiddingAdaptersConfigBuilder: biddingAdaptersConfigBuilder,
	}
}

func ptr[T any](t T) *T {
	return &t
}

func testHelperAuctionBuilder(opts ...BuilderOption) *auctionv2.Builder {
	m := testHelperDefaultAuctionBuilderMocks()

	for _, opt := range opts {
		opt(m)
	}

	return &auctionv2.Builder{
		ConfigFetcher:                m.ConfigFetcher,
		AdUnitsMatcher:               m.AdUnitsMatcher,
		BiddingBuilder:               m.BiddingBuilder,
		BiddingAdaptersConfigBuilder: m.BiddingAdaptersConfigBuilder,
	}
}

func TestBuilder_Build2(t *testing.T) {
	request := &schema.AuctionV2Request{}
	testCases := []struct {
		name    string
		builder *auctionv2.Builder
		params  *auctionv2.BuildParams
		want    *auctionv2.AuctionResult
		wantErr bool
		err     error
	}{
		{
			name: "Success",
			builder: testHelperAuctionBuilder(WithBiddingBuilder(&mocks.BiddingBuilderMock{
				HoldAuctionFunc: func(ctx context.Context, params *bidding.BuildParams) (bidding.AuctionResult, error) {
					return bidding.AuctionResult{
						RoundNumber: 0,
						Bids: []adapters.DemandResponse{
							{DemandID: "bidmachine", Bid: &adapters.BidDemandResponse{}},
						},
					}, nil
				},
			})),
			params: &auctionv2.BuildParams{Adapters: []adapter.Key{adapter.GAMKey, adapter.BidmachineKey}, MergedAuctionRequest: request, PriceFloor: 0.1},
			want: &auctionv2.AuctionResult{
				AuctionConfiguration: &auction.Config{
					ID:      1,
					Demands: []adapter.Key{adapter.GAMKey, adapter.DTExchangeKey},
					Bidding: []adapter.Key{adapter.BidmachineKey, adapter.AmazonKey, adapter.MetaKey},
					Timeout: 15000,
				},
				AdUnits: &[]auction.AdUnit{
					{
						DemandID:   "gam",
						UID:        "123_gam",
						Label:      "gam",
						PriceFloor: ptr(0.1),
						BidType:    "CPM",
						Extra:      map[string]any{"placement_id": string("123")},
					},
					{
						DemandID:   "dtexchange",
						Label:      "dtexchange",
						PriceFloor: ptr(0.01),
						UID:        "123_dtexchange",
						BidType:    schema.CPMBidType,
						Extra: map[string]any{
							"placement_id": "123",
						},
					},
				},

				CPMAdUnits: &[]auction.AdUnit{
					{
						DemandID:   "gam",
						UID:        "123_gam",
						Label:      "gam",
						PriceFloor: ptr(0.1),
						BidType:    "CPM",
						Extra:      map[string]any{"placement_id": string("123")},
					},
				},
				BiddingAuctionResult: &bidding.AuctionResult{
					Bids: []adapters.DemandResponse{
						{DemandID: "bidmachine", Bid: &adapters.BidDemandResponse{}},
					},
				},
			},
		},
		{
			name: "No Biding Adapters Matched",
			builder: testHelperAuctionBuilder(WithBiddingBuilder(&mocks.BiddingBuilderMock{
				HoldAuctionFunc: func(ctx context.Context, params *bidding.BuildParams) (bidding.AuctionResult, error) {
					return bidding.AuctionResult{}, bidding.ErrNoAdaptersMatched
				},
			})),
			params: &auctionv2.BuildParams{Adapters: []adapter.Key{adapter.GAMKey, adapter.BidmachineKey}, MergedAuctionRequest: request, PriceFloor: 0.1},
			want: &auctionv2.AuctionResult{
				AuctionConfiguration: &auction.Config{
					ID:      1,
					Demands: []adapter.Key{adapter.GAMKey, adapter.DTExchangeKey},
					Bidding: []adapter.Key{adapter.BidmachineKey, adapter.AmazonKey, adapter.MetaKey},
					Timeout: 15000,
				},
				AdUnits: &[]auction.AdUnit{
					{
						DemandID:   "gam",
						UID:        "123_gam",
						Label:      "gam",
						PriceFloor: ptr(0.1),
						BidType:    "CPM",
						Extra:      map[string]any{"placement_id": string("123")},
					},
					{
						DemandID:   "dtexchange",
						Label:      "dtexchange",
						PriceFloor: ptr(0.01),
						UID:        "123_dtexchange",
						BidType:    schema.CPMBidType,
						Extra: map[string]any{
							"placement_id": "123",
						},
					},
				},
				CPMAdUnits: &[]auction.AdUnit{
					{
						DemandID:   "gam",
						UID:        "123_gam",
						Label:      "gam",
						PriceFloor: ptr(0.1),
						BidType:    "CPM",
						Extra:      map[string]any{"placement_id": string("123")},
					},
				},
				BiddingAuctionResult: &bidding.AuctionResult{},
			},
		},
		{
			name: "No Ads Found",
			builder: testHelperAuctionBuilder(
				WithAdUnitsMatcher(&mocks.AdUnitsMatcherMock{
					MatchCachedFunc: func(ctx context.Context, params *auction.BuildParams) ([]auction.AdUnit, error) {
						return []auction.AdUnit{}, nil
					},
				}),
				WithBiddingBuilder(&mocks.BiddingBuilderMock{
					HoldAuctionFunc: func(ctx context.Context, params *bidding.BuildParams) (bidding.AuctionResult, error) {
						return bidding.AuctionResult{}, nil
					},
				}),
			),
			params:  &auctionv2.BuildParams{Adapters: []adapter.Key{adapter.GAMKey, adapter.BidmachineKey}, MergedAuctionRequest: request, PriceFloor: 0.1},
			wantErr: true,
			err:     auction.ErrNoAdsFound,
		},
		{
			name: "Has auction for AuctionKey",
			builder: testHelperAuctionBuilder(WithBiddingBuilder(&mocks.BiddingBuilderMock{
				HoldAuctionFunc: func(ctx context.Context, params *bidding.BuildParams) (bidding.AuctionResult, error) {
					return bidding.AuctionResult{
						RoundNumber: 0,
						Bids: []adapters.DemandResponse{
							{DemandID: "bidmachine", Bid: &adapters.BidDemandResponse{}},
						},
					}, nil
				},
			})),
			params: &auctionv2.BuildParams{
				AuctionKey:           "1ERNSV33K4000",
				Adapters:             []adapter.Key{adapter.GAMKey, adapter.BidmachineKey},
				MergedAuctionRequest: request,
				PriceFloor:           0.1,
			},
			want: &auctionv2.AuctionResult{
				AuctionConfiguration: &auction.Config{
					ID:      1,
					Demands: []adapter.Key{adapter.GAMKey, adapter.DTExchangeKey},
					Bidding: []adapter.Key{adapter.BidmachineKey, adapter.AmazonKey, adapter.MetaKey},
					Timeout: 15000,
				},
				AdUnits: &[]auction.AdUnit{
					{
						DemandID:   "gam",
						UID:        "123_gam",
						Label:      "gam",
						PriceFloor: ptr(0.1),
						BidType:    "CPM",
						Extra:      map[string]any{"placement_id": string("123")},
					},
					{
						DemandID:   "dtexchange",
						Label:      "dtexchange",
						PriceFloor: ptr(0.01),
						UID:        "123_dtexchange",
						BidType:    schema.CPMBidType,
						Extra: map[string]any{
							"placement_id": "123",
						},
					},
				},
				CPMAdUnits: &[]auction.AdUnit{
					{
						DemandID:   "gam",
						UID:        "123_gam",
						Label:      "gam",
						PriceFloor: ptr(0.1),
						BidType:    "CPM",
						Extra:      map[string]any{"placement_id": string("123")},
					},
				},
				BiddingAuctionResult: &bidding.AuctionResult{
					Bids: []adapters.DemandResponse{
						{DemandID: "bidmachine", Bid: &adapters.BidDemandResponse{}},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "No auction for AuctionKey",
			builder: testHelperAuctionBuilder(WithBiddingBuilder(&mocks.BiddingBuilderMock{
				HoldAuctionFunc: func(ctx context.Context, params *bidding.BuildParams) (bidding.AuctionResult, error) {
					return bidding.AuctionResult{
						RoundNumber: 0,
						Bids: []adapters.DemandResponse{
							{DemandID: "bidmachine", Bid: &adapters.BidDemandResponse{}},
						},
					}, nil
				},
			})),
			params: &auctionv2.BuildParams{
				AuctionKey:           "1F60CVMI00400",
				Adapters:             []adapter.Key{adapter.GAMKey, adapter.BidmachineKey},
				MergedAuctionRequest: request,
				PriceFloor:           0.1,
			},
			wantErr: true,
			err:     auction.InvalidAuctionKey,
		},
	}

	for _, tC := range testCases {
		got, err := tC.builder.Build(context.Background(), tC.params)

		if tC.wantErr {
			if !errors.Is(err, tC.err) {
				t.Errorf("Expected error %v, got: %v", tC.err, err)
			}
		} else {
			if err != nil {
				t.Errorf("Error Build: %v", err)
				return // Skip further checks
			}

			got.Stat = nil // Stat is not deterministic
			if diff := cmp.Diff(tC.want, got); diff != "" {
				t.Errorf("builder.Build -> %+v mismatch \n(-want, +got)\n%s", tC.name, diff)
			}
		}
	}
}
