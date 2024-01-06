package infrastructure

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"log"
)

func ConnectToDynamo() (*dynamodb.Client, error) {
	config, err := NewAwsConfig()
	if err != nil {
		log.Fatal(err)
	}
	client := dynamodb.NewFromConfig(config)
	return client, nil
}
