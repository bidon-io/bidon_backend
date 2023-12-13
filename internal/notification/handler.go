package notification

import (
	"context"
	"errors"
	"github.com/bidon-io/bidon-backend/internal/auction"
	"github.com/bidon-io/bidon-backend/internal/bidding"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
	"github.com/prebid/openrtb/v19/openrtb3"
	"golang.org/x/exp/slices"
	"log"
)

type Handler struct {
	AuctionResultRepo AuctionResultRepo
	Sender            Sender
}

//go:generate go run -mod=mod github.com/matryer/moq@latest -out mocks/mocks.go -pkg mocks . AuctionResultRepo Sender

type AuctionResultRepo interface {
	CreateOrUpdate(ctx context.Context, imp *schema.Imp, bids []Bid) error
	Find(ctx context.Context, auctionID string) (*AuctionResult, error)
}

type Sender interface {
	SendEvent(ctx context.Context, p Params)
}

// HandleBiddingRound is used to handle bidding round, it is called after all adapters have responded with bids or errors
// Results saved to redis
func (h Handler) HandleBiddingRound(ctx context.Context, imp *schema.Imp, auctionResult bidding.AuctionResult, bundle, adType string) error {
	var bids []Bid
	bidFloor := imp.GetBidFloor()

	for _, resp := range auctionResult.Bids {
		if errors.Is(resp.Error, context.DeadlineExceeded) && resp.TimeoutURL != "" {
			// Handle Timeout, currently only Meta supports this
			p := Params{
				Bundle:           bundle,
				AdType:           adType,
				AuctionID:        imp.AuctionID,
				NotificationType: "TimeoutURL",
				URL:              resp.TimeoutURL,
				Bid:              Bid{RequestID: resp.RequestID},
				Reason:           openrtb3.LossExpired,
				FirstPrice:       bidFloor,
				SecondPrice:      bidFloor,
			}

			h.Sender.SendEvent(ctx, p)
		} else if resp.IsBid() {
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

			if bid.Price >= bidFloor { // Valid Bid, use for further processing
				bids = append(bids, bid)
			} else { // Send Loss notification straight away
				go h.Sender.SendEvent(ctx, Params{
					Bundle:           bundle,
					AdType:           adType,
					AuctionID:        imp.AuctionID,
					NotificationType: "LURL",
					URL:              bid.LURL,
					Bid:              bid,
					Reason:           openrtb3.LossBelowAuctionFloor,
					FirstPrice:       bidFloor,
					SecondPrice:      bidFloor,
				})
			}
		}
	}

	return h.AuctionResultRepo.CreateOrUpdate(ctx, imp, bids)
}

// HandleStats is used to handle /stats request
// Finalize results of auction in redis
// If external_win_notification is enabled - do nothing, wait /win or /loss request
// If external_win_notification is disabled - send win/loss notifications to demands
func (h Handler) HandleStats(ctx context.Context, stats schema.Stats, config *auction.Config, bundle, adType string) {
	if config == nil {
		log.Printf("HandleStats: cannot find config: %v", stats.AuctionConfigurationID)
		return
	}

	// Disable external_win_notification until we have cancel notifications
	//if config.ExternalWinNotifications {
	//	return
	//}

	// Get AuctionResult from redis
	auctionResult, err := h.AuctionResultRepo.Find(ctx, stats.AuctionID)
	if err != nil {
		log.Printf("HandleStats: AuctionResult exception: %s", err)
		return
	}

	if auctionResult == nil {
		log.Printf("HandleStats: AuctionResult not found: %s", stats.AuctionID)
		return
	}

	var notifications []Params
	var prices []float64
	for _, round := range stats.Rounds {
		prices = append(prices, round.PriceFloor)

		for _, bid := range round.Demands {
			if bid.IsFill() {
				prices = append(prices, bid.GetPrice())
			}
		}

		for _, bid := range round.Bidding.Bids {
			if bid.IsFill() {
				prices = append(prices, bid.GetPrice())
			}
		}
	}
	slices.Sort(prices)
	var firstPrice, secondPrice float64

	if len(prices) == 0 {
		log.Printf("HandleStats: no valid prices: %s", stats.AuctionID)
		return
	} else if len(prices) == 1 {
		firstPrice = prices[0]
		secondPrice = prices[0]
	} else {
		firstPrice = prices[len(prices)-1]
		secondPrice = prices[len(prices)-2]
	}

	switch stats.Result.Status {
	case "SUCCESS": // We have winner
		// Find all bidding rounds for this auction
		for _, round := range auctionResult.Rounds {
			// Find all bids for this round
			for _, bid := range round.Bids {
				if bid.Price == firstPrice {
					notifications = append(notifications, Params{
						Bundle:           bundle,
						AdType:           adType,
						AuctionID:        stats.AuctionID,
						NotificationType: "NURL",
						URL:              bid.NURL,
						Bid:              bid,
						Reason:           openrtb3.LossWon,
						FirstPrice:       firstPrice,
						SecondPrice:      secondPrice,
					})
				} else {
					notifications = append(notifications, Params{
						Bundle:           bundle,
						AdType:           adType,
						AuctionID:        stats.AuctionID,
						NotificationType: "LURL",
						URL:              bid.LURL,
						Bid:              bid,
						Reason:           openrtb3.LossLostToHigherBid,
						FirstPrice:       firstPrice,
						SecondPrice:      secondPrice,
					})
				}
			}
		}
	case "FAIL":
		// Find all bidding rounds for this auction
		for _, round := range auctionResult.Rounds {
			// Find all bids for this round
			for _, bid := range round.Bids {
				notifications = append(notifications, Params{
					Bundle:           bundle,
					AdType:           adType,
					AuctionID:        stats.AuctionID,
					NotificationType: "LURL",
					URL:              bid.LURL,
					Bid:              bid,
					Reason:           openrtb3.LossInternalError,
					FirstPrice:       firstPrice,
					SecondPrice:      secondPrice,
				})
			}
		}
	case "AUCTION_CANCELLED":
		for _, round := range auctionResult.Rounds {
			// Find all bids for this round
			for _, bid := range round.Bids {
				notifications = append(notifications, Params{
					Bundle:           bundle,
					AdType:           adType,
					AuctionID:        stats.AuctionID,
					NotificationType: "LURL",
					URL:              bid.LURL,
					Bid:              bid,
					Reason:           openrtb3.LossExpired,
					FirstPrice:       firstPrice,
					SecondPrice:      secondPrice,
				})
			}
		}
	}

	for _, n := range notifications {
		go h.Sender.SendEvent(ctx, n)
	}
}

// HandleShow is used to handle /show request
// Send burl to demand
func (h Handler) HandleShow(ctx context.Context, impression *schema.Bid, bundle, adType string) {
	if !impression.IsBidding() {
		// Not bidding, do nothing
		return
	}

	auctionResult, err := h.AuctionResultRepo.Find(ctx, impression.AuctionID)
	if err != nil {
		log.Printf("HandleShow: AuctionResult exception: %s", err)
		return
	}

	if auctionResult == nil {
		log.Printf("HandleShow: AuctionResult not found: %s", impression.AuctionID)
		return
	}

	for _, round := range auctionResult.Rounds {
		for _, bid := range round.Bids {
			if bid.Price == impression.GetPrice() {
				go h.Sender.SendEvent(ctx, Params{
					Bundle:           bundle,
					AdType:           adType,
					AuctionID:        impression.AuctionID,
					NotificationType: "BURL",
					URL:              bid.BURL,
					Bid:              bid,
					Reason:           openrtb3.LossWon,
					FirstPrice:       impression.GetPrice(),
					SecondPrice:      0,
				})

				break
			}
		}
	}
}

// HandleWin is used to handle /win request
// If external_win_notification is enabled - send win/loss notifications to demands
// If external_win_notification is disabled - do nothing
func (h Handler) HandleWin(ctx context.Context, bid *schema.Bid) error {
	return nil
}

// HandleLoss is used to handle /loss request
// If external_win_notification is enabled - send win/loss notifications to demands
// If external_win_notification is disabled - do nothing
func (h Handler) HandleLoss(ctx context.Context, bid *schema.Bid) error {
	return nil
}
