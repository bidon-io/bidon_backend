package adminstore

import (
	"context"
	"errors"

	"github.com/bidon-io/bidon-backend/internal/db"
	"github.com/pilagod/gorm-cursor-paginator/v2/paginator"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (r *resourceRepo[Resource, ResourceAttrs, DBModel]) List(ctx context.Context, _ map[string][]string) ([]Resource, error) {
	return r.list(ctx, nil)
}

func (r *resourceRepo[Resource, ResourceAttrs, DBModel]) list(ctx context.Context, addFilters func(*gorm.DB) *gorm.DB) ([]Resource, error) {
	var dbModels []DBModel
	db := r.db.WithContext(ctx)
	for _, association := range r.associations {
		db = db.Preload(association)
	}

	if addFilters != nil {
		db = addFilters(db)
	}

	if err := db.Find(&dbModels).Error; err != nil {
		return nil, err
	}

	resources := make([]Resource, len(dbModels))
	for i := range dbModels {
		resources[i] = r.mapper.resource(&dbModels[i])
	}

	return resources, nil
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

func (r *resourceRepo[Resource, ResourceAttrs, DBModel]) listWithCursor(ctx context.Context, addFilters func(*gorm.DB) *gorm.DB, pgn *paginator.Paginator) ([]Resource, *paginator.Cursor, error) {
	if pgn == nil {
		return nil, nil, errors.New("paginator not provided")
	}

	var dbModels []DBModel

	db := r.db.WithContext(ctx)
	for _, association := range r.associations {
		db = db.Preload(association)
	}

	if addFilters != nil {
		db = addFilters(db)
	}

	db, cursor, err := pgn.Paginate(db, &dbModels)
	if err != nil {
		return nil, nil, err
	}

	resources := make([]Resource, len(dbModels))
	for i := range dbModels {
		resources[i] = r.mapper.resource(&dbModels[i])
	}

	return resources, &cursor, nil
}

func (r *resourceRepo[Resource, ResourceAttrs, DBModel]) getTotal(ctx context.Context, addFilters func(*gorm.DB) *gorm.DB) (int64, error) {
	var total int64
	db := r.db.WithContext(ctx)

	if addFilters != nil {
		db = addFilters(db)
	}

	if err := db.Model(new(DBModel)).Count(&total).Error; err != nil {
		return 0, err
	}

	return total, nil
}
