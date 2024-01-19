package main

import (
	"github.com/94DanielBrown/roc/internal/platform/db"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
)

func (app *Config) Home(c echo.Context) error {
	return c.JSON(http.StatusOK, "Home")
}

func (app *Config) ListRoasts(c echo.Context) error {
	// Logic to list roasts
	return c.JSON(http.StatusOK, []string{"Roast1", "Roast2"}) // Example response
}

// Need to call this with Roast struct
func (app *Config) CreateRoast(c echo.Context) error {
	var roastRequest db.Roast

	if err := c.Bind(&roastRequest); err != nil {
		log.Error("Error in binding request: ", err)
		return err
	}

	log.Info("Roast request received: ", roastRequest)

	response := map[string]string{
		"roastID":  roastRequest.RoastID,
		"reviewer": roastRequest.Name,
	}

	return c.JSON(http.StatusCreated, response)
}

func (app *Config) GetRoast(c echo.Context, roastID string) error {
	return c.JSON(http.StatusOK, "Roast")
}

func GetReview(c echo.Context, roastID string) {

}
