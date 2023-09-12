package auth

type LogInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogInResponse struct {
	User        User   `json:"user"`
	AccessToken string `json:"access_token"`
}

type User struct {
	ID      int64  `json:"-"`
	Email   string `json:"email"`
	IsAdmin bool   `json:"is_admin"`
}
