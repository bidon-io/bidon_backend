package auction_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/auction"
	"github.com/bidon-io/bidon-backend/internal/auction/mocks"
	"github.com/bidon-io/bidon-backend/internal/bidding"
	"github.com/bidon-io/bidon-backend/internal/bidding/adapters"
	"github.com/bidon-io/bidon-backend/internal/sdkapi"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
)

func testApp(id int64) *sdkapi.App {
	return &sdkapi.App{ID: id}
}

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
		MatchCachedFunc: func(_ context.Context, _ *auction.BuildParams) ([]auction.AdUnit, error) {
			return adUnits, nil
		},
	}
	biddingAdaptersConfigBuilder := &mocks.BiddingAdaptersConfigBuilderMock{
		BuildFunc: func(_ context.Context, _ int64, _ []adapter.Key, _ *auction.AdUnitsMap) (adapter.ProcessedConfigsMap, error) {
			return adapter.ProcessedConfigsMap{}, nil
		},
	}
	biddingBuilder := &mocks.BiddingBuilderMock{
		HoldAuctionFunc: func(_ context.Context, _ *bidding.BuildParams) (bidding.AuctionResult, error) {
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

func testHelperAuctionBuilder(opts ...BuilderOption) *auction.Builder {
	m := testHelperDefaultAuctionBuilderMocks()

	for _, opt := range opts {
		opt(m)
	}

	return &auction.Builder{
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
	request := &schema.AuctionRequest{}
	testCases := []struct {
		name    string
		builder *auction.Builder
		params  *auction.BuildParams
		want    *auction.Result
		wantErr bool
		err     error
	}{
		{
			name: "Success",
			builder: testHelperAuctionBuilder(WithBiddingBuilder(&mocks.BiddingBuilderMock{
				HoldAuctionFunc: func(_ context.Context, _ *bidding.BuildParams) (bidding.AuctionResult, error) {
					return bidding.AuctionResult{
						RoundNumber: 0,
						Bids: []adapters.DemandResponse{
							{DemandID: "bidmachine", Bid: &adapters.BidDemandResponse{}},
						},
					}, nil
				},
			})),
			params: &auction.BuildParams{
				App: testApp(1),
				Adapters:             []adapter.Key{adapter.GAMKey, adapter.BidmachineKey},
				AuctionRequest:       request,
				PriceFloor:           0.02,
				AuctionConfiguration: auctionConfig,
			},
			want: &auction.Result{
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
			name: "Uplifts PriceFloor for BM CPM",
			builder: testHelperAuctionBuilder(
				WithAdUnitsMatcher(
					&mocks.AdUnitsMatcherMock{
						MatchCachedFunc: func(_ context.Context, _ *auction.BuildParams) ([]auction.AdUnit, error) {
							return []auction.AdUnit{
								{
									DemandID:   "bidmachine",
									Label:      "BM",
									PriceFloor: ptr(0.1),
									UID:        "123_bidmachine",
									BidType:    schema.CPMBidType,
								},
							}, nil
						},
					},
				),
				WithBiddingBuilder(
					&mocks.BiddingBuilderMock{
						HoldAuctionFunc: func(_ context.Context, _ *bidding.BuildParams) (bidding.AuctionResult, error) {
							return bidding.AuctionResult{
								Bids: []adapters.DemandResponse{
									{DemandID: "meta", Bid: &adapters.BidDemandResponse{Price: 0.5}},
								},
							}, nil
						},
					},
				),
			),
			params: &auction.BuildParams{
				App: testApp(1),
				Adapters:             []adapter.Key{adapter.GAMKey, adapter.BidmachineKey},
				AuctionRequest:       request,
				PriceFloor:           0.02,
				AuctionConfiguration: auctionConfig,
			},
			want: &auction.Result{
				AuctionConfiguration: auctionConfig,
				AdUnits: &[]auction.AdUnit{
					{
						DemandID:   "bidmachine",
						Label:      "BM",
						PriceFloor: ptr(0.1),
						UID:        "123_bidmachine",
						BidType:    schema.CPMBidType,
					},
				},

				CPMAdUnits: &[]auction.AdUnit{
					{
						DemandID:   "bidmachine",
						Label:      "BM",
						PriceFloor: ptr(0.51),
						UID:        "123_bidmachine",
						BidType:    schema.CPMBidType,
					},
				},
				BiddingAuctionResult: &bidding.AuctionResult{
					Bids: []adapters.DemandResponse{
						{DemandID: "meta", Bid: &adapters.BidDemandResponse{Price: 0.5}},
					},
				},
			},
		},
		{
			name: "No Biding Adapters Matched",
			builder: testHelperAuctionBuilder(WithBiddingBuilder(&mocks.BiddingBuilderMock{
				HoldAuctionFunc: func(_ context.Context, _ *bidding.BuildParams) (bidding.AuctionResult, error) {
					return bidding.AuctionResult{}, bidding.ErrNoAdaptersMatched
				},
			})),
			params: &auction.BuildParams{
				App: testApp(1),
				Adapters:             []adapter.Key{adapter.GAMKey, adapter.BidmachineKey},
				AuctionRequest:       request,
				PriceFloor:           0.02,
				AuctionConfiguration: auctionConfig,
			},
			want: &auction.Result{
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
					MatchCachedFunc: func(_ context.Context, _ *auction.BuildParams) ([]auction.AdUnit, error) {
						return []auction.AdUnit{}, nil
					},
				}),
				WithBiddingBuilder(&mocks.BiddingBuilderMock{
					HoldAuctionFunc: func(_ context.Context, _ *bidding.BuildParams) (bidding.AuctionResult, error) {
						return bidding.AuctionResult{}, nil
					},
				}),
			),
			params: &auction.BuildParams{
				App: testApp(1),
				Adapters:             []adapter.Key{adapter.GAMKey, adapter.BidmachineKey},
				AuctionRequest:       request,
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
					HoldAuctionFunc: func(_ context.Context, _ *bidding.BuildParams) (bidding.AuctionResult, error) {
						return bidding.AuctionResult{
							RoundNumber: 0,
							Bids: []adapters.DemandResponse{
								{DemandID: "bidmachine", Bid: &adapters.BidDemandResponse{}},
							},
						}, nil
					},
				}),
			),
			params: &auction.BuildParams{
				App: testApp(1),
				Adapters:       []adapter.Key{adapter.GAMKey, adapter.BidmachineKey},
				AuctionRequest: request,
				PriceFloor:     0.01,
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
				HoldAuctionFunc: func(_ context.Context, _ *bidding.BuildParams) (bidding.AuctionResult, error) {
					return bidding.AuctionResult{
						RoundNumber: 0,
						Bids: []adapters.DemandResponse{
							{DemandID: "bidmachine", Bid: &adapters.BidDemandResponse{}},
						},
					}, nil
				},
			})),
			params: &auction.BuildParams{
				App: testApp(1),
				AuctionKey:           "1ERNSV33K4000",
				Adapters:             []adapter.Key{adapter.GAMKey, adapter.BidmachineKey},
				AuctionRequest:       request,
				PriceFloor:           0.01,
				AuctionConfiguration: auctionConfig,
			},
			want: &auction.Result{
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
		given auction.Result
		want  int64
	}{
		{
			name: "Duration is 12345 when Stat present",
			given: auction.Result{
				Stat: &auction.Stat{
					DurationTS: 12345,
				},
			},
			want: 12345,
		},
		{
			name:  "Duration is 0 when Stat is nil",
			given: auction.Result{},
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

func TestAuctionResult_GetMaxBidPrice(t *testing.T) {
	tests := []struct {
		name     string
		bids     []adapters.DemandResponse
		expected float64
	}{
		{
			name:     "no bids",
			bids:     []adapters.DemandResponse{},
			expected: 0.0,
		},
		{
			name: "single bid",
			bids: []adapters.DemandResponse{
				{
					Bid: &adapters.BidDemandResponse{
						Price: 1.5,
					},
				},
			},
			expected: 1.5,
		},
		{
			name: "multiple bids",
			bids: []adapters.DemandResponse{
				{
					Bid: &adapters.BidDemandResponse{
						Price: 1.5,
					},
				},
				{},
				{
					Bid: &adapters.BidDemandResponse{
						Price: 1,
					},
				},
			},
			expected: 1.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := bidding.AuctionResult{Bids: tt.bids}
			maxPrice := result.GetMaxBidPrice()
			if maxPrice != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, maxPrice)
			}
		})
	}
}
