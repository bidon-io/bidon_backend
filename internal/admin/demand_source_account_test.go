package admin

import (
	"context"
	"testing"

	"github.com/bidon-io/bidon-backend/internal/adapter"
)

func Test_demandSourceAccountValidator_ValidateWithContext(t *testing.T) {
	tests := []struct {
		name         string
		attrs        *DemandSourceAccountAttrs
		demandSource *DemandSource
		wantErr      bool
	}{
		{
			"valid Amazon",
			&DemandSourceAccountAttrs{
				DemandSourceID: 1,
				Extra: map[string]any{
					"price_points_map": map[string]any{
						"price_point": map[string]any{
							"name":        "name",
							"price_point": "price_point",
							"price":       1.0,
						},
					},
				},
			},
			&DemandSource{
				DemandSourceAttrs: DemandSourceAttrs{
					ApiKey: string(adapter.AmazonKey),
				},
			},
			false,
		},
		{
			"valid Applovin",
			&DemandSourceAccountAttrs{
				DemandSourceID: 1,
				Extra: map[string]any{
					"sdk_key": "f6NN9Sc0UcaPPzoTVTUUklLmU2HmkJPumt9E_6ueD2GFTkKC8XJoR4b1J8Z-2EeaafI42GFc9tmNeOg1qFgvFy",
					"foo":     "bar",
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
			"valid Bidmachine 1",
			&DemandSourceAccountAttrs{
				DemandSourceID: 1,
				Extra: map[string]any{
					"seller_id":        "154",
					"endpoint":         "x.appbaqend.com",
					"mediation_config": []string{"foo", "bar"},
					"foo":              "bar",
				},
			},
			&DemandSource{
				DemandSourceAttrs: DemandSourceAttrs{
					ApiKey: string(adapter.BidmachineKey),
				},
			},
			false,
		},
		{
			"valid Bidmachine 2",
			&DemandSourceAccountAttrs{
				DemandSourceID: 1,
				Extra: map[string]any{
					"seller_id":        "154",
					"endpoint":         "https://x.appbaqend.com",
					"mediation_config": []string{},
					"foo":              "bar",
				},
			},
			&DemandSource{
				DemandSourceAttrs: DemandSourceAttrs{
					ApiKey: string(adapter.BidmachineKey),
				},
			},
			false,
		},
		{
			"valid BigoAds",
			&DemandSourceAccountAttrs{
				DemandSourceID: 1,
				Extra: map[string]any{
					"publisher_id": "200104",
					"endpoint":     "https://api.gov-static.tech/Ad/GetUniAdS2s?id=200104",
					"foo":          "bar",
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
			"valid GAM",
			&DemandSourceAccountAttrs{
				DemandSourceID: 1,
				Extra: map[string]any{
					"network_code": "22897248651",
					"foo":          "bar",
				},
			},
			&DemandSource{
				DemandSourceAttrs: DemandSourceAttrs{
					ApiKey: string(adapter.GAMKey),
				},
			},
			false,
		},
		{
			"valid Mintegral",
			&DemandSourceAccountAttrs{
				DemandSourceID: 1,
				Extra: map[string]any{
					"app_key":      "f674fbc5dea30e35fee264960a16e5f9",
					"publisher_id": "28686",
					"foo":          "bar",
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
			&DemandSourceAccountAttrs{
				DemandSourceID: 1,
				Extra: map[string]any{
					"account_id": "627b73f7023d420018e2f038",
					"foo":        "bar",
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
			"valid nil Extra",
			&DemandSourceAccountAttrs{
				DemandSourceID: 1,
				Extra:          nil,
			},
			&DemandSource{
				DemandSourceAttrs: DemandSourceAttrs{
					ApiKey: string(adapter.VungleKey),
				},
			},
			false,
		},
		{
			"valid adapter that has no required keys",
			&DemandSourceAccountAttrs{
				DemandSourceID: 1,
				Extra: map[string]any{
					"foo": "bar",
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
			"invalid when no keys present",
			&DemandSourceAccountAttrs{
				DemandSourceID: 1,
				Extra:          map[string]any{},
			},
			&DemandSource{
				DemandSourceAttrs: DemandSourceAttrs{
					ApiKey: string(adapter.BidmachineKey),
				},
			},
			true,
		},
		{
			"invalid when values are not string",
			&DemandSourceAccountAttrs{
				DemandSourceID: 1,
				Extra: map[string]any{
					"seller_id": 154,
					"endpoint":  154,
				},
			},
			&DemandSource{
				DemandSourceAttrs: DemandSourceAttrs{
					ApiKey: string(adapter.BidmachineKey),
				},
			},
			true,
		},
		{
			"invalid Bidmachine when endpoint is not URL nor host",
			&DemandSourceAccountAttrs{
				DemandSourceID: 1,
				Extra: map[string]any{
					"seller_id":        "154",
					"endpoint":         "not host nor url",
					"mediation_config": []string{"foo", "bar"},
				},
			},
			&DemandSource{
				DemandSourceAttrs: DemandSourceAttrs{
					ApiKey: string(adapter.BidmachineKey),
				},
			},
			true,
		},
		{
			"invalid Bidmachine when mediation_config is not slice of strings",
			&DemandSourceAccountAttrs{
				DemandSourceID: 1,
				Extra: map[string]any{
					"seller_id":        "154",
					"endpoint":         "https://api.test.bidmachine.io",
					"mediation_config": []int{1, 2, 3},
				},
			},
			&DemandSource{
				DemandSourceAttrs: DemandSourceAttrs{
					ApiKey: string(adapter.BidmachineKey),
				},
			},
			true,
		},
		{
			"invalid Bidmachine when mediation_config is nil",
			&DemandSourceAccountAttrs{
				DemandSourceID: 1,
				Extra: map[string]any{
					"seller_id":        "154",
					"endpoint":         "https://api.test.bidmachine.io",
					"mediation_config": nil,
				},
			},
			&DemandSource{
				DemandSourceAttrs: DemandSourceAttrs{
					ApiKey: string(adapter.BidmachineKey),
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
			v := &demandSourceAccountValidator{
				attrs:            tt.attrs,
				demandSourceRepo: repo,
			}
			if err := v.ValidateWithContext(context.Background()); (err != nil) != tt.wantErr {
				t.Errorf("ValidateWithContext() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
