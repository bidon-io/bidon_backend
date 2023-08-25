package notification

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/bidon-io/bidon-backend/internal/auction"
	"github.com/bidon-io/bidon-backend/internal/bidding"
	"github.com/bidon-io/bidon-backend/internal/bidding/adapters"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
	"github.com/cenkalti/backoff/v4"
	"github.com/prebid/openrtb/v19/openrtb3"
)

type Notification struct {
	RoundID string
}

type Handler struct {
	AuctionResultRepo AuctionResultRepo
	HttpClient        *http.Client
}

//go:generate go run -mod=mod github.com/matryer/moq@latest -out mocks/mocks.go -pkg mocks . AuctionResultRepo

type AuctionResultRepo interface {
	CreateOrUpdate(ctx context.Context, imp *schema.Imp, bids []Bid) error
	Find(ctx context.Context, auctionID string) (*AuctionResult, error)
}

// HandleRound is used to handle bidding round, it is called after all adapters have responded with bids or errors
// Results saved to redis
func (h Handler) HandleRound(ctx context.Context, imp *schema.Imp, auctionResult bidding.AuctionResult) error {
	bids := []Bid{}
	bidFloor := imp.GetBidFloor()

	for _, resp := range auctionResult.Bids {
		if resp.IsBid() {
			bid := Bid{
				ID:        resp.Bid.ID,
				ImpID:     resp.Bid.ImpID,
				Price:     resp.Bid.Price,
				DemandID:  resp.Bid.DemandID,
				AdID:      resp.Bid.AdID,
				SeatID:    resp.Bid.SeatID,
				LURL:      resp.Bid.LURL,
				NURL:      resp.Bid.NURL,
				BURL:      resp.Bid.BURL,
				RequestID: resp.RequestID,
			}

			if bid.Price >= bidFloor {
				bids = append(bids, bid)
			} else {
				h.SendLoss(ctx, bid, openrtb3.LossBelowAuctionFloor, bidFloor, bidFloor)
			}
		} else {
			// TODO: maybe log no bid responses
			continue
		}
	}

	return h.AuctionResultRepo.CreateOrUpdate(ctx, imp, bids)
}

// HandleStats is used to handle /stats request
// Finalize results of auction in redis
// If external_win_notification is enabled - do nothing, wait /win or /loss request
// If external_win_notification is disabled - send win/loss notifications to demands
func (h Handler) HandleStats(ctx context.Context, stats schema.Stats, config auction.Config) error {
	if config.ExternalWinNotifications {
		return nil
	}

	// Get AuctionResult from redis
	auctionResult, err := h.AuctionResultRepo.Find(ctx, stats.AuctionID)
	if err != nil {
		return err
	}

	if auctionResult == nil {
		log.Printf("auction result not found: %s", stats.AuctionID)
		return nil
	}

	var winner *Bid
	loosers := []Bid{}
	lossReason := openrtb3.LossWon
	switch stats.Result.Status {
	case "SUCCESS": // We have winner
		winEcpm := stats.Result.ECPM
		lossReason = openrtb3.LossLostToHigherBid

		// Find all bidding rounds for this auction
		for _, round := range auctionResult.Rounds {
			// Find all bids for this round
			for _, bid := range round.Bids {
				if bid.Price == winEcpm {
					winner = &bid
				} else {
					loosers = append(loosers, bid)
				}
			}
		}

		fmt.Println(winner)
		fmt.Println(loosers)
	case "FAIL":
		lossReason = openrtb3.LossInternalError
		fmt.Println("FAIL")
	case "AUCTION_CANCELLED":
		lossReason = openrtb3.LossExpired
		fmt.Println("AUCTION_CANCELLED")
	}

	for _, bid := range loosers {
		h.SendLoss(ctx, bid, lossReason, 0, 0)
	}

	return nil
}

// HandleShow is used to handle impressions
// Send burl to demand
func (h Handler) HandleShow(ctx context.Context, imp *schema.Imp, responses []*adapters.DemandResponse) error {
	return nil
}

// If external_win_notification is enabled - send win/loss notifications to demands
// If external_win_notification is disabled - do nothing
func (h Handler) HandleWin(ctx context.Context, imp *schema.Imp, responses []*adapters.DemandResponse) error {
	return nil
}

// If external_win_notification is enabled - send win/loss notifications to demands
// If external_win_notification is disabled - do nothing
func (h Handler) HandleLoss(ctx context.Context, imp *schema.Imp, responses []*adapters.DemandResponse) error {
	return nil
}

func (h Handler) SendLoss(ctx context.Context, bid Bid, lossReason openrtb3.LossReason, winPrice, secondPrice float64) {
	u, err := url.Parse(bid.LURL)
	if bid.LURL == "" || err != nil {
		log.Printf("failed to parse url: %s", bid.LURL)
		return
	}
	macroses := macrosesMap(bid, lossReason, winPrice, secondPrice)
	params, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		log.Printf("failed to parse params: %s", u.RawQuery)
		return
	}
	for param := range params {
		if val, ok := macroses[params.Get(param)]; ok {
			params.Set(param, val)
		}
	}
	u.RawQuery = params.Encode()
	// TODO: write raw events to kafka
	err = backoff.Retry(func() error {
		_, err := h.HttpClient.Get(u.String())
		return err
	}, backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 3))

	if err != nil {
		log.Printf("failed to send loss notification: %s -> %s", bid.DemandID, bid.LURL)
	}
}

func macrosesMap(bid Bid, lossReason openrtb3.LossReason, winPrice, secondPrice float64) map[string]string {
	return map[string]string{
		"${AUCTION_MIN_TO_WIN}":         strconv.FormatFloat(secondPrice, 'f', -1, 64),
		"${AUCTION_MINIMUM_BID_TO_WIN}": strconv.FormatFloat(secondPrice, 'f', -1, 64),
		"${MIN_BID_TO_WIN}":             strconv.FormatFloat(secondPrice, 'f', -1, 64),
		"${AUCTION_ID}":                 bid.RequestID,
		"${AUCTION_BID_ID}":             bid.ID,
		"${AUCTION_IMP_ID}":             bid.ImpID,
		"${AUCTION_SEAT_ID}":            bid.SeatID,
		"${AUCTION_AD_ID}":              bid.AdID,
		"${AUCTION_PRICE}":              strconv.FormatFloat(winPrice, 'f', -1, 64),
		"${AUCTION_LOSS}":               fmt.Sprintf("%d", lossReason),
		"${AUCTION_CURRENCY}":           "USD",
	}
}
