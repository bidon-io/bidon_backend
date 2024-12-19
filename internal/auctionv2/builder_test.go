package auctionv2_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/auction"
	"github.com/bidon-io/bidon-backend/internal/auctionv2"
	"github.com/bidon-io/bidon-backend/internal/auctionv2/mocks"
	"github.com/bidon-io/bidon-backend/internal/bidding"
	"github.com/bidon-io/bidon-backend/internal/bidding/adapters"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
)

type BuilderMocks struct {
	AdUnitsMatcher               *mocks.AdUnitsMatcherMock
	BiddingBuilder               *mocks.BiddingBuilderMock
	BiddingAdaptersConfigBuilder *mocks.BiddingAdaptersConfigBuilderMock
}

type BuilderOption func(*BuilderMocks)

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
	adUnitsMatcher := &mocks.AdUnitsMatcherMock{
		MatchCachedFunc: func(ctx context.Context, params *auction.BuildParams) ([]auction.AdUnit, error) {
			return adUnits, nil
		},
	}
	biddingAdaptersConfigBuilder := &mocks.BiddingAdaptersConfigBuilderMock{
		BuildFunc: func(ctx context.Context, appID int64, adapterKeys []adapter.Key, adUnitsMap *auction.AdUnitsMap) (adapter.ProcessedConfigsMap, error) {
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
		AdUnitsMatcher:               m.AdUnitsMatcher,
		BiddingBuilder:               m.BiddingBuilder,
		BiddingAdaptersConfigBuilder: m.BiddingAdaptersConfigBuilder,
	}
}

func TestBuilder_Build(t *testing.T) {
	auctionConfig := &auction.Config{
		ID:        1,
		Demands:   []adapter.Key{adapter.GAMKey, adapter.DTExchangeKey},
		Bidding:   []adapter.Key{adapter.BidmachineKey, adapter.AmazonKey, adapter.MetaKey},
		AdUnitIDs: []int64{1, 2},
		Timeout:   15000,
	}
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
			params: &auctionv2.BuildParams{
				Adapters:             []adapter.Key{adapter.GAMKey, adapter.BidmachineKey},
				MergedAuctionRequest: request,
				PriceFloor:           0.01,
				AuctionConfiguration: auctionConfig,
			},
			want: &auctionv2.AuctionResult{
				AuctionConfiguration: auctionConfig,
				AdUnits: &[]auction.AdUnit{
					{
						DemandID:   "gam",
						UID:        "123_gam",
						Label:      "gam",
						PriceFloor: ptr(0.1),
						BidType:    "CPM",
						Extra:      map[string]any{"placement_id": "123"},
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
						Extra:      map[string]any{"placement_id": "123"},
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
			params: &auctionv2.BuildParams{
				Adapters:             []adapter.Key{adapter.GAMKey, adapter.BidmachineKey},
				MergedAuctionRequest: request,
				PriceFloor:           0.01,
				AuctionConfiguration: auctionConfig,
			},
			want: &auctionv2.AuctionResult{
				AuctionConfiguration: auctionConfig,
				AdUnits: &[]auction.AdUnit{
					{
						DemandID:   "gam",
						UID:        "123_gam",
						Label:      "gam",
						PriceFloor: ptr(0.1),
						BidType:    "CPM",
						Extra:      map[string]any{"placement_id": "123"},
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
						Extra:      map[string]any{"placement_id": "123"},
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
			params: &auctionv2.BuildParams{
				Adapters:             []adapter.Key{adapter.GAMKey, adapter.BidmachineKey},
				MergedAuctionRequest: request,
				PriceFloor:           0.01,
				AuctionConfiguration: auctionConfig,
			},
			wantErr: true,
			err:     auction.ErrNoAdsFound,
		},
		{
			name: "No Ads Found due to empty AdUnitIDs",
			builder: testHelperAuctionBuilder(
				WithBiddingBuilder(&mocks.BiddingBuilderMock{
					HoldAuctionFunc: func(ctx context.Context, params *bidding.BuildParams) (bidding.AuctionResult, error) {
						return bidding.AuctionResult{
							RoundNumber: 0,
							Bids: []adapters.DemandResponse{
								{DemandID: "bidmachine", Bid: &adapters.BidDemandResponse{}},
							},
						}, nil
					},
				}),
			),
			params: &auctionv2.BuildParams{
				Adapters:             []adapter.Key{adapter.GAMKey, adapter.BidmachineKey},
				MergedAuctionRequest: request,
				PriceFloor:           0.01,
				AuctionConfiguration: &auction.Config{
					ID:        1,
					Demands:   []adapter.Key{adapter.GAMKey, adapter.DTExchangeKey},
					Bidding:   []adapter.Key{adapter.BidmachineKey, adapter.AmazonKey, adapter.MetaKey},
					AdUnitIDs: []int64{},
					Timeout:   15000,
				},
			},
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
				PriceFloor:           0.01,
				AuctionConfiguration: auctionConfig,
			},
			want: &auctionv2.AuctionResult{
				AuctionConfiguration: auctionConfig,
				AdUnits: &[]auction.AdUnit{
					{
						DemandID:   "gam",
						UID:        "123_gam",
						Label:      "gam",
						PriceFloor: ptr(0.1),
						BidType:    "CPM",
						Extra:      map[string]any{"placement_id": "123"},
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
						Extra:      map[string]any{"placement_id": "123"},
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

func TestAuctionResult_GetDuration(t *testing.T) {
	tests := []struct {
		name  string
		given auctionv2.AuctionResult
		want  int64
	}{
		{
			name: "Duration is 12345 when Stat present",
			given: auctionv2.AuctionResult{
				Stat: &auctionv2.Stat{
					DurationTS: 12345,
				},
			},
			want: 12345,
		},
		{
			name:  "Duration is 0 when Stat is nil",
			given: auctionv2.AuctionResult{},
			want:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.given.GetDuration(); got != tt.want {
				t.Errorf("GetDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}
