package admin

import (
	"context"

	v8n "github.com/go-ozzo/ozzo-validation/v4"
)

// ResourceService wraps ResourceRepo and provides additional functionality like validations, etc.
// It is used to separate business logic from the store. This is what consumers should use.
type ResourceService[Resource, ResourceAttrs any] struct {
	ResourceRepo[Resource, ResourceAttrs]

	getValidator func(*ResourceAttrs) v8n.ValidatableWithContext
}

// ResourceRepo provides CRUD operations for managing a resource.
//
//go:generate go run -mod=mod github.com/matryer/moq@latest -out mocks_test.go . ResourceRepo
type ResourceRepo[Resource, ResourceAttrs any] interface {
	List(ctx context.Context) ([]Resource, error)
	Find(ctx context.Context, id int64) (*Resource, error)
	Create(ctx context.Context, attrs *ResourceAttrs) (*Resource, error)
	Update(ctx context.Context, id int64, attrs *ResourceAttrs) (*Resource, error)
	Delete(ctx context.Context, id int64) error
}

func (s *ResourceService[Resource, ResourceAttrs]) Create(ctx context.Context, attrs *ResourceAttrs) (*Resource, error) {
	if err := s.validate(ctx, attrs); err != nil {
		return nil, err
	}

	return s.ResourceRepo.Create(ctx, attrs)
}

func (s *ResourceService[Resource, ResourceAttrs]) Update(ctx context.Context, id int64, attrs *ResourceAttrs) (*Resource, error) {
	if err := s.validate(ctx, attrs); err != nil {
		return nil, err
	}

	return s.ResourceRepo.Update(ctx, id, attrs)
}

func (s *ResourceService[Resource, ResourceAttrs]) validate(ctx context.Context, attrs *ResourceAttrs) error {
	if s.getValidator != nil {
		validator := s.getValidator(attrs)
		return validator.ValidateWithContext(ctx)
	}

	return nil
}
