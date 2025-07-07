package grpcserver

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"

	"github.com/bidon-io/bidon-backend/internal/ad"
	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/auction"
	auctionmocks "github.com/bidon-io/bidon-backend/internal/auction/mocks"
	"github.com/bidon-io/bidon-backend/internal/bidding"
	"github.com/bidon-io/bidon-backend/internal/bidding/adapters"
	"github.com/bidon-io/bidon-backend/internal/sdkapi"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/event"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/event/engine"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/geocoder"
	handlersmocks "github.com/bidon-io/bidon-backend/internal/sdkapi/v2/apihandlers/mocks"
	"github.com/bidon-io/bidon-backend/internal/segment"
	segmentmocks "github.com/bidon-io/bidon-backend/internal/segment/mocks"
	v3 "github.com/bidon-io/bidon-backend/pkg/proto/com/iabtechlab/openrtb/v3"
)

type serverParams struct {
	app             *sdkapi.App
	geodata         *geocoder.GeoData
	segments        *[]segment.Segment
	auctionConfig   *auction.Config
	adUnits         *[]auction.AdUnit
	demandResponses *[]adapters.DemandResponse
}

func defaultServerParams() *serverParams {
	adUnits := DefaultAdUnits()
	config := DefaultAuctionConfig()
	demandResp := BuildDemandResponses(adUnits)
	return &serverParams{
		app:             &sdkapi.App{ID: 1},
		geodata:         &geocoder.GeoData{CountryCode: "US"},
		segments:        &[]segment.Segment{DefaultSegment()},
		auctionConfig:   &config,
		adUnits:         &adUnits,
		demandResponses: &demandResp,
	}
}

func buildServer(p *serverParams) *Server {
	adUnitsMatcher := &auctionmocks.AdUnitsMatcherMock{
		MatchCachedFunc: func(ctx context.Context, params *auction.BuildParams) ([]auction.AdUnit, error) {
			return *p.adUnits, nil
		},
	}
	appFetcher := &handlersmocks.AppFetcherMock{
		FetchCachedFunc: func(ctx context.Context, appKey string, appBundle string) (sdkapi.App, error) {
			return *p.app, nil
		},
	}
	gcoder := &handlersmocks.GeocoderMock{
		LookupFunc: func(ctx context.Context, ipString string) (geocoder.GeoData, error) {
			return *p.geodata, nil
		},
	}
	configFetcher := &handlersmocks.ConfigFetcherMock{
		MatchFunc: func(ctx context.Context, appID int64, adType ad.Type, segmentID int64, version string) (*auction.Config, error) {
			return p.auctionConfig, nil
		},
		FetchByUIDCachedFunc: func(ctx context.Context, appId int64, key string, aucUID string) *auction.Config {
			return p.auctionConfig
		},
	}
	segmentFetcher := &segmentmocks.FetcherMock{
		FetchCachedFunc: func(ctx context.Context, appID int64) ([]segment.Segment, error) {
			return *p.segments, nil
		},
	}
	segmentMatcher := &segment.Matcher{
		Fetcher: segmentFetcher,
	}
	adapterKeysFetcher := &auctionmocks.AdapterKeysFetcherMock{
		FetchEnabledAdapterKeysFunc: func(ctx context.Context, appID int64, keys []adapter.Key) ([]adapter.Key, error) {
			return keys, nil
		},
	}

	biddingAdaptersConfigBuilder := &auctionmocks.BiddingAdaptersConfigBuilderMock{
		BuildFunc: func(ctx context.Context, appID int64, adapterKeys []adapter.Key, adUnitsMap *auction.AdUnitsMap) (adapter.ProcessedConfigsMap, error) {
			return adapter.ProcessedConfigsMap{
				adapter.MetaKey: map[string]any{
					"app_id":     "123",
					"app_secret": "123",
					"seller_id":  "123",
					"tag_id":     "123",
				},
			}, nil
		},
	}
	biddingBuilder := &auctionmocks.BiddingBuilderMock{
		HoldAuctionFunc: func(ctx context.Context, params *bidding.BuildParams) (bidding.AuctionResult, error) {
			return bidding.AuctionResult{
				Bids: *p.demandResponses,
			}, nil
		},
	}
	auctionBuilderV2 := &auction.Builder{
		AdUnitsMatcher:               adUnitsMatcher,
		BiddingBuilder:               biddingBuilder,
		BiddingAdaptersConfigBuilder: biddingAdaptersConfigBuilder,
	}
	auctionService := &auction.Service{
		AdapterKeysFetcher: adapterKeysFetcher,
		ConfigFetcher:      configFetcher,
		AuctionBuilder:     auctionBuilderV2,
		SegmentMatcher:     segmentMatcher,
		EventLogger:        &event.Logger{Engine: &engine.Log{}},
	}

	return NewServer(auctionService, appFetcher, gcoder)
}

func TestServer_Bid(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		buildServer func() *Server
		input       func() *v3.Openrtb
		want        func() *v3.Openrtb
		wantErr     bool
		errorMsg    string
	}{
		{
			name: "valid request",
			buildServer: func() *Server {
				return buildServer(defaultServerParams())
			},
			input: func() *v3.Openrtb {
				return NewRequestBuilder().Build()
			},
			want: func() *v3.Openrtb {
				return NewResponseBuilder().Build()
			},
			wantErr: false,
		},
		{
			name: "empty response",
			buildServer: func() *Server {
				p := defaultServerParams()
				p.adUnits = &[]auction.AdUnit{}
				return buildServer(p)
			},
			input: func() *v3.Openrtb {
				return NewRequestBuilder().Build()
			},
			want: func() *v3.Openrtb {
				return NewResponseBuilder().WithAdUnits([]auction.AdUnit{}).Build()
			},
			wantErr: false,
		},
		{
			name: "without auction config",
			buildServer: func() *Server {
				p := defaultServerParams()
				p.auctionConfig = nil
				return buildServer(p)
			},
			input: func() *v3.Openrtb {
				return NewRequestBuilder().Build()
			},
			want: func() *v3.Openrtb {
				return &v3.Openrtb{}
			},
			wantErr:  true,
			errorMsg: "rpc error: code = InvalidArgument desc = {\"error\":{\"code\":422,\"message\":\"Invalid Auction Key\"}}",
		},
		{
			name: "all networks are disabled",
			buildServer: func() *Server {
				p := defaultServerParams()
				s := buildServer(p)

				as := s.AuctionService.(*auction.Service)
				as.AdapterKeysFetcher = &auctionmocks.AdapterKeysFetcherMock{
					FetchEnabledAdapterKeysFunc: func(ctx context.Context, appID int64, keys []adapter.Key) ([]adapter.Key, error) {
						return []adapter.Key{}, nil
					},
				}

				return s
			},
			input: func() *v3.Openrtb {
				return NewRequestBuilder().Build()
			},
			want: func() *v3.Openrtb {
				return &v3.Openrtb{}
			},
			wantErr:  true,
			errorMsg: "rpc error: code = InvalidArgument desc = {\"error\":{\"code\":422,\"message\":\"No ads found\"}}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.buildServer()
			got, err := s.Bid(ctx, tt.input())
			if (err != nil) != tt.wantErr {
				t.Errorf("Server.Bid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err.Error() != tt.errorMsg {
				t.Errorf("Server.Bid() error = %v, want %v", err.Error(), tt.errorMsg)
				return
			}

			if diff := cmp.Diff(tt.want(), got, protocmp.Transform()); diff != "" {
				t.Errorf("Server.Bid() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
