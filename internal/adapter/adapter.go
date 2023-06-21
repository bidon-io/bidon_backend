package adapter

type Key string

const (
	// Sorted alphabetically
	ApplovinKey   Key = "applovin"
	BidmachineKey Key = "bidmachine"
	DTExchangeKey Key = "dtexchange"
	UnityAdsKey   Key = "unityads"
)

func AdapterKeys() []Key {
	return []Key{
		ApplovinKey,
		BidmachineKey,
		DTExchangeKey,
		UnityAdsKey,
	}
}

func GetCommonAdapters(adapters1 []Key, adapters2 []Key) []Key {
	result := make([]Key, 0)
	hash := make(map[Key]bool)

	for _, v := range adapters1 {
		hash[v] = true
	}

	for _, v := range adapters2 {
		if _, ok := hash[v]; ok {
			result = append(result, v)
		}
	}

	return result
}
