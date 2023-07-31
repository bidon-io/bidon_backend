package store_test

import (
	"context"
	"testing"
	"time"

	"github.com/bidon-io/bidon-backend/internal/notification"
	"github.com/bidon-io/bidon-backend/internal/notification/store"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
	"github.com/go-redis/redismock/v9"
	"github.com/google/go-cmp/cmp"
)

func TestAuctionResultRepo_CreateOrUpdate(t *testing.T) {
	ctx := context.Background()
	bidFloor := 0.5
	imp := &schema.Imp{
		AuctionID: "auction-1",
		RoundID:   "round-1",
		BidFloor:  &bidFloor,
	}
	bids := []notification.Bid{
		{ID: "bid-1", ImpID: "imp-1", Price: 1.23},
		{ID: "bid-2", ImpID: "imp-1", Price: 4.56},
		{ID: "bid-3", ImpID: "imp-2", Price: 7.89},
		{ID: "bid-4", ImpID: "imp-1", Price: 0.12},
	}
	expectedAuctionResult := &notification.AuctionResult{
		AuctionID: "auction-1",
		Rounds: []notification.Round{
			{
				RoundID:  "round-1",
				Bids:     bids,
				BidFloor: 0.5,
			},
		},
	}
	rdb, mock := redismock.NewClientMock()
	mock.ExpectGet("auction-1").RedisNil()
	mock.ExpectSet("auction-1", expectedAuctionResult, 24*time.Hour).SetVal("OK")

	repo := store.AuctionResultRepo{Redis: rdb}
	err := repo.CreateOrUpdate(ctx, imp, bids)

	if mock.ExpectationsWereMet() != nil {
		t.Errorf("expectation not met: %v", mock.ExpectationsWereMet())
	}
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestAuctionResultRepo_CreateOrUpdate_DuplicateRound(t *testing.T) {
	bidFloor := 0.5
	ctx := context.Background()
	imp := &schema.Imp{
		AuctionID: "auction-1",
		RoundID:   "round-1",
		BidFloor:  &bidFloor,
	}
	bids := []notification.Bid{
		{ID: "bid-1", ImpID: "imp-1", Price: 1.23},
		{ID: "bid-2", ImpID: "imp-1", Price: 4.56},
		{ID: "bid-3", ImpID: "imp-2", Price: 7.89},
		{ID: "bid-4", ImpID: "imp-1", Price: 0.12},
	}
	existingAuctionResult := &notification.AuctionResult{
		AuctionID: "auction-1",
		Rounds: []notification.Round{
			{
				RoundID:  "round-1",
				Bids:     bids,
				BidFloor: 0.5,
			},
		},
	}
	bytes, _ := existingAuctionResult.MarshalBinary()
	rdb, mock := redismock.NewClientMock()
	mock.ExpectGet("auction-1").SetVal(string(bytes))

	repo := store.AuctionResultRepo{Redis: rdb}
	err := repo.CreateOrUpdate(ctx, imp, bids)

	if err.Error() != "round round-1 already exists" {
		t.Errorf("expectation not met: %v", err)
	}
	if err == nil {
		t.Errorf("expected error, got not errors")
	}
}

func TestAuctionResultRepo_Find(t *testing.T) {
	ctx := context.Background()
	expectedAuctionResult := &notification.AuctionResult{
		AuctionID: "auction-1",
		Rounds: []notification.Round{
			{
				RoundID:  "round-1",
				Bids:     []notification.Bid{},
				BidFloor: 0.5,
			},
		},
	}
	bytes, _ := expectedAuctionResult.MarshalBinary()
	rdb, mock := redismock.NewClientMock()
	mock.ExpectGet("auction-1").SetVal(string(bytes))

	repo := store.AuctionResultRepo{Redis: rdb}
	actualAuctionResult, err := repo.Find(ctx, "auction-1")

	if mock.ExpectationsWereMet() != nil {
		t.Errorf("expectation not met: %v", mock.ExpectationsWereMet())
	}
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if diff := cmp.Diff(expectedAuctionResult, actualAuctionResult); diff != "" {
		t.Errorf("expectedAuctionResult -> %+v mismatch \n(-want, +got)\n%s", expectedAuctionResult, diff)
	}
}

func TestAuctionResultRepo_Find_NotFound(t *testing.T) {
	ctx := context.Background()
	rdb, mock := redismock.NewClientMock()
	mock.ExpectGet("auction-1").RedisNil()

	repo := store.AuctionResultRepo{Redis: rdb}
	actualAuctionResult, err := repo.Find(ctx, "auction-1")

	if mock.ExpectationsWereMet() != nil {
		t.Errorf("expectation not met: %v", mock.ExpectationsWereMet())
	}
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if actualAuctionResult != nil {
		t.Errorf("actualAuctionResult not nil")
	}
}

func TestAuctionResult_Save(t *testing.T) {
	ctx := context.Background()
	auctionResult := &notification.AuctionResult{
		AuctionID: "auction-1",
		Rounds: []notification.Round{
			{
				RoundID:  "round-1",
				Bids:     []notification.Bid{},
				BidFloor: 0.5,
			},
		},
	}
	rdb, mock := redismock.NewClientMock()
	mock.ExpectSet("auction-1", auctionResult, 24*time.Hour).SetVal("OK")
	repo := store.AuctionResultRepo{Redis: rdb}

	err := repo.Save(ctx, auctionResult)

	if mock.ExpectationsWereMet() != nil {
		t.Errorf("expectation not met: %v", mock.ExpectationsWereMet())
	}
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
