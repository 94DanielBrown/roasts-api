package db

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"log"
	"time"
)

// WaitForDynamoDBTableCreate waits for dynamodb table to be created
func WaitForDynamoDBTableCreate(ctx context.Context, client *dynamodb.Client, tableName string) {

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
