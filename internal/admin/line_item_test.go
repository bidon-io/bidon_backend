package admin

import (
	"context"
	"testing"

	"github.com/bidon-io/bidon-backend/internal/adapter"
)

func Test_lineItemAttrsValidator_ValidateWithContext(t *testing.T) {
	tests := []struct {
		name                string
		attrs               *LineItemAttrs
		demandSourceAccount *DemandSourceAccount
		wantErr             bool
	}{
		{
			"valid Admob",
			&LineItemAttrs{
				AccountID: 1,
				Extra: map[string]any{
					"ad_unit_id": "ca-app-pub-3940256099942544/5224354917",
					"foo":        "bar",
				},
			},
			&DemandSourceAccount{
				DemandSource: DemandSource{
					DemandSourceAttrs: DemandSourceAttrs{
						ApiKey: string(adapter.AdmobKey),
					},
				},
			},
			false,
		},
		{
			"valid Amazon",
			&LineItemAttrs{
				AccountID: 1,
				Extra: map[string]any{
					"slot_uuid": "26069ec0-4151-4194-a181-7a0017efdf28",
					"format":    "VIDEO",
				},
			},
			&DemandSourceAccount{
				DemandSource: DemandSource{
					DemandSourceAttrs: DemandSourceAttrs{
						ApiKey: string(adapter.AmazonKey),
					},
				},
			},
			false,
		},
		{
			"valid Applovin",
			&LineItemAttrs{
				AccountID: 1,
				Extra: map[string]any{
					"zone_id": "bd706625a42e3413",
					"foo":     "bar",
				},
			},
			&DemandSourceAccount{
				DemandSource: DemandSource{
					DemandSourceAttrs: DemandSourceAttrs{
						ApiKey: string(adapter.ApplovinKey),
					},
				},
			},
			false,
		},
		{
			"valid BigoAds",
			&LineItemAttrs{
				AccountID: 1,
				Extra: map[string]any{
					"slot_id": "10175763-10078514",
					"foo":     "bar",
				},
			},
			&DemandSourceAccount{
				DemandSource: DemandSource{
					DemandSourceAttrs: DemandSourceAttrs{
						ApiKey: string(adapter.BigoAdsKey),
					},
				},
			},
			false,
		},
		{
			"valid DT Exchange",
			&LineItemAttrs{
				AccountID: 1,
				Extra: map[string]any{
					"spot_id": "1187213",
					"foo":     "bar",
				},
			},
			&DemandSourceAccount{
				DemandSource: DemandSource{
					DemandSourceAttrs: DemandSourceAttrs{
						ApiKey: string(adapter.DTExchangeKey),
					},
				},
			},
			false,
		},
		{
			"valid GAM",
			&LineItemAttrs{
				AccountID: 1,
				Extra: map[string]any{
					"ad_unit_id": "/111/Bidon/Interstitials/0.4 USD",
					"foo":        "bar",
				},
			},
			&DemandSourceAccount{
				DemandSource: DemandSource{
					DemandSourceAttrs: DemandSourceAttrs{
						ApiKey: string(adapter.GAMKey),
					},
				},
			},
			false,
		},
		{
			"valid Meta",
			&LineItemAttrs{
				AccountID: 1,
				Extra: map[string]any{
					"placement_id": "767803077426274_1212622446277666",
					"foo":          "bar",
				},
			},
			&DemandSourceAccount{
				DemandSource: DemandSource{
					DemandSourceAttrs: DemandSourceAttrs{
						ApiKey: string(adapter.MetaKey),
					},
				},
			},
			false,
		},
		{
			"valid Unity",
			&LineItemAttrs{
				AccountID: 1,
				Extra: map[string]any{
					"placement_id": "bidon_rv_43",
					"foo":          "bar",
				},
			},
			&DemandSourceAccount{
				DemandSource: DemandSource{
					DemandSourceAttrs: DemandSourceAttrs{
						ApiKey: string(adapter.UnityAdsKey),
					},
				},
			},
			false,
		},
		{
			"valid Vungle",
			&LineItemAttrs{
				AccountID: 1,
				Extra: map[string]any{
					"placement_id": "BANNER_TEST-8066185",
					"foo":          "bar",
				},
			},
			&DemandSourceAccount{
				DemandSource: DemandSource{
					DemandSourceAttrs: DemandSourceAttrs{
						ApiKey: string(adapter.VungleKey),
					},
				},
			},
			false,
		},
		{
			"valid MobileFuse",
			&LineItemAttrs{
				AccountID: 1,
				Extra: map[string]any{
					"placement_id": "938186",
					"foo":          "bar",
				},
			},
			&DemandSourceAccount{
				DemandSource: DemandSource{
					DemandSourceAttrs: DemandSourceAttrs{
						ApiKey: string(adapter.MobileFuseKey),
					},
				},
			},
			false,
		},
		{
			"valid Mintegral",
			&LineItemAttrs{
				AccountID: 1,
				Extra: map[string]any{
					"placement_id": "938186",
					"unit_id":      "2567735",
					"foo":          "bar",
				},
			},
			&DemandSourceAccount{
				DemandSource: DemandSource{
					DemandSourceAttrs: DemandSourceAttrs{
						ApiKey: string(adapter.MintegralKey),
					},
				},
			},
			false,
		},
		{
			"valid nil Extra",
			&LineItemAttrs{
				AccountID: 1,
				Extra:     nil,
			},
			&DemandSourceAccount{
				DemandSource: DemandSource{
					DemandSourceAttrs: DemandSourceAttrs{
						ApiKey: string(adapter.MintegralKey),
					},
				},
			},
			false,
		},
		{
			"valid adapter that has no required keys",
			&LineItemAttrs{
				AccountID: 1,
				Extra: map[string]any{
					"foo": "bar",
				},
			},
			&DemandSourceAccount{
				DemandSource: DemandSource{
					DemandSourceAttrs: DemandSourceAttrs{
						ApiKey: string(adapter.BidmachineKey),
					},
				},
			},
			false,
		},
		{
			"invalid when no keys present",
			&LineItemAttrs{
				AccountID: 1,
				Extra:     map[string]any{},
			},
			&DemandSourceAccount{
				DemandSource: DemandSource{
					DemandSourceAttrs: DemandSourceAttrs{
						ApiKey: string(adapter.MintegralKey),
					},
				},
			},
			true,
		},
		{
			"invalid when values are not string",
			&LineItemAttrs{
				AccountID: 1,
				Extra: map[string]any{
					"placement_id": 938186,
					"ad_unit_id":   2567735,
				},
			},
			&DemandSourceAccount{
				DemandSource: DemandSource{
					DemandSourceAttrs: DemandSourceAttrs{
						ApiKey: string(adapter.MintegralKey),
					},
				},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &DemandSourceAccountRepoMock{
				FindFunc: func(ctx context.Context, id int64) (*DemandSourceAccount, error) {
					if id != tt.attrs.AccountID {
						t.Errorf("Find() got = %v, want %v", id, tt.attrs.AccountID)
					}
					return tt.demandSourceAccount, nil
				},
			}
			v := &lineItemAttrsValidator{
				attrs:                   tt.attrs,
				demandSourceAccountRepo: repo,
			}
			if err := v.ValidateWithContext(context.Background()); (err != nil) != tt.wantErr {
				t.Errorf("ValidateWithContext() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
