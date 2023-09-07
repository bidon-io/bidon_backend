package tokenmgmt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtCustomClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

type TokenService struct {
	jwtSecretKey []byte
}

func NewTokenService(jwtSecretKey []byte) *TokenService {
	return &TokenService{jwtSecretKey}
}

func (s *TokenService) GenerateAccessToken(email string) (string, error) {
	claims := &JwtCustomClaims{
		email,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecretKey)
}
