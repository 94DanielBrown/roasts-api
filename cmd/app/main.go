package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/94DanielBrown/awsapp"
	s3 "github.com/94DanielBrown/awsapp/pkg/s3"
	"github.com/94DanielBrown/roasts/config"
	"github.com/94DanielBrown/roasts/internal/database"
	"github.com/94DanielBrown/roasts/internal/roasts"
	"github.com/94DanielBrown/roasts/internal/utils"
	"github.com/labstack/echo/v4"
)

const webPort = 8000

type Config struct {
	RoastModels  database.RoastModels
	ReviewModels database.ReviewModels
	UserModels   database.UserModels
	Logger       *slog.Logger
	S3           *s3.Client
}

func (app *Config) routes() *echo.Echo {
	e := echo.New()

	// Use custom middleware func to add correlationID to context to use in logging
	e.Use(utils.CorrelationID)

	e.GET("/", app.home)
	e.POST("/roast", app.createRoastHandler, roasts.CreateRoastValidator)
	e.GET("/roasts", app.getAllRoastsHandler)
	e.GET("/roast/:roastID", app.getRoastHandler)
	//Add validator
	// use request body lots of things
	e.POST("/review", app.createReviewHandler)
	e.GET("/reviews/:roastID", app.getReviewsHandler)
	e.GET("/test", func(c echo.Context) error {
		return echo.NewHTTPError(http.StatusInternalServerError, "Test error")
	})
	// creates user if not already in dynamo
	e.GET("/user/:userID", app.getUserHandler)
	// use request body lots of things
	e.POST("/saveRoast", app.saveRoastHandler)
	e.POST("/removeRoast", app.removeRoastHandler)
	e.GET("/userReviews/:userID", app.getUserReviewHandler)
	e.POST("/userSettings/:userID", app.updateUserSettingsHandler)
	e.GET("/newImage", app.uploadImage)
	return e
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	err := config.LoadEnvVariables()
	if err != nil {
		logger.Error("Unable to load env variables", "error", err)
		os.Exit(1)
	}

	// TODO - env variable for table name
	tableName := "roasts"

	ctx := context.Background()

	client, table, err := awsapp.InitDynamo(ctx, tableName)
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
	}

	e := app.routes()
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", webPort)))
}
