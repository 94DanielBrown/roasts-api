package app

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/94DanielBrown/roc/config"
	"github.com/94DanielBrown/roc/internal/platform/db"
	"github.com/94DanielBrown/roc/pkg/infrastructure"
)

const webPort = 8000

type Config struct {
	RoastModels  db.RoastModels
	ReviewModels db.ReviewModels
}

func main() {
	err := config.LoadEnvVariables()
	if err != nil {
		log.Panic(fmt.Sprintf("Unable to load env variables: %s", err))
	}

	client, err := infrastructure.ConnectToDynamo()
	if err != nil {
		log.Panicf("Error connecting to dynamodb: %v", err)
	}

	app := Config{
		RoastModels:  db.NewRoastModels(client),
		ReviewModels: db.NewReviewModels(client),
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

		//db.Waitforcreate
		log.Println("Table created successfully.")
	} else {
		log.Printf("Table with name %v already exists.", tableName)
	}
}
