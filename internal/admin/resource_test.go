package admin

import (
	"context"
	"errors"
	"testing"

	v8n "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/go-cmp/cmp"
)

type TestResource struct {
	ID int64
	TestResourceAttrs
}

type TestResourceAttrs struct {
	Name string
}

func TestResourceService_List(t *testing.T) {
	want := []TestResource{
		{ID: 1, TestResourceAttrs: TestResourceAttrs{Name: "test1"}},
		{ID: 2, TestResourceAttrs: TestResourceAttrs{Name: "test2"}},
	}

	s := ResourceService[TestResource, TestResourceAttrs]{
		policy: &resourcePolicyMock[TestResource]{
			scopeFunc: func(authCtx AuthContext) resourceScope[TestResource] {
				return &resourceScopeMock[TestResource]{
					listFunc: func(ctx context.Context) ([]TestResource, error) {
						return want, nil
					},
				}
			},
		},
	}

	resources, _ := s.List(context.Background(), nil)
	if diff := cmp.Diff(want, resources); diff != "" {
		t.Errorf("List() mismatch (-want +got):\n%s", diff)
	}
}

func TestResourceService_Find(t *testing.T) {
	want := &TestResource{
		ID:                1,
		TestResourceAttrs: TestResourceAttrs{Name: "test1"},
	}

	s := ResourceService[TestResource, TestResourceAttrs]{
		policy: &resourcePolicyMock[TestResource]{
			scopeFunc: func(authCtx AuthContext) resourceScope[TestResource] {
				return &resourceScopeMock[TestResource]{
					findFunc: func(ctx context.Context, id int64) (*TestResource, error) {
						if id != want.ID {
							t.Errorf("Find() got %d, want %d", id, want.ID)
						}
						return want, nil
					},
				}
			},
		},
	}

	resource, _ := s.Find(context.Background(), nil, want.ID)
	if diff := cmp.Diff(want, resource); diff != "" {
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
	want := &TestResource{
		ID:                1,
		TestResourceAttrs: TestResourceAttrs{Name: "test1"},
	}

	s := ResourceService[TestResource, TestResourceAttrs]{
		repo: &ResourceManipulatorMock[TestResource, TestResourceAttrs]{
			CreateFunc: func(ctx context.Context, attrs *TestResourceAttrs) (*TestResource, error) {
				if diff := cmp.Diff(&want.TestResourceAttrs, attrs); diff != "" {
					t.Errorf("Create() mismatch (-want +got):\n%s", diff)
				}
				return want, nil
			},
		},
		getValidator: func(attrs *TestResourceAttrs) v8n.ValidatableWithContext {
			if diff := cmp.Diff(&want.TestResourceAttrs, attrs); diff != "" {
				t.Errorf("getValidator() mismatch (-want +got):\n%s", diff)
			}

			return &testValidator{}
		},
	}

	resource, _ := s.Create(context.Background(), &want.TestResourceAttrs)
	if diff := cmp.Diff(want, resource); diff != "" {
		t.Errorf("Create() mismatch (-want +got):\n%s", diff)
	}
}

func TestResourceService_Create_validationError(t *testing.T) {
	testResourceAttrs := &TestResourceAttrs{Name: "test1"}

	repoMock := &ResourceManipulatorMock[TestResource, TestResourceAttrs]{}
	s := ResourceService[TestResource, TestResourceAttrs]{
		repo: repoMock,
		getValidator: func(attrs *TestResourceAttrs) v8n.ValidatableWithContext {
			if diff := cmp.Diff(testResourceAttrs, attrs); diff != "" {
				t.Errorf("getValidator() mismatch (-want +got):\n%s", diff)
			}

			return &testValidator{err: errors.New("validation error")}
		},
	}

	resource, err := s.Create(context.Background(), testResourceAttrs)
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
	want := &TestResource{
		ID:                1,
		TestResourceAttrs: TestResourceAttrs{Name: "test1"},
	}

	s := ResourceService[TestResource, TestResourceAttrs]{
		repo: &ResourceManipulatorMock[TestResource, TestResourceAttrs]{
			UpdateFunc: func(ctx context.Context, id int64, attrs *TestResourceAttrs) (*TestResource, error) {
				if id != want.ID {
					t.Errorf("Update() got %d, want %d", id, want.ID)
				}
				if diff := cmp.Diff(&want.TestResourceAttrs, attrs); diff != "" {
					t.Errorf("Update() mismatch (-want +got):\n%s", diff)
				}
				return want, nil
			},
		},
		getValidator: func(attrs *TestResourceAttrs) v8n.ValidatableWithContext {
			if diff := cmp.Diff(&want.TestResourceAttrs, attrs); diff != "" {
				t.Errorf("getValidator() mismatch (-want +got):\n%s", diff)
			}

			return &testValidator{}
		},
	}

	resource, _ := s.Update(context.Background(), want.ID, &want.TestResourceAttrs)
	if diff := cmp.Diff(want, resource); diff != "" {
		t.Errorf("Update() mismatch (-want +got):\n%s", diff)
	}
}

func TestResourceService_Update_validationError(t *testing.T) {
	testResourceAttrs := &TestResourceAttrs{Name: "test1"}

	repoMock := &ResourceManipulatorMock[TestResource, TestResourceAttrs]{}
	s := ResourceService[TestResource, TestResourceAttrs]{
		repo: repoMock,
		getValidator: func(attrs *TestResourceAttrs) v8n.ValidatableWithContext {
			if diff := cmp.Diff(testResourceAttrs, attrs); diff != "" {
				t.Errorf("getValidator() mismatch (-want +got):\n%s", diff)
			}

			return &testValidator{err: errors.New("validation error")}
		},
	}

	resource, err := s.Update(context.Background(), 1, testResourceAttrs)
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

	s := ResourceService[TestResource, TestResourceAttrs]{
		repo: &ResourceManipulatorMock[TestResource, TestResourceAttrs]{
			DeleteFunc: func(ctx context.Context, id int64) error {
				if id != want {
					t.Errorf("Delete() got %d, want %d", id, want)
				}
				return nil
			},
		},
	}

	if err := s.Delete(context.Background(), want); err != nil {
		t.Errorf("Delete() got %v, want nil", err)
	}
}
