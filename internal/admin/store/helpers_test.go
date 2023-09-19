package adminstore

import (
	"github.com/bidon-io/bidon-backend/internal/admin"
	"github.com/bidon-io/bidon-backend/internal/db"
)

func AppResource(dbModel *db.App) admin.App {
	resource := appMapper{}.resource(dbModel)

	return resource
}

func AppAttrsWithId(dbModel *db.App) admin.App {
	resource := admin.App{
		ID:       dbModel.ID,
		AppAttrs: appMapper{}.resourceAttrs(dbModel),
	}

	return resource
}

func AppDemandProfileResource(dbModel *db.AppDemandProfile) admin.AppDemandProfile {
	return appDemandProfileMapper{}.resource(dbModel)
}

func SegmentResource(dbModel *db.Segment) *admin.Segment {
	resource := segmentMapper{}.resource(dbModel)

	return &resource
}

func SegmentAttrsWithId(dbModel *db.Segment) *admin.Segment {
	resource := admin.Segment{
		ID:           dbModel.ID,
		SegmentAttrs: segmentMapper{}.resourceAttrs(dbModel),
	}

	return &resource
}

func UserResource(dbModel *db.User) *admin.User {
	resource := userMapper{}.resource(dbModel)

	return &resource
}

func DemandSourceAccountResource(dbModel *db.DemandSourceAccount) admin.DemandSourceAccount {
	resource := demandSourceAccountMapper{}.resource(dbModel)

	return resource
}

func DemandSourceAccountAttrsWithId(dbModel *db.DemandSourceAccount) admin.DemandSourceAccount {
	resource := admin.DemandSourceAccount{
		ID:                       dbModel.ID,
		DemandSourceAccountAttrs: demandSourceAccountMapper{}.resourceAttrs(dbModel),
	}

	return resource
}

func DemandSourceResource(dbModel *db.DemandSource) *admin.DemandSource {
	resource := demandSourceMapper{}.resource(dbModel)

	return &resource
}
