package notification

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/cenkalti/backoff/v4"
	"github.com/prebid/openrtb/v19/openrtb3"

	"github.com/bidon-io/bidon-backend/internal/sdkapi/event"
)

type Params struct {
	Bundle           string
	AdType           string
	AuctionID        string
	NotificationType string
	URL              string
	Bid              Bid
	Reason           openrtb3.LossReason
	FirstPrice       float64
	SecondPrice      float64
}

type EventSender struct {
	HttpClient  *http.Client
	EventLogger *event.Logger
}

func (es EventSender) SendEvent(ctx context.Context, p Params) {
	u, err := url.Parse(p.URL)
	if p.URL == "" || err != nil {
		log.Printf("SendNotificationEvent: failed to parse URL type %s: %s", p.NotificationType, p.URL)
		return
	}
	macroses := macrosesMap(p.Bid, p.Reason, p.FirstPrice, p.SecondPrice)
	params, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		log.Printf("SendNotificationEvent: failed to parse params: %s", u.RawQuery)
		return
	}
	for param := range params {
		if val, ok := macroses[params.Get(param)]; ok {
			params.Set(param, val)
		}
	}
	u.RawQuery = params.Encode()
	err = backoff.Retry(func() error {
		httpResp, err := es.HttpClient.Get(u.String())
		if err != nil {
			log.Printf("SendNotificationEvent: send failed: %v", err)
			return err
		}
		defer httpResp.Body.Close()

		return nil
	}, backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 3))

	e := event.NewNotificationEvent(event.NotificationParams{
		EventType:   p.NotificationType,
		ImpID:       p.Bid.ImpID,
		Bundle:      p.Bundle,
		AdType:      p.AdType,
		AuctionID:   p.AuctionID,
		DemandID:    string(p.Bid.DemandID),
		LossReason:  int64(p.Reason),
		Price:       p.Bid.Price,
		FirstPrice:  p.FirstPrice,
		SecondPrice: p.SecondPrice,
		URL:         u.String(),
		TemplateURL: p.URL,
		Error:       err,
	})
	es.EventLogger.Log(e, func(err error) {
		log.Printf("SendNotificationEvent: log notification event: %v", err)
	})

	if err != nil {
		log.Printf("SendNotificationEvent: failed to send loss notification: %s -> %s", p.Bid.DemandID, p.URL)
	}
}

func macrosesMap(bid Bid, lossReason openrtb3.LossReason, firstPrice, secondPrice float64) map[string]string {
	return map[string]string{
		"${AUCTION_MIN_TO_WIN}":         strconv.FormatFloat(secondPrice, 'f', -1, 64),
		"${AUCTION_MINIMUM_BID_TO_WIN}": strconv.FormatFloat(secondPrice, 'f', -1, 64),
		"${MIN_BID_TO_WIN}":             strconv.FormatFloat(secondPrice, 'f', -1, 64),
		"${AUCTION_ID}":                 bid.RequestID,
		"${AUCTION_BID_ID}":             bid.ID,
		"${AUCTION_IMP_ID}":             bid.ImpID,
		"${AUCTION_SEAT_ID}":            bid.SeatID,
		"${AUCTION_AD_ID}":              bid.AdID,
		"${AUCTION_PRICE}":              strconv.FormatFloat(firstPrice, 'f', -1, 64),
		"${AUCTION_LOSS}":               fmt.Sprintf("%d", lossReason),
		"${AUCTION_CURRENCY}":           "USD",
	}
}
