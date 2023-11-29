package notification_test

import (
	"context"
	"fmt"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/event"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/event/engine"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bidon-io/bidon-backend/internal/notification"
	"github.com/google/go-cmp/cmp"
	"github.com/prebid/openrtb/v19/openrtb3"
)

func TestHandler_SendEvent(t *testing.T) {
	// Create a test context and input data
	ctx := context.Background()

	// Create a mock HTTP server to handle the LURL request
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("unexpected method: %s", r.Method)
		}

		params := r.URL.Query()
		if diff := cmp.Diff("request-1", params.Get("id")); diff != "" {
			t.Errorf("mismatched id (-want, +got)\n%s", diff)
		}
		if diff := cmp.Diff("4.56", params.Get("auction_price")); diff != "" {
			t.Errorf("mismatched winprice (-want, +got)\n%s", diff)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	bid := notification.Bid{
		RequestID: "request-1",
		ImpID:     "imp-1",
		Price:     1.23,
		LURL:      fmt.Sprintf("%s/lurl?auction_price=${AUCTION_PRICE}&id=${AUCTION_ID}", server.URL),
	}

	// Create a Handler instance with the mock HTTP client and server
	sender := notification.EventSender{
		HttpClient:  server.Client(),
		EventLogger: &event.Logger{Engine: &engine.Log{}},
	}

	p := notification.Params{
		NotificationType: "LURL",
		URL:              bid.LURL,
		Bid:              bid,
		Reason:           openrtb3.LossBelowAuctionFloor,
		FirstPrice:       4.56,
		SecondPrice:      3.00,
	}

	// Call the SendNotificationEvent method with the test context and input data
	sender.SendEvent(ctx, p)
}
