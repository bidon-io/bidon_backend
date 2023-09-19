package admin

//go:generate go run -mod=mod github.com/matryer/moq@latest -out user_mocks_test.go . UserRepo

import (
	"context"

	v8n "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type User struct {
	ID        int64  `json:"id"`
	PublicUID string `json:"public_uid"`
	IsAdmin   *bool  `json:"is_admin"`
	Email     string `json:"email"`
}

type UserAttrs struct {
	Email    string `json:"email"`
	IsAdmin  *bool  `json:"is_admin"`
	Password string `json:"password"`
}

type UserService = ResourceService[User, UserAttrs]

func NewUserService(store Store) *UserService {
	s := &UserService{
		repo:   store.Users(),
		policy: newUserPolicy(store),
	}

	s.getValidator = func(attrs *UserAttrs) v8n.ValidatableWithContext {
		return &userAttrsValidator{
			attrs: attrs,
		}
	}

	return s
}

type UserRepo interface {
	AllResourceQuerier[User]
	ResourceManipulator[User, UserAttrs]
}

type userPolicy struct {
	repo UserRepo
}

func newUserPolicy(store Store) *userPolicy {
	return &userPolicy{
		repo: store.Users(),
	}
}

func (p *userPolicy) getReadScope(authCtx AuthContext) resourceScope[User] {
	return &privateResourceScope[User]{
		repo:    p.repo,
		authCtx: authCtx,
	}
}

func (p *userPolicy) getManageScope(authCtx AuthContext) resourceScope[User] {
	return &privateResourceScope[User]{
		repo:    p.repo,
		authCtx: authCtx,
	}
}

func (p *userPolicy) authorizeCreate(_ context.Context, authCtx AuthContext, _ *UserAttrs) error {
	if !authCtx.IsAdmin() {
		return ErrActionForbidden
	}

	return nil
}

func (p *userPolicy) authorizeUpdate(_ context.Context, _ AuthContext, _ *User, _ *UserAttrs) error {
	return nil
}

func (p *userPolicy) authorizeDelete(_ context.Context, _ AuthContext, _ *User) error {
	return nil
}

type userAttrsValidator struct {
	attrs    *UserAttrs
	userRepo UserRepo
}

func (v *userAttrsValidator) ValidateWithContext(_ context.Context) error {
	return v8n.ValidateStruct(v.attrs,
		v8n.Field(&v.attrs.Email, is.EmailFormat),
		v8n.Field(&v.attrs.Password, v8n.Length(6, 50)),
	)
}
