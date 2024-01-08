package db

import "github.com/aws/aws-sdk-go-v2/service/dynamodb"

var client *dynamodb.Client

type RoastModels struct {
	Roast Roast
}

type ReviewModels struct {
	Review Review
}

type Roast struct {
}

type Review struct {
}

func NewRoastModels(dynamo *dynamodb.Client) RoastModels {
	client = dynamo
	return RoastModels{}
}

func NewReviewModels(dynamo *dynamodb.Client) ReviewModels {
	client = dynamo
	return ReviewModels{}
}
