package adminstore

import (
	"context"
	"fmt"
	"github.com/bidon-io/bidon-backend/internal/admin"
	"github.com/bidon-io/bidon-backend/internal/admin/auth"
	"github.com/bidon-io/bidon-backend/internal/admin/resource"
	"github.com/bidon-io/bidon-backend/internal/db"
	"github.com/gofrs/uuid/v5"
	"time"
)

type APIKeyRepo struct {
	db *db.DB
}

func NewAPIKeyRepo(d *db.DB) *APIKeyRepo {
	return &APIKeyRepo{db: d}
}

func (r *APIKeyRepo) ListOwnedByUser(ctx context.Context, userID int64) (*resource.Collection[admin.APIKeyShort], error) {
	var dbKeys []db.APIKey
	db := r.db.WithContext(ctx).Where("user_id = ?", userID)

	if err := db.Find(&dbKeys).Error; err != nil {
		return nil, err
	}

	keys := make([]admin.APIKeyShort, len(dbKeys))
	for i := range dbKeys {
		keys[i] = admin.APIKeyShort{
			ID: dbKeys[i].ID.String(),
		}
		if !dbKeys[i].LastAccessedAt.IsZero() {
			keys[i].LastAccessedAt = &dbKeys[i].LastAccessedAt
		}
	}

	collection := &resource.Collection[admin.APIKeyShort]{
		Items: keys,
		Meta: resource.CollectionMeta{
			TotalCount: int64(len(keys)),
		},
	}

	return collection, nil
}

func (r *APIKeyRepo) FindOwnedByUser(ctx context.Context, userID int64, idStr string) (*admin.APIKeyFull, error) {
	var dbKey db.APIKey
	db := r.db.WithContext(ctx).Where("user_id = ?", userID)

	id, err := uuid.FromString(idStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse API key ID: %v", err)
	}

	if err := db.First(&dbKey, id).Error; err != nil {
		return nil, err
	}

	key := &admin.APIKeyFull{
		ID:    dbKey.ID.String(),
		Value: dbKey.Value,
	}
	if !dbKey.LastAccessedAt.IsZero() {
		key.LastAccessedAt = &dbKey.LastAccessedAt
	}

	return key, nil
}

func (r *APIKeyRepo) Create(ctx context.Context, userID int64) (*admin.APIKeyFull, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("failed to generate API key ID: %v", err)
	}

	value, err := auth.NewAPIKey(id)
	if err != nil {
		return nil, fmt.Errorf("failed to generate API key value: %v", err)
	}

	dbKey := &db.APIKey{
		ID:     id,
		Value:  value,
		UserID: userID,
	}

	if err := r.db.WithContext(ctx).Create(dbKey).Error; err != nil {
		return nil, err
	}

	key := &admin.APIKeyFull{
		ID:    dbKey.ID.String(),
		Value: dbKey.Value,
	}
	return key, nil
}

func (r *APIKeyRepo) Delete(ctx context.Context, idStr string) error {
	var dbKey db.APIKey

	id, err := uuid.FromString(idStr)
	if err != nil {
		return fmt.Errorf("failed to parse API key ID: %v", err)
	}

	return r.db.WithContext(ctx).Delete(&dbKey, id).Error
}

func (r *APIKeyRepo) Access(ctx context.Context, id uuid.UUID) (auth.APIKey, error) {
	var dbKey db.APIKey

	if err := r.db.WithContext(ctx).Preload("User").First(&dbKey, id).Error; err != nil {
		return auth.APIKey{}, err
	}

	key := auth.APIKey{
		ID: dbKey.ID,
		User: auth.User{
			ID:      dbKey.User.ID,
			Email:   dbKey.User.Email,
			IsAdmin: *dbKey.User.IsAdmin,
		},
		PreviousAccessedAt: dbKey.LastAccessedAt,
	}

	if err := r.db.WithContext(ctx).Model(&dbKey).Update("last_accessed_at", time.Now()).Error; err != nil {
		return auth.APIKey{}, err
	}

	return key, nil
}
