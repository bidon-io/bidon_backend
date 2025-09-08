package notification

import (
	"context"
	"errors"
	"log"

	"github.com/prebid/openrtb/v19/openrtb3"
	"golang.org/x/exp/slices"

	"github.com/bidon-io/bidon-backend/internal/auction"
	"github.com/bidon-io/bidon-backend/internal/bidding"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
)

type Handler struct {
	AuctionResultRepo AuctionResultRepo
	Sender            Sender
	ConfigFetcher     ConfigFetcher
}

//go:generate go run -mod=mod github.com/matryer/moq@v0.5.3 -out mocks/mocks.go -pkg mocks . AuctionResultRepo Sender ConfigFetcher

type AuctionResultRepo interface {
	CreateOrUpdate(ctx context.Context, adObject *schema.AdObject, bids []Bid) error
	Find(ctx context.Context, auctionID string) (*AuctionResult, error)
}

type Sender interface {
	SendEvent(ctx context.Context, p Params)
}

type ConfigFetcher interface {
	FetchByUIDCached(ctx context.Context, appID int64, id, uid string) *auction.Config
}

// HandleBiddingRound is used to handle bidding round, it is called after all adapters have responded with bids or errors
// Results saved to redis
func (h Handler) HandleBiddingRound(ctx context.Context, adObject *schema.AdObject, auctionResult bidding.AuctionResult, bundle, adType string) error {
	var bids []Bid
	bidFloor := adObject.GetBidFloor()

	for _, resp := range auctionResult.Bids {
		if errors.Is(resp.Error, context.DeadlineExceeded) && resp.TimeoutURL != "" {
			// Handle Timeout, currently only Meta supports this
			p := Params{
				Bundle:           bundle,
				AdType:           adType,
				AuctionID:        adObject.AuctionID,
				NotificationType: "TimeoutURL",
				URL:              resp.TimeoutURL,
				Bid:              Bid{RequestID: resp.RequestID, DemandID: resp.DemandID},
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
					AuctionID:        adObject.AuctionID,
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

	if len(bids) > 0 {
		return h.AuctionResultRepo.CreateOrUpdate(ctx, adObject, bids)
	}

	return nil
}

// HandleStats is used to handle v2/stats request
// Finalize results of auction in redis
// If external_win_notification is enabled - do nothing, wait /win or /loss request
// If external_win_notification is disabled - send win/loss notifications to demands
func (h Handler) HandleStats(ctx context.Context, stats schema.Stats, config *auction.Config, bundle, adType string) {
	if config == nil {
		log.Printf("HandleStats: cannot find config: %v", stats.AuctionConfigurationID)
		return
	}

	// If external_win_notification is enabled, do nothing - wait for /win or /loss request
	if config.ExternalWinNotifications {
		return
	}

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

	prices = append(prices, stats.AuctionPricefloor)

	for _, adUnit := range stats.AdUnits {
		if adUnit.IsFill() {
			prices = append(prices, adUnit.GetPrice())
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
		for _, bid := range auctionResult.Bids {
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
	case "FAIL":
		for _, bid := range auctionResult.Bids {
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
	case "AUCTION_CANCELLED":
		for _, bid := range auctionResult.Bids {
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

	for _, bid := range auctionResult.Bids {
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

// HandleWin is used to handle /win request
// If external_win_notification is enabled - send win/loss notifications to demands
// If external_win_notification is disabled - do nothing
func (h Handler) HandleWin(ctx context.Context, bid *schema.Bid, config *auction.Config, bundle, adType string) error {
	if config == nil {
		log.Printf("HandleWin: cannot find config: %v", bid.AuctionConfigurationID)
		return nil
	}

	// If external_win_notification is disabled, do nothing (notifications already sent in HandleStats)
	if !config.ExternalWinNotifications {
		return nil
	}

	// Get AuctionResult from redis
	auctionResult, err := h.AuctionResultRepo.Find(ctx, bid.AuctionID)
	if err != nil {
		log.Printf("HandleWin: AuctionResult exception: %s", err)
		return err
	}

	if auctionResult == nil {
		log.Printf("HandleWin: AuctionResult not found: %s", bid.AuctionID)
		return nil
	}

	// Find the winning bid and send notifications
	winningPrice := bid.GetPrice()
	var prices []float64
	prices = append(prices, winningPrice)

	// Collect all bid prices to determine second price
	for _, auctionBid := range auctionResult.Bids {
		prices = append(prices, auctionBid.Price)
	}

	slices.Sort(prices)
	var firstPrice, secondPrice float64

	if len(prices) == 0 {
		log.Printf("HandleWin: no valid prices: %s", bid.AuctionID)
		return nil
	} else if len(prices) == 1 {
		firstPrice = prices[0]
		secondPrice = prices[0]
	} else {
		firstPrice = prices[len(prices)-1]
		secondPrice = prices[len(prices)-2]
	}

	// Send notifications for all bids stored in auctionResult, regardless of incoming bid type
	for _, auctionBid := range auctionResult.Bids {
		if auctionBid.Price == winningPrice {
			// Send win notification
			go h.Sender.SendEvent(ctx, Params{
				Bundle:           bundle,
				AdType:           adType,
				AuctionID:        bid.AuctionID,
				NotificationType: "NURL",
				URL:              auctionBid.NURL,
				Bid:              auctionBid,
				Reason:           openrtb3.LossWon,
				FirstPrice:       firstPrice,
				SecondPrice:      secondPrice,
			})
		} else {
			// Send loss notification
			go h.Sender.SendEvent(ctx, Params{
				Bundle:           bundle,
				AdType:           adType,
				AuctionID:        bid.AuctionID,
				NotificationType: "LURL",
				URL:              auctionBid.LURL,
				Bid:              auctionBid,
				Reason:           openrtb3.LossLostToHigherBid,
				FirstPrice:       firstPrice,
				SecondPrice:      secondPrice,
			})
		}
	}

	return nil
}

// HandleLoss is used to handle /loss request
// If external_win_notification is enabled - send win/loss notifications to demands
// If external_win_notification is disabled - do nothing
func (h Handler) HandleLoss(ctx context.Context, bid *schema.Bid, externalWinner *schema.ExternalWinner, config *auction.Config, bundle, adType string) error {
	if config == nil {
		log.Printf("HandleLoss: cannot find config: %v", bid.AuctionConfigurationID)
		return nil
	}

	// If external_win_notification is disabled, do nothing (notifications already sent in HandleStats)
	if !config.ExternalWinNotifications {
		return nil
	}

	// Get AuctionResult from redis
	auctionResult, err := h.AuctionResultRepo.Find(ctx, bid.AuctionID)
	if err != nil {
		log.Printf("HandleLoss: AuctionResult exception: %s", err)
		return err
	}

	if auctionResult == nil {
		log.Printf("HandleLoss: AuctionResult not found: %s", bid.AuctionID)
		return nil
	}

	// Calculate prices including external winner
	var prices []float64
	prices = append(prices, bid.GetPrice())

	// Add external winner price if available
	if externalWinner != nil {
		externalPrice := externalWinner.GetPrice()
		if externalPrice > 0 {
			prices = append(prices, externalPrice)
		}
	}

	// Collect all bid prices
	for _, auctionBid := range auctionResult.Bids {
		prices = append(prices, auctionBid.Price)
	}

	slices.Sort(prices)
	var firstPrice, secondPrice float64

	if len(prices) == 0 {
		log.Printf("HandleLoss: no valid prices: %s", bid.AuctionID)
		return nil
	} else if len(prices) == 1 {
		firstPrice = prices[0]
		secondPrice = prices[0]
	} else {
		firstPrice = prices[len(prices)-1]
		secondPrice = prices[len(prices)-2]
	}

	// Send loss notifications for all bids stored in auctionResult
	// (since external winner won)
	for _, auctionBid := range auctionResult.Bids {
		go h.Sender.SendEvent(ctx, Params{
			Bundle:           bundle,
			AdType:           adType,
			AuctionID:        bid.AuctionID,
			NotificationType: "LURL",
			URL:              auctionBid.LURL,
			Bid:              auctionBid,
			Reason:           openrtb3.LossLostToHigherBid,
			FirstPrice:       firstPrice,
			SecondPrice:      secondPrice,
		})
	}

	return nil
}
