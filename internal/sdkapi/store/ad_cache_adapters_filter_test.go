package store

import (
	"testing"

	"github.com/bidon-io/bidon-backend/internal/ad"
	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
	"github.com/stretchr/testify/assert"
)

func TestFilterCached(t *testing.T) {
	filter := NewAdCacheAdaptersFilter()

	tests := []struct {
		name     string
		OS       ad.OS
		adType   ad.Type
		adapters []adapter.Key
		adCache  []schema.AdCacheObject
		expected []adapter.Key
	}{
		{
			name:     "Empty adCache includes all adapters",
			OS:       ad.AndroidOS,
			adType:   ad.RewardedType,
			adapters: []adapter.Key{adapter.AdmobKey, adapter.ApplovinKey, adapter.UnityAdsKey, adapter.IronSourceKey},
			adCache:  []schema.AdCacheObject{},
			expected: []adapter.Key{adapter.AdmobKey, adapter.ApplovinKey, adapter.UnityAdsKey, adapter.IronSourceKey},
		},
		{
			name:     "Excludes adapters with exceed max_cache_count",
			OS:       ad.AndroidOS,
			adType:   ad.RewardedType,
			adapters: []adapter.Key{adapter.AdmobKey, adapter.ApplovinKey, adapter.UnityAdsKey},
			adCache: []schema.AdCacheObject{
				{DemandID: string(adapter.AdmobKey), Price: 3.37},
				{DemandID: string(adapter.UnityAdsKey), Price: 4.3},
			},
			expected: []adapter.Key{adapter.ApplovinKey},
		},
		{
			name:     "Default max_cache_count is applied for missing settings",
			OS:       ad.AndroidOS,
			adType:   ad.RewardedType,
			adapters: []adapter.Key{adapter.MintegralKey},
			adCache: []schema.AdCacheObject{
				{DemandID: string(adapter.MintegralKey), Price: 3.37},
			},
			expected: []adapter.Key{adapter.MintegralKey},
		},
		{
			name:     "Excludes adapters with exceed default max_cache_count",
			OS:       ad.AndroidOS,
			adType:   ad.RewardedType,
			adapters: []adapter.Key{adapter.MintegralKey},
			adCache: []schema.AdCacheObject{
				{DemandID: string(adapter.MintegralKey), Price: 3.37},
				{DemandID: string(adapter.MintegralKey), Price: 4.37},
				{DemandID: string(adapter.MintegralKey), Price: 5.37},
			},
			expected: []adapter.Key{},
		},
		{
			name:     "Platform-specific settings are applied",
			OS:       ad.IOSOS,
			adType:   ad.BannerType,
			adapters: []adapter.Key{adapter.IronSourceKey, adapter.AdmobKey},
			adCache: []schema.AdCacheObject{
				{DemandID: string(adapter.IronSourceKey), Price: 3.37},
				{DemandID: string(adapter.AdmobKey), Price: 4.0},
			},
			expected: []adapter.Key{adapter.AdmobKey},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filter.Filter(tt.OS, tt.adType, tt.adapters, tt.adCache)
			assert.ElementsMatch(t, tt.expected, result)
		})
	}
}
