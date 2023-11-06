package schema

type BidType string

const (
	EmptyBidType BidType = ""
	RTBBidType   BidType = "RTB"
	CPMBidType   BidType = "CPM"
)

func (b BidType) String() string {
	return string(b)
}
