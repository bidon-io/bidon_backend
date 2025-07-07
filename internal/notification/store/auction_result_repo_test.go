package store_test

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/google/go-cmp/cmp"

	"github.com/bidon-io/bidon-backend/internal/notification"
	"github.com/bidon-io/bidon-backend/internal/notification/store"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
)

func TestAuctionResultRepo_CreateOrUpdate(t *testing.T) {
	ctx := context.Background()
	bidFloor := 0.5
	imp := &schema.AdObject{
		AuctionID:  "auction-1",
		PriceFloor: bidFloor,
	}
	bids := []notification.Bid{
		{ID: "bid-1", ImpID: "imp-1", Price: 1.23},
		{ID: "bid-2", ImpID: "imp-1", Price: 4.56},
		{ID: "bid-3", ImpID: "imp-2", Price: 7.89},
		{ID: "bid-4", ImpID: "imp-1", Price: 0.12},
	}
	expectedAuctionResultV2 := &notification.AuctionResult{
		AuctionID: "auction-1",
		Bids:      bids,
	}
	rdb, mock := redismock.NewClusterMock()
	mock.ExpectGet("auction-1").RedisNil()
	mock.ExpectSet("auction-1", expectedAuctionResultV2, 4*time.Hour).SetVal("OK")

	repo := store.AuctionResultRepo{Redis: rdb}
	err := repo.CreateOrUpdate(ctx, imp, bids)

	if mock.ExpectationsWereMet() != nil {
		t.Errorf("expectation not met: %v", mock.ExpectationsWereMet())
	}
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestAuctionResultRepo_Find(t *testing.T) {
	ctx := context.Background()
	expectedAuctionResultV2 := &notification.AuctionResult{
		AuctionID: "auction-1",
		Bids:      []notification.Bid{},
	}
	bytes, _ := expectedAuctionResultV2.MarshalBinary()
	rdb, mock := redismock.NewClusterMock()
	mock.ExpectGet("auction-1").SetVal(string(bytes))

	repo := store.AuctionResultRepo{Redis: rdb}
	actualAuctionResultV2, err := repo.Find(ctx, "auction-1")

	if mock.ExpectationsWereMet() != nil {
		t.Errorf("expectation not met: %v", mock.ExpectationsWereMet())
	}
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if diff := cmp.Diff(expectedAuctionResultV2, actualAuctionResultV2); diff != "" {
		t.Errorf("expectedAuctionResultV2 -> %+v mismatch \n(-want, +got)\n%s", expectedAuctionResultV2, diff)
	}
}

func TestAuctionResultRepo_Find_NotFound(t *testing.T) {
	ctx := context.Background()
	rdb, mock := redismock.NewClusterMock()
	mock.ExpectGet("auction-1").RedisNil()

	repo := store.AuctionResultRepo{Redis: rdb}
	actualAuctionResultV2, err := repo.Find(ctx, "auction-1")

	if mock.ExpectationsWereMet() != nil {
		t.Errorf("expectation not met: %v", mock.ExpectationsWereMet())
	}
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if actualAuctionResultV2 != nil {
		t.Errorf("actualAuctionResultV2 not nil")
	}
}

func TestAuctionResultV2_Save(t *testing.T) {
	ctx := context.Background()
	AuctionResultV2 := &notification.AuctionResult{
		AuctionID: "auction-1",
		Bids:      []notification.Bid{},
	}
	rdb, mock := redismock.NewClusterMock()
	mock.ExpectSet("auction-1", AuctionResultV2, 4*time.Hour).SetVal("OK")
	repo := store.AuctionResultRepo{Redis: rdb}

	err := repo.Save(ctx, AuctionResultV2)

	if mock.ExpectationsWereMet() != nil {
		t.Errorf("expectation not met: %v", mock.ExpectationsWereMet())
	}
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
