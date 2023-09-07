package auth

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	User        PublicUser `json:"user"`
	AccessToken string     `json:"access_token"`
}

type PublicUser struct {
	Email   string `json:"email"`
	IsAdmin *bool  `json:"is_admin"`
}
