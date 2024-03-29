package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/94DanielBrown/roasts/internal/database"
	"github.com/94DanielBrown/roasts/internal/ratings"
	"github.com/94DanielBrown/roasts/internal/reviews"
	"github.com/94DanielBrown/roasts/internal/utils"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type CustomClaims struct {
	UserID    string `json:"userID"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	jwt.RegisteredClaims
}

// home returns a simple welcome message to check if the server is running properly
func (app *Config) home(c echo.Context) error {
	return c.JSON(http.StatusOK, "home")
}

// listRoasts returns a dummy list atm
func (app *Config) listRoasts(c echo.Context) error {
	// Logic to list roasts
	return c.JSON(http.StatusOK, []string{"Roast1", "Roast2"}) // Example response
}

// getRoastHandlers gets an individual roast by roastID
func (app *Config) getRoastHandler(c echo.Context) error {
	correlationId := c.Get("correlationID")
	roastID := c.Param("roastID")
	roastPrefix := "ROAST#" + roastID

	// roast is a pointer here to deal with nil values being returned
	roast, err := app.RoastModels.GetRoastByPrefix(roastPrefix)
	if err != nil {
		errMsg := "error getting roast"
		app.Logger.Error(errMsg, "error", err, "correlationID", correlationId)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": errMsg})
	}

	app.Logger.Info("roast returned", "correlationID", correlationId)
	return c.JSON(http.StatusOK, roast)
}

// createRoastHandler adds the new roast to DynamoDB
func (app *Config) createRoastHandler(c echo.Context) error {
	correlationId := c.Get("correlationID")
	var newRoast database.Roast

	if err := c.Bind(&newRoast); err != nil {
		errMsg := "error in binding request"
		slog.Error(errMsg, "err", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": errMsg})
	}

	RoastID := utils.ToPascalCase(newRoast.Name)
	newRoast.RoastID = RoastID
	newRoast.RoastKey = "ROAST#" + RoastID
	newRoast.SK = "PROFILE#" + time.Now().Format("02042006")

	app.Logger.Info("Roast request received: ", "payload", newRoast, "correlationID", correlationId)

	if err := app.RoastModels.CreateRoast(newRoast); err != nil {
		errMsg := "error creating roast"
		app.Logger.Error(errMsg, "err", err, "correlationID", correlationId)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": errMsg})
	}

	app.Logger.Info("roast created", "correlationID", correlationId)
	return c.JSON(http.StatusOK, newRoast)
}

// getRoast just returns a string atm, may not be required
func (app *Config) getRoast(c echo.Context, roastID string) error {
	return c.JSON(http.StatusOK, "Roast")
}

// getAllRoastsHandler gets and returns all roasts from DynamoDB
func (app *Config) getAllRoastsHandler(c echo.Context) error {
	correlationId := c.Get("correlationID")
	allRoasts, err := app.RoastModels.GetAllRoasts()
	if err != nil {
		errMsg := "error getting all roasts from dynamodb"
		slog.Error(errMsg, "err", err, "correlationID", correlationId)
		return c.JSON(http.StatusInternalServerError, "error getting all roasts from dynamodb")
	}

	app.Logger.Info("all roasts returned", "correlationID", correlationId)
	return c.JSON(http.StatusOK, allRoasts)
}

// createReviewHandler adds the review to DynamoDB
func (app *Config) createReviewHandler(c echo.Context) error {
	correlationId := c.Get("correlationID")
	var newReview database.Review

	// TODO - auth supabase jwt token
	// Check header and get jwt token if present
	authHeader := c.Request().Header.Get("Authorization")
	fmt.Println("Authorization", authHeader)
	var tokenString string
	if len(authHeader) > 7 && strings.ToUpper(authHeader[0:7]) == "BEARER " {
		tokenString = authHeader[7:]
	}
	if tokenString == "" {
		errMsg := "jwt token is missing in the authorization header"
		app.Logger.Error(errMsg)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": errMsg})
	}

	// Parse the JWT token and extract the claims
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Ensure the token's algorithm matches your expected signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, http.ErrNotSupported
		}
		return []byte("qwertyuiopasdfghjklzxcvbnm123456"), nil // Replace with non temp key
	})

	if err != nil {
		errMsg := "error parsing JWT token"
		app.Logger.Error(errMsg, "error", err, "correlationID", correlationId)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": errMsg})
	}

	// Map claims from token so can be used in review
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		newReview.UserID = claims.UserID
		newReview.FirstName = claims.FirstName
		newReview.LastName = claims.LastName
	} else {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "invalid jwt token"})
	}

	if err := c.Bind(&newReview); err != nil {
		errMsg := "error in binding request"
		slog.Error(errMsg, "err", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": errMsg})
	}

	RoastID := utils.ToPascalCase(newReview.RoastName)
	newReview.RoastID = RoastID
	newReview.RoastKey = "ROAST#" + RoastID
	newReview.SK = "REVIEW#" + reviews.GenerateID()

	app.Logger.Info("review request received: ", "payload", newReview, "correlationID", correlationId)

	if err := app.ReviewModels.CreateReview(newReview); err != nil {
		errMsg := "error creating review"
		app.Logger.Error(errMsg, "err", err, "correlationID", correlationId)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": errMsg})
	}

	// TODO - Send to queue for retry or in the meantime just have a scheduled job to rectify inconsistencies
	go func() {
		err := ratings.UpdateAverages(app.RoastModels, newReview)
		if err != nil {
			app.Logger.Error("error updating average rating", "error", err, "correlationID", correlationId)
		}
	}()

	app.Logger.Info("review created", "correlationID", correlationId)
	return c.JSON(http.StatusOK, newReview)
}

// getReviewsHandler gets all reviews for a roast by roastID
func (app *Config) getReviewsHandler(c echo.Context) error {
	correlationId := c.Get("correlationID")
	roastID := c.Param("roastID")
	roastKey := "ROAST#" + roastID

	roastReviews, err := app.ReviewModels.GetReviewsByRoast(roastKey)
	if err != nil {
		errMsg := "error getting roast"
		app.Logger.Error(errMsg, "err", err, "correlationID", correlationId)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": errMsg})
	}

	if roastReviews == nil {
		app.Logger.Info("no roast reviews returned due to no reviews", "correlationID", correlationId)
		return c.JSON(http.StatusNotFound, map[string]string{"message": "reviews not found"})
	}

	app.Logger.Info("reviews returned", "roastID", roastID, "correlationID", correlationId)
	return c.JSON(http.StatusOK, roastReviews)
}
