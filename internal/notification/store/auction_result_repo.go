package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/bidon-io/bidon-backend/internal/notification"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
)

var TTL = 4 * time.Hour

type AuctionResultRepo struct {
	Redis *redis.ClusterClient
}

func (r AuctionResultRepo) CreateOrUpdate(ctx context.Context, adObject *schema.AdObject, bids []notification.Bid) error {
	auctionResult, err := r.Find(ctx, adObject.AuctionID)
	if err != nil {
		return err
	}

	if auctionResult != nil {
		auctionResult.Bids = bids
	} else {
		auctionResult = &notification.AuctionResult{
			AuctionID: adObject.AuctionID,
			Bids:      bids,
		}
	}

	err = r.Save(ctx, auctionResult)
	if err != nil {
		return err
	}

	return nil
}

func (r AuctionResultRepo) FinalizeResult(ctx context.Context, statsRequest *schema.Stats) error {
	if !statsRequest.Result.IsSuccess() {
		return nil
	}

	winningPrice := statsRequest.Result.Price
	fmt.Println(winningPrice)
	auctionResult, err := r.Find(ctx, statsRequest.AuctionID)
	if err != nil {
		return err
	}
	fmt.Println(auctionResult)

	return nil
}

func (r AuctionResultRepo) Find(ctx context.Context, auctionID string) (*notification.AuctionResult, error) {
	auctionResult := &notification.AuctionResult{}
	err := r.Redis.Get(ctx, auctionID).Scan(auctionResult)
	switch {
	case errors.Is(err, redis.Nil): // Key does not exist
		return nil, nil
	case err == nil:
		return auctionResult, nil
	default:
		return nil, err
	}
}

func (r AuctionResultRepo) Save(ctx context.Context, a *notification.AuctionResult) error {
	err := r.Redis.Set(ctx, a.AuctionID, a, TTL).Err()
	if err != nil {
		return err
	}

	return nil
}
