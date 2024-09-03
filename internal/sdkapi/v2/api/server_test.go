package api_test

import (
	"github.com/bidon-io/bidon-backend/internal/sdkapi/v2/api"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/v2/api/mocks"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer_GetAuction(t *testing.T) {
	e := echo.New()

	auctionHandlerMock := &mocks.HandlerMock{
		HandleFunc: func(c echo.Context) error {
			return nil
		},
	}

	srv := &api.Server{
		AuctionHandler: auctionHandlerMock,
	}

	tests := []struct {
		name        string
		handler     func(c echo.Context, _ api.GetAuctionParamsAdType) error
		method      string
		url         string
		adType      api.GetAuctionParamsAdType
		mockHandler *mocks.HandlerMock
	}{
		{
			name:        "GetAuction",
			handler:     srv.GetAuction,
			method:      http.MethodPost,
			url:         "/v2/auction/banner",
			adType:      "banner",
			mockHandler: auctionHandlerMock,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.url, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			if err := tt.handler(c, tt.adType); err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if calls := len(tt.mockHandler.HandleCalls()); calls != 1 {
				t.Errorf("expected Handle to be called once, got %d calls", calls)
			}
		})
	}
}

func TestServer_GetConfig(t *testing.T) {
	e := echo.New()

	configHandlerMock := &mocks.HandlerMock{
		HandleFunc: func(c echo.Context) error {
			return nil
		},
	}

	srv := &api.Server{
		ConfigHandler: configHandlerMock,
	}

	tests := []struct {
		name        string
		handler     func(c echo.Context) error
		method      string
		url         string
		mockHandler *mocks.HandlerMock
	}{
		{
			name:        "GetConfig",
			handler:     srv.GetConfig,
			method:      http.MethodGet,
			url:         "/v2/config",
			mockHandler: configHandlerMock,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.url, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			if err := tt.handler(c); err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if calls := len(tt.mockHandler.HandleCalls()); calls != 1 {
				t.Errorf("expected Handle to be called once, got %d calls", calls)
			}
		})
	}
}

func TestServer_PostStats(t *testing.T) {
	e := echo.New()

	statsHandlerMock := &mocks.HandlerMock{
		HandleFunc: func(c echo.Context) error {
			return nil
		},
	}

	srv := &api.Server{
		StatsHandler: statsHandlerMock,
	}

	tests := []struct {
		name        string
		handler     func(c echo.Context, _ api.PostStatsParamsAdType) error
		method      string
		url         string
		adType      api.PostStatsParamsAdType
		mockHandler *mocks.HandlerMock
	}{
		{
			name:        "PostStats",
			handler:     srv.PostStats,
			method:      http.MethodPost,
			url:         "/v2/stats/banner",
			adType:      "banner",
			mockHandler: statsHandlerMock,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.url, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			if err := tt.handler(c, tt.adType); err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if calls := len(tt.mockHandler.HandleCalls()); calls != 1 {
				t.Errorf("expected Handle to be called once, got %d calls", calls)
			}
		})
	}
}

func TestServer_PostShow(t *testing.T) {
	e := echo.New()

	showHandlerMock := &mocks.HandlerMock{
		HandleFunc: func(c echo.Context) error {
			return nil
		},
	}

	srv := &api.Server{
		ShowHandler: showHandlerMock,
	}

	tests := []struct {
		name        string
		handler     func(c echo.Context, _ api.PostShowParamsAdType) error
		method      string
		url         string
		adType      api.PostShowParamsAdType
		mockHandler *mocks.HandlerMock
	}{
		{
			name:        "PostShow",
			handler:     srv.PostShow,
			method:      http.MethodPost,
			url:         "/v2/show/banner",
			adType:      "banner",
			mockHandler: showHandlerMock,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.url, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			if err := tt.handler(c, tt.adType); err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if calls := len(tt.mockHandler.HandleCalls()); calls != 1 {
				t.Errorf("expected Handle to be called once, got %d calls", calls)
			}
		})
	}
}

func TestServer_PostClick(t *testing.T) {
	e := echo.New()

	clickHandlerMock := &mocks.HandlerMock{
		HandleFunc: func(c echo.Context) error {
			return nil
		},
	}

	srv := &api.Server{
		ClickHandler: clickHandlerMock,
	}

	tests := []struct {
		name        string
		handler     func(c echo.Context, _ api.PostClickParamsAdType) error
		method      string
		url         string
		adType      api.PostClickParamsAdType
		mockHandler *mocks.HandlerMock
	}{
		{
			name:        "PostClick",
			handler:     srv.PostClick,
			method:      http.MethodPost,
			url:         "/v2/click/banner",
			adType:      "banner",
			mockHandler: clickHandlerMock,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.url, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			if err := tt.handler(c, tt.adType); err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if calls := len(tt.mockHandler.HandleCalls()); calls != 1 {
				t.Errorf("expected Handle to be called once, got %d calls", calls)
			}
		})
	}
}

func TestServer_PostReward(t *testing.T) {
	e := echo.New()

	rewardHandlerMock := &mocks.HandlerMock{
		HandleFunc: func(c echo.Context) error {
			return nil
		},
	}

	srv := &api.Server{
		RewardHandler: rewardHandlerMock,
	}

	tests := []struct {
		name        string
		handler     func(c echo.Context, _ api.PostRewardParamsAdType) error
		method      string
		url         string
		adType      api.PostRewardParamsAdType
		mockHandler *mocks.HandlerMock
	}{
		{
			name:        "PostReward",
			handler:     srv.PostReward,
			method:      http.MethodPost,
			url:         "/v2/reward/banner",
			adType:      "banner",
			mockHandler: rewardHandlerMock,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.url, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			if err := tt.handler(c, tt.adType); err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if calls := len(tt.mockHandler.HandleCalls()); calls != 1 {
				t.Errorf("expected Handle to be called once, got %d calls", calls)
			}
		})
	}
}

func TestServer_PostWin(t *testing.T) {
	e := echo.New()

	winHandlerMock := &mocks.HandlerMock{
		HandleFunc: func(c echo.Context) error {
			return nil
		},
	}

	srv := &api.Server{
		WinHandler: winHandlerMock,
	}

	tests := []struct {
		name        string
		handler     func(c echo.Context, _ api.PostWinParamsAdType) error
		method      string
		url         string
		adType      api.PostWinParamsAdType
		mockHandler *mocks.HandlerMock
	}{
		{
			name:        "PostWin",
			handler:     srv.PostWin,
			method:      http.MethodPost,
			url:         "/v2/win/banner",
			adType:      "banner",
			mockHandler: winHandlerMock,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.url, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			if err := tt.handler(c, tt.adType); err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if calls := len(tt.mockHandler.HandleCalls()); calls != 1 {
				t.Errorf("expected Handle to be called once, got %d calls", calls)
			}
		})
	}
}
