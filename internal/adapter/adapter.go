package adapter

type Key string

type Config map[Key]map[string]any

const (
	// Sorted alphabetically
	AdmobKey      Key = "admob"
	ApplovinKey   Key = "applovin"
	BidmachineKey Key = "bidmachine"
	DTExchangeKey Key = "dtexchange"
	MetaKey       Key = "meta"
	MintegralKey  Key = "mintegral"
	MobileFuseKey Key = "mobilefuse"
	UnityAdsKey   Key = "unityads"
	VungleKey     Key = "vungle"
)

var Keys = []Key{
	AdmobKey,
	ApplovinKey,
	BidmachineKey,
	DTExchangeKey,
	MetaKey,
	MintegralKey,
	MobileFuseKey,
	UnityAdsKey,
	VungleKey,
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
