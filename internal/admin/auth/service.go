package auth

import (
	"context"
	"crypto/subtle"

	"github.com/golang-jwt/jwt/v5"
)

type Service struct {
	userRepo UserRepo
	config   Config
}

type UserRepo interface {
	FindByEmailAndPassword(ctx context.Context, email, password string) (User, error)
}

type Config struct {
	SecretKey         []byte
	SuperUserLogin    []byte
	SuperUserPassword []byte
}

func NewAuthService(userRepo UserRepo, config Config) *Service {
	return &Service{
		userRepo: userRepo,
		config:   config,
	}
}

func (s *Service) LogIn(ctx context.Context, r LogInRequest) (*LogInResponse, error) {
	user, err := s.userRepo.FindByEmailAndPassword(ctx, r.Email, r.Password)
	if err != nil {
		return nil, err
	}

	claims := newJWTClaims(user)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString(s.GetSecretKey())
	if err != nil {
		return nil, err
	}

	return &LogInResponse{
		User:        user,
		AccessToken: accessToken,
	}, nil
}

func (s *Service) GetSecretKey() []byte {
	return s.config.SecretKey
}

func (s *Service) IsSuperUser(username, password string) bool {
	return subtle.ConstantTimeCompare([]byte(username), s.config.SuperUserLogin) == 1 &&
		subtle.ConstantTimeCompare([]byte(password), s.config.SuperUserPassword) == 1
}
