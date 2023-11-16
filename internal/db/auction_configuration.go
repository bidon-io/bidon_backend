package db

import (
	"errors"

	"gorm.io/gorm"
)

func (ac *AuctionConfiguration) BeforeSave(tx *gorm.DB) (err error) {
	// Check if the combination of app_id, ad_type, and segment_id is already taken
	var count int64

	query := tx.Model(&AuctionConfiguration{}).
		Where("app_id = ? AND ad_type = ?", ac.AppID, ac.AdType)

	if ac.SegmentID != nil && ac.SegmentID.Valid {
		query = query.Where("segment_id = ?", ac.SegmentID.Int64)
	} else {
		query = query.Where("segment_id IS NULL")
	}

	query = query.Not(ac.ID).Count(&count)

	if query.Error != nil {
		return query.Error
	}

	if count > 0 {
		return errors.New("the combination of app_id, ad_type, and segment_id already exists")
	}

	return nil
}
