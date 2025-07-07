package v2

import (
	"github.com/labstack/echo/v4"

	adapterstore "github.com/bidon-io/bidon-backend/internal/adapter/store"
	"github.com/bidon-io/bidon-backend/internal/auction"
	auctionstore "github.com/bidon-io/bidon-backend/internal/auction/store"
	"github.com/bidon-io/bidon-backend/internal/bidding"
	"github.com/bidon-io/bidon-backend/internal/notification"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/event"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/geocoder"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
	sdkapistore "github.com/bidon-io/bidon-backend/internal/sdkapi/store"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/v2/api"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/v2/apihandlers"
	"github.com/bidon-io/bidon-backend/internal/segment"
)

type Router struct {
	ConfigFetcher             *auctionstore.ConfigFetcher
	AppFetcher                *sdkapistore.AppFetcher
	SegmentMatcher            *segment.Matcher
	AdUnitsMatcher            *auctionstore.AdUnitsMatcher
	NotificationHandler       notification.Handler
	GeoCoder                  *geocoder.Geocoder
	EventLogger               *event.Logger
	AdapterInitConfigsFetcher *sdkapistore.AdapterInitConfigsFetcher
	ConfigurationFetcher      *adapterstore.ConfigurationFetcher
	BiddingBuilder            *bidding.Builder
	AuctionService            *auction.Service
	AdUnitLookup              *sdkapistore.AdUnitLookup
}

func (r *Router) RegisterRoutes(g *echo.Group) {
	auctionHandler := apihandlers.AuctionHandler{
		BaseHandler: &apihandlers.BaseHandler[schema.AuctionRequest, *schema.AuctionRequest]{
			AppFetcher:    r.AppFetcher,
			ConfigFetcher: r.ConfigFetcher,
			Geocoder:      r.GeoCoder,
		},
		AuctionService: r.AuctionService,
	}
	statsHandler := apihandlers.StatsHandler{
		BaseHandler: &apihandlers.BaseHandler[schema.StatsRequest, *schema.StatsRequest]{
			AppFetcher:    r.AppFetcher,
			ConfigFetcher: r.ConfigFetcher,
			Geocoder:      r.GeoCoder,
		},
		EventLogger:         r.EventLogger,
		NotificationHandler: r.NotificationHandler,
	}
	configHandler := apihandlers.ConfigHandler{
		BaseHandler: &apihandlers.BaseHandler[schema.ConfigRequest, *schema.ConfigRequest]{
			AppFetcher:    r.AppFetcher,
			ConfigFetcher: r.ConfigFetcher,
			Geocoder:      r.GeoCoder,
		},
		SegmentMatcher:            r.SegmentMatcher,
		AdapterInitConfigsFetcher: r.AdapterInitConfigsFetcher,
		EventLogger:               r.EventLogger,
	}
	showHandler := apihandlers.ShowHandler{
		BaseHandler: &apihandlers.BaseHandler[schema.ShowRequest, *schema.ShowRequest]{
			AppFetcher:    r.AppFetcher,
			ConfigFetcher: r.ConfigFetcher,
			Geocoder:      r.GeoCoder,
		},
		EventLogger:         r.EventLogger,
		NotificationHandler: r.NotificationHandler,
		AdUnitLookup:        r.AdUnitLookup,
	}
	clickHandler := apihandlers.ClickHandler{
		BaseHandler: &apihandlers.BaseHandler[schema.ClickRequest, *schema.ClickRequest]{
			AppFetcher:    r.AppFetcher,
			ConfigFetcher: r.ConfigFetcher,
			Geocoder:      r.GeoCoder,
		},
		EventLogger: r.EventLogger,
	}
	rewardHandler := apihandlers.RewardHandler{
		BaseHandler: &apihandlers.BaseHandler[schema.RewardRequest, *schema.RewardRequest]{
			AppFetcher:    r.AppFetcher,
			ConfigFetcher: r.ConfigFetcher,
			Geocoder:      r.GeoCoder,
		},
		EventLogger: r.EventLogger,
	}
	lossHandler := apihandlers.LossHandler{
		BaseHandler: &apihandlers.BaseHandler[schema.LossRequest, *schema.LossRequest]{
			AppFetcher:    r.AppFetcher,
			ConfigFetcher: r.ConfigFetcher,
			Geocoder:      r.GeoCoder,
		},
		EventLogger:         r.EventLogger,
		NotificationHandler: r.NotificationHandler,
	}
	winHandler := apihandlers.WinHandler{
		BaseHandler: &apihandlers.BaseHandler[schema.WinRequest, *schema.WinRequest]{
			AppFetcher:    r.AppFetcher,
			ConfigFetcher: r.ConfigFetcher,
			Geocoder:      r.GeoCoder,
		},
		EventLogger:         r.EventLogger,
		NotificationHandler: r.NotificationHandler,
	}

	api.RegisterHandlers(g, &api.Server{
		AuctionHandler: &auctionHandler,
		ClickHandler:   &clickHandler,
		ConfigHandler:  &configHandler,
		LossHandler:    &lossHandler,
		StatsHandler:   &statsHandler,
		ShowHandler:    &showHandler,
		RewardHandler:  &rewardHandler,
		WinHandler:     &winHandler,
	})
}
