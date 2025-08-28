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
	ChartboostKey Key = "chartboost"
	DTExchangeKey Key = "dtexchange"
	GAMKey        Key = "gam"
	InmobiKey     Key = "inmobi"
	IronSourceKey Key = "ironsource"
	MetaKey       Key = "meta"
	MintegralKey  Key = "mintegral"
	MobileFuseKey Key = "mobilefuse"
	MolocoKey     Key = "moloco"
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
	ChartboostKey,
	DTExchangeKey,
	GAMKey,
	InmobiKey,
	IronSourceKey,
	MetaKey,
	MintegralKey,
	MobileFuseKey,
	MolocoKey,
	UnityAdsKey,
	VKAdsKey,
	VungleKey,
	YandexKey,
}

var CustomAdapters = [...]string{"max", "level_play"}

var itemExists = struct{}{}

// GetCommonAdapters returns the intersection of all the slices passed as arguments.
func GetCommonAdapters(slices ...[]Key) []Key {
	if len(slices) == 0 {
		return []Key{}
	}

	elementCount := make(map[Key]int)
	for _, slice := range slices {
		uniqueElements := make(map[Key]struct{})
		for _, elem := range slice {
			if _, ok := uniqueElements[elem]; !ok {
				elementCount[elem]++
				uniqueElements[elem] = itemExists
			}
		}
	}

	intersection := make([]Key, 0)
	for elem, count := range elementCount {
		if count == len(slices) {
			intersection = append(intersection, elem)
		}
	}

	return intersection
}

func IsDisabledForCOPPA(adapter Key) bool {
	return adapter == ApplovinKey
}
