package admin_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/bidon-io/bidon-backend/internal/admin"
	"github.com/google/go-cmp/cmp"
)

func TestAppService_Meta(t *testing.T) {
	users := []admin.User{
		{ID: 1, IsAdmin: ptr(true)},
		{ID: 2, IsAdmin: ptr(false)},
	}

	store := &admin.StoreMock{
		AppsFunc: func() admin.AppRepo {
			return new(admin.AppRepoMock)
		},
		UsersFunc: func() admin.UserRepo {
			return new(admin.UserRepoMock)
		},
	}

	appService := admin.NewAppService(store)

	tests := map[string]struct {
		authCtx admin.AuthContext
		want    admin.ResourceMeta
	}{
		"admin user requests app metadata": {
			authCtx: userContext{user: users[0]},
			want: admin.ResourceMeta{
				Key: admin.AppResourceKey,
				Permissions: admin.ResourcePermissions{
					Read:   true,
					Create: true,
				},
			},
		},
		"non-admin user requests app metadata": {
			authCtx: userContext{user: users[1]},
			want: admin.ResourceMeta{
				Key: admin.AppResourceKey,
				Permissions: admin.ResourcePermissions{
					Read:   true,
					Create: true,
				},
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := appService.Meta(context.Background(), tt.authCtx)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("%v: AppService.Meta() mismatch (-want +got):\n%s", name, diff)
			}
		})
	}
}

func TestAppService_List(t *testing.T) {
	users := []admin.User{
		{ID: 1, IsAdmin: ptr(true)},
		{ID: 2, IsAdmin: ptr(false)},
		{ID: 3, IsAdmin: ptr(false)},
	}

	apps := []admin.App{
		{ID: 1, AppAttrs: admin.AppAttrs{UserID: users[1].ID}, User: users[1]},
		{ID: 2, AppAttrs: admin.AppAttrs{UserID: users[1].ID}, User: users[1]},

		{ID: 3, AppAttrs: admin.AppAttrs{UserID: users[2].ID}, User: users[2]},
		{ID: 4, AppAttrs: admin.AppAttrs{UserID: users[2].ID}, User: users[2]},
	}

	store := &admin.StoreMock{
		AppsFunc: func() admin.AppRepo {
			return &admin.AppRepoMock{
				ListFunc: func(_ context.Context) ([]admin.App, error) {
					return apps, nil
				},
				ListOwnedByUserFunc: func(_ context.Context, userID int64) ([]admin.App, error) {
					userApps := make([]admin.App, 0)
					for _, app := range apps {
						if app.UserID == userID {
							userApps = append(userApps, app)
						}
					}

					return userApps, nil
				},
			}
		},
		UsersFunc: func() admin.UserRepo {
			return new(admin.UserRepoMock)
		},
	}

	appService := admin.NewAppService(store)

	tests := map[string]struct {
		authCtx  admin.AuthContext
		want     []admin.AppResource
		checkErr func(error) error
	}{
		"admin user lists apps": {
			authCtx: userContext{user: users[0]},
			want: func() []admin.AppResource {
				resources := make([]admin.AppResource, 0)
				for _, app := range apps {
					app := app
					resources = append(resources, admin.AppResource{
						App: &app,
						Permissions: admin.ResourceInstancePermissions{
							Update: true,
							Delete: true,
						},
					})
				}

				return resources
			}(),
			checkErr: func(err error) error {
				if err != nil {
					return fmt.Errorf("got err = %q, want nil err", err)
				}

				return nil
			},
		},
		"non-admin user lists apps": {
			authCtx: userContext{user: users[1]},
			want: func() []admin.AppResource {
				resources := make([]admin.AppResource, 0)
				for _, app := range apps[0:2] {
					app := app
					resources = append(resources, admin.AppResource{
						App: &app,
						Permissions: admin.ResourceInstancePermissions{
							Update: true,
							Delete: true,
						},
					})
				}
				return resources
			}(),
			checkErr: func(err error) error {
				if err != nil {
					return fmt.Errorf("got err = %q, want nil err", err)
				}

				return nil
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := appService.List(context.Background(), tt.authCtx)
			if err := tt.checkErr(err); err != nil {
				t.Errorf("%v: AppService.List() %v", name, err)
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("%v: AppService.List() mismatch (-want +got):\n%s", name, diff)
			}
		})
	}
}

func TestAppService_Find(t *testing.T) {
	users := []admin.User{
		{ID: 1, IsAdmin: ptr(true)},
		{ID: 2, IsAdmin: ptr(false)},
		{ID: 3, IsAdmin: ptr(false)},
	}

	apps := []admin.App{
		{ID: 1, AppAttrs: admin.AppAttrs{UserID: users[1].ID}, User: users[1]},
		{ID: 2, AppAttrs: admin.AppAttrs{UserID: users[1].ID}, User: users[1]},

		{ID: 3, AppAttrs: admin.AppAttrs{UserID: users[2].ID}, User: users[2]},
		{ID: 4, AppAttrs: admin.AppAttrs{UserID: users[2].ID}, User: users[2]},
	}

	store := &admin.StoreMock{
		AppsFunc: func() admin.AppRepo {
			return &admin.AppRepoMock{
				FindFunc: func(_ context.Context, id int64) (*admin.App, error) {
					for _, app := range apps {
						if app.ID == id {
							return &app, nil
						}
					}

					return nil, errors.New("not found")
				},
				FindOwnedByUserFunc: func(_ context.Context, userID, id int64) (*admin.App, error) {
					for _, app := range apps {
						if app.ID == id && app.UserID == userID {
							return &app, nil
						}
					}

					return nil, errors.New("not found")
				},
			}
		},
		UsersFunc: func() admin.UserRepo {
			return new(admin.UserRepoMock)
		},
	}

	appService := admin.NewAppService(store)

	tests := map[string]struct {
		authCtx  admin.AuthContext
		id       int64
		want     *admin.AppResource
		checkErr func(error) error
	}{
		"admin user finds any app": {
			authCtx: userContext{user: users[0]},
			id:      apps[0].ID,
			want: &admin.AppResource{
				App: &apps[0],
				Permissions: admin.ResourceInstancePermissions{
					Update: true,
					Delete: true,
				},
			},
			checkErr: func(err error) error {
				if err != nil {
					return fmt.Errorf("got err = %q, want nil err", err)
				}

				return nil
			},
		},
		"non-admin user finds owned app": {
			authCtx: userContext{user: users[1]},
			id:      apps[0].ID,
			want: &admin.AppResource{
				App: &apps[0],
				Permissions: admin.ResourceInstancePermissions{
					Update: true,
					Delete: true,
				},
			},
			checkErr: func(err error) error {
				if err != nil {
					return fmt.Errorf("got err = %q, want nil err", err)
				}

				return nil
			},
		},
		"non-admin user finds not owned app": {
			authCtx: userContext{user: users[1]},
			id:      apps[2].ID,
			want:    nil,
			checkErr: func(err error) error {
				if err == nil {
					return errors.New("got nil err, want err")
				}

				return nil
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := appService.Find(context.Background(), tt.authCtx, tt.id)
			if err := tt.checkErr(err); err != nil {
				t.Errorf("%v: AppService.Find() %v", name, err)
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("%v: AppService.Find() mismatch (-want +got):\n%s", name, diff)
			}
		})
	}
}

func TestAppService_Create(t *testing.T) {
	users := []admin.User{
		{ID: 1, IsAdmin: ptr(true)},
		{ID: 2, IsAdmin: ptr(false)},
	}

	store := &admin.StoreMock{
		AppsFunc: func() admin.AppRepo {
			return &admin.AppRepoMock{
				CreateFunc: func(_ context.Context, attrs *admin.AppAttrs) (*admin.App, error) {
					app := new(admin.App)

					app.ID = 1
					app.AppAttrs = *attrs

					for _, u := range users {
						if u.ID == attrs.UserID {
							app.User = u
							break
						}
					}

					return app, nil
				},
			}
		},
		UsersFunc: func() admin.UserRepo {
			return &admin.UserRepoMock{
				FindFunc: func(_ context.Context, id int64) (*admin.User, error) {
					for _, u := range users {
						if u.ID == id {
							return &u, nil
						}
					}

					return nil, errors.New("not found")
				},
			}
		},
	}

	appService := admin.NewAppService(store)

	tests := map[string]struct {
		authCtx  admin.AuthContext
		attrs    admin.AppAttrs
		want     *admin.App
		checkErr func(error) error
	}{
		"admin user creates app for themselves": {
			authCtx: userContext{user: users[0]},
			attrs: admin.AppAttrs{
				UserID: users[0].ID,
			},
			want: &admin.App{
				ID: 1,
				AppAttrs: admin.AppAttrs{
					UserID: users[0].ID,
				},
				User: users[0],
			},
			checkErr: func(err error) error {
				if err != nil {
					return fmt.Errorf("got err = %q, want nil err", err)
				}

				return nil
			},
		},
		"admin user creates app with no userID set": {
			authCtx: userContext{user: users[0]},
			attrs:   admin.AppAttrs{},
			want: &admin.App{
				ID: 1,
				AppAttrs: admin.AppAttrs{
					UserID: users[0].ID,
				},
				User: users[0],
			},
			checkErr: func(err error) error {
				if err != nil {
					return fmt.Errorf("got err = %q, want nil err", err)
				}

				return nil
			},
		},
		"admin user creates app for another user": {
			authCtx: userContext{user: users[0]},
			attrs: admin.AppAttrs{
				UserID: users[1].ID,
			},
			want: &admin.App{
				ID: 1,
				AppAttrs: admin.AppAttrs{
					UserID: users[1].ID,
				},
				User: users[1],
			},
			checkErr: func(err error) error {
				if err != nil {
					return fmt.Errorf("got err = %q, want nil err", err)
				}

				return nil
			},
		},
		"admin user creates app for non-existent user": {
			authCtx: userContext{user: users[0]},
			attrs: admin.AppAttrs{
				UserID: users[len(users)-1].ID + 1,
			},
			want: nil,
			checkErr: func(err error) error {
				if err == nil {
					return errors.New("got nil err, want err")
				}

				return nil
			},
		},
		"non-admin user creates app for themselves": {
			authCtx: userContext{user: users[1]},
			attrs: admin.AppAttrs{
				UserID: users[1].ID,
			},
			want: &admin.App{
				ID: 1,
				AppAttrs: admin.AppAttrs{
					UserID: users[1].ID,
				},
				User: users[1],
			},
			checkErr: func(err error) error {
				if err != nil {
					return fmt.Errorf("got err = %q, want nil err", err)
				}

				return nil
			},
		},
		"non-admin user creates app with no userID set": {
			authCtx: userContext{user: users[1]},
			attrs:   admin.AppAttrs{},
			want: &admin.App{
				ID: 1,
				AppAttrs: admin.AppAttrs{
					UserID: users[1].ID,
				},
				User: users[1],
			},
			checkErr: func(err error) error {
				if err != nil {
					return fmt.Errorf("got err = %q, want nil err", err)
				}

				return nil
			},
		},
		"non-admin user creates app for another user": {
			authCtx: userContext{user: users[1]},
			attrs: admin.AppAttrs{
				UserID: users[0].ID,
			},
			want: nil,
			checkErr: func(err error) error {
				if err == nil {
					return errors.New("got nil error, want err")
				}

				return nil
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := appService.Create(context.Background(), tt.authCtx, &tt.attrs)
			if err := tt.checkErr(err); err != nil {
				t.Errorf("%v: AppService.Create() %v", name, err)
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("%v: AppService.Create() mismatch (-want +got):\n%s", name, diff)
			}
		})
	}
}

func TestAppService_Update(t *testing.T) {
	users := []admin.User{
		{ID: 1, IsAdmin: ptr(true)},
		{ID: 2, IsAdmin: ptr(false)},
	}
	adminUser := users[0]
	nonAdminUser := users[1]

	apps := []admin.App{
		{ID: 1, AppAttrs: admin.AppAttrs{UserID: users[0].ID}, User: users[0]},
		{ID: 2, AppAttrs: admin.AppAttrs{UserID: users[1].ID}, User: users[1]},
	}
	adminApps := apps[:1]
	nonAdminApps := apps[1:]

	store := &admin.StoreMock{
		AppsFunc: func() admin.AppRepo {
			return &admin.AppRepoMock{
				FindFunc: func(_ context.Context, id int64) (*admin.App, error) {
					for _, app := range apps {
						if app.ID == id {
							return &app, nil
						}
					}

					return nil, errors.New("not found")
				},
				FindOwnedByUserFunc: func(ctx context.Context, userID int64, id int64) (*admin.App, error) {
					for _, app := range apps {
						if app.ID == id && app.UserID == userID {
							return &app, nil
						}
					}

					return nil, errors.New("not found")
				},
				UpdateFunc: func(_ context.Context, id int64, attrs *admin.AppAttrs) (*admin.App, error) {
					app := new(admin.App)
					for _, *app = range apps {
						if app.ID == id {
							break
						}
					}

					if attrs.HumanName != "" {
						app.HumanName = attrs.HumanName
					}
					if attrs.UserID != 0 {
						app.UserID = attrs.UserID

						for _, app.User = range users {
							if app.User.ID == attrs.UserID {
								break
							}
						}
					}

					return app, nil
				},
			}
		},
		UsersFunc: func() admin.UserRepo {
			return &admin.UserRepoMock{
				FindFunc: func(_ context.Context, id int64) (*admin.User, error) {
					for _, u := range users {
						if u.ID == id {
							return &u, nil
						}
					}

					return nil, errors.New("not found")
				},
			}
		},
	}

	appService := admin.NewAppService(store)

	tests := map[string]struct {
		authCtx  admin.AuthContext
		id       int64
		attrs    admin.AppAttrs
		want     *admin.App
		checkErr func(error) error
	}{
		"admin user updates their own app": {
			authCtx: userContext{user: adminUser},
			id:      adminApps[0].ID,
			attrs: admin.AppAttrs{
				HumanName: "new name",
			},
			want: &admin.App{
				ID: adminApps[0].ID,
				AppAttrs: admin.AppAttrs{
					UserID:    adminApps[0].UserID,
					HumanName: "new name",
				},
				User: adminApps[0].User,
			},
			checkErr: func(err error) error {
				if err != nil {
					return fmt.Errorf("got err = %q, want nil err", err)
				}

				return nil
			},
		},
		"admin user updates another user's app": {
			authCtx: userContext{user: adminUser},
			id:      nonAdminApps[0].ID,
			attrs: admin.AppAttrs{
				HumanName: "new name",
			},
			want: &admin.App{
				ID: nonAdminApps[0].ID,
				AppAttrs: admin.AppAttrs{
					UserID:    nonAdminApps[0].UserID,
					HumanName: "new name",
				},
				User: nonAdminApps[0].User,
			},
			checkErr: func(err error) error {
				if err != nil {
					return fmt.Errorf("got err = %q, want nil err", err)
				}

				return nil
			},
		},
		"admin user updates app owner": {
			authCtx: userContext{user: adminUser},
			id:      adminApps[0].ID,
			attrs: admin.AppAttrs{
				UserID: nonAdminUser.ID,
			},
			want: &admin.App{
				ID: adminApps[0].ID,
				AppAttrs: admin.AppAttrs{
					UserID: nonAdminUser.ID,
				},
				User: nonAdminUser,
			},
			checkErr: func(err error) error {
				if err != nil {
					return fmt.Errorf("got err = %q, want nil err", err)
				}

				return nil
			},
		},
		"admin user updates app owner to non-existent user": {
			authCtx: userContext{user: adminUser},
			id:      adminApps[0].ID,
			attrs: admin.AppAttrs{
				UserID: users[len(users)-1].ID + 1,
			},
			want: nil,
			checkErr: func(err error) error {
				if err == nil {
					return errors.New("got nil err, want err")
				}

				return nil
			},
		},
		"non-admin user updates their own app": {
			authCtx: userContext{user: nonAdminUser},
			id:      nonAdminApps[0].ID,
			attrs: admin.AppAttrs{
				HumanName: "new name",
			},
			want: &admin.App{
				ID: nonAdminApps[0].ID,
				AppAttrs: admin.AppAttrs{
					UserID:    nonAdminApps[0].UserID,
					HumanName: "new name",
				},
				User: nonAdminApps[0].User,
			},
			checkErr: func(err error) error {
				if err != nil {
					return fmt.Errorf("got err = %q, want nil err", err)
				}

				return nil
			},
		},
		"non-admin user updates another user's app": {
			authCtx: userContext{user: nonAdminUser},
			id:      adminApps[0].ID,
			attrs: admin.AppAttrs{
				HumanName: "new name",
			},
			want: nil,
			checkErr: func(err error) error {
				if err == nil {
					return errors.New("got nil err, want err")
				}

				return nil
			},
		},
		"non-admin user updates app owner": {
			authCtx: userContext{user: nonAdminUser},
			id:      nonAdminApps[0].ID,
			attrs: admin.AppAttrs{
				UserID: adminUser.ID,
			},
			want: nil,
			checkErr: func(err error) error {
				if err == nil {
					return errors.New("got nil err, want err")
				}

				return nil
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := appService.Update(context.Background(), tt.authCtx, tt.id, &tt.attrs)
			if err := tt.checkErr(err); err != nil {
				t.Errorf("%v: AppService.Update() %v", name, err)
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("%v: AppService.Update() mismatch (-want +got):\n%s", name, diff)
			}
		})
	}
}

func TestAppService_Delete(t *testing.T) {
	users := []admin.User{
		{ID: 1, IsAdmin: ptr(true)},
		{ID: 2, IsAdmin: ptr(false)},
	}
	adminUser := users[0]
	nonAdminUser := users[1]

	apps := []admin.App{
		{ID: 1, AppAttrs: admin.AppAttrs{UserID: users[0].ID}, User: users[0]},
		{ID: 2, AppAttrs: admin.AppAttrs{UserID: users[1].ID}, User: users[1]},
	}
	adminApps := apps[:1]
	nonAdminApps := apps[1:]

	store := &admin.StoreMock{
		AppsFunc: func() admin.AppRepo {
			return &admin.AppRepoMock{
				FindFunc: func(_ context.Context, id int64) (*admin.App, error) {
					for _, app := range apps {
						if app.ID == id {
							return &app, nil
						}
					}

					return nil, errors.New("not found")
				},
				FindOwnedByUserFunc: func(ctx context.Context, userID int64, id int64) (*admin.App, error) {
					for _, app := range apps {
						if app.ID == id && app.UserID == userID {
							return &app, nil
						}
					}

					return nil, errors.New("not found")
				},
				DeleteFunc: func(_ context.Context, id int64) error {
					return nil
				},
			}
		},
		UsersFunc: func() admin.UserRepo {
			return new(admin.UserRepoMock)
		},
	}

	appService := admin.NewAppService(store)

	tests := map[string]struct {
		authCtx  admin.AuthContext
		id       int64
		checkErr func(error) error
	}{
		"admin user deletes their own app": {
			authCtx: userContext{user: adminUser},
			id:      adminApps[0].ID,
			checkErr: func(err error) error {
				if err != nil {
					return fmt.Errorf("got err = %q, want nil err", err)
				}
				return nil
			},
		},
		"admin user deletes another user's app": {
			authCtx: userContext{user: adminUser},
			id:      nonAdminApps[0].ID,
			checkErr: func(err error) error {
				if err != nil {
					return fmt.Errorf("got err = %q, want nil err", err)
				}
				return nil
			},
		},
		"admin user deletes non-existent app": {
			authCtx: userContext{user: adminUser},
			id:      apps[len(apps)-1].ID + 1,
			checkErr: func(err error) error {
				if err == nil {
					return errors.New("got nil err, want err")
				}
				return nil
			},
		},
		"non-admin user deletes their own app": {
			authCtx: userContext{user: nonAdminUser},
			id:      nonAdminApps[0].ID,
			checkErr: func(err error) error {
				if err != nil {
					return fmt.Errorf("got err = %q, want nil err", err)
				}
				return nil
			},
		},
		"non-admin user deletes another user's app": {
			authCtx: userContext{user: nonAdminUser},
			id:      adminApps[0].ID,
			checkErr: func(err error) error {
				if err == nil {
					return errors.New("got nil err, want err")
				}
				return nil
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			err := appService.Delete(context.Background(), tt.authCtx, tt.id)
			if err := tt.checkErr(err); err != nil {
				t.Errorf("%v: AppService.Delete() %v", name, err)
			}
		})
	}
}
