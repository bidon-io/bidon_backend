package admin

import (
	"context"
	"errors"

	"github.com/pilagod/gorm-cursor-paginator/v2/paginator"
)

//go:generate go run -mod=mod github.com/matryer/moq@latest -out resource_scopes_mocks_test.go . AllResourceQuerier OwnedResourceQuerier OwnedOrSharedResourceQuerier

// AllResourceQuerier defines the interface for querying all resources from persistence layer.
// Resource repositories implement this interface.
type AllResourceQuerier[Resource any] interface {
	List(context.Context, map[string][]string) ([]Resource, error)
	Find(ctx context.Context, id int64) (*Resource, error)
}

// AllResourceQuerierWithCursor defines the interface for querying all resources from persistence layer with cursor pagination.
// Resource repositories implement this interface.
// TODO: merge with AllResourceQuerier once all repositories are updated to use cursor pagination.
type AllResourceQuerierWithCursor[Resource any] interface {
	ListWithCursor(ctx context.Context, qParams map[string][]string) ([]Resource, *paginator.Cursor, error)
	GetTotal(ctx context.Context, qParams map[string][]string) (int64, error)
}

// OwnedResourceQuerier defines the interface for querying resources owned by a user from persistence layer.
// Resource repositories implement this interface.
type OwnedResourceQuerier[Resource any] interface {
	ListOwnedByUser(ctx context.Context, userID int64, qParams map[string][]string) ([]Resource, error)
	FindOwnedByUser(ctx context.Context, userID, id int64) (*Resource, error)
}

// OwnedResourceQuerierWithCursor defines the interface for querying resources owned by a user from persistence layer with cursor pagination.
// Resource repositories implement this interface.
// TODO: merge with OwnedResourceQuerier once all repositories are updated to use cursor pagination.
type OwnedResourceQuerierWithCursor[Resource any] interface {
	ListOwnedByUserWithCursor(ctx context.Context, userID int64, qParams map[string][]string) ([]Resource, *paginator.Cursor, error)
	GetTotalOwnedByUser(ctx context.Context, userID int64, qParams map[string][]string) (int64, error)
}

// OwnedOrSharedResourceQuerier defines the interface for querying resources owned by a user or shared with a user from persistence layer.
// Resource repositories implement this interface.
type OwnedOrSharedResourceQuerier[Resource any] interface {
	ListOwnedByUserOrShared(ctx context.Context, userID int64) ([]Resource, error)
	FindOwnedByUserOrShared(ctx context.Context, userID, id int64) (*Resource, error)
}

// publicResourceScope is a resource scope that allows access to all resources.
type publicResourceScope[Resource any] struct {
	repo AllResourceQuerier[Resource]
}

func (s *publicResourceScope[Resource]) list(ctx context.Context, _ map[string][]string) ([]Resource, error) {
	return s.repo.List(ctx, nil)
}

func (s *publicResourceScope[Resource]) listWithCursor(ctx context.Context, _ map[string][]string) ([]Resource, *paginator.Cursor, error) {
	return nil, nil, errors.New("cursor pagination not supported")
}

func (s *publicResourceScope[Resource]) find(ctx context.Context, id int64) (*Resource, error) {
	return s.repo.Find(ctx, id)
}

// privateResourceScope is a resource scope that allows access to all resources only for admin users.
type privateResourceScope[Resource any] struct {
	repo AllResourceQuerier[Resource]

	authCtx AuthContext
}

func (s *privateResourceScope[Resource]) list(ctx context.Context, qParams map[string][]string) ([]Resource, error) {
	if s.authCtx.IsAdmin() {
		return s.repo.List(ctx, qParams)
	}

	return nil, errors.New("unauthorized")
}

func (s *privateResourceScope[Resource]) listWithCursor(ctx context.Context, _ map[string][]string) ([]Resource, *paginator.Cursor, error) {
	return nil, nil, errors.New("cursor pagination not supported")
}

func (s *privateResourceScope[Resource]) find(ctx context.Context, id int64) (*Resource, error) {
	if s.authCtx.IsAdmin() {
		if id == 0 {
			return nil, nil
		} else {
			return s.repo.Find(ctx, id)
		}
	}

	return nil, errors.New("unauthorized")
}

// ownedResourceScope is a resource scope that allows access to resources owned by a user and all resources for admin users.
type ownedResourceScope[Resource any] struct {
	repo interface {
		AllResourceQuerier[Resource]
		OwnedResourceQuerier[Resource]
	}

	authCtx AuthContext
}

func (s *ownedResourceScope[Resource]) list(ctx context.Context, qParams map[string][]string) ([]Resource, error) {
	if s.authCtx.IsAdmin() {
		return s.repo.List(ctx, qParams)
	}

	return s.repo.ListOwnedByUser(ctx, s.authCtx.UserID(), qParams)
}

func (s *ownedResourceScope[Resource]) listWithCursor(ctx context.Context, _ map[string][]string) ([]Resource, *paginator.Cursor, error) {
	return nil, nil, errors.New("cursor pagination not supported")
}

func (s *ownedResourceScope[Resource]) find(ctx context.Context, id int64) (*Resource, error) {
	if s.authCtx.IsAdmin() {
		return s.repo.Find(ctx, id)
	}

	return s.repo.FindOwnedByUser(ctx, s.authCtx.UserID(), id)
}

type ownedResourceScopeWithCursor[Resource any] struct {
	*ownedResourceScope[Resource]
	repo interface {
		AllResourceQuerierWithCursor[Resource]
		OwnedResourceQuerierWithCursor[Resource]
	}
}

func (s *ownedResourceScopeWithCursor[Resource]) listWithCursor(ctx context.Context, qParams map[string][]string) ([]Resource, *paginator.Cursor, error) {
	if s.authCtx.IsAdmin() {
		return s.repo.ListWithCursor(ctx, qParams)
	}

	return s.repo.ListOwnedByUserWithCursor(ctx, s.authCtx.UserID(), qParams)
}

// ownedOrSharedResourceScope is a resource scope that allows access to resources owned by a user or shared with a user and all resources for admin users.
type ownedOrSharedResourceScope[Resource any] struct {
	repo interface {
		AllResourceQuerier[Resource]
		OwnedOrSharedResourceQuerier[Resource]
	}

	authCtx AuthContext
}

func (s *ownedOrSharedResourceScope[Resource]) list(ctx context.Context, qParams map[string][]string) ([]Resource, error) {
	if s.authCtx.IsAdmin() {
		return s.repo.List(ctx, qParams)
	}

	return s.repo.ListOwnedByUserOrShared(ctx, s.authCtx.UserID())
}

func (s *ownedOrSharedResourceScope[Resource]) listWithCursor(ctx context.Context, _ map[string][]string) ([]Resource, *paginator.Cursor, error) {
	return nil, nil, errors.New("cursor pagination not supported")
}

func (s *ownedOrSharedResourceScope[Resource]) find(ctx context.Context, id int64) (*Resource, error) {
	if s.authCtx.IsAdmin() {
		return s.repo.Find(ctx, id)
	}

	return s.repo.FindOwnedByUserOrShared(ctx, s.authCtx.UserID(), id)
}
