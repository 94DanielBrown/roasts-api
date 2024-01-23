package main

import (
	"github.com/94DanielBrown/roasts/internal/platform/db"
	"github.com/94DanielBrown/roasts/pkg/utils"
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

// CreateRoastHandler adds the new roast to DynamoDB
func (app *Config) CreateRoastHandler(c echo.Context) error {
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

func GetReviewHandler(c echo.Context, roastID string) {

}

// GetAverageRatings returns a map of roast IDs to average ratings for each roast
func (app *Config) GetAverageRatings(c echo.Context) error {
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
