package grpcserver

import (
	"context"
	"fmt"
	"log"

	"github.com/bidon-io/bidon-backend/internal/auctionv2"
	"github.com/bidon-io/bidon-backend/internal/sdkapi"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/geocoder"
	v3 "github.com/bidon-io/bidon-backend/pkg/proto/com/iabtechlab/openrtb/v3"
	pb "github.com/bidon-io/bidon-backend/pkg/proto/org/bidon/proto/v1"
)

type Server struct {
	pb.UnimplementedBiddingServiceServer
	AuctionService AuctionService
	AppFetcher     AppFetcher
	GeoCoder       Geocoder
}

func NewServer(auctionService AuctionService, appFetcher AppFetcher, geoCoder Geocoder) *Server {
	return &Server{
		AuctionService: auctionService,
		AppFetcher:     appFetcher,
		GeoCoder:       geoCoder,
	}
}

//go:generate go run -mod=mod github.com/matryer/moq@latest -out mocks/mocks.go -pkg mocks . AppFetcher AuctionService Geocoder

type AppFetcher interface {
	FetchCached(ctx context.Context, appKey, appBundle string) (sdkapi.App, error)
}

type Geocoder interface {
	Lookup(ctx context.Context, ipString string) (geocoder.GeoData, error)
}

type AuctionService interface {
	Run(ctx context.Context, params *auctionv2.ExecutionParams) (*auctionv2.Response, error)
}

func (s *Server) Bid(ctx context.Context, o *v3.Openrtb) (*v3.Openrtb, error) {
	adapter := NewAuctionAdapter()
	ar, err := adapter.OpenRTBToAuctionRequest(o)
	if err != nil {
		return &v3.Openrtb{}, err
	}

	app, err := s.AppFetcher.FetchCached(ctx, ar.App.Key, ar.App.Bundle)
	if err != nil {
		return &v3.Openrtb{}, err
	}

	geo, err := s.GeoCoder.Lookup(ctx, ar.Device.IP)
	if err != nil {
		return &v3.Openrtb{}, fmt.Errorf("failed to lookup ip: %w", err)
	}

	params := &auctionv2.ExecutionParams{
		Req:     ar,
		AppID:   app.ID,
		Country: geo.CountryCode,
		GeoData: geo,
		Log: func(s string) {
			log.Print(s)
		},
		LogErr: func(err error) {
			log.Print(err)
		},
	}

	result, err := s.AuctionService.Run(ctx, params)
	if err != nil {
		return &v3.Openrtb{}, err
	}

	response, err := adapter.AuctionResponseToOpenRTB(result)
	if err != nil {
		return &v3.Openrtb{}, err
	}

	return response, nil
}
