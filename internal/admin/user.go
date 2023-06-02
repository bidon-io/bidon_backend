package admin

type UserService = resourceService[User, UserAttrs]

type User struct {
	ID int64 `json:"id"`
	UserAttrs
}

type UserAttrs struct {
	Email string `json:"email"`
}
