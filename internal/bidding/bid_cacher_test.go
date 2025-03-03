package bidding_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/bidding"
	"github.com/bidon-io/bidon-backend/internal/bidding/adapters"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
	"github.com/bidon-io/bidon-backend/pkg/clock"
)

func TestBidCache_ApplyBidCache(t *testing.T) {
	redisClient, mock := redismock.NewClusterMock()
	mockTime := clock.NewMock()
	mockTime.Set(time.Now())
	bidCache := &bidding.BidCache{Redis: redisClient, Clock: mockTime}

	ctx := context.Background()
	br := &schema.BiddingRequest{
		BaseRequest: schema.BaseRequest{
			Session: schema.Session{ID: "session1"},
			Ext:     "{\"ext\":{\"bid_cache\": true}}",
		},
		AdType: "banner",
	}
	br.NormalizeValues()

	tests := []struct {
		name     string
		bids     []adapters.DemandResponse
		cacheGet bidding.Cache
		cacheSet bidding.Cache
		want     []adapters.DemandResponse
	}{
		{
			name:     "no cache, no bids",
			bids:     []adapters.DemandResponse{},
			cacheGet: bidding.Cache{},
			cacheSet: bidding.Cache{},
			want:     []adapters.DemandResponse{},
		},
		{
			name: "no cache, has bids",
			bids: []adapters.DemandResponse{
				{DemandID: adapter.BidmachineKey, Bid: &adapters.BidDemandResponse{DemandID: adapter.BidmachineKey, Price: 1.0}},
				{DemandID: adapter.ApplovinKey, Bid: &adapters.BidDemandResponse{DemandID: adapter.ApplovinKey, Price: 2.0}},
				{DemandID: adapter.MetaKey, Bid: &adapters.BidDemandResponse{DemandID: adapter.MetaKey, Price: 3.0}},
			},
			cacheGet: bidding.Cache{},
			cacheSet: bidding.Cache{
				Bids: map[adapter.Key]bidding.CacheEntry{
					adapter.ApplovinKey: {
						Bid:       bidding.CachedBid{DemandID: adapter.ApplovinKey, Price: 2.0},
						CreatedAt: mockTime.Now(),
						AuctionID: "session1",
					},
				},
			},
			want: []adapters.DemandResponse{
				{DemandID: adapter.BidmachineKey, Bid: &adapters.BidDemandResponse{DemandID: adapter.BidmachineKey, Price: 1.0}},
				{DemandID: adapter.MetaKey, Bid: &adapters.BidDemandResponse{DemandID: adapter.MetaKey, Price: 3.0}},
			},
		},
		{
			name: "no bids, has cache",
			bids: []adapters.DemandResponse{},
			cacheGet: bidding.Cache{
				Bids: map[adapter.Key]bidding.CacheEntry{
					adapter.ApplovinKey: {
						Bid:       bidding.CachedBid{DemandID: adapter.ApplovinKey, Price: 2.0},
						CreatedAt: mockTime.Now(),
						AuctionID: "session1",
					},
					adapter.VKAdsKey: {
						Bid:       bidding.CachedBid{DemandID: adapter.VKAdsKey, Price: 1.0},
						CreatedAt: mockTime.Now(),
						AuctionID: "session1",
					},
				},
			},
			cacheSet: bidding.Cache{
				Bids: map[adapter.Key]bidding.CacheEntry{
					adapter.VKAdsKey: {
						Bid:       bidding.CachedBid{DemandID: adapter.VKAdsKey, Price: 1.0},
						CreatedAt: mockTime.Now(),
						AuctionID: "session1",
					},
				},
			},
			want: []adapters.DemandResponse{
				{DemandID: adapter.ApplovinKey, Bid: &adapters.BidDemandResponse{DemandID: adapter.ApplovinKey, Price: 2.0}},
			},
		},
		{
			name: "no bids, has expired cache",
			bids: []adapters.DemandResponse{},
			cacheGet: bidding.Cache{
				Bids: map[adapter.Key]bidding.CacheEntry{
					adapter.ApplovinKey: {
						Bid:       bidding.CachedBid{DemandID: adapter.ApplovinKey, Price: 2.0},
						CreatedAt: mockTime.Now().Add(-6 * time.Minute), // Highest bid, but expired
						AuctionID: "session1",
					},
					adapter.VKAdsKey: {
						Bid:       bidding.CachedBid{DemandID: adapter.VKAdsKey, Price: 1.0},
						CreatedAt: mockTime.Now(),
						AuctionID: "session1",
					},
				},
			},
			cacheSet: bidding.Cache{},
			want: []adapters.DemandResponse{
				{DemandID: adapter.VKAdsKey, Bid: &adapters.BidDemandResponse{DemandID: adapter.VKAdsKey, Price: 1.0}},
			},
		},
		{
			name: "no valid bids, has cache",
			bids: []adapters.DemandResponse{
				{DemandID: adapter.BigoAdsKey, Bid: nil, Error: errors.New("some error")},
			},
			cacheGet: bidding.Cache{
				Bids: map[adapter.Key]bidding.CacheEntry{
					adapter.ApplovinKey: {
						Bid:       bidding.CachedBid{DemandID: adapter.ApplovinKey, Price: 2.0},
						CreatedAt: mockTime.Now().Add(-6 * time.Minute), // Highest bid, but expired
						AuctionID: "session1",
					},
					adapter.VKAdsKey: {
						Bid:       bidding.CachedBid{DemandID: adapter.VKAdsKey, Price: 1.0},
						CreatedAt: mockTime.Now(),
						AuctionID: "session1",
					},
					adapter.MetaKey: {
						Bid:       bidding.CachedBid{DemandID: adapter.MetaKey, Price: 1.5},
						CreatedAt: mockTime.Now(),
						AuctionID: "session1",
					},
				},
			},
			cacheSet: bidding.Cache{
				Bids: map[adapter.Key]bidding.CacheEntry{
					adapter.VKAdsKey: {
						Bid:       bidding.CachedBid{DemandID: adapter.VKAdsKey, Price: 1.0},
						CreatedAt: mockTime.Now(),
						AuctionID: "session1",
					},
				},
			},
			want: []adapters.DemandResponse{
				{DemandID: adapter.BigoAdsKey, Bid: nil, Error: errors.New("some error")},
				{DemandID: adapter.MetaKey, Bid: &adapters.BidDemandResponse{DemandID: adapter.MetaKey, Price: 1.5}},
			},
		},
		{
			name: "has bids, has cache",
			bids: []adapters.DemandResponse{
				{DemandID: adapter.BidmachineKey, Bid: &adapters.BidDemandResponse{DemandID: adapter.BidmachineKey, Price: 1.0}},
				{DemandID: adapter.VungleKey, Bid: &adapters.BidDemandResponse{DemandID: adapter.VungleKey, Price: 2.0}},
				{DemandID: adapter.MetaKey, Bid: &adapters.BidDemandResponse{DemandID: adapter.MetaKey, Price: 2.5}},
			},
			cacheGet: bidding.Cache{
				Bids: map[adapter.Key]bidding.CacheEntry{
					adapter.VungleKey: {
						Bid:       bidding.CachedBid{DemandID: adapter.VungleKey, Price: 3.0},
						CreatedAt: mockTime.Now().Add(-1 * time.Minute),
						AuctionID: "session1",
					},
					adapter.VKAdsKey: {
						Bid:       bidding.CachedBid{DemandID: adapter.VKAdsKey, Price: 1.0},
						CreatedAt: mockTime.Now().Add(-3 * time.Minute),
						AuctionID: "session1",
					},
				},
			},
			cacheSet: bidding.Cache{
				Bids: map[adapter.Key]bidding.CacheEntry{
					adapter.VKAdsKey: {
						Bid:       bidding.CachedBid{DemandID: adapter.VKAdsKey, Price: 1.0},
						CreatedAt: mockTime.Now().Add(-3 * time.Minute),
						AuctionID: "session1",
					},
					adapter.MetaKey: {
						Bid:       bidding.CachedBid{DemandID: adapter.MetaKey, Price: 2.5},
						CreatedAt: mockTime.Now(),
						AuctionID: "session1",
					},
				},
			},
			want: []adapters.DemandResponse{
				{DemandID: adapter.BidmachineKey, Bid: &adapters.BidDemandResponse{DemandID: adapter.BidmachineKey, Price: 1.0}},
				{DemandID: adapter.VungleKey, Bid: &adapters.BidDemandResponse{DemandID: adapter.VungleKey, Price: 3.0}},
			},
		},
		{
			name: "has bids, has cheap cache",
			bids: []adapters.DemandResponse{
				{DemandID: adapter.VungleKey, Bid: &adapters.BidDemandResponse{DemandID: adapter.VungleKey, Price: 3.0}},
			},
			cacheGet: bidding.Cache{
				Bids: map[adapter.Key]bidding.CacheEntry{
					adapter.VungleKey: {
						Bid:       bidding.CachedBid{DemandID: adapter.VungleKey, Price: 2.0},
						CreatedAt: mockTime.Now().Add(-1 * time.Minute),
						AuctionID: "session1",
					},
				},
			},
			cacheSet: bidding.Cache{
				Bids: map[adapter.Key]bidding.CacheEntry{},
			},
			want: []adapters.DemandResponse{
				{DemandID: adapter.VungleKey, Bid: &adapters.BidDemandResponse{DemandID: adapter.VungleKey, Price: 3.0}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aucResult := &bidding.AuctionResult{
				Bids: tt.bids,
			}
			bytes, _ := tt.cacheGet.MarshalBinary()
			if len(tt.cacheGet.Bids) > 0 {
				mock.ExpectGetDel("bidding:session1:banner").SetVal(string(bytes))
			} else {
				mock.ExpectGetDel("bidding:session1:banner").RedisNil()
			}
			if len(tt.cacheSet.Bids) > 0 {
				bytes, _ := tt.cacheSet.MarshalBinary()
				mock.ExpectSet("bidding:session1:banner", string(bytes), bidding.TTL).SetVal("OK")
			}

			got := bidCache.ApplyBidCache(ctx, br, aucResult)

			if diff := cmp.Diff(tt.want, got, cmpopts.IgnoreFields(adapters.DemandResponse{}, "Error")); diff != "" {
				t.Errorf("Create() mismatch (-want +got):\n%s", diff)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unmet redis expectations: %+v", err)
			}
		})
	}
}
