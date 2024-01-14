package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/94DanielBrown/roc/config"
	"github.com/94DanielBrown/roc/internal/platform/db"
	"github.com/94DanielBrown/roc/pkg/infrastructure"
)

const webPort = 8000

type Config struct {
	RoastModels  db.RoastModels
	ReviewModels db.ReviewModels
	Logger       *slog.Logger
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	err := config.LoadEnvVariables()
	if err != nil {
		logger.Error("Unable to load env variables", "error", err)
		os.Exit(1)
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
	tableName := "roc"

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
