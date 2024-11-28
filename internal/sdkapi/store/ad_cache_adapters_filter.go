package store

import (
	"github.com/bidon-io/bidon-backend/internal/ad"
	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
)

type AdCacheAdaptersFilter struct {
	DefaultMaxCacheCount int
	Settings             map[ad.OS]map[adapter.Key]map[ad.Type]int
}

func NewAdCacheAdaptersFilter() *AdCacheAdaptersFilter {
	defaultSettings := map[ad.OS]map[adapter.Key]map[ad.Type]int{
		ad.AndroidOS: {
			adapter.AdmobKey:      {ad.RewardedType: 1},
			adapter.ApplovinKey:   {ad.InterstitialType: 1, ad.RewardedType: 1},
			adapter.GAMKey:        {ad.RewardedType: 1},
			adapter.IronSourceKey: {ad.InterstitialType: 1, ad.RewardedType: 1, ad.BannerType: 1},
			adapter.UnityAdsKey:   {ad.InterstitialType: 1, ad.RewardedType: 1},
		},
		ad.IOSOS: {
			adapter.AdmobKey:      {ad.InterstitialType: 1, ad.RewardedType: 1},
			adapter.ApplovinKey:   {ad.InterstitialType: 1, ad.BannerType: 1},
			adapter.GAMKey:        {ad.InterstitialType: 1, ad.RewardedType: 1},
			adapter.IronSourceKey: {ad.InterstitialType: 1, ad.RewardedType: 1, ad.BannerType: 1},
			adapter.UnityAdsKey:   {ad.InterstitialType: 1, ad.RewardedType: 1},
			adapter.DTExchangeKey: {ad.InterstitialType: 1, ad.RewardedType: 1, ad.BannerType: 1},
			adapter.MintegralKey:  {ad.RewardedType: 1},
		},
	}

	return &AdCacheAdaptersFilter{
		DefaultMaxCacheCount: 3,
		Settings:             defaultSettings,
	}
}

func (f *AdCacheAdaptersFilter) Filter(OS ad.OS, adType ad.Type, adapters []adapter.Key, adCache []schema.AdCacheObject) []adapter.Key {
	demandCount := make(map[adapter.Key]int)

	for _, entry := range adCache {
		demandCount[adapter.Key(entry.DemandID)]++
	}

	filteredAdapters := make([]adapter.Key, 0)
	for _, adapterKey := range adapters {
		maxCount := f.DefaultMaxCacheCount
		if platformSettings, platformExists := f.Settings[OS]; platformExists {
			if adTypeSettings, demandExists := platformSettings[adapterKey]; demandExists {
				if specificMax, exists := adTypeSettings[adType]; exists {
					maxCount = specificMax
				}
			}
		}

		if count, found := demandCount[adapterKey]; !found || count < maxCount {
			filteredAdapters = append(filteredAdapters, adapterKey)
		}
	}

	return filteredAdapters
}
