package auth

import (
	"context"

	"github.com/alexedwards/scs/v2"
)

type sessionAuthContext struct {
	sm  *scs.SessionManager
	ctx context.Context
}

func (c *sessionAuthContext) UserID() int64 {
	return c.sm.GetInt64(c.ctx, "user_id")
}

func (c *sessionAuthContext) IsAdmin() bool {
	return c.sm.GetBool(c.ctx, "is_admin")
}
