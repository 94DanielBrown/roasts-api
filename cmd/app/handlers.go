package main

import (
	"fmt"
	"github.com/94DanielBrown/roasts/internal/platform/db"
	"github.com/94DanielBrown/roasts/pkg/utils"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
	"time"
)

func (app *Config) home(c echo.Context) error {
	return c.JSON(http.StatusOK, "home")
}

func (app *Config) listRoasts(c echo.Context) error {
	// Logic to list roasts
	return c.JSON(http.StatusOK, []string{"Roast1", "Roast2"}) // Example response
}

// createRoastHandler adds the new roast to DynamoDB
func (app *Config) createRoastHandler(c echo.Context) error {
	correlationId := c.Get("correlationID")
	var newRoast db.Roast

	if err := c.Bind(&newRoast); err != nil {
		errMsg := "Error in binding request"
		slog.Error(errMsg, "err", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": errMsg})
	}

	newRoast.RoastID = "ROAST#" + utils.ToPascalCase(newRoast.Name)
	newRoast.SK = "#PROFILE" + time.Now().Format("02042006")

	app.Logger.Info("Roast request received: ", "payload", newRoast, "correlationID", correlationId)

	if err := app.RoastModels.CreateRoast(newRoast); err != nil {
		errMsg := "Error creating roast"
		app.Logger.Error(errMsg, "err", err, "correlationID", correlationId)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": errMsg})
	}

	app.Logger.Info("Roast created", "correlationID", correlationId)
	return c.JSON(http.StatusOK, newRoast)
}

func (app *Config) getRoast(c echo.Context, roastID string) error {
	return c.JSON(http.StatusOK, "Roast")
}

func (app *Config) getAllRoastsHandler(c echo.Context) error {
	correlationId := c.Get("correlationID")
	allRoasts, err := app.RoastModels.GetAllRoasts()
	if err != nil {
		errMsg := "Error getting all roasts from dynamodb"
		slog.Error(errMsg, "err", err, "correlationID", correlationId)
		return c.JSON(http.StatusInternalServerError, "Error getting all roasts from dynamodb")
	}

	fmt.Println(allRoasts)

	app.Logger.Info("All roasts returned", "correlationID", correlationId)
	return c.JSON(http.StatusOK, allRoasts)
}

func GetReviewHandler(c echo.Context, roastID string) {

}

// getAverageRatings returns a map of roast IDs to average ratings for each roast
func (app *Config) getAverageRatings(c echo.Context) error {
	// Fetch all roasts
	roasts, err := app.RoastModels.GetAllRoasts()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch roasts"})
	}

	// Prepare a map to hold the average ratings
	averageRatings := make(map[string]float64)

	// Calculate average rating for each roast
	for _, roast := range roasts {
		reviews, err := app.ReviewModels.GetReviewsByRoast(roast.RoastID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch reviews for roast " + roast.RoastID})
		}

		var totalRating float64
		for _, review := range reviews {
			totalRating += review.Rating
		}

		if len(reviews) > 0 {
			averageRatings[roast.RoastID] = totalRating / float64(len(reviews))
		} else {
			averageRatings[roast.RoastID] = 0 // Handle case with no reviews
		}
	}

	// Return the average ratings
	return c.JSON(http.StatusOK, averageRatings)
}
