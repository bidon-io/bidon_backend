package adminecho

import "github.com/pilagod/gorm-cursor-paginator/v2/paginator"

type Pagination struct {
	Cursor *paginator.Cursor `json:"cursor"`
}

type Meta struct {
	Total      int        `json:"total"`
	Pagination Pagination `json:"pagination"`
}

type ResourceListResponse[Resource any] struct {
	Data []Resource `json:"data"`
	Meta Meta       `json:"meta"`
}
