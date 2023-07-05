package store

import (
	"github.com/bidon-io/bidon-backend/internal/admin"
	"github.com/bidon-io/bidon-backend/internal/db"
)

func UserResource(dbModel *db.User) *admin.User {
	resource := userMapper{}.resource(dbModel)

	return &resource
}
