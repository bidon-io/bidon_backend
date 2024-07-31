package adapter

type Key string

type (
	RawConfigsMap       map[Key]Config
	ProcessedConfigsMap map[Key]map[string]any
)

type Config struct {
	AccountExtra map[string]any
	AppData      map[string]any
}

const (
	// Sorted alphabetically
	AdmobKey      Key = "admob"
	AmazonKey     Key = "amazon"
	ApplovinKey   Key = "applovin"
	BidmachineKey Key = "bidmachine"
	BigoAdsKey    Key = "bigoads"
	DTExchangeKey Key = "dtexchange"
	GAMKey        Key = "gam"
	InmobiKey     Key = "inmobi"
	IronSourceKey Key = "ironsource"
	MetaKey       Key = "meta"
	MintegralKey  Key = "mintegral"
	MobileFuseKey Key = "mobilefuse"
	UnityAdsKey   Key = "unityads"
	VKAdsKey      Key = "vkads"
	VungleKey     Key = "vungle"
	YandexKey     Key = "yandex"
)

var Keys = []Key{
	AdmobKey,
	AmazonKey,
	ApplovinKey,
	BidmachineKey,
	BigoAdsKey,
	DTExchangeKey,
	GAMKey,
	InmobiKey,
	IronSourceKey,
	MetaKey,
	MintegralKey,
	MobileFuseKey,
	UnityAdsKey,
	VKAdsKey,
	VungleKey,
	YandexKey,
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
