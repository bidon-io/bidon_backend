package admin

//go:generate go run -mod=mod github.com/matryer/moq@latest -out app_mocks_test.go . AppRepo

import (
	"context"
)

type App struct {
	ID        int64  `json:"id"`
	PublicUID string `json:"public_uid"`
	AppAttrs
	User User `json:"user"`
}

type AppAttrs struct {
	PlatformID  PlatformID `json:"platform_id"`
	HumanName   string     `json:"human_name"`
	PackageName string     `json:"package_name"`
	UserID      int64      `json:"user_id"`
	AppKey      string     `json:"app_key"`
}

type PlatformID string

const (
	UnknownPlatformID PlatformID = ""
	IOSPlatformID     PlatformID = "ios"
	AndroidPlatformID PlatformID = "android"
)

type AppService = ResourceService[App, AppAttrs]

func NewAppService(store Store) *AppService {
	return &AppService{
		repo:   store.Apps(),
		policy: newAppPolicy(store),

		prepareCreateAttrs: func(authCtx AuthContext, attrs *AppAttrs) {
			if attrs.UserID == 0 {
				attrs.UserID = authCtx.UserID()
			}
		},
	}
}

type AppRepo interface {
	AllResourceQuerier[App]
	OwnedResourceQuerier[App]
	ResourceManipulator[App, AppAttrs]
}

type appPolicy struct {
	repo AppRepo

	userPolicy *userPolicy
}

func newAppPolicy(store Store) *appPolicy {
	return &appPolicy{
		repo: store.Apps(),

		userPolicy: newUserPolicy(store),
	}
}

func (p *appPolicy) getReadScope(authCtx AuthContext) resourceScope[App] {
	return &ownedResourceScope[App]{
		repo:    p.repo,
		authCtx: authCtx,
	}
}

func (p *appPolicy) getManageScope(authCtx AuthContext) resourceScope[App] {
	return &ownedResourceScope[App]{
		repo:    p.repo,
		authCtx: authCtx,
	}
}

func (p *appPolicy) authorizeCreate(ctx context.Context, authCtx AuthContext, attrs *AppAttrs) error {
	// If user is not the owner, check if user can manage the owner.
	if attrs.UserID != authCtx.UserID() {
		_, err := p.userPolicy.getManageScope(authCtx).find(ctx, attrs.UserID)
		return err
	}

	return nil
}

func (p *appPolicy) authorizeUpdate(ctx context.Context, authCtx AuthContext, app *App, attrs *AppAttrs) error {
	// If user tries to change the owner and owner is not the same as before, check if user can manage the new owner.
	if attrs.UserID != 0 && attrs.UserID != app.UserID {
		_, err := p.userPolicy.getManageScope(authCtx).find(ctx, attrs.UserID)
		return err
	}

	return nil
}

func (p *appPolicy) authorizeDelete(_ context.Context, _ AuthContext, _ *App) error {
	return nil
}
