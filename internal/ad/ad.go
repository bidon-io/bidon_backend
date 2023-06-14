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

func (f Format) IsBannerFormat() bool {
	return f == BannerFormat ||
		f == LeaderboardFormat ||
		f == MRECFormat ||
		f == AdaptiveFormat
}
