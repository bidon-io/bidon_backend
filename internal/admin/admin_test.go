package admin_test

import "github.com/bidon-io/bidon-backend/internal/admin"

type userContext struct {
	user admin.User
}

func (c userContext) UserID() int64 {
	return c.user.ID
}

func (c userContext) IsAdmin() bool {
	return c.user.IsAdmin != nil && *c.user.IsAdmin
}

func ptr[T any](t T) *T {
	return &t
}
