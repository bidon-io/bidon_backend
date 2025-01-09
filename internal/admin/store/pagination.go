package adminstore

import (
	"github.com/pilagod/gorm-cursor-paginator/v2/paginator"
	"strconv"
)

const (
	MaxPaginationLimit     = 100
	DefaultPaginationLimit = 12
)

func queryToPaginator(qParams map[string][]string) *paginator.Paginator {
	p := paginator.New(
		paginator.WithKeys("ID"),
		paginator.WithOrder(paginator.ASC),
	)

	limit := DefaultPaginationLimit
	if v, ok := qParams["limit"]; ok {
		parsedLimit, err := strconv.Atoi(v[0])
		if err == nil && parsedLimit > 0 && parsedLimit <= MaxPaginationLimit {
			limit = parsedLimit
		}
	}
	p.SetLimit(limit)

	if v, ok := qParams["cursor_before"]; ok {
		p.SetBeforeCursor(v[0])
	}

	if v, ok := qParams["cursor_after"]; ok {
		p.SetAfterCursor(v[0])
	}

	return p
}
