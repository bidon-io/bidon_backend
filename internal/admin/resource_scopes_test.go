package admin

import (
	"context"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPublicResourceScope(t *testing.T) {
	want := []TestResourceData{
		{ID: 1, TestResourceAttrs: TestResourceAttrs{Name: "test1"}},
		{ID: 2, TestResourceAttrs: TestResourceAttrs{Name: "test1"}},
	}

	type repo struct {
		*AllResourceQuerierMock[TestResourceData]
	}

	s := &publicResourceScope[TestResourceData]{
		repo: &repo{
			&AllResourceQuerierMock[TestResourceData]{
				ListFunc: func(_ context.Context) ([]TestResourceData, error) {
					return want, nil
				},
				FindFunc: func(_ context.Context, id int64) (*TestResourceData, error) {
					if id != want[0].ID {
						t.Errorf("Find() got = %v, want %v", id, want[0].ID)
					}

					return &want[0], nil
				},
			},
		},
	}

	gotList, _ := s.list(context.Background())
	if diff := cmp.Diff(gotList, want); diff != "" {
		t.Errorf("list() mismatch (-want +got):\n%s", diff)
	}

	gotFind, _ := s.find(context.Background(), want[0].ID)
	if diff := cmp.Diff(gotFind, &want[0]); diff != "" {
		t.Errorf("find() mismatch (-want +got):\n%s", diff)
	}
}

func TestPrivateResourceScope(t *testing.T) {
	want := []TestResourceData{
		{ID: 1, TestResourceAttrs: TestResourceAttrs{Name: "test1"}},
		{ID: 2, TestResourceAttrs: TestResourceAttrs{Name: "test1"}},
	}

	type repo struct {
		*AllResourceQuerierMock[TestResourceData]
	}

	s := &privateResourceScope[TestResourceData]{
		repo: &repo{
			&AllResourceQuerierMock[TestResourceData]{
				ListFunc: func(_ context.Context) ([]TestResourceData, error) {
					return want, nil
				},
				FindFunc: func(_ context.Context, id int64) (*TestResourceData, error) {
					if id != want[0].ID {
						t.Errorf("Find() got = %v, want %v", id, want[0].ID)
					}

					return &want[0], nil
				},
			},
		},
	}

	tests := []struct {
		name    string
		authCtx AuthContext
		wantErr bool
	}{
		{
			"admin user",
			&AuthContextMock{
				IsAdminFunc: func() bool {
					return true
				},
			},
			false,
		},
		{
			"non-admin user",
			&AuthContextMock{
				IsAdminFunc: func() bool {
					return false
				},
			},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s.authCtx = tt.authCtx

			gotList, err := s.list(context.Background())
			if tt.wantErr {
				if err == nil {
					t.Errorf("list() error = %v, wantErr %v", err, tt.wantErr)
				}
			} else if diff := cmp.Diff(gotList, want); diff != "" {
				t.Errorf("list() mismatch (-want +got):\n%s", diff)
			}

			gotFind, err := s.find(context.Background(), want[0].ID)
			if tt.wantErr {
				if err == nil {
					t.Errorf("find() error = %v, wantErr %v", err, tt.wantErr)
				}
			} else if diff := cmp.Diff(gotFind, &want[0]); diff != "" {
				t.Errorf("find() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestOwnedResourceScope_list(t *testing.T) {
	testResources := []TestResourceData{
		{ID: 1, TestResourceAttrs: TestResourceAttrs{Name: "test1"}},
		{ID: 2, TestResourceAttrs: TestResourceAttrs{Name: "test2"}},
	}

	type repo struct {
		*AllResourceQuerierMock[TestResourceData]
		*OwnedResourceQuerierMock[TestResourceData]
	}

	s := &ownedResourceScope[TestResourceData]{
		repo: &repo{
			&AllResourceQuerierMock[TestResourceData]{
				ListFunc: func(_ context.Context) ([]TestResourceData, error) {
					return testResources, nil
				},
			},
			&OwnedResourceQuerierMock[TestResourceData]{
				ListOwnedByUserFunc: func(_ context.Context, userID int64) ([]TestResourceData, error) {
					if userID != 1 {
						t.Errorf("ListOwnedByUser() got = %v, want %v", userID, 1)
					}

					return testResources[0 : len(testResources)/2], nil
				},
			},
		},
	}

	tests := []struct {
		name    string
		authCtx AuthContext
		want    []TestResourceData
	}{
		{
			"admin user",
			&AuthContextMock{
				IsAdminFunc: func() bool {
					return true
				},
			},
			[]TestResourceData{
				{ID: 1, TestResourceAttrs: TestResourceAttrs{Name: "test1"}},
				{ID: 2, TestResourceAttrs: TestResourceAttrs{Name: "test2"}},
			},
		},
		{
			"non-admin user",
			&AuthContextMock{
				IsAdminFunc: func() bool {
					return false
				},
				UserIDFunc: func() int64 {
					return 1
				},
			},
			[]TestResourceData{
				{ID: 1, TestResourceAttrs: TestResourceAttrs{Name: "test1"}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s.authCtx = tt.authCtx

			got, _ := s.list(context.Background())
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("list() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestOwnedResourceScope_find(t *testing.T) {
	testResources := []TestResourceData{
		{ID: 1, TestResourceAttrs: TestResourceAttrs{Name: "test1"}},
		{ID: 2, TestResourceAttrs: TestResourceAttrs{Name: "test2"}},
	}

	type repo struct {
		*AllResourceQuerierMock[TestResourceData]
		*OwnedResourceQuerierMock[TestResourceData]
	}

	s := &ownedResourceScope[TestResourceData]{
		repo: &repo{
			&AllResourceQuerierMock[TestResourceData]{
				FindFunc: func(_ context.Context, id int64) (*TestResourceData, error) {
					for i := range testResources {
						r := &testResources[i]
						if r.ID == id {
							return r, nil
						}
					}

					return nil, errors.New("not found")
				},
			},
			&OwnedResourceQuerierMock[TestResourceData]{
				FindOwnedByUserFunc: func(_ context.Context, userID, id int64) (*TestResourceData, error) {
					for i := range testResources[:len(testResources)/2] {
						r := &testResources[i]
						if r.ID == id {
							return r, nil
						}
					}

					return nil, errors.New("not found")
				},
			},
		},
	}

	tests := []struct {
		name    string
		authCtx AuthContext
		id      int64
		want    *TestResourceData
		wantErr bool
	}{
		{
			"admin user",
			&AuthContextMock{
				IsAdminFunc: func() bool {
					return true
				},
			},
			1,
			&TestResourceData{ID: 1, TestResourceAttrs: TestResourceAttrs{Name: "test1"}},
			false,
		},
		{
			"non-admin user with owned resource",
			&AuthContextMock{
				IsAdminFunc: func() bool {
					return false
				},
				UserIDFunc: func() int64 {
					return 1
				},
			},
			1,
			&TestResourceData{ID: 1, TestResourceAttrs: TestResourceAttrs{Name: "test1"}},
			false,
		},
		{
			"non-admin user with not-owned resource",
			&AuthContextMock{
				IsAdminFunc: func() bool {
					return false
				},
				UserIDFunc: func() int64 {
					return 1
				},
			},
			2,
			nil,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s.authCtx = tt.authCtx

			got, err := s.find(context.Background(), tt.id)
			if tt.wantErr {
				if err == nil {
					t.Errorf("find() error = %v, wantErr %v", err, tt.wantErr)
				}
			} else if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("find() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestOwnedOrSharedResourceScope_list(t *testing.T) {
	testResources := []TestResourceData{
		{ID: 1, TestResourceAttrs: TestResourceAttrs{Name: "test1"}},
		{ID: 2, TestResourceAttrs: TestResourceAttrs{Name: "test2"}},
	}

	type repo struct {
		*AllResourceQuerierMock[TestResourceData]
		*OwnedOrSharedResourceQuerierMock[TestResourceData]
	}

	s := &ownedOrSharedResourceScope[TestResourceData]{
		repo: &repo{
			&AllResourceQuerierMock[TestResourceData]{
				ListFunc: func(_ context.Context) ([]TestResourceData, error) {
					return testResources, nil
				},
			},
			&OwnedOrSharedResourceQuerierMock[TestResourceData]{
				ListOwnedByUserOrSharedFunc: func(_ context.Context, userID int64) ([]TestResourceData, error) {
					if userID != 1 {
						t.Errorf("ListOwnedByUserOrShared() got = %v, want %v", userID, 1)
					}

					return testResources[0 : len(testResources)/2], nil
				},
			},
		},
	}

	tests := []struct {
		name    string
		authCtx AuthContext
		want    []TestResourceData
	}{
		{
			"admin user",
			&AuthContextMock{
				IsAdminFunc: func() bool {
					return true
				},
			},
			[]TestResourceData{
				{ID: 1, TestResourceAttrs: TestResourceAttrs{Name: "test1"}},
				{ID: 2, TestResourceAttrs: TestResourceAttrs{Name: "test2"}},
			},
		},
		{
			"non-admin user",
			&AuthContextMock{
				IsAdminFunc: func() bool {
					return false
				},
				UserIDFunc: func() int64 {
					return 1
				},
			},
			[]TestResourceData{
				{ID: 1, TestResourceAttrs: TestResourceAttrs{Name: "test1"}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s.authCtx = tt.authCtx

			got, _ := s.list(context.Background())
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("list() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestOwnedOrSharedResourceScope_find(t *testing.T) {
	testResources := []TestResourceData{
		{ID: 1, TestResourceAttrs: TestResourceAttrs{Name: "test1"}},
		{ID: 2, TestResourceAttrs: TestResourceAttrs{Name: "test2"}},
	}

	type repo struct {
		*AllResourceQuerierMock[TestResourceData]
		*OwnedOrSharedResourceQuerierMock[TestResourceData]
	}

	s := &ownedOrSharedResourceScope[TestResourceData]{
		repo: &repo{
			&AllResourceQuerierMock[TestResourceData]{
				FindFunc: func(_ context.Context, id int64) (*TestResourceData, error) {
					for i := range testResources {
						r := &testResources[i]
						if r.ID == id {
							return r, nil
						}
					}

					return nil, errors.New("not found")
				},
			},
			&OwnedOrSharedResourceQuerierMock[TestResourceData]{
				FindOwnedByUserOrSharedFunc: func(_ context.Context, userID, id int64) (*TestResourceData, error) {
					for i := range testResources[:len(testResources)/2] {
						r := &testResources[i]
						if r.ID == id {
							return r, nil
						}
					}

					return nil, errors.New("not found")
				},
			},
		},
	}

	tests := []struct {
		name    string
		authCtx AuthContext
		id      int64
		want    *TestResourceData
		wantErr bool
	}{
		{
			"admin user",
			&AuthContextMock{
				IsAdminFunc: func() bool {
					return true
				},
			},
			1,
			&TestResourceData{ID: 1, TestResourceAttrs: TestResourceAttrs{Name: "test1"}},
			false,
		},
		{
			"non-admin user with owned resource",
			&AuthContextMock{
				IsAdminFunc: func() bool {
					return false
				},
				UserIDFunc: func() int64 {
					return 1
				},
			},
			1,
			&TestResourceData{ID: 1, TestResourceAttrs: TestResourceAttrs{Name: "test1"}},
			false,
		},
		{
			"non-admin user with not-owned resource",
			&AuthContextMock{
				IsAdminFunc: func() bool {
					return false
				},
				UserIDFunc: func() int64 {
					return 1
				},
			},
			2,
			nil,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s.authCtx = tt.authCtx

			got, err := s.find(context.Background(), tt.id)
			if tt.wantErr {
				if err == nil {
					t.Errorf("find() error = %v, wantErr %v", err, tt.wantErr)
				}
			} else if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("find() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
