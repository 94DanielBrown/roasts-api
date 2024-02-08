package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/94DanielBrown/roasts/config"
	"github.com/94DanielBrown/roasts/internal/platform/db"
	"github.com/94DanielBrown/roasts/pkg/infrastructure"
)

const webPort = 8000

type Config struct {
	RoastModels  db.RoastModels
	ReviewModels db.ReviewModels
	Logger       *slog.Logger
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	err := config.LoadEnvVariables()
	if err != nil {
		logger.Error("Unable to load env variables", "error", err)
		os.Exit(1)
	}

	for _, env := range os.Environ() {
		fmt.Println(env)
	}

	client, err := infrastructure.ConnectToDynamo()
	if err != nil {
		logger.Error("Error connecting to dynamodb", "error", err)
		os.Exit(1)
	}

	app := Config{
		RoastModels:  db.NewRoastModels(client),
		ReviewModels: db.NewReviewModels(client),
		Logger:       logger,
	}

	ctx := context.Background()
	tableName := "roasts"

	exists, err := db.TableExists(ctx, client, tableName)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	if !exists {
		err = db.CreateDynamoDBTable(ctx, client, tableName)
		if err != nil {
			log.Fatalf("Error creating Dynamoodb table: %v", err)
		}

		log.Println("Table created successfully.")
	} else {
		log.Printf("Table with name %v already exists.", tableName)
	}

	e := app.routes()

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", webPort)))
}
