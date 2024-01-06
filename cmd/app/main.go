package app

import (
	"fmt"
	"github.com/94DanielBrown/roc/pkg/infrastructure"
	"log"

	"github.com/94DanielBrown/roc/config"
)

const webPort = 8000

type Config struct {
	Models data.Models
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
		Models: data.New(client),
	}
}
