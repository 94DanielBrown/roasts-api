package dynamo

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/94DanielBrown/roasts/pkg/awsconfig"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// Connect to dynamodb
func Connect() (*dynamodb.Client, error) {
	config, err := awsconfig.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	client := dynamodb.NewFromConfig(config)
	return client, nil
}

// Create a dynamodb table
func Create(ctx context.Context, client *dynamodb.Client, tableName string) error {

	_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("id"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("timestamp"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("id"),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String("timestamp"),
				KeyType:       types.KeyTypeRange,
			},
		},
		TableName:   aws.String(tableName),
		BillingMode: types.BillingModePayPerRequest,
	})

	if err != nil {
		return err
	}

	return nil
}

// Exists checks if dynamodb table exists or not
func Exists(ctx context.Context, client *dynamodb.Client, tableName string) (bool, error) {
	p := dynamodb.NewListTablesPaginator(client, nil, func(o *dynamodb.ListTablesPaginatorOptions) {
		o.StopOnDuplicateToken = true
	})

	for p.HasMorePages() {
		out, err := p.NextPage(ctx)
		if err != nil {
			return false, err
		}

		for _, tn := range out.TableNames {
			if tn == tableName {
				return true, nil
			}
		}
	}
	return false, nil
}

// Wait for dynamodb table to be created
func Wait(ctx context.Context, client *dynamodb.Client, tableName string) {

	waiter := dynamodb.NewTableExistsWaiter(client, func(t *dynamodb.TableExistsWaiterOptions) {
		t.MinDelay = 5 * time.Second
		t.MaxDelay = 30 * time.Second
	})

	maxWait := 150 * time.Second

	ti := dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	}
	err := waiter.Wait(ctx, &ti, maxWait)
	if err != nil {
		log.Panic(fmt.Sprintf("time out waiting for table %s to be created: %v", tableName, err))
	}
}
