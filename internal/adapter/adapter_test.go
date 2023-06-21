package adapter_test

import (
	"testing"

	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/google/go-cmp/cmp"
)

func TestAdapter_GetCommonAdapters(t *testing.T) {
	testCases := []struct {
		name   string
		params [][]adapter.Key
		want   []adapter.Key
	}{
		{
			name:   "One element overlap",
			params: [][]adapter.Key{{adapter.UnityAdsKey, adapter.BidmachineKey}, {adapter.BidmachineKey}},
			want:   []adapter.Key{adapter.BidmachineKey},
		},
		{
			name:   "Two elements overlap",
			params: [][]adapter.Key{{adapter.UnityAdsKey, adapter.BidmachineKey}, {adapter.UnityAdsKey, adapter.BidmachineKey}},
			want:   []adapter.Key{adapter.UnityAdsKey, adapter.BidmachineKey},
		},
		{
			name:   "No elements overlap",
			params: [][]adapter.Key{{adapter.UnityAdsKey, adapter.BidmachineKey}, {adapter.DTExchangeKey, adapter.ApplovinKey}},
			want:   []adapter.Key{},
		},
		{
			name:   "Empty input",
			params: [][]adapter.Key{{}, {}},
			want:   []adapter.Key{},
		},
	}

	for _, tC := range testCases {
		got := adapter.GetCommonAdapters(tC.params[0], tC.params[1])

		if diff := cmp.Diff(tC.want, got); diff != "" {
			t.Errorf("builder.Build -> %+v mismatch \n(-want, +got)\n%s", tC.name, diff)
		}
	}
}
