package admin

type User struct {
	ID int64 `json:"id"`
	UserAttrs
}

type UserAttrs struct {
	Email string `json:"email"`
}

type UserRepo = ResourceRepo[User, UserAttrs]

type UserService = ResourceService[User, UserAttrs]

func NewUserService(store Store) *UserService {
	return &UserService{
		ResourceRepo: store.Users(),
	}
}
