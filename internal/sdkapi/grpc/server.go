package grpcserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/bidon-io/bidon-backend/config"
	"github.com/bidon-io/bidon-backend/internal/auction"
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
	Run(ctx context.Context, params *auction.ExecutionParams) (*auction.Response, error)
}

func (s *Server) Bid(ctx context.Context, o *v3.Openrtb) (*v3.Openrtb, error) {
	adapter := NewAuctionAdapter()
	ar, err := adapter.OpenRTBToAuctionRequest(o)
	if err != nil {
		return &v3.Openrtb{}, err2GrpcStatus(err)
	}

	app, err := s.AppFetcher.FetchCached(ctx, ar.App.Key, ar.App.Bundle)
	if err != nil {
		return &v3.Openrtb{}, err2GrpcStatus(err)
	}

	geo, err := s.GeoCoder.Lookup(ctx, ar.Device.IP)
	if err != nil {
		return &v3.Openrtb{}, fmt.Errorf("failed to lookup ip: %w", err)
	}

	logger := ctxzap.Extract(ctx)
	params := &auction.ExecutionParams{
		Req:     ar,
		AppID:   app.ID,
		Country: geo.CountryCode,
		GeoData: geo,
		Log: func(s string) {
			logger.Info(s)
		},
		LogErr: func(err error) {
			logger.Error(err.Error())
		},
	}

	result, err := s.AuctionService.Run(ctx, params)
	if err != nil {
		return &v3.Openrtb{}, err2GrpcStatus(err)
	}

	response, err := adapter.AuctionResponseToOpenRTB(result)
	if err != nil {
		return &v3.Openrtb{}, err2GrpcStatus(err)
	}

	return response, nil
}

func err2GrpcStatus(err error) error {
	httpErr := echoError(err)
	response := map[string]any{
		"error": map[string]any{
			"code":    httpErr.Code,
			"message": httpErr.Message,
		},
	}
	message, _ := json.Marshal(response)
	code := httpCodeToGRPCCode(httpErr.Code)
	return status.Error(code, string(message))
}

func echoError(err error) *echo.HTTPError {
	var x *echo.HTTPError

	if errors.As(err, &x) {
		return x
	}

	var message string
	if config.Debug() {
		message = err.Error()
	} else {
		message = http.StatusText(http.StatusInternalServerError)
	}

	return echo.NewHTTPError(http.StatusInternalServerError, message)
}

func httpCodeToGRPCCode(code int) codes.Code {
	switch code {
	case http.StatusBadRequest:
	case http.StatusUnprocessableEntity:
		return codes.InvalidArgument
	}

	return codes.Internal
}
