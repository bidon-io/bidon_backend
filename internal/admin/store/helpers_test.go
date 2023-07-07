package store

import (
	"github.com/bidon-io/bidon-backend/internal/admin"
	"github.com/bidon-io/bidon-backend/internal/db"
)

func AppResource(dbModel *db.App) *admin.App {
	resource := appMapper{}.resource(dbModel)

	return &resource
}

func SegmentResource(dbModel *db.Segment) *admin.Segment {
	resource := segmentMapper{}.resource(dbModel)

	return &resource
}

func UserResource(dbModel *db.User) *admin.User {
	resource := userMapper{}.resource(dbModel)

	return &resource
}
