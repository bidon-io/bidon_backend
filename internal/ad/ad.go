package ad

type Type string

const (
	UnknownType      Type = ""
	BannerType       Type = "banner"
	InterstitialType Type = "interstitial"
	RewardedType     Type = "rewarded"
)

type Format string

const (
	EmptyFormat       Format = ""
	BannerFormat      Format = "BANNER"
	LeaderboardFormat Format = "LEADERBOARD"
	MRECFormat        Format = "MREC"
	AdaptiveFormat    Format = "ADAPTIVE"
)

var BannerFormats = []Format{BannerFormat, LeaderboardFormat, MRECFormat, AdaptiveFormat}
