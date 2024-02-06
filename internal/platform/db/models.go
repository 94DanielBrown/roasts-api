package db

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"strconv"
)

var client *dynamodb.Client

type RoastModels struct {
	client *dynamodb.Client
}

type ReviewModels struct {
	client *dynamodb.Client
}

type Roast struct {
	// Exclude id from JSON as it's generated in API from the name
	RoastID    string `dynamodbav:"PK" json:"-"`
	SK         string `dynamodbav:"SK" json:"-"`
	Name       string `dynamodbav:"Name" json:"name"`
	ImageUrl   string `dynamodbav:"ImageUrl" json:"imageUrl"`
	PriceRange string `dynamodbav:"PriceRange" json:"priceRange"`
	// Average rating of 0 is omitted, frontend should take no result as an indication to display that there's no reviews yet
	AverageRating float64 `dynamodbav:"AverageRating" json:"averageRating,omitempty"`
}

type Review struct {
	RoastID  string  `dynamodbav:"PK" json:"-"`
	SortKey  string  `dynamodbav:"SK" json:"-"`
	UserID   string  `dynamodbav:"UserID"`
	Rating   float64 `dynamodbav:"Rating"`
	Comment  string  `dynamodbav:"Comment,omitempty"`
	ReviewID string  `dynamodbav:"ReviewID"`
	ImageUrl string  `dynamodbav:"ImageUrl" json:"imageUrl"`
}

func NewRoastModels(dynamo *dynamodb.Client) RoastModels {
	return RoastModels{client: dynamo}
}

func NewReviewModels(dynamo *dynamodb.Client) ReviewModels {
	return ReviewModels{client: dynamo}
}

func (rm *RoastModels) CreateRoast(roast Roast) error {
	av, err := attributevalue.MarshalMap(roast)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item: av,
		// TODO - Don't hardcode table name
		TableName: aws.String("roasts"),
	}

	_, err = rm.client.PutItem(context.Background(), input)
	return err
}

func (rm *RoastModels) UpdateAverageRating(roastID string, newAverage float64) error {

	// Construct the update input
	input := &dynamodb.UpdateItemInput{
		// TODO - Don't hardcode table name
		TableName: aws.String("roasts"),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: roastID},
			"SK": &types.AttributeValueMemberS{Value: "#PROFILE"},
		},
		UpdateExpression: aws.String("set AverageRating = :r"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":r": &types.AttributeValueMemberN{Value: strconv.FormatFloat(newAverage, 'f', 2, 64)},
		},
	}

	// Execute the update
	_, err := rm.client.UpdateItem(context.Background(), input)
	return err
}

// GetRoastsById gets
func (rm *RoastModels) GetRoastByPrefix(roastPrefix string) (*Roast, error) {
	input := &dynamodb.QueryInput{
		// TODO - Don't hardcode table name
		TableName:              aws.String("roasts"),
		KeyConditionExpression: aws.String("PK = :pkval and begins_with(SK, :skval)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pkval": &types.AttributeValueMemberS{Value: roastPrefix},
			":skval": &types.AttributeValueMemberS{Value: "#PROFILE"},
		},
	}

	// TODO - Should probably pass ctx through rather than use background
	result, err := rm.client.Query(context.Background(), input)
	if err != nil {
		return nil, err
	}

	if len(result.Items) == 0 {
		return nil, nil
	}

	var roast Roast
	err = attributevalue.UnmarshalMap(result.Items[0], &roast)
	if err != nil {
		return nil, err
	}

	return &roast, err
}

// GetAllRoasts performs a scan of dynamodb to get all roasts
func (rm *RoastModels) GetAllRoasts() ([]Roast, error) {
	input := &dynamodb.ScanInput{
		// TODO - Don't hardcode table name
		TableName:        aws.String("roasts"),
		FilterExpression: aws.String("begins_with(PK, :pkval)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pkval": &types.AttributeValueMemberS{Value: "ROAST#"},
		},
	}

	result, err := rm.client.Scan(context.Background(), input)
	if err != nil {
		return nil, err
	}

	var roasts []Roast
	err = attributevalue.UnmarshalListOfMaps(result.Items, &roasts)
	return roasts, err
}

func (rm *ReviewModels) GetReviewsByRoast(roastID string) ([]Review, error) {
	input := &dynamodb.QueryInput{
		// TODO - Don't hardcode table name
		TableName:              aws.String("roasts"),
		KeyConditionExpression: aws.String("PK = :pkval and begins_with(SK, :skval)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pkval": &types.AttributeValueMemberS{Value: roastID},
			":skval": &types.AttributeValueMemberS{Value: "#REVIEW#"},
		},
	}

	result, err := rm.client.Query(context.Background(), input)
	if err != nil {
		return nil, err
	}

	var reviews []Review
	err = attributevalue.UnmarshalListOfMaps(result.Items, &reviews)
	return reviews, err
}
