package admin

import (
	"context"
	"testing"

	"github.com/bidon-io/bidon-backend/internal/adapter"
)

func Test_appDemandProfileAttrsValidator_ValidateWithContext(t *testing.T) {
	tests := []struct {
		name         string
		attrs        *AppDemandProfileAttrs
		demandSource *DemandSource
		wantErr      bool
	}{
		{
			"valid Admob",
			&AppDemandProfileAttrs{
				DemandSourceID: 1,
				Data: map[string]any{
					"app_id": "ca-app-pub-7174718190807894~2828867145",
					"foo":    "bar",
				},
			},
			&DemandSource{
				DemandSourceAttrs: DemandSourceAttrs{
					ApiKey: string(adapter.AdmobKey),
				},
			},
			false,
		},
		{
			"valid BigoAds",
			&AppDemandProfileAttrs{
				DemandSourceID: 1,
				Data: map[string]any{
					"app_id": "10182906",
					"foo":    "bar",
				},
			},
			&DemandSource{
				DemandSourceAttrs: DemandSourceAttrs{
					ApiKey: string(adapter.BigoAdsKey),
				},
			},
			false,
		},
		{
			"valid DT Exchange",
			&AppDemandProfileAttrs{
				DemandSourceID: 1,
				Data: map[string]any{
					"app_id": "147573",
					"foo":    "bar",
				},
			},
			&DemandSource{
				DemandSourceAttrs: DemandSourceAttrs{
					ApiKey: string(adapter.DTExchangeKey),
				},
			},
			false,
		},
		{
			"valid Mintegral",
			&AppDemandProfileAttrs{
				DemandSourceID: 1,
				Data: map[string]any{
					"app_id": "223817",
					"foo":    "bar",
				},
			},
			&DemandSource{
				DemandSourceAttrs: DemandSourceAttrs{
					ApiKey: string(adapter.MintegralKey),
				},
			},
			false,
		},
		{
			"valid Vungle",
			&AppDemandProfileAttrs{
				DemandSourceID: 1,
				Data: map[string]any{
					"app_id": "64afd303f5edf073b3bd24a7",
					"foo":    "bar",
				},
			},
			&DemandSource{
				DemandSourceAttrs: DemandSourceAttrs{
					ApiKey: string(adapter.VungleKey),
				},
			},
			false,
		},
		{
			"valid Meta",
			&AppDemandProfileAttrs{
				DemandSourceID: 1,
				Data: map[string]any{
					"app_id":     "767803077426274",
					"app_secret": "2457f83ef75ab78d249d13df9b74f45b",
					"foo":        "bar",
				},
			},
			&DemandSource{
				DemandSourceAttrs: DemandSourceAttrs{
					ApiKey: string(adapter.MetaKey),
				},
			},
			false,
		},
		{
			"valid Unity",
			&AppDemandProfileAttrs{
				DemandSourceID: 1,
				Data: map[string]any{
					"game_id": "3716005",
				},
			},
			&DemandSource{
				DemandSourceAttrs: DemandSourceAttrs{
					ApiKey: string(adapter.UnityAdsKey),
				},
			},
			false,
		},
		{
			"valid nil Data",
			&AppDemandProfileAttrs{
				DemandSourceID: 1,
				Data:           nil,
			},
			&DemandSource{
				DemandSourceAttrs: DemandSourceAttrs{
					ApiKey: string(adapter.UnityAdsKey),
				},
			},
			false,
		},
		{
			"valid adapter that has no required keys",
			&AppDemandProfileAttrs{
				DemandSourceID: 1,
				Data: map[string]any{
					"foo": "bar",
				},
			},
			&DemandSource{
				DemandSourceAttrs: DemandSourceAttrs{
					ApiKey: string(adapter.ApplovinKey),
				},
			},
			false,
		},
		{
			"invalid when no keys present",
			&AppDemandProfileAttrs{
				DemandSourceID: 1,
				Data:           map[string]any{},
			},
			&DemandSource{
				DemandSourceAttrs: DemandSourceAttrs{
					ApiKey: string(adapter.MetaKey),
				},
			},
			true,
		},
		{
			"invalid when values are not string",
			&AppDemandProfileAttrs{
				DemandSourceID: 1,
				Data: map[string]any{
					"app_id":     123,
					"app_secret": 321,
				},
			},
			&DemandSource{
				DemandSourceAttrs: DemandSourceAttrs{
					ApiKey: string(adapter.MetaKey),
				},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &DemandSourceRepoMock{
				FindFunc: func(ctx context.Context, id int64) (*DemandSource, error) {
					if id != tt.attrs.DemandSourceID {
						t.Errorf("Find() got = %v, want %v", id, tt.attrs.DemandSourceID)
					}
					return tt.demandSource, nil
				},
			}
			v := &appDemandProfileAttrsValidator{
				attrs:            tt.attrs,
				demandSourceRepo: repo,
			}
			if err := v.ValidateWithContext(context.Background()); (err != nil) != tt.wantErr {
				t.Errorf("ValidateWithContext() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
