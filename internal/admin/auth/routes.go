package auth

import (
	"github.com/bidon-io/bidon-backend/internal/admin/auth/tokenmgmt"
	"github.com/bidon-io/bidon-backend/internal/admin/auth/usermgmt"
	"github.com/bidon-io/bidon-backend/internal/db"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func SetUpRoutes(e *echo.Echo, db *db.DB, jwtSecretKey []byte) {
	userService := usermgmt.NewUserService(db)
	tokenService := tokenmgmt.NewTokenService(jwtSecretKey)
	authService := NewAuthService(userService, tokenService)

	e.POST("/auth/login", authService.LogIn)
}

func ConfigureJWT(g *echo.Group, jwtSecretKey []byte) {
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(tokenmgmt.JwtCustomClaims)
		},
		SigningKey: jwtSecretKey,
	}
	g.Use(echojwt.WithConfig(config))
}
