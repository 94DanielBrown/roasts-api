package db

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var client *dynamodb.Client

type RoastModels struct {
	client *dynamodb.Client
}

type ReviewModels struct {
	Review Review
}

type Roast struct {
	// Exclude id from JSON as it's generated in API from the name
	RoastID    string `dynamodbav:"PK" json:"-"`
	PriceRange string `dynamodbav:"SK" json:"priceRange"`
	Name       string `dynamodbav:"Name" json:"name"`
	ImageUrl   string `dynamodbav:"ImageUrl" json:"imageUrl"`
}

type Review struct {
}

func NewRoastModels(dynamo *dynamodb.Client) RoastModels {
	return RoastModels{client: dynamo}
}

func NewReviewModels(dynamo *dynamodb.Client) ReviewModels {
	client = dynamo
	return ReviewModels{}
}

func (rm *RoastModels) CreateRoast(roast Roast) error {
	av, err := attributevalue.MarshalMap(roast)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("roc"),
	}

	_, err = rm.client.PutItem(context.Background(), input)
	return err
}
