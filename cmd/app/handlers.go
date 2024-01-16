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

// Need to call this with Roast struct
func (app *Config) CreateRoast(c echo.Context) error {
	// Logic to create a roast
	//roastID := c.FormValue("roastID")
	//reviewer := c.FormValue("reviewer")
	//// Image will be uploaded on frontend and then url passed through to store in dynamo
	//imageUrl := c.FormValue("imageUrl")
	return c.String(http.StatusCreated, "Roast created")
}

func (app *Config) GetRoast(c echo.Context, roastID string) error {
	return c.JSON(http.StatusOK, "Roast")
}

func GetReview(c echo.Context, roastID string) {

}
