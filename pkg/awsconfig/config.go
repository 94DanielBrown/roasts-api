package awsconfig

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

func NewConfig() (aws.Config, error) {
	awsAccessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	awsRegion := os.Getenv("AWS_REGION")

	credProvider := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(
		awsAccessKey,
		awsSecretAccessKey,
		"",
	))

	conf, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion(awsRegion),
		config.WithCredentialsProvider(credProvider),
	)
	if err != nil {
		panic(err)
	}
	return conf, nil
}
