package store

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type dbModel[Resource, ResourceAttrs, T any] interface {
	initFromResourceAttrs(*ResourceAttrs)
	toResource() Resource
	*T
}

// A resourceRepo is a generic basic repository for API resources that map directly to database models.
// It implements [admin.ResourceRepo]
type resourceRepo[Resource, ResourceAttrs, DBModel any, DBModelP dbModel[Resource, ResourceAttrs, DBModel]] struct {
	DB *gorm.DB
}

func (r *resourceRepo[Resource, ResourceAttrs, DBModel, DBModelP]) List(ctx context.Context) ([]Resource, error) {
	var dbRecords []DBModel
	if err := r.DB.WithContext(ctx).Find(&dbRecords).Error; err != nil {
		return nil, err
	}

	resources := make([]Resource, len(dbRecords))
	for i := range dbRecords {
		resources[i] = DBModelP(&dbRecords[i]).toResource()
	}

	return resources, nil
}

func (r *resourceRepo[Resource, ResourceAttrs, DBModel, DBModelP]) Find(ctx context.Context, id int64) (*Resource, error) {
	var dbRecord DBModel
	if err := r.DB.WithContext(ctx).First(&dbRecord, id).Error; err != nil {
		return nil, err
	}

	resource := DBModelP(&dbRecord).toResource()
	return &resource, nil
}

func (r *resourceRepo[Resource, ResourceAttrs, DBModel, DBModelP]) Create(ctx context.Context, attrs *ResourceAttrs) (*Resource, error) {
	var dbRecord DBModel
	DBModelP(&dbRecord).initFromResourceAttrs(attrs)

	if err := r.DB.WithContext(ctx).Create(&dbRecord).Error; err != nil {
		return nil, err
	}

	resource := DBModelP(&dbRecord).toResource()
	return &resource, nil
}

func (r *resourceRepo[Resource, ResourceAttrs, DBModel, DBModelP]) Update(ctx context.Context, id int64, attrs *ResourceAttrs) (*Resource, error) {
	var dbRecord DBModel
	DBModelP(&dbRecord).initFromResourceAttrs(attrs)

	if err := r.DB.WithContext(ctx).Model(&dbRecord).Where("id = ?", id).Clauses(clause.Returning{}).Updates(&dbRecord).Error; err != nil {
		return nil, err
	}

	resource := DBModelP(&dbRecord).toResource()
	return &resource, nil
}

func (r *resourceRepo[Resource, ResourceAttrs, DBModel, DBModelP]) Delete(ctx context.Context, id int64) error {
	var dbRecord DBModel

	return r.DB.WithContext(ctx).Delete(&dbRecord, id).Error
}
