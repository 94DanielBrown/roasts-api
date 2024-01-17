package main

import (
	"github.com/labstack/echo/v4"
)

func (app *Config) routes() *echo.Echo {
	e := echo.New()

	e.GET("/", app.Home)
	e.GET("/roasts", app.ListRoasts)
	e.POST("/roast", app.CreateRoast, app.CreateRoastValidator)

	return e
}
