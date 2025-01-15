package admin

import (
	"context"
	"errors"
	"testing"

	"github.com/bidon-io/bidon-backend/internal/admin/resource"

	"github.com/google/go-cmp/cmp"
)

func TestPublicResourceScope(t *testing.T) {
	want := &resource.Collection[TestResourceData]{
		Items: []TestResourceData{
			{ID: 1, TestResourceAttrs: TestResourceAttrs{Name: "test1"}},
			{ID: 2, TestResourceAttrs: TestResourceAttrs{Name: "test1"}},
		},
		Meta: resource.CollectionMeta{
			TotalCount: 2,
		},
	}

	type repo struct {
		*AllResourceQuerierMock[TestResourceData]
	}

	s := &publicResourceScope[TestResourceData]{
		repo: &repo{
			&AllResourceQuerierMock[TestResourceData]{
				ListFunc: func(_ context.Context, qParams map[string][]string) (*resource.Collection[TestResourceData], error) {
					return want, nil
				},
				FindFunc: func(_ context.Context, id int64) (*TestResourceData, error) {
					if id != want.Items[0].ID {
						t.Errorf("Find() got = %v, want %v", id, want.Items[0].ID)
					}

					return &want.Items[0], nil
				},
			},
		},
	}

	gotList, _ := s.list(context.Background(), nil)
	if diff := cmp.Diff(gotList, want); diff != "" {
		t.Errorf("list() mismatch (-want +got):\n%s", diff)
	}

	gotFind, _ := s.find(context.Background(), want.Items[0].ID)
	if diff := cmp.Diff(gotFind, &want.Items[0]); diff != "" {
		t.Errorf("find() mismatch (-want +got):\n%s", diff)
	}
}

func TestPrivateResourceScope(t *testing.T) {
	want := &resource.Collection[TestResourceData]{
		Items: []TestResourceData{
			{ID: 1, TestResourceAttrs: TestResourceAttrs{Name: "test1"}},
			{ID: 2, TestResourceAttrs: TestResourceAttrs{Name: "test1"}},
		},
		Meta: resource.CollectionMeta{
			TotalCount: 2,
		},
	}

	type repo struct {
		*AllResourceQuerierMock[TestResourceData]
	}

	s := &privateResourceScope[TestResourceData]{
		repo: &repo{
			&AllResourceQuerierMock[TestResourceData]{
				ListFunc: func(_ context.Context, qParams map[string][]string) (*resource.Collection[TestResourceData], error) {
					return want, nil
				},
				FindFunc: func(_ context.Context, id int64) (*TestResourceData, error) {
					if id != want.Items[0].ID {
						t.Errorf("Find() got = %v, want %v", id, want.Items[0].ID)
					}

					return &want.Items[0], nil
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

			gotList, err := s.list(context.Background(), nil)
			if tt.wantErr {
				if err == nil {
					t.Errorf("list() error = %v, wantErr %v", err, tt.wantErr)
				}
			} else if diff := cmp.Diff(gotList, want); diff != "" {
				t.Errorf("list() mismatch (-want +got):\n%s", diff)
			}

			gotFind, err := s.find(context.Background(), want.Items[0].ID)
			if tt.wantErr {
				if err == nil {
					t.Errorf("find() error = %v, wantErr %v", err, tt.wantErr)
				}
			} else if diff := cmp.Diff(gotFind, &want.Items[0]); diff != "" {
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
	testCollection := &resource.Collection[TestResourceData]{
		Items: testResources,
		Meta: resource.CollectionMeta{
			TotalCount: 2,
		},
	}

	type repo struct {
		*AllResourceQuerierMock[TestResourceData]
		*OwnedResourceQuerierMock[TestResourceData]
	}

	s := &ownedResourceScope[TestResourceData]{
		repo: &repo{
			&AllResourceQuerierMock[TestResourceData]{
				ListFunc: func(_ context.Context, qParams map[string][]string) (*resource.Collection[TestResourceData], error) {
					return testCollection, nil
				},
			},
			&OwnedResourceQuerierMock[TestResourceData]{
				ListOwnedByUserFunc: func(_ context.Context, userID int64, qParams map[string][]string) (*resource.Collection[TestResourceData], error) {
					if userID != 1 {
						t.Errorf("ListOwnedByUser() got = %v, want %v", userID, 1)
					}

					items := testResources[0 : len(testResources)/2]
					collection := &resource.Collection[TestResourceData]{
						Items: testResources[0 : len(testResources)/2],
						Meta: resource.CollectionMeta{
							TotalCount: int64(len(items)),
						},
					}

					return collection, nil
				},
			},
		},
	}

	tests := []struct {
		name    string
		authCtx AuthContext
		want    *resource.Collection[TestResourceData]
	}{
		{
			"admin user",
			&AuthContextMock{
				IsAdminFunc: func() bool {
					return true
				},
			},
			&resource.Collection[TestResourceData]{
				Items: []TestResourceData{
					{ID: 1, TestResourceAttrs: TestResourceAttrs{Name: "test1"}},
					{ID: 2, TestResourceAttrs: TestResourceAttrs{Name: "test2"}},
				},
				Meta: resource.CollectionMeta{
					TotalCount: 2,
				},
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
			&resource.Collection[TestResourceData]{
				Items: []TestResourceData{
					{ID: 1, TestResourceAttrs: TestResourceAttrs{Name: "test1"}},
				},
				Meta: resource.CollectionMeta{
					TotalCount: 1,
				},
			},
		}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s.authCtx = tt.authCtx

			got, _ := s.list(context.Background(), nil)
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
	testCollection := &resource.Collection[TestResourceData]{
		Items: testResources,
		Meta: resource.CollectionMeta{
			TotalCount: int64(len(testResources)),
		},
	}

	type repo struct {
		*AllResourceQuerierMock[TestResourceData]
		*OwnedOrSharedResourceQuerierMock[TestResourceData]
	}

	s := &ownedOrSharedResourceScope[TestResourceData]{
		repo: &repo{
			&AllResourceQuerierMock[TestResourceData]{
				ListFunc: func(_ context.Context, qParams map[string][]string) (*resource.Collection[TestResourceData], error) {
					return testCollection, nil
				},
			},
			&OwnedOrSharedResourceQuerierMock[TestResourceData]{
				ListOwnedByUserOrSharedFunc: func(_ context.Context, userID int64) (*resource.Collection[TestResourceData], error) {
					if userID != 1 {
						t.Errorf("ListOwnedByUserOrShared() got = %v, want %v", userID, 1)
					}

					items := testResources[0 : len(testResources)/2]
					collection := &resource.Collection[TestResourceData]{
						Items: testResources[0 : len(testResources)/2],
						Meta: resource.CollectionMeta{
							TotalCount: int64(len(items)),
						},
					}

					return collection, nil
				},
			},
		},
	}

	tests := []struct {
		name    string
		authCtx AuthContext
		want    *resource.Collection[TestResourceData]
	}{
		{
			"admin user",
			&AuthContextMock{
				IsAdminFunc: func() bool {
					return true
				},
			},
			&resource.Collection[TestResourceData]{
				Items: []TestResourceData{
					{ID: 1, TestResourceAttrs: TestResourceAttrs{Name: "test1"}},
					{ID: 2, TestResourceAttrs: TestResourceAttrs{Name: "test2"}},
				},
				Meta: resource.CollectionMeta{
					TotalCount: 2,
				},
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
			&resource.Collection[TestResourceData]{
				Items: []TestResourceData{
					{ID: 1, TestResourceAttrs: TestResourceAttrs{Name: "test1"}},
				},
				Meta: resource.CollectionMeta{
					TotalCount: 1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s.authCtx = tt.authCtx

			got, _ := s.list(context.Background(), nil)
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
