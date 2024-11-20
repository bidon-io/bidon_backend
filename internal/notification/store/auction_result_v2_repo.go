package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/bidon-io/bidon-backend/internal/notification"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
	"github.com/redis/go-redis/v9"
)

type AuctionResultV2Repo struct {
	Redis *redis.Client
}

func (r AuctionResultV2Repo) CreateOrUpdate(ctx context.Context, imp *schema.Imp, bids []notification.Bid) error {
	auctionResult, err := r.Find(ctx, imp.AuctionID)
	if err != nil {
		return err
	}

	if auctionResult != nil {
		auctionResult.Bids = bids
	} else {
		auctionResult = &notification.AuctionResult{
			AuctionID: imp.AuctionID,
			Bids:      bids,
		}
	}

	err = r.Save(ctx, auctionResult)
	if err != nil {
		return err
	}

	return nil
}

func (r AuctionResultV2Repo) FinalizeResult(ctx context.Context, statsRequest *schema.Stats) error {
	if !statsRequest.Result.IsSuccess() {
		return nil
	}

	winningPrice := statsRequest.Result.ECPM
	fmt.Println(winningPrice)
	auctionResult, err := r.Find(ctx, statsRequest.AuctionID)
	if err != nil {
		return err
	}
	fmt.Println(auctionResult)

	return nil
}

func (r AuctionResultV2Repo) Find(ctx context.Context, auctionID string) (*notification.AuctionResult, error) {
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

func (r AuctionResultV2Repo) Save(ctx context.Context, a *notification.AuctionResult) error {
	err := r.Redis.Set(ctx, a.AuctionID, a, TTL).Err()
	if err != nil {
		return err
	}

	return nil
}
