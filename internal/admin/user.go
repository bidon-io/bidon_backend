package admin

type User struct {
	ID int64 `json:"id"`
	UserAttrs
}

type UserAttrs struct {
	Email string `json:"email"`
}

type UserService = ResourceService[User, UserAttrs]

func NewUserService(store Store) *UserService {
	return &UserService{
		repo: store.Users(),
		policy: &userPolicy{
			repo: store.Users(),
		},
	}
}

type UserRepo interface {
	AllResourceQuerier[User]
	ResourceManipulator[User, UserAttrs]
}

type userPolicy struct {
	repo UserRepo
}

func (p *userPolicy) scope(authCtx AuthContext) resourceScope[User] {
	return &privateResourceScope[User]{
		repo:    p.repo,
		authCtx: authCtx,
	}
}
