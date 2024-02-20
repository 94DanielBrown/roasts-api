package db

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"os"
	"strconv"
)

var client *dynamodb.Client

type RoastModels struct {
	client    *dynamodb.Client
	tableName string
}

type ReviewModels struct {
	client    *dynamodb.Client
	tableName string
}

type Roast struct {
	// Exclude id from JSON as it's generated in API from the name
	RoastID string `dynamodbav:"PK" json:"-"`
	// Using date created as SK
	SK         string `dynamodbav:"SK" json:"-"`
	Name       string `dynamodbav:"Name" json:"name"`
	ImageUrl   string `dynamodbav:"ImageUrl" json:"imageUrl"`
	PriceRange string `dynamodbav:"PriceRange" json:"priceRange"`
	// Average rating of 0 is omitted, frontend should take no result as an indication to display that there's no reviews yet
	AverageRating float64 `dynamodbav:"AverageRating" json:"averageRating,omitempty"`
}

type Review struct {
	RoastID string `dynamodbav:"PK" json:"-"`
	// Using unique ID as SK generated from epoch time
	SK        string `dynamodbav:"SK" json:"-"`
	UserID    string `dynamodbav:"userID"`
	Rating    int    `dynamodbav:"rating"`
	Comment   string `dynamodbav:"Comment,omitempty"`
	RoastName string `dynamodbav:"RoastName" json:"roastName"`
	ImageUrl  string `dynamodbav:"ImageUrl" json:"imageUrl"`
}

func NewRoastModels(dynamo *dynamodb.Client) RoastModels {
	tn := os.Getenv("TABLE_NAME")
	return RoastModels{client: dynamo, tableName: tn}
}

func NewReviewModels(dynamo *dynamodb.Client) ReviewModels {
	tn := os.Getenv("TABLE_NAME")
	return ReviewModels{client: dynamo, tableName: tn}
}

func (rm *RoastModels) CreateRoast(roast Roast) error {
	fmt.Println("tablename: ", rm.tableName)
	av, err := attributevalue.MarshalMap(roast)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(rm.tableName),
	}

	_, err = rm.client.PutItem(context.Background(), input)
	return err
}

func (rm *RoastModels) UpdateAverageRating(roastID string, newAverage float64) error {

	// Construct the update input
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(rm.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: roastID},
			"SK": &types.AttributeValueMemberS{Value: "PROFILE"},
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
		TableName:              aws.String(rm.tableName),
		KeyConditionExpression: aws.String("PK = :pkval and begins_with(SK, :skval)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pkval": &types.AttributeValueMemberS{Value: roastPrefix},
			":skval": &types.AttributeValueMemberS{Value: "PROFILE"},
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
		TableName:        aws.String(rm.tableName),
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

func (rm *ReviewModels) CreateReview(review Review) error {
	fmt.Println("tablename: ", rm.tableName)
	av, err := attributevalue.MarshalMap(review)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(rm.tableName),
	}

	_, err = rm.client.PutItem(context.Background(), input)
	return err
}

func (rm *ReviewModels) GetReviewsByRoast(roastID string) ([]Review, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(rm.tableName),
		KeyConditionExpression: aws.String("PK = :pkval and begins_with(SK, :skval)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pkval": &types.AttributeValueMemberS{Value: roastID},
			":skval": &types.AttributeValueMemberS{Value: "REVIEW#"},
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
