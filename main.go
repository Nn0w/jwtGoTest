package main

import (
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	initDB()

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	h := &handler{}
	e.POST("/api/auth/login", h.login)

	e.POST("/api/auth/refresh", h.refreshToken)

	// Restricted group
	r := e.Group("/api/restricted")
	{
		r.Use(echojwt.WithConfig(jwtConfig))

		r.GET("", h.restricted)
	}

	e.Logger.Fatal(e.Start(":" + loadedConfig.Port))
}
