package main

import (
	"github.com/94DanielBrown/roc/internal/platform/db"
	"github.com/94DanielBrown/roc/pkg/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"log/slog"
	"net/http"
)

func (app *Config) Home(c echo.Context) error {
	return c.JSON(http.StatusOK, "Home")
}

func (app *Config) ListRoasts(c echo.Context) error {
	// Logic to list roasts
	return c.JSON(http.StatusOK, []string{"Roast1", "Roast2"}) // Example response
}

// CreateRoast adds the new roast to DynamoDB
func (app *Config) CreateRoast(c echo.Context) error {
	var newRoast db.Roast

	if err := c.Bind(&newRoast); err != nil {
		errMsg := "Error in binding request"
		slog.Error(errMsg, "err", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": errMsg})
	}

	newRoast.RoastID = "Roast#" + utils.ToPascalCase(newRoast.Name)

	log.Info("Roast request received: ", newRoast)

	if err := app.RoastModels.CreateRoast(newRoast); err != nil {
		errMsg := "Error creating roast"
		slog.Error(errMsg, "err", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": errMsg})
	}

	return c.JSON(http.StatusOK, newRoast)
}

func (app *Config) GetRoast(c echo.Context, roastID string) error {
	return c.JSON(http.StatusOK, "Roast")
}

func GetReview(c echo.Context, roastID string) {

}
