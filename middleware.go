package main

import (
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

// Configure middleware jwt with the custom claims type
var jwtConfig = echojwt.Config{
	NewClaimsFunc: func(c echo.Context) jwt.Claims {
		return new(jwtAccessClaims)
	},
	SigningKey:    []byte(loadedConfig.AccessTokenHSKey),
	SigningMethod: "HS512",
}
