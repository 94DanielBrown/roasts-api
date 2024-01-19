package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
	"strings"

	"github.com/94DanielBrown/roc/internal/platform/db"
)

func (app *Config) CreateRoastValidator(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var roastReq db.Roast

		// Bind the JSON payload to the struct
		if err := c.Bind(&roastReq); err != nil {
			log.Error(err)
			return err
		}

		// Check if roastID is lowercase
		if roastReq.RoastID != strings.ToLower(roastReq.RoastID) {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "roastId must be in lowercase"})
		}

		c.Set("roastReq", roastReq)

		// Call handler
		return next(c)
	}
}
