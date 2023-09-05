package admin

import (
	"context"

	v8n "github.com/go-ozzo/ozzo-validation/v4"
)

//go:generate go run -mod=mod github.com/matryer/moq@latest -out resource_mocks_test.go . ResourceManipulator AuthContext resourcePolicy resourceScope

// ResourceService provides CRUD operations for managing a resource. It handles validation, authorization, and persistence.
type ResourceService[Resource, ResourceAttrs any] struct {
	repo   ResourceManipulator[Resource, ResourceAttrs]
	policy resourcePolicy[Resource]

	getValidator func(*ResourceAttrs) v8n.ValidatableWithContext
}

// ResourceManipulator abstracts the persistence layer for a resource. Resource repositories implement this interface.
type ResourceManipulator[Resource, ResourceAttrs any] interface {
	Create(ctx context.Context, attrs *ResourceAttrs) (*Resource, error)
	Update(ctx context.Context, id int64, attrs *ResourceAttrs) (*Resource, error)
	Delete(ctx context.Context, id int64) error
}

// AuthContext defines the authorization context when accessing a resource.
type AuthContext interface {
	UserID() int64
	IsAdmin() bool
}

// resourcePolicy defines the authorization policy for a resource.
type resourcePolicy[Resource any] interface {
	scope(AuthContext) resourceScope[Resource]
}

// resourceScope handles visibility of resource.
type resourceScope[Resource any] interface {
	list(context.Context) ([]Resource, error)
	find(context.Context, int64) (*Resource, error)
}

func (s *ResourceService[Resource, ResourceAttrs]) List(ctx context.Context, authCtx AuthContext) ([]Resource, error) {
	scope := s.policy.scope(authCtx)

	return scope.list(ctx)
}

func (s *ResourceService[Resource, ResourceAttrs]) Find(ctx context.Context, authCtx AuthContext, id int64) (*Resource, error) {
	scope := s.policy.scope(authCtx)

	return scope.find(ctx, id)
}

func (s *ResourceService[Resource, ResourceAttrs]) Create(ctx context.Context, attrs *ResourceAttrs) (*Resource, error) {
	if err := s.validate(ctx, attrs); err != nil {
		return nil, err
	}

	return s.repo.Create(ctx, attrs)
}

func (s *ResourceService[Resource, ResourceAttrs]) Update(ctx context.Context, id int64, attrs *ResourceAttrs) (*Resource, error) {
	if err := s.validate(ctx, attrs); err != nil {
		return nil, err
	}

	return s.repo.Update(ctx, id, attrs)
}

func (s *ResourceService[Resource, ResourceAttrs]) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *ResourceService[Resource, ResourceAttrs]) validate(ctx context.Context, attrs *ResourceAttrs) error {
	if s.getValidator != nil {
		validator := s.getValidator(attrs)
		return validator.ValidateWithContext(ctx)
	}

	return nil
}
