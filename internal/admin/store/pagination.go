package adminstore

import (
	"strconv"

	"gorm.io/gorm"
)

func NewPagination[DBModel any](limit, page int, sort string) *Pagination[DBModel] {
	return &Pagination[DBModel]{
		Limit: limit,
		Page:  page,
		Sort:  sort,
	}
}

func PaginationFromQueryParams[DBModel any](queryParams map[string][]string) *Pagination[DBModel] {
	var limit int
	if v, ok := queryParams["limit"]; ok {
		if parsed, err := strconv.Atoi(v[0]); err == nil {
			limit = parsed
		}
	}

	var page int
	if v, ok := queryParams["page"]; ok {
		if parsed, err := strconv.Atoi(v[0]); err == nil {
			page = parsed
		}
	}

	// If page and limit are not provided, return nil
	if page == 0 && limit == 0 {
		return nil
	}

	var sort string
	if v, ok := queryParams["sort"]; ok {
		sort = v[0]
	}

	return NewPagination[DBModel](limit, page, sort)
}

type Pagination[DBModel any] struct {
	Limit      int       `json:"limit,omitempty;query:limit"`
	Page       int       `json:"page,omitempty;query:page"`
	Sort       string    `json:"sort,omitempty;query:sort"`
	TotalCount int64     `json:"total_count"`
	Rows       []DBModel `json:"rows"`
}

func (p *Pagination[DBModel]) Paginate(records *[]DBModel, db *gorm.DB) error {
	db.Model(new(DBModel)).Count(&p.TotalCount)

	if err := db.Scopes(p.PageScope()).Find(records).Error; err != nil {
		return err
	}
	p.Rows = *records

	return nil
}

func (p *Pagination[DBModel]) PageScope() func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(p.GetOffset()).Limit(p.GetLimit()).Order(p.GetSort())
	}
}

func (p *Pagination[DBModel]) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p *Pagination[DBModel]) GetLimit() int {
	if p.Limit == 0 {
		p.Limit = 100
	}
	return p.Limit
}

func (p *Pagination[DBModel]) GetPage() int {
	if p.Page == 0 {
		p.Page = 1
	}
	return p.Page
}

func (p *Pagination[DBModel]) GetSort() string {
	switch p.Sort {
	case "desc":
		return "id desc"
	default:
		return "id asc"
	}
	return p.Sort
}
