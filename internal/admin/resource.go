package admin

import (
	"context"
	"errors"

	v8n "github.com/go-ozzo/ozzo-validation/v4"
)

var ErrActionForbidden = errors.New("action forbidden")

//go:generate go run -mod=mod github.com/matryer/moq@latest -out resource_mocks_test.go . ResourceManipulator AuthContext resourcePolicy resourceScope

type ResourceMeta struct {
	Key         string              `json:"key"`
	Permissions ResourcePermissions `json:"permissions"`
}

// ResourceService provides CRUD operations for managing a resource. It handles validation, authorization, and persistence.
type ResourceService[Resource, ResourceData, ResourceAttrs any] struct {
	resourceKey string

	repo   ResourceManipulator[ResourceData, ResourceAttrs]
	policy resourcePolicy[ResourceData, ResourceAttrs]

	prepareResource    func(authCtx AuthContext, data *ResourceData) Resource
	prepareCreateAttrs func(authCtx AuthContext, attrs *ResourceAttrs)
	getValidator       func(*ResourceAttrs) v8n.ValidatableWithContext
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
type resourcePolicy[Resource, ResourceAttrs any] interface {
	// getReadScope returns the resource scope for reading resources.
	getReadScope(AuthContext) resourceScope[Resource]
	// getManageScope returns the resource scope for managing resources.
	getManageScope(AuthContext) resourceScope[Resource]

	// authorizeCreate checks for authorization to create a resource.
	authorizeCreate(ctx context.Context, authCtx AuthContext, attrs *ResourceAttrs) error
	// authorizeUpdate checks for authorization to update a resource. It is called with resource from getManageScope.
	authorizeUpdate(ctx context.Context, authCtx AuthContext, resource *Resource, attrs *ResourceAttrs) error
	// authorizeDelete checks for authorization to delete a resource. It is called with resource from getManageScope.
	authorizeDelete(ctx context.Context, authCtx AuthContext, resource *Resource) error

	// permissions returns the permissions for a resource as a whole.
	permissions(AuthContext) ResourcePermissions
	// instancePermissions returns the permissions for a resource instance.
	instancePermissions(AuthContext, *Resource) ResourceInstancePermissions
}

// resourceScope handles visibility of resource.
type resourceScope[Resource any] interface {
	list(context.Context) ([]Resource, error)
	find(context.Context, int64) (*Resource, error)
}

func (s *ResourceService[Resource, ResourceData, ResourceAttrs]) Meta(_ context.Context, authContext AuthContext) ResourceMeta {
	return ResourceMeta{
		Key:         s.resourceKey,
		Permissions: s.policy.permissions(authContext),
	}
}

func (s *ResourceService[Resource, ResourceData, ResourceAttrs]) List(ctx context.Context, authCtx AuthContext) ([]Resource, error) {
	data, err := s.policy.getReadScope(authCtx).list(ctx)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, len(data))
	for i, r := range data {
		resources[i] = s.prepareResource(authCtx, &r)
	}

	return resources, nil
}

func (s *ResourceService[Resource, ResourceData, ResourceAttrs]) Find(ctx context.Context, authCtx AuthContext, id int64) (*Resource, error) {
	data, err := s.policy.getReadScope(authCtx).find(ctx, id)
	if err != nil {
		return nil, err
	}

	resource := s.prepareResource(authCtx, data)

	return &resource, nil
}

func (s *ResourceService[Resource, ResourceData, ResourceAttrs]) Create(ctx context.Context, authCtx AuthContext, attrs *ResourceAttrs) (*ResourceData, error) {
	if s.prepareCreateAttrs != nil {
		s.prepareCreateAttrs(authCtx, attrs)
	}

	if err := s.policy.authorizeCreate(ctx, authCtx, attrs); err != nil {
		return nil, err
	}

	if err := s.validate(ctx, attrs); err != nil {
		return nil, err
	}

	return s.repo.Create(ctx, attrs)
}

func (s *ResourceService[Resource, ResourceData, ResourceAttrs]) Update(ctx context.Context, authCtx AuthContext, id int64, attrs *ResourceAttrs) (*ResourceData, error) {
	scope := s.policy.getManageScope(authCtx)

	resource, err := scope.find(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := s.policy.authorizeUpdate(ctx, authCtx, resource, attrs); err != nil {
		return nil, err
	}

	if err := s.validate(ctx, attrs); err != nil {
		return nil, err
	}

	return s.repo.Update(ctx, id, attrs)
}

func (s *ResourceService[Resource, ResourceData, ResourceAttrs]) Delete(ctx context.Context, authCtx AuthContext, id int64) error {
	scope := s.policy.getManageScope(authCtx)

	resource, err := scope.find(ctx, id)
	if err != nil {
		return err
	}

	if err := s.policy.authorizeDelete(ctx, authCtx, resource); err != nil {
		return err
	}

	return s.repo.Delete(ctx, id)
}

func (s *ResourceService[Resource, ResourceData, ResourceAttrs]) validate(ctx context.Context, attrs *ResourceAttrs) error {
	if s.getValidator != nil {
		validator := s.getValidator(attrs)
		return validator.ValidateWithContext(ctx)
	}

	return nil
}
