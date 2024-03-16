package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/94DanielBrown/roasts/internal/database"
	"github.com/94DanielBrown/roasts/internal/reviews"
	"github.com/94DanielBrown/roasts/internal/roasts"
	"github.com/94DanielBrown/roasts/internal/utils"
	"github.com/94DanielBrown/roasts/pkg/dynamo"
	"github.com/labstack/echo/v4"

	"github.com/94DanielBrown/roasts/config"
)

const webPort = 8000

type Config struct {
	RoastModels  database.RoastModels
	ReviewModels database.ReviewModels
	Logger       *slog.Logger
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	err := config.LoadEnvVariables()
	if err != nil {
		logger.Error("Unable to load env variables", "error", err)
		os.Exit(1)
	}

	ctx := context.Background()

	if err := run(ctx, logger); err != nil {
		logger.Error("Failed to startup", "error", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, logger *slog.Logger) error {
	client, err := dynamo.Connect()
	if err != nil {
		logger.Error("Error connecting to dynamodb", "error", err)
		os.Exit(1)
	}

	app := Config{
		RoastModels:  database.NewRoastModels(client),
		ReviewModels: database.NewReviewModels(client),
		Logger:       logger,
	}

	tableName := "roasts"

	exists, err := dynamo.Exists(ctx, client, tableName)
	if err != nil {
		return fmt.Errorf("error checking if dynamodb table exists: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	if !exists {
		err = dynamo.Create(ctx, client, tableName)
		if err != nil {
			return fmt.Errorf("error creating dynamodb table: %w", err)
		}
		logger.Info("table created successfully")
	} else {
		logger.Info("table already exists", "tableName", tableName)
	}

	e := app.routes()
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", webPort)))

	return nil
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
	e.POST("/review", app.createReviewHandler, reviews.CreateReviewValidator)
	e.GET("/reviews/:roastID", app.getReviewsHandler)
	e.GET("/test", func(c echo.Context) error {
		return echo.NewHTTPError(http.StatusInternalServerError, "Test error")
	})

	return e
}
