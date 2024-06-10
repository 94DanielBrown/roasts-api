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
	ImageURL    string `dynamodbav:"ImageURL" json:"imageURL"`
	PriceRange  int    `dynamodbav:"PriceRange" json:"priceRange"`
	Location    string `dynamodbav:"Location" json:"location"`
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
	RoastKey string `dynamodbav:"PK" json:"-"`
	// Using unique RoastID as SK generated from epoch time
	SK             string `dynamodbav:"SK" json:"-"`
	RoastID        string `dynamodbav:"RoastID" json:"roastID"`
	OverallRating  int    `dynamodbav:"OverallRating" json:"overallRating"`
	MeatRating     int    `dynamodbav:"MeatRating" json:"meatRating"`
	PotatoesRating int    `dynamodbav:"PotatoesRating" json:"potatoesRating"`
	VegRating      int    `dynamodbav:"VegRating" json:"vegRating"`
	GravyRating    int    `dynamodbav:"GravyRating" json:"gravyRating"`
	Comment        string `dynamodbav:"Comment,omitempty" json:"comment,omitempty"`
	RoastName      string `dynamodbav:"RoastName" json:"roastName"`
	// When you update your image it needs to update it on all of the users reviews?
	// Like wise if they want to change their displayname ......
	ImageURL    string `dynamodbav:"ImageURL" json:"imageURL"`
	UserID      string `dynamodbav:"UserID" json:"userID"`
	DisplayName string `dynamodbav:"Name" json:"displayName,omitempty"`
	FirstName   string `dynamodbav:"FirstName" json:"firstName,omitempty"`
	LastName    string `dynamodbav:"LastName" json:"lastName,omitempty"`
	DateAdded   int    `dynamodbav:"DateAdded" json:"dateAdded"`
}

type User struct {
	UserKey         string   `dynamodbav:"PK" json:"userKey"`
	SK              string   `dynamodbav:"SK" json:"-"`
	ProfilePhotoUrl string   `dynamodbav:"ProfilePhotoUrl" json:"profilePhotoUrl,omitempty"`
	SavedRoasts     []string `dynamodbav:"SavedRoasts" json:"savedRoasts,omitempty"`
	FirstName       string   `dynamodbav:"FirstName" json:"firstName,omitempty"`
	LastName        string   `dynamodbav:"LastName" json:"lastName,omitempty"`
	DisplayName     string   `dynamodbav:"DisplayName" json:"displayName,omitempty"`
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

func (rm *ReviewModels) RemoveReview(roastKey, reviewKey string) error {
	fmt.Println("PK: ", roastKey)
	fmt.Println("SK: ", reviewKey)
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(rm.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: roastKey},
			"SK": &types.AttributeValueMemberS{Value: reviewKey},
		},
	}
	_, err := rm.client.DeleteItem(context.Background(), input)
	if err != nil {
		return err
	}
	return nil
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
func (um *UserModels) CreateUser(user User) error {
	av, err := attributevalue.MarshalMap(user)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(um.tableName),
		Item:      av,
	}

	_, err = um.client.PutItem(context.Background(), input)
	if err != nil {
		return err
	}

	return nil
}

// UpdateUser updates a user item in the database
func (um *UserModels) UpdateUser(user User) error {
	av, err := attributevalue.MarshalMap(user)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(um.tableName),
		Item:      av,
	}

	_, err = um.client.PutItem(context.Background(), input)
	return err
}

// UpdateSavedRoasts updates the SavedRoasts array for a user identified by userID
func (um *UserModels) UpdateSavedRoasts(userID, roastID string) error {
	// Retrieve the user by userID
	fmt.Println("userID: ", userID)
	userKey := "USER#" + userID
	user, err := um.GetUserByPrefix(userKey)
	if err != nil {
		return fmt.Errorf("error retrieving user: %w", err)
	}
	if user == nil {
		return fmt.Errorf("user not found with userID: %s", userID)
	}
	// Append the roastID to the SavedRoasts array
	user.SavedRoasts = append(user.SavedRoasts, roastID)
	// test logging
	fmt.Println("user.SavedRoasts: ", user.SavedRoasts)
	// Update the user item in the database
	if err := um.UpdateUser(*user); err != nil {
		return fmt.Errorf("error updating users SavedRoasts: %w", err)
	}
	return nil
}

// RemoveSavedRoast removed roastID from the SavedRoasts array for a user identified by userID
func (um *UserModels) RemoveSavedRoast(userID, roastID string) error {
	// Retrieve the user by userID
	fmt.Println("userID: ", userID)
	userKey := "USER#" + userID
	user, err := um.GetUserByPrefix(userKey)
	if err != nil {
		return fmt.Errorf("error retrieving user: %w", err)
	}
	if user == nil {
		return fmt.Errorf("user not found with userID: %s", userID)
	}
	fmt.Println("saved roasts before", user.SavedRoasts)
	// Get the index of the roastID in the SavedRoasts array
	index := -1
	for i, id := range user.SavedRoasts {
		if id == roastID {
			index = i
			fmt.Println("index", index)
			break
		}
	}
	if index == -1 {
		return fmt.Errorf("roastID not found in users SavedRoasts")
	}
	// Remove the roastID from the SavedRoasts array
	user.SavedRoasts = append(user.SavedRoasts[:index], user.SavedRoasts[index+1:]...)
	// Update the user item in the database
	if err := um.UpdateUser(*user); err != nil {
		return fmt.Errorf("error updating users SavedRoasts: %w", err)
	}
	return nil
}

// GetUserReviews retrieves all reviews a user has made from dynamoDB
func (rm *UserModels) GetUserReviews(userID string) ([]Review, error) {
	// TODO - query doesn't work would need secondary index if scales to avoid scanning
	//input := &dynamodb.QueryInput{
	//	TableName:              aws.String(rm.tableName),
	//	KeyConditionExpression: aws.String("PK = :pkval and begins_with(SK, :skval)"),
	//	ExpressionAttributeValues: map[string]types.AttributeValue{
	//		":pkval":  &types.AttributeValueMemberS{Value: "ROAST#"},
	//		":skval":  &types.AttributeValueMemberS{Value: "REVIEW#"},
	//		":userID": &types.AttributeValueMemberS{Value: userID},
	//	},
	//	FilterExpression: aws.String("userID = :userID"),
	//}

	input := &dynamodb.ScanInput{
		TableName:        aws.String(rm.tableName),
		FilterExpression: aws.String("UserID = :userID"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":userID": &types.AttributeValueMemberS{Value: userID},
		},
	}

	// TODO - Should probably pass ctx through rather than use background
	result, err := rm.client.Scan(context.Background(), input)
	if err != nil {
		return nil, err
	}

	if len(result.Items) == 0 {
		return nil, nil
	}

	var reviews []Review
	for _, item := range result.Items {
		var review Review
		err = attributevalue.UnmarshalMap(item, &review)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}
	fmt.Println("reviews: ", reviews)

	return reviews, err
}

// UpdateSettings retrieves all reviews a user has made from dynamoDB
func (um *UserModels) UpdateSettings(userID, displayName, firstName, lastName string) error {
	userKey := "USER#" + userID
	user, err := um.GetUserByPrefix(userKey)
	if err != nil {
		return fmt.Errorf("error retrieving user: %w", err)
	}
	if user == nil {
		return fmt.Errorf("user not found with userID: %s", userID)
	}
	user.DisplayName = displayName
	user.FirstName = firstName
	user.LastName = lastName
	if err := um.UpdateUser(*user); err != nil {
		return fmt.Errorf("error updating users settings: %w", err)
	}
	return nil
}
