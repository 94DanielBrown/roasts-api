package main

import (
	"github.com/94DanielBrown/roasts/pkg/middleware"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (app *Config) routes() *echo.Echo {
	e := echo.New()

	// Use custom middleware func to add correlationID to context to use in logging
	e.Use(middleware.CorrelationIDMiddleware)

	e.GET("/", app.home)
	e.GET("/roast/:roastID", app.getRoastHandler)
	e.GET("/roasts", app.getAllRoastsHandler)
	e.POST("/roast", app.createRoastHandler, app.CreateRoastValidator)
	e.GET("/test", func(c echo.Context) error {
		return echo.NewHTTPError(http.StatusInternalServerError, "Test error")
	})

	return e
}
