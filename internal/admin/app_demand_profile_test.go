package admin_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/bidon-io/bidon-backend/internal/admin"
	"github.com/google/go-cmp/cmp"
)

func TestAppDemandProfileService_Meta(t *testing.T) {
	users := []admin.User{
		{ID: 1, IsAdmin: ptr(true)},
		{ID: 2, IsAdmin: ptr(false)},
	}

	store := &admin.StoreMock{
		AppDemandProfilesFunc: func() admin.AppDemandProfileRepo {
			return new(admin.AppDemandProfileRepoMock)
		},
		AppsFunc: func() admin.AppRepo {
			return new(admin.AppRepoMock)
		},
		UsersFunc: func() admin.UserRepo {
			return new(admin.UserRepoMock)
		},
		DemandSourceAccountsFunc: func() admin.DemandSourceAccountRepo {
			return new(admin.DemandSourceAccountRepoMock)
		},
		DemandSourcesFunc: func() admin.DemandSourceRepo {
			return new(admin.DemandSourceRepoMock)
		},
	}

	profileService := admin.NewAppDemandProfileService(store)

	tests := map[string]struct {
		authCtx admin.AuthContext
		want    admin.ResourceMeta
	}{
		"admin user requests app demand profile metadata": {
			authCtx: userContext{user: users[0]},
			want: admin.ResourceMeta{
				Key: admin.AppDemandProfileResourceKey,
				Permissions: admin.ResourcePermissions{
					Read:   true,
					Create: true,
				},
			},
		},
		"non-admin user requests app demand profile metadata": {
			authCtx: userContext{user: users[1]},
			want: admin.ResourceMeta{
				Key: admin.AppDemandProfileResourceKey,
				Permissions: admin.ResourcePermissions{
					Read:   true,
					Create: true,
				},
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := profileService.Meta(context.Background(), tt.authCtx)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("%v: AppDemandProfileService.Meta() mismatch (-want +got):\n%s", name, diff)
			}
		})
	}
}

func TestAppDemandProfileService_List(t *testing.T) {
	users := []admin.User{
		{ID: 1, IsAdmin: ptr(true)},
		{ID: 2, IsAdmin: ptr(false)},
		{ID: 3, IsAdmin: ptr(false)},
	}
	adminUser := users[0]
	firstUser := users[1]
	secondUser := users[2]

	apps := []admin.App{
		{ID: 1, AppAttrs: admin.AppAttrs{UserID: firstUser.ID}, User: firstUser},
		{ID: 2, AppAttrs: admin.AppAttrs{UserID: secondUser.ID}, User: secondUser},
	}
	firstUserApps := apps[0:1]
	secondUserApps := apps[1:2]

	profiles := []admin.AppDemandProfile{
		{ID: 1, AppDemandProfileAttrs: admin.AppDemandProfileAttrs{AppID: firstUserApps[0].ID}, App: firstUserApps[0]},
		{ID: 2, AppDemandProfileAttrs: admin.AppDemandProfileAttrs{AppID: firstUserApps[0].ID}, App: firstUserApps[0]},

		{ID: 3, AppDemandProfileAttrs: admin.AppDemandProfileAttrs{AppID: secondUserApps[0].ID}, App: secondUserApps[0]},
		{ID: 4, AppDemandProfileAttrs: admin.AppDemandProfileAttrs{AppID: secondUserApps[0].ID}, App: secondUserApps[0]},
	}
	firstUserProfiles := profiles[0:2]

	store := &admin.StoreMock{
		AppDemandProfilesFunc: func() admin.AppDemandProfileRepo {
			return &admin.AppDemandProfileRepoMock{
				ListFunc: func(_ context.Context) ([]admin.AppDemandProfile, error) {
					return profiles, nil
				},
				ListOwnedByUserFunc: func(_ context.Context, userID int64) ([]admin.AppDemandProfile, error) {
					userProfiles := make([]admin.AppDemandProfile, 0)
					for _, profile := range profiles {
						if profile.App.UserID == userID {
							userProfiles = append(userProfiles, profile)
						}
					}
					return userProfiles, nil
				},
			}
		},
		AppsFunc: func() admin.AppRepo {
			return new(admin.AppRepoMock)
		},
		UsersFunc: func() admin.UserRepo {
			return new(admin.UserRepoMock)
		},
		DemandSourceAccountsFunc: func() admin.DemandSourceAccountRepo {
			return new(admin.DemandSourceAccountRepoMock)
		},
		DemandSourcesFunc: func() admin.DemandSourceRepo {
			return new(admin.DemandSourceRepoMock)
		},
	}

	profileService := admin.NewAppDemandProfileService(store)

	tests := map[string]struct {
		authCtx  admin.AuthContext
		want     []admin.AppDemandProfileResource
		checkErr func(error) error
	}{
		"admin user lists profiles": {
			authCtx: userContext{user: adminUser},
			want: func() []admin.AppDemandProfileResource {
				resources := make([]admin.AppDemandProfileResource, 0)
				for _, profile := range profiles {
					profile := profile
					resources = append(resources, admin.AppDemandProfileResource{
						AppDemandProfile: &profile,
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
		"non-admin user lists profiles": {
			authCtx: userContext{user: firstUser},
			want: func() []admin.AppDemandProfileResource {
				resources := make([]admin.AppDemandProfileResource, 0)
				for _, profile := range firstUserProfiles {
					profile := profile
					resources = append(resources, admin.AppDemandProfileResource{
						AppDemandProfile: &profile,
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
			got, err := profileService.List(context.Background(), tt.authCtx)
			if err := tt.checkErr(err); err != nil {
				t.Errorf("%v: AppDemandProfileService.List() %v", name, err)
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("%v: AppDemandProfileService.List() mismatch (-want +got):\n%s", name, diff)
			}
		})
	}
}

func TestAppDemandProfileService_Find(t *testing.T) {
	users := []admin.User{
		{ID: 1, IsAdmin: ptr(true)},
		{ID: 2, IsAdmin: ptr(false)},
		{ID: 3, IsAdmin: ptr(false)},
	}
	adminUser := users[0]
	firstUser := users[1]
	secondUser := users[2]

	apps := []admin.App{
		{ID: 1, AppAttrs: admin.AppAttrs{UserID: firstUser.ID}, User: firstUser},
		{ID: 2, AppAttrs: admin.AppAttrs{UserID: secondUser.ID}, User: secondUser},
	}
	firstUserApps := apps[0:1]
	secondUserApps := apps[1:2]

	profiles := []admin.AppDemandProfile{
		{ID: 1, AppDemandProfileAttrs: admin.AppDemandProfileAttrs{AppID: firstUserApps[0].ID}, App: firstUserApps[0]},
		{ID: 2, AppDemandProfileAttrs: admin.AppDemandProfileAttrs{AppID: firstUserApps[0].ID}, App: firstUserApps[0]},

		{ID: 3, AppDemandProfileAttrs: admin.AppDemandProfileAttrs{AppID: secondUserApps[0].ID}, App: secondUserApps[0]},
		{ID: 4, AppDemandProfileAttrs: admin.AppDemandProfileAttrs{AppID: secondUserApps[0].ID}, App: secondUserApps[0]},
	}
	firstUserProfiles := profiles[0:2]
	secondUserProfiles := profiles[2:4]

	store := &admin.StoreMock{
		AppDemandProfilesFunc: func() admin.AppDemandProfileRepo {
			return &admin.AppDemandProfileRepoMock{
				FindFunc: func(_ context.Context, id int64) (*admin.AppDemandProfile, error) {
					for _, profile := range profiles {
						if profile.ID == id {
							return &profile, nil
						}
					}
					return nil, errors.New("not found")
				},
				FindOwnedByUserFunc: func(_ context.Context, userID, id int64) (*admin.AppDemandProfile, error) {
					for _, profile := range firstUserProfiles {
						if profile.ID == id {
							return &profile, nil
						}
					}
					return nil, errors.New("not found")
				},
			}
		},
		AppsFunc: func() admin.AppRepo {
			return new(admin.AppRepoMock)
		},
		UsersFunc: func() admin.UserRepo {
			return new(admin.UserRepoMock)
		},
		DemandSourceAccountsFunc: func() admin.DemandSourceAccountRepo {
			return new(admin.DemandSourceAccountRepoMock)
		},
		DemandSourcesFunc: func() admin.DemandSourceRepo {
			return new(admin.DemandSourceRepoMock)
		},
	}

	profileService := admin.NewAppDemandProfileService(store)

	tests := map[string]struct {
		authCtx  admin.AuthContext
		id       int64
		want     *admin.AppDemandProfileResource
		checkErr func(error) error
	}{
		"admin user finds any profile": {
			authCtx: userContext{user: adminUser},
			id:      profiles[0].ID,
			want: &admin.AppDemandProfileResource{
				AppDemandProfile: &profiles[0],
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
		"non-admin user finds own profile": {
			authCtx: userContext{user: firstUser},
			id:      firstUserProfiles[0].ID,
			want: &admin.AppDemandProfileResource{
				AppDemandProfile: &firstUserProfiles[0],
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
		"non-admin user finds other user's profile": {
			authCtx: userContext{user: firstUser},
			id:      secondUserProfiles[0].ID,
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
			got, err := profileService.Find(context.Background(), tt.authCtx, tt.id)
			if err := tt.checkErr(err); err != nil {
				t.Errorf("%v: AppDemandProfileService.Find() %v", name, err)
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("%v: AppDemandProfileService.Find() mismatch (-want +got):\n%s", name, diff)
			}
		})
	}
}

func TestAppDemandProfileService_Create(t *testing.T) {
	users := []admin.User{
		{ID: 1, IsAdmin: ptr(true)},
		{ID: 2, IsAdmin: ptr(false)},
		{ID: 3, IsAdmin: ptr(false)},
	}
	adminUser := users[0]
	firstUser := users[1]
	secondUser := users[2]

	apps := []admin.App{
		{ID: 1, AppAttrs: admin.AppAttrs{UserID: adminUser.ID}, User: adminUser},
		{ID: 2, AppAttrs: admin.AppAttrs{UserID: firstUser.ID}, User: secondUser},
		{ID: 3, AppAttrs: admin.AppAttrs{UserID: secondUser.ID}, User: secondUser},
	}
	adminApps := apps[0:1]
	//nonAdminApps := apps[1:2]

	demandSources := []admin.DemandSource{
		{ID: 1},
	}

	demandSourceAccounts := []admin.DemandSourceAccount{
		{
			ID: 1,
			DemandSourceAccountAttrs: admin.DemandSourceAccountAttrs{
				DemandSourceID: demandSources[0].ID,
				UserID:         adminUser.ID,
			},
			DemandSource: demandSources[0],
			User:         adminUser,
		},
		{
			ID: 2,
			DemandSourceAccountAttrs: admin.DemandSourceAccountAttrs{
				DemandSourceID: demandSources[0].ID,
				UserID:         firstUser.ID,
			},
			DemandSource: demandSources[0],
			User:         firstUser,
		},
		{
			ID: 3,
			DemandSourceAccountAttrs: admin.DemandSourceAccountAttrs{
				DemandSourceID: demandSources[0].ID,
				UserID:         secondUser.ID,
			},
			DemandSource: demandSources[0],
			User:         secondUser,
		},
	}

	store := &admin.StoreMock{
		AppDemandProfilesFunc: func() admin.AppDemandProfileRepo {
			return &admin.AppDemandProfileRepoMock{
				CreateFunc: func(_ context.Context, attrs *admin.AppDemandProfileAttrs) (*admin.AppDemandProfile, error) {
					profile := new(admin.AppDemandProfile)

					profile.ID = 1
					profile.AppDemandProfileAttrs = *attrs

					for _, profile.App = range apps {
						if profile.App.ID == attrs.AppID {
							break
						}
					}
					for _, profile.Account = range demandSourceAccounts {
						if profile.Account.ID == attrs.AccountID {
							break
						}
					}
					for _, profile.DemandSource = range demandSources {
						if profile.DemandSource.ID == attrs.DemandSourceID {
							break
						}
					}

					return profile, nil
				},
			}
		},
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
		DemandSourceAccountsFunc: func() admin.DemandSourceAccountRepo {
			return &admin.DemandSourceAccountRepoMock{
				FindFunc: func(_ context.Context, id int64) (*admin.DemandSourceAccount, error) {
					for _, account := range demandSourceAccounts {
						if account.ID == id {
							return &account, nil
						}
					}
					return nil, errors.New("not found")
				},
				FindOwnedByUserOrSharedFunc: func(_ context.Context, userID int64, id int64) (*admin.DemandSourceAccount, error) {
					for _, account := range demandSourceAccounts {
						if account.ID == id && (account.UserID == userID || account.UserID == adminUser.ID) {
							return &account, nil
						}
					}
					return nil, errors.New("not found")
				},
			}
		},
		DemandSourcesFunc: func() admin.DemandSourceRepo {
			return &admin.DemandSourceRepoMock{
				FindFunc: func(_ context.Context, id int64) (*admin.DemandSource, error) {
					for _, source := range demandSources {
						if source.ID == id {
							return &source, nil
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

	profileService := admin.NewAppDemandProfileService(store)

	tests := map[string]struct {
		authCtx  admin.AuthContext
		attrs    admin.AppDemandProfileAttrs
		want     *admin.AppDemandProfile
		checkErr func(error) error
	}{
		"admin creates profile for their app": {
			authCtx: userContext{user: adminUser},
			attrs: admin.AppDemandProfileAttrs{
				DemandSourceID: demandSources[0].ID,
				AccountID:      demandSourceAccounts[0].ID,
				AppID:          adminApps[0].ID,
			},
			want: &admin.AppDemandProfile{
				ID: 1,
				AppDemandProfileAttrs: admin.AppDemandProfileAttrs{
					DemandSourceID: demandSources[0].ID,
					AccountID:      demandSourceAccounts[0].ID,
					AppID:          adminApps[0].ID,
				},
				DemandSource: demandSources[0],
				Account:      demandSourceAccounts[0],
				App:          adminApps[0],
			},
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
			got, err := profileService.Create(context.Background(), tt.authCtx, &tt.attrs)
			if err := tt.checkErr(err); err != nil {
				t.Errorf("%v: AppDemandProfileService.Create() %v", name, err)
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("%v: AppDemandProfileService.Create() mismatch (-want +got):\n%s", name, diff)
			}
		})
	}
}
