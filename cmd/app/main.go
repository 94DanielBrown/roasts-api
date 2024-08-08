package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/94DanielBrown/awsapp"
	s3 "github.com/94DanielBrown/awsapp/pkg/s3"
	_ "github.com/94DanielBrown/roasts-api/cmd/app/docs"
	"github.com/94DanielBrown/roasts-api/config"
	"github.com/94DanielBrown/roasts-api/internal/database"
	"github.com/94DanielBrown/roasts-api/internal/utils"
	"github.com/94DanielBrown/roasts-api/pkg/apikey"
	"github.com/94DanielBrown/roasts-api/pkg/firebase"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

const webPort = 8000

type Config struct {
	RoastModels  database.RoastModels
	ReviewModels database.ReviewModels
	UserModels   database.UserModels
	Logger       *slog.Logger
	S3           *s3.Client
	ImageBucket  string
}

func (app *Config) routes() *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Use custom middleware func to add correlationID to context to use in logging
	e.Use(utils.CorrelationID)

	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.POST("/roast", app.createRoastHandler, apikey.Validate())
	e.POST("/deleteRoast", app.deleteRoastHandler, apikey.Validate())
	e.GET("/roasts", app.getAllRoastsHandler)
	e.GET("/roast/:roastID", app.getRoastHandler, firebase.FirebaseJWTMiddleware())
	e.POST("/saveRoast", app.saveRoastHandler, firebase.FirebaseJWTMiddleware())
	e.POST("/removeRoast", app.removeRoastHandler, firebase.FirebaseJWTMiddleware())
	//Add validator
	// use request body lots of things
	e.POST("/review", app.createReviewHandler, firebase.FirebaseJWTMiddleware())
	e.GET("/reviews/:roastID", app.getReviewsHandler)
	e.POST("/removeReview", app.removeReviewHandler)
	// creates user if not already in dynamo
	e.GET("/user/:userID", app.getUserHandler)
	// use request body lots of things
	e.GET("/userReviews/:userID", app.getUserReviewsHandler)
	e.POST("/userSettings/:userID", app.updateUserSettingsHandler)
	e.GET("/newImage", app.uploadImage)
	return e
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	env, err := config.LoadEnvVariables()
	if err != nil {
		logger.Error("Unable to load env variables", "error", err)
		os.Exit(1)
	}

	ctx := context.Background()

	client, table, err := awsapp.InitDynamo(ctx, env.TableName)
	if err != nil {
		logger.Error("error setting up dynamo for app", "error", err)
		os.Exit(1)
	} else {
		logger.Info(table)
	}

	s3Client, err := s3.Connect()
	if err != nil {
		logger.Error("error setting up s3 for app", "error", err)
		os.Exit(1)
	}

	app := Config{
		RoastModels:  database.NewRoastModels(client),
		ReviewModels: database.NewReviewModels(client),
		UserModels:   database.NewUserModels(client),
		Logger:       logger,
		S3:           s3Client,
		ImageBucket:  env.ImageBucket,
	}

	e := app.routes()
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", env.WebPort)))
}
