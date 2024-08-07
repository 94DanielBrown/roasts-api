package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/94DanielBrown/roasts-api/internal/database"
	"github.com/94DanielBrown/roasts-api/internal/ratings"
	"github.com/94DanielBrown/roasts-api/internal/reviews"
	"github.com/94DanielBrown/roasts-api/internal/utils"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type CustomClaims struct {
	UserID    string `json:"userID"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	jwt.RegisteredClaims
}

type message struct {
	Message string `json:"message"`
}

// @Summary get a roast
// @ID get-roast
// @Tags roasts
// @Produce json
// @Success 200 {object} database.Roast
// @Failure 404 {object} message
// @Failure 500 {object} message
// @Router /roast/{roastID} [get]
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
	if roast == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "roast not found"})
	}

	app.Logger.Info("roast returned", "correlationID", correlationId)
	return c.JSON(http.StatusOK, roast)
}

// @Summary create a roast
// @ID create-roast
// @Tags roasts
// @Accept json
// @Produce json
// @Param data body database.Roast true "Roast object that needs to be created"
// @Success 200 {object} database.Roast
// @Failure 400 {object} message
// @Failure 500 {object} message
// @Router /roast/{roastID} [post]
func (app *Config) createRoastHandler(c echo.Context) error {
	correlationId := c.Get("correlationID")
	var newRoast database.Roast

	if err := c.Bind(&newRoast); err != nil {
		errMsg := "Error in binding request"
		app.Logger.Error(errMsg, "err", err, "correlationID", correlationId)
		return c.JSON(http.StatusBadRequest, message{Message: errMsg})
	}

	RoastID := utils.ToPascalCase(newRoast.Name)
	newRoast.RoastID = RoastID
	newRoast.RoastKey = "ROAST#" + RoastID
	newRoast.SK = "PROFILE#" + time.Now().Format("02042006")

	app.Logger.Info("Roast request received", "payload", newRoast, "correlationID", correlationId)

	if err := app.RoastModels.CreateRoast(newRoast); err != nil {
		errMsg := "Error creating roast"
		app.Logger.Error(errMsg, "err", err, "correlationID", correlationId)
		return c.JSON(http.StatusInternalServerError, message{Message: errMsg})
	}

	app.Logger.Info("Roast created", "correlationID", correlationId)
	return c.JSON(http.StatusOK, newRoast)
}

// @Summary delete a roast
// @ID delete-roast
// @Tags roasts
// @Accept json
// @Produce json
// @Param data body database.Roast true "Roast object that needs to be created"
// @Success 200 {object} database.Roast
// @Failure 400 {object} message
// @Failure 500 {object} message
// @Router /roast/{roastID} [post]
func (app *Config) deleteRoastHandler(c echo.Context) error {
	correlationId := c.Get("correlationID")
	roastName := c.Request().Header.Get("Roast-Name")
	app.Logger.Info("Roast deletion request received", "roast", roastName, "correlationID", correlationId)

	if err := app.RoastModels.DeleteRoast(roastName); err != nil {
		errMsg := "Error delete roast"
		app.Logger.Error(errMsg, "err", err, "correlationID", correlationId)
		return c.JSON(http.StatusInternalServerError, message{Message: errMsg})
	}

	app.Logger.Info("Roast deleted", "correlationID", correlationId)
	return c.JSON(http.StatusOK, fmt.Sprintf("%s deleted", roastName))
}

// @Summary get all roasts
// @ID  get-all-roasts
// @Tags roasts
// @Produce json
// @Success 200 {object} []database.Roast
// @Failure 400 {object} message
// @Failure 500 {object} message
// @Router /roasts [get]
func (app *Config) getAllRoastsHandler(c echo.Context) error {
	correlationId := c.Get("correlationID")
	allRoasts, err := app.RoastModels.GetAllRoasts()
	if err != nil {
		errMsg := "Error getting all roasts from dynamodb"
		app.Logger.Error(errMsg, "err", err, "correlationID", correlationId)
		return c.JSON(http.StatusInternalServerError, message{Message: errMsg})
	}
	if allRoasts == nil {
		allRoasts = []database.Roast{}
	}

	app.Logger.Info("all roasts returned", "correlationID", correlationId)
	return c.JSON(http.StatusOK, allRoasts)
}

// @Summary save roast
// @ID save-roast
// @Tags roasts
// @Accept json
// @Produce json
// @Success 200 {object} message
// @Failure 400 {object} message
// @Failure 500 {object} message
// @Router /saveRoast [post]
func (app *Config) saveRoastHandler(c echo.Context) error {
	correlationId := c.Get("correlationID")
	userID := c.Get("userID")
	var requestData struct {
		RoastID string `json:"roastID"`
		UserID  string `json:"userID"`
	}
	if err := c.Bind(&requestData); err != nil {
		errMsg := "error binding request"
		app.Logger.Error(errMsg, "error", err, "correlationID", correlationId)
		return c.JSON(http.StatusBadRequest, message{Message: errMsg})
	}
	if userID != requestData.UserID {
		return fmt.Errorf("uid in jwt doesn't match request data")
	}
	err := app.UserModels.UpdateSavedRoasts(requestData.UserID, requestData.RoastID)
	if err != nil {
		errMsg := "Error saving roast"
		app.Logger.Error(errMsg, "err", err, "correlationID", correlationId)
		return c.JSON(http.StatusInternalServerError, message{Message: errMsg})
	}
	return c.JSON(http.StatusOK, requestData.RoastID)
}

// @Summary delete roast
// @ID delete-roast
// @Tags roasts
// @Accept json
// @Produce json
// @Success 200 {object} message
// @Failure 400 {object} message
// @Failure 500 {object} message
// @Router /removeRoast/{roastID} [post]
func (app *Config) removeRoastHandler(c echo.Context) error {
	correlationId := c.Get("correlationID")
	userID := c.Get("userID")
	var requestData struct {
		RoastID string `json:"roastID"`
		UserID  string `json:"userID"`
	}
	if err := c.Bind(&requestData); err != nil {
		errMsg := "Error binding request"
		app.Logger.Error(errMsg, "error", err, "correlationID", correlationId)
		return c.JSON(http.StatusBadRequest, message{Message: errMsg})
	}

	if userID != requestData.UserID {
		return fmt.Errorf("uid in jwt doesn't match request data")
	}
	err := app.UserModels.RemoveSavedRoast(requestData.UserID, requestData.RoastID)
	if err != nil {
		errMsg := "Error removing roast"
		app.Logger.Error(errMsg, "err", err, "correlationID", correlationId)
		return c.JSON(http.StatusInternalServerError, message{Message: errMsg})
	}
	return c.JSON(http.StatusOK, requestData.RoastID)
}

// @Summary create review
// @ID create-review
// @Tags reviews
// @Accept json
// @Produce json
// @Success 200 {object} database.Review
// @Failure 400 {object} message
// @Failure 500 {object} message
// @Router /review [post]
func (app *Config) createReviewHandler(c echo.Context) error {
	correlationId := c.Get("correlationID")
	userID := c.Get("userID")
	var newReview database.Review

	if err := c.Bind(&newReview); err != nil {
		errMsg := "error in binding request"
		slog.Error(errMsg, "err", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": errMsg})
	}

	fmt.Println("newReview", newReview)

	if newReview.ReviewKey == "" {
		newReview.RoastKey = "ROAST#" + newReview.RoastID
		newReview.ReviewKey = "REVIEW#" + reviews.GenerateID()
	}

	newReview.RoastKey = "ROAST#" + newReview.RoastID
	app.Logger.Info("review request received: ", "payload", newReview, "correlationID", correlationId)

	fmt.Println("context UserID", userID)
	fmt.Println("new review UserID", newReview.UserID)

	if userID != newReview.UserID {
		return fmt.Errorf("uid in jwt doesn't match request data")
	}
	if err := app.ReviewModels.CreateReview(newReview); err != nil {
		errMsg := "error creating review"
		app.Logger.Error(errMsg, "err", err, "correlationID", correlationId)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": errMsg})
	}

	app.Logger.Info("review created", "correlationID", correlationId)
	err := ratings.UpdateAverages(app.RoastModels, newReview, "plusCount")
	if err != nil {
		errMsg := "error updating averages"
		app.Logger.Error(errMsg, "err", err, "correlationID", correlationId)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": errMsg})
	}
	return c.JSON(http.StatusOK, newReview)
}

// @Summary get reviews for a roast
// @ID  get-roast-reviews
// @Tags reviews
// @Produce json
// @Success 200 {object} []database.Review
// @Failure 404 {object} message
// @Failure 500 {object} message
// @Router /reviews/{roastID} [get]
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

	// TODO - we should really return a 400 for invalid roast ID?
	if roastReviews == nil {
		app.Logger.Info("no roast reviews returned due to no reviews", "correlationID", correlationId)
		return c.JSON(http.StatusNotFound, map[string]string{"message": "reviews not found"})
	}

	app.Logger.Info("reviews returned", "roastID", roastID, "correlationID", correlationId)
	return c.JSON(http.StatusOK, roastReviews)
}

// @Summary delete a review
// @ID delete-review
// @Tags reviews
// @Produce json
// @Success 200 {object} database.Review.reviewKey
// @Failure 500 {object} message
// @Router /removeReview [post]
func (app *Config) removeReviewHandler(c echo.Context) error {
	correlationId := c.Get("correlationID")
	var requestData struct {
		RoastID   string `json:"roastID"`
		ReviewKey string `json:"reviewKey"`
	}
	if err := c.Bind(&requestData); err != nil {
		errMsg := "error binding request"
		app.Logger.Error(errMsg, "error", err, "correlationID", correlationId)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": errMsg})

	}
	roastKey := "ROAST#" + requestData.RoastID
	oldReview, err := app.ReviewModels.GetReviewByKey(roastKey, requestData.ReviewKey)
	if err != nil {
		errMsg := "error getting review"
		app.Logger.Error(errMsg, "error", err, "correlationID", correlationId)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": errMsg})
	}

	err = app.ReviewModels.RemoveReview(roastKey, requestData.ReviewKey)
	if err != nil {
		errMsg := "error removing review"
		app.Logger.Error(errMsg, "error", err, "correlationID", correlationId)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": errMsg})
	}

	err = ratings.UpdateAverages(app.RoastModels, *oldReview, "minusCount")
	if err != nil {
		errMsg := "error updating averages"
		app.Logger.Error(errMsg, "err", err, "correlationID", correlationId)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": errMsg})
	}
	app.Logger.Info("review removed", "correlationID", correlationId)
	return c.JSON(http.StatusOK, requestData.ReviewKey)
}

// getUserHandler retrieves the user's information from DynamoDB or otherwise creates a new user
func (app *Config) getUserHandler(c echo.Context) error {
	correlationId := c.Get("correlationID")
	userID := c.Param("userID")
	app.Logger.Info("User request received", "userID", userID, "correlationID", correlationId)
	userPrefix := "USER#" + userID
	user, err := app.UserModels.GetUserByPrefix(userPrefix)
	if err != nil {
		errMsg := "error retrieving user"
		app.Logger.Error(errMsg, "err", err, "userID", userID, "correlationID", correlationId)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": errMsg})
	}
	if user == nil {
		// User not found, so let's create one
		app.Logger.Info("User not found, creating user", "userID", userID, "correlationID", correlationId)
		newUser := database.User{
			UserKey: userPrefix,
			// Can use SK for something else in future if needed
			SK: "PROFILE#" + userID,
		}
		if err := app.UserModels.CreateUser(newUser); err != nil {
			errMsg := "error creating user"
			app.Logger.Error(errMsg, "err", err, "userID", userID, "correlationID", correlationId)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": errMsg})
		}
		return c.JSON(http.StatusOK, newUser)
	}
	app.Logger.Info("User retrieved", "user", user, "correlationID", correlationId)
	fmt.Print("user: ", user)
	return c.JSON(http.StatusOK, user)
}

// @Summary get a users reviews
// @ID get-user-reviews
// @Tags reviews
// @Produce json
// @Success 200 {object} []database.Review
// @Failure 404 {object} message
// @Failure 500 {object} message
// @Router /userReviews/{userID} [get]
func (app *Config) getUserReviewsHandler(c echo.Context) error {
	correlationId := c.Get("correlationID")
	userID := c.Param("userID")
	app.Logger.Info("user review request received", "userID", userID, "correlationID", correlationId)
	userReviews, err := app.UserModels.GetUserReviews(userID)
	if err != nil {
		errMsg := "error retrieving user"
		app.Logger.Error(errMsg, "err", err, "userID", userID, "correlationID", correlationId)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": errMsg})
	}
	if userReviews == nil {
		app.Logger.Info("no reviews found for user", "correlationID", correlationId)
		return c.JSON(http.StatusNotFound, map[string]string{"message": "reviews not found for user"})
	}
	app.Logger.Info("user reviews returned", "user", userID, "correlationID", correlationId)
	return c.JSON(http.StatusOK, userReviews)
}

// TODO - validate if names are valid and not empty
// updateUserSettingsHandler retrieves the user's reviews from DynamoDB
func (app *Config) updateUserSettingsHandler(c echo.Context) error {
	fmt.Println("test")
	correlationId := c.Get("correlationID")
	userID := c.Param("userID")
	app.Logger.Info("user settings update request received", "userID", userID, "correlationID", correlationId)
	var requestData struct {
		DisplayName string `json:"displayName"`
		FirstName   string `json:"firstName"`
		LastName    string `json:"lastName"`
	}
	if err := c.Bind(&requestData); err != nil {
		app.Logger.Error("error binding request", "error", err, "correlationID", correlationId)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}
	err := app.UserModels.UpdateSettings(userID, requestData.DisplayName, requestData.FirstName, requestData.LastName)
	if err != nil {
		errMsg := "error updating user settings"
		app.Logger.Error(errMsg, "error", err, "userID", userID, "correlationID", correlationId)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": errMsg})
	}
	app.Logger.Info("user settings updated", "user", userID, "correlationID", correlationId)
	return c.JSON(http.StatusOK, "okay")
}

func (app *Config) uploadImage(c echo.Context) error {
	correlationId := c.Get("correlationID")
	bucketName := app.ImageBucket
	objectKey := fmt.Sprintf("upload/%d", time.Now().Unix())
	expiry := 30 * time.Minute
	imageType := "image/jpeg"

	presignedURL, err := app.S3.GeneratePresignedURL(bucketName, objectKey, imageType, expiry)
	if err != nil {
		errMsg := "error creating presigned URL"
		app.Logger.Error("failed to generate presigned URL", "error", err, "correlationID", correlationId)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": errMsg})
	}

	response := map[string]string{
		"presignedURL": presignedURL,
		"objectKey":    objectKey,
	}
	fmt.Println("presigned URL", presignedURL)

	return c.JSON(http.StatusOK, response)
}
