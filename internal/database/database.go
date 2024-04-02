package database

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type RoastModels struct {
	client    *dynamodb.Client
	tableName string
}

type ReviewModels struct {
	client    *dynamodb.Client
	tableName string
}

type UserModels struct {
	client    *dynamodb.Client
	tableName string
}

type Roast struct {
	RoastKey string `dynamodbav:"PK" json:"-"`
	// Using date created as SK
	SK          string `dynamodbav:"SK" json:"-"`
	RoastID     string `dynamodbav:"RoastID" json:"id"`
	Name        string `dynamodbav:"Name" json:"name"`
	ImageUri    string `dynamodbav:"ImageUrl" json:"imageUrl"`
	PriceRange  string `dynamodbav:"PriceRange" json:"priceRange"`
	ReviewCount int    `dynamodbav:"ReviewCount" json:"reviewCount"`
	// Average rating of 0 is omitted, frontend should take no result as an indication to display that there's no reviews yet
	OverallRating           float64 `dynamodbav:"OverallRating" json:"overallRating,omitempty"`
	MeatRating              float64 `dynamodbav:"meatRating" json:"meatRating,omitempty"`
	PotatoesRating          float64 `dynamodbav:"PotatoesRating" json:"potatoesRating,omitempty"`
	VegRating               float64 `dynamodbav:"VegRating" json:"vegRating,omitempty"`
	GravyRating             float64 `dynamodbav:"GravyRating" json:"gravyRating,omitempty"`
	MeatPotatoesRating      float64 `dynamodbav:"MeatPotatoesRating" json:"meatPotatoesRating,omitempty"`
	MeatVegRating           float64 `dynamodbav:"MeatVegRating" json:"meatVegRating,omitempty"`
	MeatGravyRating         float64 `dynamodbav:"MeatGravyRating" json:"meatGravyRating,omitempty"`
	PotatoesVegRating       float64 `dynamodbav:"PotatoesVegRating" json:"potatoesVegRating,omitempty"`
	PotatoesGravyRating     float64 `dynamodbav:"PotatoesGravyRating" json:"potatoesGravyRating,omitempty"`
	VegGravyRating          float64 `dynamodbav:"VegGravyRating" json:"vegGravyRating,omitempty"`
	MeatPotatoesVegRating   float64 `dynamodbav:"MeatPotatoesVegRating" json:"meatPotatoesVegRating,omitempty"`
	MeatPotatoesGravyRating float64 `dynamodbav:"MeatPotatoesGravyRating" json:"meatPotatoesGravyRating,omitempty"`
	MeatVegGravyRating      float64 `dynamodbav:"MeatVegGravyRating" json:"meatVegGravyRating,omitempty"`
}

type Review struct {
	// TODO - will need a unique RoastID returned per review for FlatList key
	RoastKey string `dynamodbav:"PK" json:"-"`
	// Using unique RoastID as SK generated from epoch time
	SK             string `dynamodbav:"SK" json:"-"`
	RoastID        string `dynamodbav:"RoastID" json:"id"`
	OverallRating  int    `dynamodbav:"OverallRating" json:"overallRating"`
	MeatRating     int    `dynamodbav:"MeatRating" json:"meatRating"`
	PotatoesRating int    `dynamodbav:"PotatoesRating" json:"potatoesRating"`
	VegRating      int    `dynamodbav:"VegRating" json:"vegRating"`
	GravyRating    int    `dynamodbav:"GravyRating" json:"gravyRating"`
	Comment        string `dynamodbav:"Comment,omitempty" json:"comment,omitempty"`
	RoastName      string `dynamodbav:"RoastName" json:"roastName"`
	ImageUrl       string `dynamodbav:"ImageUrl" json:"imageUrl"`
	UserID         string `dynamodbav:"userID" json:"userID"`
	FirstName      string `dynamodbav:"FirstName" json:"firstName"`
	LastName       string `dynamodbav:"LastName" json:"lastName"`
}

type User struct {
	UserKey         string   `dynamodbav:"PK" json:"userKey"`
	SK              string   `dynamodbav:"SK" json:"-"`
	ProfilePhotoUrl string   `dynamodbav:"ProfilePhotoUrl" json:"profilePhotoUrl,omitempty"`
	SavedRoasts     []string `dynamodbav:"SavedRoasts" json:"savedRoasts,omitempty"`
}

func NewRoastModels(dynamo *dynamodb.Client) RoastModels {
	tn := os.Getenv("TABLE_NAME")
	return RoastModels{client: dynamo, tableName: tn}
}

func NewReviewModels(dynamo *dynamodb.Client) ReviewModels {
	tn := os.Getenv("TABLE_NAME")
	return ReviewModels{client: dynamo, tableName: tn}
}

func NewUserModels(dynamo *dynamodb.Client) UserModels {
	tn := os.Getenv("TABLE_NAME")
	return UserModels{client: dynamo, tableName: tn}
}

func (rm *RoastModels) CreateRoast(roast Roast) error {
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

func (rm *RoastModels) UpdateRoast(roast *Roast) error {
	updateExpr := "set OverallRating = :or, MeatRating = :mr, PotatoesRating = :pr, VegRating = :vr, GravyRating = :gr, ReviewCount = :rc, " +
		"MeatPotatoesRating = :mpr, MeatVegRating = :mvr, MeatGravyRating = :mgr, PotatoesVegRating = :pvr, " +
		"PotatoesGravyRating = :pgr, VegGravyRating = :vgr, MeatPotatoesVegRating = :mprv, MeatPotatoesGravyRating = :mpgr, MeatVegGravyRating = :mvg"

	exprAttrValues, err := attributevalue.MarshalMap(map[string]interface{}{
		":or":   roast.OverallRating,
		":mr":   roast.MeatRating,
		":pr":   roast.PotatoesRating,
		":vr":   roast.VegRating,
		":gr":   roast.GravyRating,
		":mpr":  roast.MeatPotatoesRating,
		":mvr":  roast.MeatVegRating,
		":mgr":  roast.MeatGravyRating,
		":pvr":  roast.PotatoesVegRating,
		":pgr":  roast.PotatoesGravyRating,
		":vgr":  roast.VegGravyRating,
		":mprv": roast.MeatPotatoesVegRating,
		":mpgr": roast.MeatPotatoesGravyRating,
		":mvg":  roast.MeatVegGravyRating,
		":rc":   roast.ReviewCount,
	})
	// TODO - Wrap errors up stack
	if err != nil {
		return fmt.Errorf("error marshalling attribute values for update: %w", err)
	}

	// Construct the input for the UpdateItem operation
	input := &dynamodb.UpdateItemInput{
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: roast.RoastKey},
			"SK": &types.AttributeValueMemberS{Value: roast.SK},
		},
		TableName:                 aws.String(rm.tableName),
		UpdateExpression:          aws.String(updateExpr),
		ExpressionAttributeValues: exprAttrValues,
		ReturnValues:              types.ReturnValueUpdatedNew,
	}

	_, err = rm.client.UpdateItem(context.Background(), input)
	if err != nil {
		return err
	}
	return nil
}

// GetRoastByPrefix retrieves a roast by its prefix
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

	var items []Roast
	err = attributevalue.UnmarshalListOfMaps(result.Items, &items)
	if err != nil {
		return nil, err
	}

	var roasts []Roast
	for _, item := range items {
		if strings.HasPrefix(item.SK, "PROFILE") {
			roasts = append(roasts, item)
		}
	}
	return roasts, nil
}

func (rm *ReviewModels) CreateReview(review Review) error {
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

func (rm *ReviewModels) GetReviewsByRoast(roastKey string) ([]Review, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(rm.tableName),
		KeyConditionExpression: aws.String("PK = :pkval and begins_with(SK, :skval)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pkval": &types.AttributeValueMemberS{Value: roastKey},
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

// GetUserByPrefix retrieves a user through userID
func (rm *UserModels) GetUserByPrefix(userPrefix string) (*User, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(rm.tableName),
		KeyConditionExpression: aws.String("PK = :pkval and begins_with(SK, :skval)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pkval": &types.AttributeValueMemberS{Value: userPrefix},
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

	var user User
	err = attributevalue.UnmarshalMap(result.Items[0], &user)
	if err != nil {
		return nil, err
	}

	return &user, err
}

// CreateUser creates a new user in DynamoDB
func (rm *UserModels) CreateUser(user User) error {
	av, err := attributevalue.MarshalMap(user)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(rm.tableName),
		Item:      av,
	}

	_, err = rm.client.PutItem(context.Background(), input)
	if err != nil {
		return err
	}

	return nil
}
