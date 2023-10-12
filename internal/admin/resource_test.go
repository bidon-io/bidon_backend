package admin

import (
	"context"
	"errors"
	"testing"

	v8n "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/go-cmp/cmp"
)

type TestResource struct {
	*TestResourceData
	Permissions ResourceInstancePermissions
}

type TestResourceData struct {
	ID int64
	TestResourceAttrs
}

type TestResourceAttrs struct {
	Name string
}

func TestResourceService_List(t *testing.T) {
	data := []TestResourceData{
		{ID: 1, TestResourceAttrs: TestResourceAttrs{Name: "test1"}},
		{ID: 2, TestResourceAttrs: TestResourceAttrs{Name: "test2"}},
	}

	s := ResourceService[TestResource, TestResourceData, TestResourceAttrs]{
		policy: &resourcePolicyMock[TestResourceData, TestResourceAttrs]{
			getReadScopeFunc: func(authCtx AuthContext) resourceScope[TestResourceData] {
				return &resourceScopeMock[TestResourceData]{
					listFunc: func(ctx context.Context) ([]TestResourceData, error) {
						return data, nil
					},
				}
			},
		},
		prepareResource: func(authCtx AuthContext, data *TestResourceData) TestResource {
			return TestResource{
				TestResourceData: data,
				Permissions: ResourceInstancePermissions{
					Update: true,
					Delete: true,
				},
			}
		},
	}

	want := make([]TestResource, len(data))
	for i := range data {
		want[i] = s.prepareResource(nil, &data[i])
	}

	resources, _ := s.List(context.Background(), nil)
	if diff := cmp.Diff(want, resources); diff != "" {
		t.Errorf("List() mismatch (-want +got):\n%s", diff)
	}
}

func TestResourceService_Find(t *testing.T) {
	data := &TestResourceData{
		ID:                1,
		TestResourceAttrs: TestResourceAttrs{Name: "test1"},
	}

	s := ResourceService[TestResource, TestResourceData, TestResourceAttrs]{
		policy: &resourcePolicyMock[TestResourceData, TestResourceAttrs]{
			getReadScopeFunc: func(authCtx AuthContext) resourceScope[TestResourceData] {
				return &resourceScopeMock[TestResourceData]{
					findFunc: func(ctx context.Context, id int64) (*TestResourceData, error) {
						if id != data.ID {
							t.Errorf("Find() got %d, want %d", id, data.ID)
						}
						return data, nil
					},
				}
			},
		},
		prepareResource: func(authCtx AuthContext, data *TestResourceData) TestResource {
			return TestResource{
				TestResourceData: data,
				Permissions: ResourceInstancePermissions{
					Update: true,
					Delete: true,
				},
			}
		},
	}

	want := s.prepareResource(nil, data)

	resource, _ := s.Find(context.Background(), nil, data.ID)
	if diff := cmp.Diff(&want, resource); diff != "" {
		t.Errorf("Find() mismatch (-want +got):\n%s", diff)
	}
}

type testValidator struct {
	err error
}

func (v *testValidator) ValidateWithContext(_ context.Context) error {
	return v.err
}

func TestResourceService_Create(t *testing.T) {
	want := &TestResourceData{
		ID:                1,
		TestResourceAttrs: TestResourceAttrs{Name: "test1"},
	}

	s := ResourceService[TestResource, TestResourceData, TestResourceAttrs]{
		repo: &ResourceManipulatorMock[TestResourceData, TestResourceAttrs]{
			CreateFunc: func(ctx context.Context, attrs *TestResourceAttrs) (*TestResourceData, error) {
				if diff := cmp.Diff(&want.TestResourceAttrs, attrs); diff != "" {
					t.Errorf("Create() mismatch (-want +got):\n%s", diff)
				}
				return want, nil
			},
		},
		policy: &resourcePolicyMock[TestResourceData, TestResourceAttrs]{
			authorizeCreateFunc: func(ctx context.Context, authCtx AuthContext, attrs *TestResourceAttrs) error {
				if diff := cmp.Diff(&want.TestResourceAttrs, attrs); diff != "" {
					t.Errorf("authorizeCreate() mismatch (-want +got):\n%s", diff)
				}
				return nil
			},
		},
		getValidator: func(attrs *TestResourceAttrs) v8n.ValidatableWithContext {
			if diff := cmp.Diff(&want.TestResourceAttrs, attrs); diff != "" {
				t.Errorf("getValidator() mismatch (-want +got):\n%s", diff)
			}

			return &testValidator{}
		},
	}

	resource, _ := s.Create(context.Background(), nil, &want.TestResourceAttrs)
	if diff := cmp.Diff(want, resource); diff != "" {
		t.Errorf("Create() mismatch (-want +got):\n%s", diff)
	}
}

func TestResourceService_Create_validationError(t *testing.T) {
	testResourceAttrs := &TestResourceAttrs{Name: "test1"}

	repoMock := &ResourceManipulatorMock[TestResourceData, TestResourceAttrs]{}
	s := ResourceService[TestResource, TestResourceData, TestResourceAttrs]{
		repo: repoMock,
		policy: &resourcePolicyMock[TestResourceData, TestResourceAttrs]{
			authorizeCreateFunc: func(ctx context.Context, authCtx AuthContext, attrs *TestResourceAttrs) error {
				if diff := cmp.Diff(testResourceAttrs, attrs); diff != "" {
					t.Errorf("authorizeCreate() mismatch (-want +got):\n%s", diff)
				}
				return nil
			},
		},
		getValidator: func(attrs *TestResourceAttrs) v8n.ValidatableWithContext {
			if diff := cmp.Diff(testResourceAttrs, attrs); diff != "" {
				t.Errorf("getValidator() mismatch (-want +got):\n%s", diff)
			}

			return &testValidator{err: errors.New("validation error")}
		},
	}

	resource, err := s.Create(context.Background(), nil, testResourceAttrs)
	if err == nil {
		t.Errorf("Create() got %v, want error", resource)
	} else if resource != nil {
		t.Errorf("Create() got %v, want nil resource with error", resource)
	}

	if calls := len(repoMock.CreateCalls()); calls != 0 {
		t.Errorf("Create() got %d calls, want 0", calls)
	}
}

func TestResourceService_Update(t *testing.T) {
	want := &TestResourceData{
		ID:                1,
		TestResourceAttrs: TestResourceAttrs{Name: "test1"},
	}

	s := ResourceService[TestResource, TestResourceData, TestResourceAttrs]{
		repo: &ResourceManipulatorMock[TestResourceData, TestResourceAttrs]{
			UpdateFunc: func(ctx context.Context, id int64, attrs *TestResourceAttrs) (*TestResourceData, error) {
				if id != want.ID {
					t.Errorf("Update() got %d, want %d", id, want.ID)
				}
				if diff := cmp.Diff(&want.TestResourceAttrs, attrs); diff != "" {
					t.Errorf("Update() mismatch (-want +got):\n%s", diff)
				}
				return want, nil
			},
		},
		policy: &resourcePolicyMock[TestResourceData, TestResourceAttrs]{
			getManageScopeFunc: func(authCtx AuthContext) resourceScope[TestResourceData] {
				return &resourceScopeMock[TestResourceData]{
					findFunc: func(ctx context.Context, id int64) (*TestResourceData, error) {
						if id != want.ID {
							t.Errorf("Find() got %d, want %d", id, want.ID)
						}
						return want, nil
					},
				}
			},
			authorizeUpdateFunc: func(ctx context.Context, authCtx AuthContext, resource *TestResourceData, attrs *TestResourceAttrs) error {
				if diff := cmp.Diff(want, resource); diff != "" {
					t.Errorf("authorizeUpdate() mismatch (-want +got):\n%s", diff)
				}
				if diff := cmp.Diff(&want.TestResourceAttrs, attrs); diff != "" {
					t.Errorf("authorizeUpdate() mismatch (-want +got):\n%s", diff)
				}
				return nil
			},
		},
		getValidator: func(attrs *TestResourceAttrs) v8n.ValidatableWithContext {
			if diff := cmp.Diff(&want.TestResourceAttrs, attrs); diff != "" {
				t.Errorf("getValidator() mismatch (-want +got):\n%s", diff)
			}

			return &testValidator{}
		},
	}

	resource, _ := s.Update(context.Background(), nil, want.ID, &want.TestResourceAttrs)
	if diff := cmp.Diff(want, resource); diff != "" {
		t.Errorf("Update() mismatch (-want +got):\n%s", diff)
	}
}

func TestResourceService_Update_validationError(t *testing.T) {
	testResourceAttrs := &TestResourceAttrs{Name: "test1"}
	want := &TestResourceData{
		ID:                1,
		TestResourceAttrs: *testResourceAttrs,
	}

	repoMock := &ResourceManipulatorMock[TestResourceData, TestResourceAttrs]{}
	s := ResourceService[TestResource, TestResourceData, TestResourceAttrs]{
		repo: repoMock,
		policy: &resourcePolicyMock[TestResourceData, TestResourceAttrs]{
			getManageScopeFunc: func(authCtx AuthContext) resourceScope[TestResourceData] {
				return &resourceScopeMock[TestResourceData]{
					findFunc: func(ctx context.Context, id int64) (*TestResourceData, error) {
						if id != want.ID {
							t.Errorf("Find() got %d, want %d", id, want.ID)
						}
						return want, nil
					},
				}
			},
			authorizeUpdateFunc: func(ctx context.Context, authCtx AuthContext, resource *TestResourceData, attrs *TestResourceAttrs) error {
				if diff := cmp.Diff(want, resource); diff != "" {
					t.Errorf("authorizeUpdate() mismatch (-want +got):\n%s", diff)
				}
				if diff := cmp.Diff(&want.TestResourceAttrs, attrs); diff != "" {
					t.Errorf("authorizeUpdate() mismatch (-want +got):\n%s", diff)
				}
				return nil
			},
		},
		getValidator: func(attrs *TestResourceAttrs) v8n.ValidatableWithContext {
			if diff := cmp.Diff(testResourceAttrs, attrs); diff != "" {
				t.Errorf("getValidator() mismatch (-want +got):\n%s", diff)
			}

			return &testValidator{err: errors.New("validation error")}
		},
	}

	resource, err := s.Update(context.Background(), nil, 1, testResourceAttrs)
	if err == nil {
		t.Errorf("Update() got %v, want error", resource)
	} else if resource != nil {
		t.Errorf("Update() got %v, want nil resource with error", resource)
	}

	if calls := len(repoMock.UpdateCalls()); calls != 0 {
		t.Errorf("Update() got %d calls, want 0", calls)
	}
}

func TestResourceService_Delete(t *testing.T) {
	want := int64(1)

	s := ResourceService[TestResource, TestResourceData, TestResourceAttrs]{
		repo: &ResourceManipulatorMock[TestResourceData, TestResourceAttrs]{
			DeleteFunc: func(ctx context.Context, id int64) error {
				if id != want {
					t.Errorf("Delete() got %d, want %d", id, want)
				}
				return nil
			},
		},
		policy: &resourcePolicyMock[TestResourceData, TestResourceAttrs]{
			getManageScopeFunc: func(authCtx AuthContext) resourceScope[TestResourceData] {
				return &resourceScopeMock[TestResourceData]{
					findFunc: func(ctx context.Context, id int64) (*TestResourceData, error) {
						if id != want {
							t.Errorf("Find() got %d, want %d", id, id)
						}
						return &TestResourceData{ID: id}, nil
					},
				}
			},
			authorizeDeleteFunc: func(ctx context.Context, authCtx AuthContext, resource *TestResourceData) error {
				return nil
			},
		},
	}

	if err := s.Delete(context.Background(), nil, want); err != nil {
		t.Errorf("Delete() got %v, want nil", err)
	}
}
