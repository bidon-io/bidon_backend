package resource

type Collection[Resource any] struct {
	Items []Resource     `json:"items"`
	Meta  CollectionMeta `json:"meta"`
}

type CollectionMeta struct {
	TotalCount int64 `json:"total_count"`
}
