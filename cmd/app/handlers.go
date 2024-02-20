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

func (app *Config) getRoastHandler(c echo.Context) error {
	correlationId := c.Get("correlationID")
	roastID := c.Param("roastID")
	roastPrefix := "ROAST#" + roastID

	// roast is a pointer here to deal with nil values being returned
	roast, err := app.RoastModels.GetRoastByPrefix(roastPrefix)
	if err != nil {
		errMsg := "Error getting roast"
		app.Logger.Error(errMsg, "err", err, "correlationID", correlationId)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": errMsg})
	}

	if roast == nil {
		app.Logger.Info("No roast returned", "correlationID", correlationId)
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Roast not found"})
	}

	app.Logger.Info("Roast returned", "correlationID", correlationId)
	return c.JSON(http.StatusOK, roast)
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
	newRoast.SK = "PROFILE#" + time.Now().Format("02042006")

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

// createReviewHandler adds the review to DynamoDB
func (app *Config) createReviewHandler(c echo.Context) error {
	correlationId := c.Get("correlationID")
	var newReview db.Review

	if err := c.Bind(&newReview); err != nil {
		errMsg := "Error in binding request"
		slog.Error(errMsg, "err", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": errMsg})
	}

	newReview.RoastID = "ROAST#" + utils.ToPascalCase(newReview.RoastName)
	newReview.SK = "REVIEW#" + utils.GenerateReviewID()

	app.Logger.Info("Review request received: ", "payload", newReview, "correlationID", correlationId)

	if err := app.ReviewModels.CreateReview(newReview); err != nil {
		errMsg := "Error creating review"
		app.Logger.Error(errMsg, "err", err, "correlationID", correlationId)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": errMsg})
	}

	app.Logger.Info("Review created", "correlationID", correlationId)
	return c.JSON(http.StatusOK, newReview)
}

// getReviewHandler gets all reviews for a roast by roastID
func (app *Config) getReviewHandler(c echo.Context) error {
	correlationId := c.Get("correlationID")
	roastID := c.Param("roastID")
	roastPrefix := "ROAST#" + roastID

	// roast is a pointer here to deal with nil values being returned
	reviews, err := app.ReviewModels.GetReviewsByRoast(roastPrefix)
	if err != nil {
		errMsg := "Error getting roast"
		app.Logger.Error(errMsg, "err", err, "correlationID", correlationId)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": errMsg})
	}

	if reviews == nil {
		app.Logger.Info("No reviews returned", "correlationID", correlationId)
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Reviews not found"})
	}

	app.Logger.Info("Reviews returned", "correlationID", correlationId)
	return c.JSON(http.StatusOK, reviews)
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
			totalRating += float64(review.Rating)
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
