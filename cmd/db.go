package cmd

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const (
	TableName = "url_store"
)

var db *dynamodb.Client

func InitDB() error {
	db = dynamodb.NewFromConfig(aws.Config{
		BaseEndpoint: aws.String("http://localhost:5100"),
	})

	if !tableExists(context.Background(), TableName) {
		_, err := db.CreateTable(context.Background(), &dynamodb.CreateTableInput{
			AttributeDefinitions: []types.AttributeDefinition{{
				AttributeName: aws.String("url"),
				AttributeType: types.ScalarAttributeTypeS,
			}, {
				AttributeName: aws.String("shortUrl"),
				AttributeType: types.ScalarAttributeTypeS,
			}},
			KeySchema: []types.KeySchemaElement{
				{
					AttributeName: aws.String("shortUrl"),
					KeyType:       types.KeyTypeHash,
				},
			},
			TableName:   aws.String(TableName),
			BillingMode: types.BillingModePayPerRequest,
		})
		if err != nil {
			log.Printf("Couldn't create table %v. Here's why: %v\n", TableName, err)
		} else {
			waiter := dynamodb.NewTableExistsWaiter(db)
			err = waiter.Wait(context.Background(), &dynamodb.DescribeTableInput{
				TableName: aws.String(TableName),
			}, 5*time.Minute)
			if err != nil {
				log.Printf("Wait for table exists failed. Here's why: %v\n", err)
			}
		}
		return err
	}

	return nil
}

func GetDB() *dynamodb.Client {
	return db
}

func tableExists(ctx context.Context, tableName string) bool {
	exists := true
	_, err := db.DescribeTable(
		ctx, &dynamodb.DescribeTableInput{TableName: aws.String(tableName)},
	)
	if err != nil {
		var notFoundEx *types.ResourceNotFoundException
		if errors.As(err, &notFoundEx) {
			log.Printf("Table %v does not exist.\n", tableName)
			err = nil
		} else {
			log.Printf("Couldn't determine existence of table %v. Here's why: %v\n", tableName, err)
		}
		exists = false
	}
	return exists
}
