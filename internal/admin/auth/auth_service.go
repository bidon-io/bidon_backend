package auth

import (
	"fmt"
	"net/http"

	"github.com/bidon-io/bidon-backend/internal/db"

	"github.com/labstack/echo/v4"
)

type AuthService struct {
	userService  UserService
	tokenService TokenService
}

//go:generate go run -mod=mod github.com/matryer/moq@latest -out mocks/mocks.go -pkg mocks . UserService TokenService

type UserService interface {
	GetUserByEmail(email string) (*db.User, error)
	ComparePassword(storedPasswordHash, password string) bool
}

type TokenService interface {
	GenerateAccessToken(email string) (string, error)
}

func NewAuthService(userService UserService, tokenService TokenService) *AuthService {
	return &AuthService{userService: userService, tokenService: tokenService}
}

func (s *AuthService) LogIn(c echo.Context) error {
	body := &AuthRequest{}
	if err := c.Bind(body); err != nil {
		return fmt.Errorf("failed to bind: %v", err)
	}

	user, err := s.userService.GetUserByEmail(body.Email)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, fmt.Sprintf("failed to get user: %v", err))
	}
	if !s.userService.ComparePassword(user.PasswordHash, body.Password) {
		return echo.NewHTTPError(http.StatusUnauthorized, "wrong password")
	}

	token, err := s.tokenService.GenerateAccessToken(body.Email)
	if err != nil {
		return fmt.Errorf("failed generating tokens: %v", err)
	}
	publicUser := PublicUser{Email: body.Email, IsAdmin: user.IsAdmin}
	return c.JSON(http.StatusOK, AuthResponse{User: publicUser, AccessToken: token})
}
