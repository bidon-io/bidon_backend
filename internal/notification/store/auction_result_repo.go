package store

import (
	"context"
	"fmt"
	"time"

	"github.com/bidon-io/bidon-backend/internal/notification"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
	"github.com/redis/go-redis/v9"
)

type AuctionResultRepo struct {
	Redis *redis.Client
}

func (r AuctionResultRepo) CreateOrUpdate(ctx context.Context, imp *schema.Imp, bids []notification.Bid) error {
	auctionResult, err := r.Find(ctx, imp.AuctionID)
	if err != nil {
		return err
	}

	round := notification.Round{
		RoundID:  imp.RoundID,
		Bids:     bids,
		BidFloor: imp.GetBidFloor(),
	}

	if auctionResult != nil {
		for _, existingRound := range auctionResult.Rounds {
			if existingRound.RoundID == imp.RoundID {
				return fmt.Errorf("round %s already exists", imp.RoundID)
			}
		}
		// This is can be potentially a problem place if we have 2 concurrent requests. Lock should be added
		auctionResult.Rounds = append(auctionResult.Rounds, round)
	} else {
		auctionResult = &notification.AuctionResult{
			AuctionID: imp.AuctionID,
			Rounds:    []notification.Round{round},
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

	winningPrice := statsRequest.Result.ECPM
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
	switch err {
	case redis.Nil: // Key does not exist
		return nil, nil
	case nil:
		return auctionResult, nil
	default:
		return nil, err
	}
}

var TTL time.Duration = 4 * time.Hour

func (r AuctionResultRepo) Save(ctx context.Context, a *notification.AuctionResult) error {
	err := r.Redis.Set(ctx, a.AuctionID, a, TTL).Err()
	if err != nil {
		return err
	}

	return nil
}
