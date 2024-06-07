package v1

import (
	adapterstore "github.com/bidon-io/bidon-backend/internal/adapter/store"
	"github.com/bidon-io/bidon-backend/internal/auction"
	auctionstore "github.com/bidon-io/bidon-backend/internal/auction/store"
	"github.com/bidon-io/bidon-backend/internal/bidding"
	"github.com/bidon-io/bidon-backend/internal/bidding/adapters_builder"
	"github.com/bidon-io/bidon-backend/internal/notification"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/event"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/geocoder"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
	sdkapistore "github.com/bidon-io/bidon-backend/internal/sdkapi/store"
	apihandlersv1 "github.com/bidon-io/bidon-backend/internal/sdkapi/v1/apihandlers"
	apihandlersv2 "github.com/bidon-io/bidon-backend/internal/sdkapi/v2/apihandlers"
	"github.com/bidon-io/bidon-backend/internal/segment"
	"github.com/labstack/echo/v4"
)

type Router struct {
	ConfigFetcher             *auctionstore.ConfigFetcher
	AppFetcher                *sdkapistore.AppFetcher
	SegmentMatcher            *segment.Matcher
	AdUnitsMatcher            *auctionstore.AdUnitsMatcher
	NotificationHandler       notification.Handler
	GeoCoder                  *geocoder.Geocoder
	EventLogger               *event.Logger
	LineItemsMatcher          *auctionstore.LineItemsMatcher
	AdapterInitConfigsFetcher *sdkapistore.AdapterInitConfigsFetcher
	ConfigurationFetcher      *adapterstore.ConfigurationFetcher
	BiddingBuilder            *bidding.Builder
	BiddingAdaptersCfgBuilder *adapters_builder.AdaptersConfigBuilder
}

func (r *Router) RegisterRoutes(g *echo.Group) {
	auctionHandler := apihandlersv1.AuctionHandler{
		BaseHandler: &apihandlersv1.BaseHandler[schema.AuctionRequest, *schema.AuctionRequest]{
			AppFetcher:    r.AppFetcher,
			ConfigFetcher: r.ConfigFetcher,
			Geocoder:      r.GeoCoder,
		},
		SegmentMatcher: r.SegmentMatcher,
		AuctionBuilder: &auction.Builder{
			ConfigFetcher:    r.ConfigFetcher,
			LineItemsMatcher: r.LineItemsMatcher,
		},
		AuctionBuilderV2: &auction.BuilderV2{
			ConfigFetcher:  r.ConfigFetcher,
			AdUnitsMatcher: r.AdUnitsMatcher,
		},
		EventLogger: r.EventLogger,
	}
	configHandler := apihandlersv1.ConfigHandler{
		BaseHandler: &apihandlersv1.BaseHandler[schema.ConfigRequest, *schema.ConfigRequest]{
			AppFetcher:    r.AppFetcher,
			ConfigFetcher: r.ConfigFetcher,
			Geocoder:      r.GeoCoder,
		},
		SegmentMatcher:            r.SegmentMatcher,
		AdapterInitConfigsFetcher: r.AdapterInitConfigsFetcher,
		EventLogger:               r.EventLogger,
	}
	biddingHandler := apihandlersv1.BiddingHandler{
		BaseHandler: &apihandlersv1.BaseHandler[schema.BiddingRequest, *schema.BiddingRequest]{
			AppFetcher:    r.AppFetcher,
			ConfigFetcher: r.ConfigFetcher,
			Geocoder:      r.GeoCoder,
		},
		BiddingBuilder:        r.BiddingBuilder,
		AdaptersConfigBuilder: r.BiddingAdaptersCfgBuilder,
		AdUnitsMatcher:        r.AdUnitsMatcher,
		EventLogger:           r.EventLogger,
	}
	statsHandler := apihandlersv1.StatsHandler{
		BaseHandler: &apihandlersv1.BaseHandler[schema.StatsRequest, *schema.StatsRequest]{
			AppFetcher:    r.AppFetcher,
			ConfigFetcher: r.ConfigFetcher,
			Geocoder:      r.GeoCoder,
		},
		EventLogger:         r.EventLogger,
		NotificationHandler: r.NotificationHandler,
	}
	showHandler := apihandlersv2.ShowHandler{
		BaseHandler: &apihandlersv2.BaseHandler[schema.ShowRequest, *schema.ShowRequest]{
			AppFetcher:    r.AppFetcher,
			ConfigFetcher: r.ConfigFetcher,
			Geocoder:      r.GeoCoder,
		},
		EventLogger:         r.EventLogger,
		NotificationHandler: r.NotificationHandler,
	}
	clickHandler := apihandlersv2.ClickHandler{
		BaseHandler: &apihandlersv2.BaseHandler[schema.ClickRequest, *schema.ClickRequest]{
			AppFetcher:    r.AppFetcher,
			ConfigFetcher: r.ConfigFetcher,
			Geocoder:      r.GeoCoder,
		},
		EventLogger: r.EventLogger,
	}
	rewardHandler := apihandlersv2.RewardHandler{
		BaseHandler: &apihandlersv2.BaseHandler[schema.RewardRequest, *schema.RewardRequest]{
			AppFetcher:    r.AppFetcher,
			ConfigFetcher: r.ConfigFetcher,
			Geocoder:      r.GeoCoder,
		},
		EventLogger: r.EventLogger,
	}
	lossHandler := apihandlersv2.LossHandler{
		BaseHandler: &apihandlersv2.BaseHandler[schema.LossRequest, *schema.LossRequest]{
			AppFetcher:    r.AppFetcher,
			ConfigFetcher: r.ConfigFetcher,
			Geocoder:      r.GeoCoder,
		},
		EventLogger:         r.EventLogger,
		NotificationHandler: r.NotificationHandler,
	}
	winHandler := apihandlersv2.WinHandler{
		BaseHandler: &apihandlersv2.BaseHandler[schema.WinRequest, *schema.WinRequest]{
			AppFetcher:    r.AppFetcher,
			ConfigFetcher: r.ConfigFetcher,
			Geocoder:      r.GeoCoder,
		},
		EventLogger:         r.EventLogger,
		NotificationHandler: r.NotificationHandler,
	}

	g.POST("/config", configHandler.Handle)
	g.POST("/auction/:ad_type", auctionHandler.Handle)
	g.POST("/bidding/:ad_type", biddingHandler.Handle)
	g.POST("/stats/:ad_type", statsHandler.Handle)
	g.POST("/show/:ad_type", showHandler.Handle)
	g.POST("/click/:ad_type", clickHandler.Handle)
	g.POST("/reward/:ad_type", rewardHandler.Handle)
	g.POST("/loss/:ad_type", lossHandler.Handle)
	g.POST("/win/:ad_type", winHandler.Handle)

	// Legacy endpoints
	g.POST("/:ad_type/auction", auctionHandler.Handle)
	g.POST("/:ad_type/stats", statsHandler.Handle)
	g.POST("/:ad_type/show", showHandler.Handle)
	g.POST("/:ad_type/click", clickHandler.Handle)
	g.POST("/:ad_type/reward", rewardHandler.Handle)
}
