package admin

import (
	"context"
	"time"

	"github.com/bidon-io/bidon-backend/internal/admin/resource"
)

type APIKeyShort struct {
	ID             string     `json:"id"`
	LastAccessedAt *time.Time `json:"last_accessed_at"`
}

type APIKeyShortResource struct {
	*APIKeyShort
	Permissions ResourceInstancePermissions `json:"_permissions"`
}

type APIKeyFull struct {
	ID             string     `json:"id"`
	Value          string     `json:"value"`
	LastAccessedAt *time.Time `json:"last_accessed_at"`
}

type APIKeyFullResource struct {
	*APIKeyFull
	Permissions ResourceInstancePermissions `json:"_permissions"`
}

type APIKeyRepo interface {
	ListOwnedByUser(ctx context.Context, userID int64) (*resource.Collection[APIKeyShort], error)
	FindOwnedByUser(ctx context.Context, userID int64, id string) (*APIKeyFull, error)
	Create(ctx context.Context, userID int64) (*APIKeyFull, error)
	Delete(ctx context.Context, id string) error
}
type APIKeyService struct {
	repo APIKeyRepo
}

func NewAPIKeyService(store Store) *APIKeyService {
	return &APIKeyService{
		repo: store.APIKeys(),
	}
}

var apiKeyInstancePermissions = ResourceInstancePermissions{
	Update: false,
	Delete: true,
}

func (s *APIKeyService) Meta(_ context.Context, _ AuthContext) ResourceMeta {
	return ResourceMeta{
		Key: "api_key",
		Permissions: ResourcePermissions{
			Read:   true,
			Create: true,
		},
	}
}

func (s *APIKeyService) List(ctx context.Context, authCtx AuthContext) (*resource.Collection[APIKeyShortResource], error) {
	keys, err := s.repo.ListOwnedByUser(ctx, authCtx.UserID())
	if err != nil {
		return nil, err
	}

	resources := make([]APIKeyShortResource, len(keys.Items))
	for i := range keys.Items {
		resources[i] = APIKeyShortResource{
			APIKeyShort: &keys.Items[i],
			Permissions: apiKeyInstancePermissions,
		}
	}

	return &resource.Collection[APIKeyShortResource]{
		Items: resources,
		Meta:  keys.Meta,
	}, nil
}

func (s *APIKeyService) Find(ctx context.Context, authCtx AuthContext, id string) (*APIKeyFullResource, error) {
	key, err := s.repo.FindOwnedByUser(ctx, authCtx.UserID(), id)
	if err != nil {
		return nil, err
	}

	return &APIKeyFullResource{
		APIKeyFull:  key,
		Permissions: apiKeyInstancePermissions,
	}, nil
}

func (s *APIKeyService) Create(ctx context.Context, authCtx AuthContext) (*APIKeyFullResource, error) {
	key, err := s.repo.Create(ctx, authCtx.UserID())
	if err != nil {
		return nil, err
	}

	return &APIKeyFullResource{
		APIKeyFull:  key,
		Permissions: apiKeyInstancePermissions,
	}, nil
}

func (s *APIKeyService) Delete(ctx context.Context, authCtx AuthContext, id string) error {
	_, err := s.repo.FindOwnedByUser(ctx, authCtx.UserID(), id)
	if err != nil {
		return err
	}

	return s.repo.Delete(ctx, id)
}
