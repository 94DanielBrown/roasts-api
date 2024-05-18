package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/94DanielBrown/roasts-api/internal/database"
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
	fmt.Println("request body: ", c.Request().Body)
	fmt.Println("create review handler called")
	correlationId := c.Get("correlationID")
	var newReview database.Review

	// TODO - auth supabase jwt token and authorize request
	// Check header and get jwt token if present
	//authHeader := c.Request().Header.Get("Authorization")
	//fmt.Println("Authorization", authHeader)
	//var tokenString string
	//if len(authHeader) > 7 && strings.ToUpper(authHeader[0:7]) == "BEARER " {
	//	tokenString = authHeader[7:]
	//}
	//if tokenString == "" {
	//	errMsg := "jwt token is missing in the authorization header"
	//	app.Logger.Error(errMsg)
	//	return c.JSON(http.StatusBadRequest, map[string]string{"error": errMsg})
	//}
	//// Parse the JWT token and extract the claims
	//token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
	//	// Ensure the token's algorithm matches your expected signing method
	//	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
	//		return nil, http.ErrNotSupported
	//	}
	//	return []byte("qwertyuiopasdfghjklzxcvbnm123456"), nil // Replace with non temp key
	//})
	//if err != nil {
	//	errMsg := "error parsing JWT token"
	//	app.Logger.Error(errMsg, "error", err, "correlationID", correlationId)
	//	return c.JSON(http.StatusBadRequest, map[string]string{"error": errMsg})
	//}
	//// Map claims from token so can be used in review
	//if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
	//	newReview.UserID = claims.UserID
	//	newReview.FirstName = claims.FirstName
	//	newReview.LastName = claims.LastName
	//} else {
	//	return c.JSON(http.StatusInternalServerError, map[string]string{"error": "invalid jwt token"})
	//}

	if err := c.Bind(&newReview); err != nil {
		errMsg := "error in binding request"
		slog.Error(errMsg, "err", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": errMsg})
	}
	fmt.Println("date addeed", newReview.DateAdded)

	newReview.RoastKey = "ROAST#" + newReview.RoastID
	newReview.SK = "REVIEW#" + reviews.GenerateID()

	app.Logger.Info("review request received: ", "payload", newReview, "correlationID", correlationId)

	if err := app.ReviewModels.CreateReview(newReview); err != nil {
		errMsg := "error creating review"
		app.Logger.Error(errMsg, "err", err, "correlationID", correlationId)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": errMsg})
	}

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

// saveRoastHandler binds data from request body to save roast to database
func (app *Config) saveRoastHandler(c echo.Context) error {
	correlationId := c.Get("correlationID")
	var requestData struct {
		RoastID string `json:"roastID"`
		UserID  string `json:"userID"`
	}
	if err := c.Bind(&requestData); err != nil {
		app.Logger.Error("error binding request", "error", err, "correlationID", correlationId)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}
	err := app.UserModels.UpdateSavedRoasts(requestData.UserID, requestData.RoastID)
	if err != nil {
		app.Logger.Error("error saving roast", "error", err, "correlationID", correlationId)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "error saving roast"})
	}
	return c.JSON(http.StatusOK, requestData.RoastID)
}

// removeRoastHandler binds data from request body to save roast to database
func (app *Config) removeRoastHandler(c echo.Context) error {
	fmt.Println("remove roast handler called")
	correlationId := c.Get("correlationID")
	var requestData struct {
		RoastID string `json:"roastID"`
		UserID  string `json:"userID"`
	}
	if err := c.Bind(&requestData); err != nil {
		app.Logger.Error("error binding request", "error", err, "correlationID", correlationId)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}
	err := app.UserModels.RemoveSavedRoast(requestData.UserID, requestData.RoastID)
	if err != nil {
		app.Logger.Error("error removing roast", "error", err, "correlationID", correlationId)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "error removing roast"})
	}
	return c.JSON(http.StatusOK, requestData.RoastID)
}

// getUserReviewHandler retrieves the user's reviews from DynamoDB
func (app *Config) getUserReviewHandler(c echo.Context) error {
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
		errMsg := "no reviews found for user"
		app.Logger.Error(errMsg, "err", err, "userID", userID, "correlationID", correlationId)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": errMsg})
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
	fmt.Println("test")
	correlationId := c.Get("correlationID")
	bucketName := "images-testing-dev-roast"
	objectKey := fmt.Sprintf("upload/%d", time.Now().Unix())
	expiry := 30 * time.Minute

	presignedURL, err := app.S3.GeneratePresignedURL(bucketName, objectKey, expiry)
	if err != nil {
		errMsg := "error creating presigned URL"
		app.Logger.Error("failed to generate presigned URL", "error", err, "correlationID", correlationId)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": errMsg})
	}

	response := map[string]string{
		"presignedURL": presignedURL,
		"objectKey":    objectKey,
	}

	return c.JSON(http.StatusOK, response)
}
