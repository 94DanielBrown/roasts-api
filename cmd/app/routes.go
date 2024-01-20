package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (app *Config) routes() *echo.Echo {
	e := echo.New()

	e.GET("/", app.Home)
	e.GET("/roasts", app.ListRoasts)
	e.POST("/roast", app.CreateRoast, app.CreateRoastValidator)
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"https://localhost", "https://localhost"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	return e
}
