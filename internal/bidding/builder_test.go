package bidding_test

import (
	"context"
	"errors"
	"testing"

	"github.com/bidon-io/bidon-backend/internal/ad"
	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/auction"
	"github.com/bidon-io/bidon-backend/internal/bidding"
	"github.com/bidon-io/bidon-backend/internal/device"
)

type mockConfigMatcher struct {
	config *auction.Config
	err    error
}

func (m *mockConfigMatcher) Match(ctx context.Context, appID int64, adType ad.Type, segmentID int64) (*auction.Config, error) {
	return m.config, m.err
}

func TestBuilder_Build(t *testing.T) {
	tests := []struct {
		name           string
		configMatcher  bidding.ConfigMatcher
		buildParams    *bidding.BuildParams
		expectedResult *bidding.DemandResponse
		expectedError  error
	}{
		{
			name: "successful build",
			configMatcher: &mockConfigMatcher{
				config: &auction.Config{
					Rounds: []auction.RoundConfig{
						{
							ID:      "ROUND_1",
							Demands: []adapter.Key{adapter.ApplovinKey, adapter.BidmachineKey},
							Timeout: 15000,
						},
						{
							ID:      "ROUND_2",
							Demands: []adapter.Key{adapter.UnityAdsKey},
							Bidding: []adapter.Key{adapter.BidmachineKey},
							Timeout: 15000,
						},
					},
				},
				err: nil,
			},
			buildParams: &bidding.BuildParams{
				AppID:      1,
				AdType:     ad.BannerType,
				AdFormat:   ad.BannerFormat,
				DeviceType: device.PhoneType,
				Adapters:   []adapter.Key{adapter.ApplovinKey},
				SegmentID:  1,
			},
			expectedResult: &bidding.DemandResponse{
				Price: 0,
			},
			expectedError: nil,
		},
		{
			name: "config matcher error",
			configMatcher: &mockConfigMatcher{
				config: nil,
				err:    errors.New("config matcher error"),
			},
			buildParams: &bidding.BuildParams{
				AppID:      1,
				AdType:     ad.BannerType,
				AdFormat:   ad.BannerFormat,
				DeviceType: device.PhoneType,
				Adapters:   []adapter.Key{adapter.ApplovinKey},
				SegmentID:  1,
			},
			expectedResult: nil,
			expectedError:  errors.New("config matcher error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := &bidding.Builder{
				ConfigMatcher: tt.configMatcher,
			}

			result, err := builder.Build(context.Background(), tt.buildParams)

			if err != nil && tt.expectedError == nil {
				t.Errorf("unexpected error: %v", err)
			}

			if err == nil && tt.expectedError != nil {
				t.Errorf("expected error: %v, but got nil", tt.expectedError)
			}

			if result != nil && tt.expectedResult == nil {
				t.Errorf("unexpected result: %v", result)
			}

			if result == nil && tt.expectedResult != nil {
				t.Errorf("expected result: %v, but got nil", tt.expectedResult)
			}
		})
	}
}
