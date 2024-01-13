package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (app *Config) Home(c echo.Context) error {
	return c.JSON(http.StatusOK, "Home")
}

func (app *Config) ListRoasts(c echo.Context) error {
	// Logic to list roasts
	return c.JSON(http.StatusOK, []string{"Roast1", "Roast2"}) // Example response
}

func (app *Config) CreateRoast(c echo.Context) error {
	// Logic to create a roast
	return c.String(http.StatusCreated, "Roast created")
}
