package adminstore

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/bidon-io/bidon-backend/internal/admin/resource"
	"github.com/bidon-io/bidon-backend/internal/db"
)

// A resourceRepo is a generic basic repository for API resources that map directly to database models.
// It implements [admin.ResourceRepo]
type resourceRepo[Resource, ResourceAttrs, DBModel any] struct {
	db           *db.DB
	mapper       resourceMapper[Resource, ResourceAttrs, DBModel]
	associations []string
}

// resourceMapper maps resources with corresponding DB model, and vice versa
type resourceMapper[Resource, ResourceAttrs, DBModel any] interface {
	dbModel(*ResourceAttrs, int64) *DBModel
	resource(*DBModel) Resource
}

func (r *resourceRepo[Resource, ResourceAttrs, DBModel]) List(ctx context.Context, _ map[string][]string) (*resource.Collection[Resource], error) {
	return r.list(ctx, nil, nil)
}

func (r *resourceRepo[Resource, ResourceAttrs, DBModel]) list(ctx context.Context, addFilters func(*gorm.DB) *gorm.DB, pgn *Pagination[DBModel]) (*resource.Collection[Resource], error) {
	var dbModels []DBModel
	db := r.db.WithContext(ctx)
	for _, association := range r.associations {
		db = db.Preload(association)
	}

	if addFilters != nil {
		db = addFilters(db)
	}

	if pgn != nil {
		if err := pgn.Paginate(&dbModels, db); err != nil {
			return nil, err
		}
		dbModels = pgn.Rows
	} else {
		if err := db.Find(&dbModels).Error; err != nil {
			return nil, err
		}
	}

	resources := make([]Resource, len(dbModels))
	for i := range dbModels {
		resources[i] = r.mapper.resource(&dbModels[i])
	}

	totalCount := int64(len(dbModels))
	if pgn != nil {
		totalCount = pgn.TotalCount
	}

	collection := &resource.Collection[Resource]{
		Items: resources,
		Meta: resource.CollectionMeta{
			TotalCount: totalCount,
		},
	}

	return collection, nil
}

func (r *resourceRepo[Resource, ResourceAttrs, DBModel]) Find(ctx context.Context, id int64) (*Resource, error) {
	return r.find(ctx, id, nil)
}

func (r *resourceRepo[Resource, ResourceAttrs, DBModel]) find(ctx context.Context, id int64, addFilters func(*gorm.DB) *gorm.DB) (*Resource, error) {
	var dbModel DBModel
	db := r.db.WithContext(ctx)
	for _, association := range r.associations {
		db = db.Preload(association)
	}

	if addFilters != nil {
		db = addFilters(db)
	}

	if err := db.First(&dbModel, id).Error; err != nil {
		return nil, err
	}

	resource := r.mapper.resource(&dbModel)
	return &resource, nil
}

func (r *resourceRepo[Resource, ResourceAttrs, DBModel]) Create(ctx context.Context, attrs *ResourceAttrs) (*Resource, error) {
	dbModel := r.mapper.dbModel(attrs, 0)

	if err := r.db.WithContext(ctx).Create(dbModel).Error; err != nil {
		return nil, err
	}

	resource := r.mapper.resource(dbModel)
	return &resource, nil
}

func (r *resourceRepo[Resource, ResourceAttrs, DBModel]) Update(ctx context.Context, id int64, attrs *ResourceAttrs) (*Resource, error) {
	dbModel := r.mapper.dbModel(attrs, id)

	if err := r.db.WithContext(ctx).Model(dbModel).Where("id = ?", id).Clauses(clause.Returning{}).Updates(&dbModel).Error; err != nil {
		return nil, err
	}

	resource := r.mapper.resource(dbModel)
	return &resource, nil
}

func (r *resourceRepo[Resource, ResourceAttrs, DBModel]) Delete(ctx context.Context, id int64) error {
	var dbModel DBModel

	return r.db.WithContext(ctx).Delete(&dbModel, id).Error
}
