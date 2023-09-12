package auth

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	jwt.RegisteredClaims
	Email string `json:"email"`
	Admin bool   `json:"admin"`
}

func newJWTClaims(user User) JWTClaims {
	return JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.Itoa(int(user.ID)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
		Email: user.Email,
		Admin: user.IsAdmin,
	}
}

func (c JWTClaims) UserID() int64 {
	userID, err := strconv.Atoi(c.Subject)
	if err != nil {
		// We can panic here because we set `sub` on server and sign the token.
		// If it's not an id, either we have a bug or something is very wrong.
		panic(fmt.Errorf("JWT `sub` is not user id: %v", err))
	}

	return int64(userID)
}

func (c JWTClaims) IsAdmin() bool {
	return c.Admin
}
