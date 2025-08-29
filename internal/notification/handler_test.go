package notification_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/bidon-io/bidon-backend/internal/auction"
	"github.com/bidon-io/bidon-backend/internal/notification"
	"github.com/bidon-io/bidon-backend/internal/notification/mocks"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
)

func TestHandler_HandleStats_WinBid(t *testing.T) {
	ctx := context.Background()
	imp := schema.Stats{
		AuctionID:               "f26af577-869e-41cb-909e-4d3eba57a28b",
		AuctionPricefloor:       7.0,
		AuctionConfigurationID:  10,
		AuctionConfigurationUID: "1701972528521547776",
		Result: schema.AuctionResult{
			Status:            "SUCCESS",
			BidType:           schema.RTBBidType,
			Price:             7.89,
			WinnerDemandID:    "vungle",
			WinnerAdUnitUID:   "1633824270331281408",
			WinnerAdUnitLabel: "vungle_inter_mergeblock_ios_3",
		},
		AdUnits: []schema.AuctionAdUnitResult{
			{
				Price:       7.30,
				DemandID:    "applovin",
				BidType:     schema.CPMBidType,
				AdUnitUID:   "1633833116256829440",
				AdUnitLabel: "applovin_inter_mergeblock_ios_6",
				Status:      "WIN",
			},
			{
				Price:       1.23,
				DemandID:    "bigoads",
				BidType:     schema.RTBBidType,
				AdUnitUID:   "1633833116256829140",
				AdUnitLabel: "bigoads_inter_mergeblock_ios_5",
				Status:      "LOSS",
			},
			{
				DemandID: "meta",
				BidType:  schema.RTBBidType,
				Status:   "NO_BID",
			},
			{
				Price:       7.89,
				DemandID:    "vungle",
				BidType:     schema.RTBBidType,
				AdUnitUID:   "1633824270331281408",
				AdUnitLabel: "vungle_inter_mergeblock_ios_3",
				Status:      "WIN",
			},
		},
	}
	result := notification.AuctionResult{
		Bids: []notification.Bid{
			{ID: "bid-1", ImpID: "imp-1", Price: 1.23},
			{ID: "bid-2", ImpID: "imp-1", Price: 4.56},
			{ID: "bid-3", ImpID: "imp-2", Price: 7.89},
			{ID: "bid-4", ImpID: "imp-1", Price: 0.12},
		},
	}
	config := auction.Config{ExternalWinNotifications: false}
	mockRepo := &mocks.AuctionResultRepoMock{}
	mockRepo.FindFunc = func(ctx context.Context, id string) (*notification.AuctionResult, error) {
		return &result, nil
	}
	wg := &sync.WaitGroup{}
	wg.Add(4) // We have 4 bid in the auction result, wait all

	sender := &mocks.SenderMock{SendEventFunc: func(_ context.Context, p notification.Params) {
		defer wg.Done()

		if p.FirstPrice != 7.89 {
			t.Errorf("expected first price 7.89, got %f", p.FirstPrice)
		}

		if p.SecondPrice != 7.30 {
			t.Errorf("expected second price 7.30, got %f", p.SecondPrice)
		}
	}}

	handlerV2 := notification.Handler{AuctionResultRepo: mockRepo, Sender: sender}

	handlerV2.HandleStats(ctx, imp, &config, "bundle-1", "banner")

	if waitTimeout(wg, 1*time.Second) {
		t.Errorf("timeout waiting for events, sent event lower than expected")
	}
}

func TestHandler_HandleStats_Loss(t *testing.T) {
	ctx := context.Background()

	imp := schema.Stats{
		AuctionID:               "f26af577-869e-41cb-909e-4d3eba57a28b",
		AuctionPricefloor:       7.0,
		AuctionConfigurationID:  10,
		AuctionConfigurationUID: "1701972528521547776",
		Result: schema.AuctionResult{
			Status:            "SUCCESS",
			BidType:           schema.CPMBidType,
			Price:             7.89,
			WinnerDemandID:    "vungle",
			WinnerAdUnitUID:   "1633824270331281408",
			WinnerAdUnitLabel: "vungle_inter_mergeblock_ios_3",
		},
		AdUnits: []schema.AuctionAdUnitResult{
			{
				Price:       7.89,
				DemandID:    "dtexchange",
				BidType:     schema.CPMBidType,
				AdUnitUID:   "1633833116256829440",
				AdUnitLabel: "dtexchange_inter_mergeblock_ios_6",
				Status:      "WIN",
			},
			{
				Price:       7.6,
				DemandID:    "vungle",
				BidType:     schema.RTBBidType,
				AdUnitUID:   "1633824270331281408",
				AdUnitLabel: "vungle_inter_mergeblock_ios_3",
				Status:      "LOSS",
			},
		},
	}

	config := auction.Config{ExternalWinNotifications: false}
	repoMock := &mocks.AuctionResultRepoMock{
		FindFunc: func(ctx context.Context, auctionID string) (*notification.AuctionResult, error) {
			return &notification.AuctionResult{
				Bids: []notification.Bid{{
					ID:    "bid-2",
					ImpID: "imp-1",
					Price: 7.6,
				}},
			}, nil
		},
	}

	wg := &sync.WaitGroup{}
	wg.Add(1) // We have 1 bid in the auction result, should wait single goroutine to finish

	sender := &mocks.SenderMock{SendEventFunc: func(_ context.Context, p notification.Params) {
		wg.Done()
	}}

	handlerV2 := notification.Handler{
		AuctionResultRepo: repoMock,
		Sender:            sender,
	}

	handlerV2.HandleStats(ctx, imp, &config, "bundle-1", "banner")

	if waitTimeout(wg, 1*time.Second) {
		t.Errorf("timeout waiting for events, sent event lower than expected")
	}
}

func TestHandler_HandleWin_ExternalNotificationsEnabled(t *testing.T) {
	ctx := context.Background()

	bid := &schema.Bid{
		AuctionID:               "test-auction-id",
		AuctionConfigurationID:  123,
		AuctionConfigurationUID: "test-config-uid",
		DemandID:                "test-demand",
		Price:                   5.0,
		AuctionPriceFloor:       1.0,
		BidType:                 schema.RTBBidType,
	}

	config := &auction.Config{ExternalWinNotifications: true}

	auctionResult := &notification.AuctionResult{
		AuctionID: "test-auction-id",
		Bids: []notification.Bid{
			{ID: "bid-1", Price: 5.0, NURL: "http://win-url-1", LURL: "http://loss-url-1"},
			{ID: "bid-2", Price: 3.0, NURL: "http://win-url-2", LURL: "http://loss-url-2"},
			{ID: "bid-3", Price: 2.0, NURL: "http://win-url-3", LURL: "http://loss-url-3"},
		},
	}

	mockRepo := &mocks.AuctionResultRepoMock{
		FindFunc: func(_ context.Context, auctionID string) (*notification.AuctionResult, error) {
			if auctionID != "test-auction-id" {
				t.Errorf("expected auction ID 'test-auction-id', got '%s'", auctionID)
			}
			return auctionResult, nil
		},
	}

	wg := &sync.WaitGroup{}
	wg.Add(3) // Expect 3 notifications (1 win + 2 loss)

	var sentEvents []notification.Params
	sender := &mocks.SenderMock{
		SendEventFunc: func(_ context.Context, p notification.Params) {
			defer wg.Done()
			sentEvents = append(sentEvents, p)
		},
	}

	handler := notification.Handler{
		AuctionResultRepo: mockRepo,
		Sender:            sender,
	}

	err := handler.HandleWin(ctx, bid, config, "test-bundle", "banner")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if waitTimeout(wg, 1*time.Second) {
		t.Errorf("timeout waiting for events")
	}

	// Verify we got the expected notifications
	if len(sentEvents) != 3 {
		t.Errorf("expected 3 events, got %d", len(sentEvents))
	}

	// Check that one event is NURL (win) and two are LURL (loss)
	winCount := 0
	lossCount := 0
	for _, event := range sentEvents {
		switch event.NotificationType {
		case "NURL":
			winCount++
		case "LURL":
			lossCount++
		}
	}

	if winCount != 1 {
		t.Errorf("expected 1 win notification, got %d", winCount)
	}
	if lossCount != 2 {
		t.Errorf("expected 2 loss notifications, got %d", lossCount)
	}
}

func TestHandler_HandleWin_ExternalNotificationsDisabled(t *testing.T) {
	ctx := context.Background()

	bid := &schema.Bid{
		AuctionID:               "test-auction-id",
		AuctionConfigurationID:  123,
		AuctionConfigurationUID: "test-config-uid",
		DemandID:                "test-demand",
		Price:                   5.0,
		BidType:                 schema.RTBBidType,
	}

	config := &auction.Config{ExternalWinNotifications: false}

	mockRepo := &mocks.AuctionResultRepoMock{}
	sender := &mocks.SenderMock{
		SendEventFunc: func(_ context.Context, p notification.Params) {
			t.Errorf("no events should be sent when external notifications are disabled")
		},
	}

	handler := notification.Handler{
		AuctionResultRepo: mockRepo,
		Sender:            sender,
	}

	err := handler.HandleWin(ctx, bid, config, "test-bundle", "banner")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Give some time to ensure no events are sent
	time.Sleep(100 * time.Millisecond)
}

func TestHandler_HandleWin_NonRTBBid(t *testing.T) {
	ctx := context.Background()

	bid := &schema.Bid{
		AuctionID:               "test-auction-id",
		AuctionConfigurationID:  123,
		AuctionConfigurationUID: "test-config-uid",
		DemandID:                "test-demand",
		Price:                   5.0,
		AuctionPriceFloor:       1.0,
		BidType:                 schema.CPMBidType, // Non-RTB bid
	}

	config := &auction.Config{ExternalWinNotifications: true}

	auctionResult := &notification.AuctionResult{
		AuctionID: "test-auction-id",
		Bids: []notification.Bid{
			{ID: "bid-1", Price: 5.0, NURL: "http://win-url-1", LURL: "http://loss-url-1"},
			{ID: "bid-2", Price: 3.0, NURL: "http://win-url-2", LURL: "http://loss-url-2"},
		},
	}

	mockRepo := &mocks.AuctionResultRepoMock{
		FindFunc: func(_ context.Context, auctionID string) (*notification.AuctionResult, error) {
			return auctionResult, nil
		},
	}

	wg := &sync.WaitGroup{}
	wg.Add(2) // Expect 2 notifications (1 win + 1 loss) for all bids regardless of incoming bid type

	var sentEvents []notification.Params
	sender := &mocks.SenderMock{
		SendEventFunc: func(_ context.Context, p notification.Params) {
			defer wg.Done()
			sentEvents = append(sentEvents, p)
		},
	}

	handler := notification.Handler{
		AuctionResultRepo: mockRepo,
		Sender:            sender,
	}

	err := handler.HandleWin(ctx, bid, config, "test-bundle", "banner")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if waitTimeout(wg, 1*time.Second) {
		t.Errorf("timeout waiting for events")
	}

	// Verify we got the expected notifications
	if len(sentEvents) != 2 {
		t.Errorf("expected 2 events, got %d", len(sentEvents))
	}

	// Check that one event is NURL (win) and one is LURL (loss)
	winCount := 0
	lossCount := 0
	for _, event := range sentEvents {
		switch event.NotificationType {
		case "NURL":
			winCount++
		case "LURL":
			lossCount++
		}
	}

	if winCount != 1 {
		t.Errorf("expected 1 win notification, got %d", winCount)
	}
	if lossCount != 1 {
		t.Errorf("expected 1 loss notification, got %d", lossCount)
	}
}

func TestHandler_HandleLoss_ExternalNotificationsEnabled(t *testing.T) {
	ctx := context.Background()

	bid := &schema.Bid{
		AuctionID:               "test-auction-id",
		AuctionConfigurationID:  123,
		AuctionConfigurationUID: "test-config-uid",
		DemandID:                "test-demand",
		Price:                   3.0,
		AuctionPriceFloor:       1.0,
		BidType:                 schema.RTBBidType,
	}

	externalWinner := &schema.ExternalWinner{
		DemandID: "external-demand",
		Price:    &[]float64{5.0}[0],
	}

	config := &auction.Config{ExternalWinNotifications: true}

	auctionResult := &notification.AuctionResult{
		AuctionID: "test-auction-id",
		Bids: []notification.Bid{
			{ID: "bid-1", Price: 3.0, LURL: "http://loss-url-1"},
			{ID: "bid-2", Price: 2.0, LURL: "http://loss-url-2"},
		},
	}

	mockRepo := &mocks.AuctionResultRepoMock{
		FindFunc: func(ctx context.Context, auctionID string) (*notification.AuctionResult, error) {
			return auctionResult, nil
		},
	}

	wg := &sync.WaitGroup{}
	wg.Add(2) // Expect 2 loss notifications

	var sentEvents []notification.Params
	sender := &mocks.SenderMock{
		SendEventFunc: func(_ context.Context, p notification.Params) {
			defer wg.Done()
			sentEvents = append(sentEvents, p)

			// Verify it's a loss notification
			if p.NotificationType != "LURL" {
				t.Errorf("expected LURL notification, got %s", p.NotificationType)
			}

			// Verify first price includes external winner
			if p.FirstPrice != 5.0 {
				t.Errorf("expected first price 5.0, got %f", p.FirstPrice)
			}
		},
	}

	handler := notification.Handler{
		AuctionResultRepo: mockRepo,
		Sender:            sender,
	}

	err := handler.HandleLoss(ctx, bid, externalWinner, config, "test-bundle", "banner")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if waitTimeout(wg, 1*time.Second) {
		t.Errorf("timeout waiting for events")
	}

	if len(sentEvents) != 2 {
		t.Errorf("expected 2 events, got %d", len(sentEvents))
	}
}

func TestHandler_HandleLoss_ExternalNotificationsDisabled(t *testing.T) {
	ctx := context.Background()

	bid := &schema.Bid{
		AuctionID:               "test-auction-id",
		AuctionConfigurationID:  123,
		AuctionConfigurationUID: "test-config-uid",
		DemandID:                "test-demand",
		Price:                   3.0,
		BidType:                 schema.RTBBidType,
	}

	externalWinner := &schema.ExternalWinner{
		DemandID: "external-demand",
		Price:    &[]float64{5.0}[0],
	}

	config := &auction.Config{ExternalWinNotifications: false}

	mockRepo := &mocks.AuctionResultRepoMock{}
	sender := &mocks.SenderMock{
		SendEventFunc: func(_ context.Context, p notification.Params) {
			t.Errorf("no events should be sent when external notifications are disabled")
		},
	}

	handler := notification.Handler{
		AuctionResultRepo: mockRepo,
		Sender:            sender,
	}

	err := handler.HandleLoss(ctx, bid, externalWinner, config, "test-bundle", "banner")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Give some time to ensure no events are sent
	time.Sleep(100 * time.Millisecond)
}

func TestHandler_HandleLoss_NonRTBBid(t *testing.T) {
	ctx := context.Background()

	bid := &schema.Bid{
		AuctionID:               "test-auction-id",
		AuctionConfigurationID:  123,
		AuctionConfigurationUID: "test-config-uid",
		DemandID:                "test-demand",
		Price:                   3.0,
		AuctionPriceFloor:       1.0,
		BidType:                 schema.CPMBidType, // Non-RTB bid
	}

	externalWinner := &schema.ExternalWinner{
		DemandID: "external-demand",
		Price:    &[]float64{5.0}[0],
	}

	config := &auction.Config{ExternalWinNotifications: true}

	auctionResult := &notification.AuctionResult{
		AuctionID: "test-auction-id",
		Bids: []notification.Bid{
			{ID: "bid-1", Price: 3.0, LURL: "http://loss-url-1"},
			{ID: "bid-2", Price: 2.0, LURL: "http://loss-url-2"},
		},
	}

	mockRepo := &mocks.AuctionResultRepoMock{
		FindFunc: func(_ context.Context, auctionID string) (*notification.AuctionResult, error) {
			return auctionResult, nil
		},
	}

	wg := &sync.WaitGroup{}
	wg.Add(2) // Expect 2 loss notifications for all bids regardless of incoming bid type

	var sentEvents []notification.Params
	sender := &mocks.SenderMock{
		SendEventFunc: func(_ context.Context, p notification.Params) {
			defer wg.Done()
			sentEvents = append(sentEvents, p)

			// Verify it's a loss notification
			if p.NotificationType != "LURL" {
				t.Errorf("expected LURL notification, got %s", p.NotificationType)
			}

			// Verify first price includes external winner
			if p.FirstPrice != 5.0 {
				t.Errorf("expected first price 5.0, got %f", p.FirstPrice)
			}
		},
	}

	handler := notification.Handler{
		AuctionResultRepo: mockRepo,
		Sender:            sender,
	}

	err := handler.HandleLoss(ctx, bid, externalWinner, config, "test-bundle", "banner")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if waitTimeout(wg, 1*time.Second) {
		t.Errorf("timeout waiting for events")
	}

	if len(sentEvents) != 2 {
		t.Errorf("expected 2 events, got %d", len(sentEvents))
	}
}

func TestHandler_HandleWin_NilConfig(t *testing.T) {
	ctx := context.Background()

	bid := &schema.Bid{
		AuctionID:               "test-auction-id",
		AuctionConfigurationID:  123,
		AuctionConfigurationUID: "test-config-uid",
		DemandID:                "test-demand",
		Price:                   5.0,
		BidType:                 schema.RTBBidType,
	}

	mockRepo := &mocks.AuctionResultRepoMock{}
	sender := &mocks.SenderMock{
		SendEventFunc: func(_ context.Context, p notification.Params) {
			t.Errorf("no events should be sent when config is nil")
		},
	}

	handler := notification.Handler{
		AuctionResultRepo: mockRepo,
		Sender:            sender,
	}

	err := handler.HandleWin(ctx, bid, nil, "test-bundle", "banner")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestHandler_HandleWin_AuctionResultNotFound(t *testing.T) {
	ctx := context.Background()

	bid := &schema.Bid{
		AuctionID:               "test-auction-id",
		AuctionConfigurationID:  123,
		AuctionConfigurationUID: "test-config-uid",
		DemandID:                "test-demand",
		Price:                   5.0,
		BidType:                 schema.RTBBidType,
	}

	config := &auction.Config{ExternalWinNotifications: true}

	mockRepo := &mocks.AuctionResultRepoMock{
		FindFunc: func(ctx context.Context, auctionID string) (*notification.AuctionResult, error) {
			return nil, nil // Auction result not found
		},
	}

	sender := &mocks.SenderMock{
		SendEventFunc: func(_ context.Context, p notification.Params) {
			t.Errorf("no events should be sent when auction result is not found")
		},
	}

	handler := notification.Handler{
		AuctionResultRepo: mockRepo,
		Sender:            sender,
	}

	err := handler.HandleWin(ctx, bid, config, "test-bundle", "banner")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestHandler_HandleLoss_NilConfig(t *testing.T) {
	ctx := context.Background()

	bid := &schema.Bid{
		AuctionID:               "test-auction-id",
		AuctionConfigurationID:  123,
		AuctionConfigurationUID: "test-config-uid",
		DemandID:                "test-demand",
		Price:                   3.0,
		BidType:                 schema.RTBBidType,
	}

	externalWinner := &schema.ExternalWinner{
		DemandID: "external-demand",
		Price:    &[]float64{5.0}[0],
	}

	mockRepo := &mocks.AuctionResultRepoMock{}
	sender := &mocks.SenderMock{
		SendEventFunc: func(_ context.Context, p notification.Params) {
			t.Errorf("no events should be sent when config is nil")
		},
	}

	handler := notification.Handler{
		AuctionResultRepo: mockRepo,
		Sender:            sender,
	}

	err := handler.HandleLoss(ctx, bid, externalWinner, nil, "test-bundle", "banner")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestHandler_HandleLoss_AuctionResultNotFound(t *testing.T) {
	ctx := context.Background()

	bid := &schema.Bid{
		AuctionID:               "test-auction-id",
		AuctionConfigurationID:  123,
		AuctionConfigurationUID: "test-config-uid",
		DemandID:                "test-demand",
		Price:                   3.0,
		BidType:                 schema.RTBBidType,
	}

	externalWinner := &schema.ExternalWinner{
		DemandID: "external-demand",
		Price:    &[]float64{5.0}[0],
	}

	config := &auction.Config{ExternalWinNotifications: true}

	mockRepo := &mocks.AuctionResultRepoMock{
		FindFunc: func(ctx context.Context, auctionID string) (*notification.AuctionResult, error) {
			return nil, nil // Auction result not found
		},
	}

	sender := &mocks.SenderMock{
		SendEventFunc: func(_ context.Context, p notification.Params) {
			t.Errorf("no events should be sent when auction result is not found")
		},
	}

	handler := notification.Handler{
		AuctionResultRepo: mockRepo,
		Sender:            sender,
	}

	err := handler.HandleLoss(ctx, bid, externalWinner, config, "test-bundle", "banner")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestHandler_HandleStats_ExternalNotificationsEnabled(t *testing.T) {
	ctx := context.Background()

	stats := schema.Stats{
		AuctionID:               "test-auction-id",
		AuctionPricefloor:       1.0,
		AuctionConfigurationID:  123,
		AuctionConfigurationUID: "test-config-uid",
		Result: schema.AuctionResult{
			Status: "SUCCESS",
			Price:  5.0,
		},
		AdUnits: []schema.AuctionAdUnitResult{
			{Price: 5.0, Status: "WIN"},
		},
	}

	config := &auction.Config{ExternalWinNotifications: true}

	// Mock repo should not be called when external notifications are enabled
	mockRepo := &mocks.AuctionResultRepoMock{
		FindFunc: func(ctx context.Context, auctionID string) (*notification.AuctionResult, error) {
			t.Errorf("Find should not be called when external notifications are enabled")
			return nil, nil
		},
	}

	sender := &mocks.SenderMock{
		SendEventFunc: func(_ context.Context, p notification.Params) {
			t.Errorf("no events should be sent when external notifications are enabled")
		},
	}

	handler := notification.Handler{
		AuctionResultRepo: mockRepo,
		Sender:            sender,
	}

	handler.HandleStats(ctx, stats, config, "test-bundle", "banner")

	// Give some time to ensure no events are sent
	time.Sleep(100 * time.Millisecond)
}

func TestHandler_HandleWin_RepoError(t *testing.T) {
	ctx := context.Background()

	bid := &schema.Bid{
		AuctionID:               "test-auction-id",
		AuctionConfigurationID:  123,
		AuctionConfigurationUID: "test-config-uid",
		DemandID:                "test-demand",
		Price:                   5.0,
		BidType:                 schema.RTBBidType,
	}

	config := &auction.Config{ExternalWinNotifications: true}

	mockRepo := &mocks.AuctionResultRepoMock{
		FindFunc: func(ctx context.Context, auctionID string) (*notification.AuctionResult, error) {
			return nil, errors.New("database error")
		},
	}

	sender := &mocks.SenderMock{
		SendEventFunc: func(_ context.Context, p notification.Params) {
			t.Errorf("no events should be sent when repo returns error")
		},
	}

	handler := notification.Handler{
		AuctionResultRepo: mockRepo,
		Sender:            sender,
	}

	err := handler.HandleWin(ctx, bid, config, "test-bundle", "banner")
	if err == nil {
		t.Errorf("expected error from repo, got nil")
	}
}

func TestHandler_HandleLoss_RepoError(t *testing.T) {
	ctx := context.Background()

	bid := &schema.Bid{
		AuctionID:               "test-auction-id",
		AuctionConfigurationID:  123,
		AuctionConfigurationUID: "test-config-uid",
		DemandID:                "test-demand",
		Price:                   3.0,
		BidType:                 schema.RTBBidType,
	}

	externalWinner := &schema.ExternalWinner{
		DemandID: "external-demand",
		Price:    &[]float64{5.0}[0],
	}

	config := &auction.Config{ExternalWinNotifications: true}

	mockRepo := &mocks.AuctionResultRepoMock{
		FindFunc: func(ctx context.Context, auctionID string) (*notification.AuctionResult, error) {
			return nil, errors.New("database error")
		},
	}

	sender := &mocks.SenderMock{
		SendEventFunc: func(_ context.Context, p notification.Params) {
			t.Errorf("no events should be sent when repo returns error")
		},
	}

	handler := notification.Handler{
		AuctionResultRepo: mockRepo,
		Sender:            sender,
	}

	err := handler.HandleLoss(ctx, bid, externalWinner, config, "test-bundle", "banner")
	if err == nil {
		t.Errorf("expected error from repo, got nil")
	}
}

func TestHandler_HandleWin_EmptyPrices(t *testing.T) {
	ctx := context.Background()

	bid := &schema.Bid{
		AuctionID:               "test-auction-id",
		AuctionConfigurationID:  123,
		AuctionConfigurationUID: "test-config-uid",
		DemandID:                "test-demand",
		Price:                   0.0, // Zero price
		AuctionPriceFloor:       0.0, // Zero floor
		BidType:                 schema.RTBBidType,
	}

	config := &auction.Config{ExternalWinNotifications: true}

	auctionResult := &notification.AuctionResult{
		AuctionID: "test-auction-id",
		Bids:      []notification.Bid{}, // Empty bids
	}

	mockRepo := &mocks.AuctionResultRepoMock{
		FindFunc: func(ctx context.Context, auctionID string) (*notification.AuctionResult, error) {
			return auctionResult, nil
		},
	}

	sender := &mocks.SenderMock{
		SendEventFunc: func(_ context.Context, p notification.Params) {
			t.Errorf("no events should be sent when no valid prices")
		},
	}

	handler := notification.Handler{
		AuctionResultRepo: mockRepo,
		Sender:            sender,
	}

	err := handler.HandleWin(ctx, bid, config, "test-bundle", "banner")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestHandler_HandleLoss_NilExternalWinner(t *testing.T) {
	ctx := context.Background()

	bid := &schema.Bid{
		AuctionID:               "test-auction-id",
		AuctionConfigurationID:  123,
		AuctionConfigurationUID: "test-config-uid",
		DemandID:                "test-demand",
		Price:                   3.0,
		AuctionPriceFloor:       1.0,
		BidType:                 schema.RTBBidType,
	}

	config := &auction.Config{ExternalWinNotifications: true}

	auctionResult := &notification.AuctionResult{
		AuctionID: "test-auction-id",
		Bids: []notification.Bid{
			{ID: "bid-1", Price: 3.0, LURL: "http://loss-url-1"},
		},
	}

	mockRepo := &mocks.AuctionResultRepoMock{
		FindFunc: func(ctx context.Context, auctionID string) (*notification.AuctionResult, error) {
			return auctionResult, nil
		},
	}

	wg := &sync.WaitGroup{}
	wg.Add(1) // Expect 1 loss notification

	sender := &mocks.SenderMock{
		SendEventFunc: func(_ context.Context, p notification.Params) {
			defer wg.Done()
			if p.NotificationType != "LURL" {
				t.Errorf("expected LURL notification, got %s", p.NotificationType)
			}
		},
	}

	handler := notification.Handler{
		AuctionResultRepo: mockRepo,
		Sender:            sender,
	}

	err := handler.HandleLoss(ctx, bid, nil, config, "test-bundle", "banner") // nil external winner
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if waitTimeout(wg, 1*time.Second) {
		t.Errorf("timeout waiting for events")
	}
}

func waitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return false // completed normally
	case <-time.After(timeout):
		return true // timed out
	}
}
