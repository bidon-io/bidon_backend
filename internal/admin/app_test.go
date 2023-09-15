package admin_test

import (
	"context"
	"errors"
	"testing"

	"github.com/bidon-io/bidon-backend/internal/admin"
	"github.com/google/go-cmp/cmp"
)

func TestAppService_Create(t *testing.T) {
	users := []admin.User{
		{ID: 1, IsAdmin: ptr(true)},
		{ID: 2, IsAdmin: ptr(false)},
	}
	apps := []admin.App{
		{ID: 1, AppAttrs: admin.AppAttrs{UserID: users[0].ID}, User: users[0]},
		{ID: 2, AppAttrs: admin.AppAttrs{UserID: users[0].ID}, User: users[0]},
		{ID: 3, AppAttrs: admin.AppAttrs{UserID: users[1].ID}, User: users[1]},
		{ID: 4, AppAttrs: admin.AppAttrs{UserID: users[1].ID}, User: users[1]},
	}

	store := &admin.StoreMock{
		AppsFunc: func() admin.AppRepo {
			return &admin.AppRepoMock{
				CreateFunc: func(_ context.Context, attrs *admin.AppAttrs) (*admin.App, error) {
					app := new(admin.App)

					app.ID = apps[len(apps)-1].ID + 1
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

	tests := []struct {
		name     string
		authCtx  admin.AuthContext
		attrs    admin.AppAttrs
		want     *admin.App
		checkErr func(error)
	}{
		{
			name:    "admin creates app for themselves",
			authCtx: userContext{user: users[0]},
			attrs: admin.AppAttrs{
				UserID: users[0].ID,
			},
			want: &admin.App{
				ID: 5,
				AppAttrs: admin.AppAttrs{
					UserID: users[0].ID,
				},
				User: users[0],
			},
			checkErr: func(err error) {
				if err != nil {
					t.Errorf("Create() error = %v, wantErr %v", err, false)
				}
			},
		},
		{
			name:    "admin creates app with no userID set",
			authCtx: userContext{user: users[0]},
			attrs:   admin.AppAttrs{},
			want: &admin.App{
				ID: 5,
				AppAttrs: admin.AppAttrs{
					UserID: users[0].ID,
				},
				User: users[0],
			},
			checkErr: func(err error) {
				if err != nil {
					t.Errorf("Create() error = %v, wantErr %v", err, false)
				}
			},
		},
		{
			name:    "admin creates app for another user",
			authCtx: userContext{user: users[0]},
			attrs: admin.AppAttrs{
				UserID: users[1].ID,
			},
			want: &admin.App{
				ID: 5,
				AppAttrs: admin.AppAttrs{
					UserID: users[1].ID,
				},
				User: users[1],
			},
			checkErr: func(err error) {
				if err != nil {
					t.Errorf("Create() error = %v, wantErr %v", err, false)
				}
			},
		},
		{
			name:    "admin creates app for non-existent user",
			authCtx: userContext{user: users[0]},
			attrs: admin.AppAttrs{
				UserID: users[len(users)-1].ID + 1,
			},
			want: nil,
			checkErr: func(err error) {
				if err == nil {
					t.Errorf("Create() error = %v, wantErr %v", err, true)
				}
			},
		},
		{
			name:    "non-admin user creates app for themselves",
			authCtx: userContext{user: users[1]},
			attrs: admin.AppAttrs{
				UserID: users[1].ID,
			},
			want: &admin.App{
				ID: 5,
				AppAttrs: admin.AppAttrs{
					UserID: users[1].ID,
				},
				User: users[1],
			},
			checkErr: func(err error) {
				if err != nil {
					t.Errorf("Create() error = %v, wantErr %v", err, false)
				}
			},
		},
		{
			name:    "non-admin user creates app with no userID set",
			authCtx: userContext{user: users[1]},
			attrs:   admin.AppAttrs{},
			want: &admin.App{
				ID: 5,
				AppAttrs: admin.AppAttrs{
					UserID: users[1].ID,
				},
				User: users[1],
			},
			checkErr: func(err error) {
				if err != nil {
					t.Errorf("Create() error = %v, wantErr %v", err, false)
				}
			},
		},
		{
			name:    "non-admin user creates app for another user",
			authCtx: userContext{user: users[1]},
			attrs: admin.AppAttrs{
				UserID: users[0].ID,
			},
			want: nil,
			checkErr: func(err error) {
				if err == nil {
					t.Errorf("Create() error = %v, wantErr %v", err, true)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := appService.Create(context.Background(), tt.authCtx, &tt.attrs)
			tt.checkErr(err)

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("Create() mismatch (-want +got):\n%s", diff)
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
		{ID: 2, AppAttrs: admin.AppAttrs{UserID: users[0].ID}, User: users[0]},
		{ID: 3, AppAttrs: admin.AppAttrs{UserID: users[1].ID}, User: users[1]},
		{ID: 4, AppAttrs: admin.AppAttrs{UserID: users[1].ID}, User: users[1]},
	}
	adminApps := apps[:2]
	nonAdminApps := apps[2:]

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
					var app admin.App
					for _, a := range apps {
						if a.ID == id {
							app = a
							break
						}
					}

					if attrs.HumanName != "" {
						app.HumanName = attrs.HumanName
					}
					if attrs.UserID != 0 {
						app.UserID = attrs.UserID

						for _, u := range users {
							if u.ID == attrs.UserID {
								app.User = u
								break
							}
						}
					}

					return &app, nil
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

	tests := []struct {
		name     string
		authCtx  admin.AuthContext
		id       int64
		attrs    admin.AppAttrs
		want     *admin.App
		checkErr func(error)
	}{
		{
			name:    "admin updates their own app",
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
			checkErr: func(err error) {
				if err != nil {
					t.Errorf("Update() error = %v, wantErr %v", err, false)
				}
			},
		},
		{
			name:    "admin updates another user's app",
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
			checkErr: func(err error) {
				if err != nil {
					t.Errorf("Update() error = %v, wantErr %v", err, false)
				}
			},
		},
		{
			name:    "admin updates app owner",
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
			checkErr: func(err error) {
				if err != nil {
					t.Errorf("Create() error = %v, wantErr %v", err, false)
				}
			},
		},
		{
			name:    "admin updates app owner to non-existent user",
			authCtx: userContext{user: adminUser},
			id:      adminApps[0].ID,
			attrs: admin.AppAttrs{
				UserID: users[len(users)-1].ID + 1,
			},
			want: nil,
			checkErr: func(err error) {
				if err == nil {
					t.Errorf("Create() error = %v, wantErr %v", err, true)
				}
			},
		},
		{
			name:    "non-admin user updates their own app",
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
			checkErr: func(err error) {
				if err != nil {
					t.Errorf("Update() error = %v, wantErr %v", err, false)
				}
			},
		},
		{
			name:    "non-admin user updates another user's app",
			authCtx: userContext{user: nonAdminUser},
			id:      adminApps[0].ID,
			attrs: admin.AppAttrs{
				HumanName: "new name",
			},
			want: nil,
			checkErr: func(err error) {
				if err == nil {
					t.Errorf("Update() error = %v, wantErr %v", err, true)
				}
			},
		},
		{
			name:    "non-admin user updates app owner",
			authCtx: userContext{user: nonAdminUser},
			id:      nonAdminApps[0].ID,
			attrs: admin.AppAttrs{
				UserID: adminUser.ID,
			},
			want: nil,
			checkErr: func(err error) {
				if err == nil {
					t.Errorf("Update() error = %v, wantErr %v", err, true)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := appService.Update(context.Background(), tt.authCtx, tt.id, &tt.attrs)
			tt.checkErr(err)

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("Update() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

type userContext struct {
	user admin.User
}

func (c userContext) UserID() int64 {
	return c.user.ID
}

func (c userContext) IsAdmin() bool {
	return c.user.IsAdmin != nil && *c.user.IsAdmin
}

func ptr[T any](t T) *T {
	return &t
}
