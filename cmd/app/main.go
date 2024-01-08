package app

import (
	"fmt"
	"log"

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
}
