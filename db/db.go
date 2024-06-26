package db

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type activity_types struct {
	userId        string `json:"userId"`
	activity_type string `json:"activity_type"`
}

type ActivityRecord struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	Date         time.Time      `json:"date"`
	TotalMinutes int            `json:"totalMinutes"`
	CreatorID    string         `json:"creatorID"`
	Type         activity_types `json:"activity_type"`
}

type User struct {
	id       string `json:"id"`
	password string `json:"password"`
}

func createDatabase() {
	localdb, err := newclient("local")

	if err != nil {
		log.Printf("Wait for table exists failed. Here's why: %v\n", err)
		return
	}

	local := TableBasics{
		localdb,
		"Activities",
	}

	local.CreateActivitiyRecordTable()
}

func newclient(profile string) (*dynamodb.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("localhost"),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: "http://localhost:8000"}, nil
			})),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID: "abcd", SecretAccessKey: "a1b2c3", SessionToken: "",
				Source: "Mock credentials used above for local instance",
			},
		}),
	)
	if err != nil {
		return nil, err
	}

	c := dynamodb.NewFromConfig(cfg)
	return c, nil
}

type TableBasics struct {
	DynamoDbClient *dynamodb.Client
	TableName      string
}

func (basics TableBasics) CreateUserTable() (*types.TableDescription, error) {
	var tableDesc *types.TableDescription
	table, err :=
		basics.DynamoDbClient.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
			AttributeDefinitions: []types.AttributeDefinition{{
				AttributeName: aws.String("date"),
				AttributeType: types.ScalarAttributeTypeN,
			}},
			KeySchema: []types.KeySchemaElement{{
				AttributeName: aws.String("date"),
				KeyType:       types.KeyTypeHash,
			}},
			TableName: aws.String(basics.TableName),
			ProvisionedThroughput: &types.ProvisionedThroughput{
				ReadCapacityUnits:  aws.Int64(10),
				WriteCapacityUnits: aws.Int64(10),
			},
		})
	if err != nil {
		log.Printf("Couldn't create table %v. Here's why: %v\n", basics.TableName, err)
	} else {
		waiter := dynamodb.NewTableExistsWaiter(basics.DynamoDbClient)
		err = waiter.Wait(context.TODO(), &dynamodb.DescribeTableInput{
			TableName: aws.String(basics.TableName)}, 5*time.Minute)
		if err != nil {
			log.Printf("Wait for table exists failed. Here's why: %v\n", err)
		}
		tableDesc = table.TableDescription
	}
	return tableDesc, err
}

func (basics TableBasics) CreateActivitiyRecordTable() (*types.TableDescription, error) {
	var tableDesc *types.TableDescription
	table, err := basics.DynamoDbClient.CreateTable(context.TODO(),
		&dynamodb.CreateTableInput{
			AttributeDefinitions: []types.AttributeDefinition{{
				AttributeName: aws.String("id"),
				AttributeType: types.ScalarAttributeTypeS,
			},
				{
					AttributeName: aws.String("date"),
					AttributeType: types.ScalarAttributeTypeS, // data type descriptor: S == string
				},
			},
			KeySchema: []types.KeySchemaElement{
				{
					AttributeName: aws.String("id"),
					KeyType:       types.KeyTypeHash,
				},
				{
					AttributeName: aws.String("date"),
					KeyType:       types.KeyTypeRange,
				},
			},
			TableName: aws.String(basics.TableName),
			ProvisionedThroughput: &types.ProvisionedThroughput{
				ReadCapacityUnits:  aws.Int64(10),
				WriteCapacityUnits: aws.Int64(10),
			},
		})
	if err != nil {
		log.Printf("Couldn't create table %v. Here's why: %v\n", basics.TableName, err)
	} else {
		waiter := dynamodb.NewTableExistsWaiter(basics.DynamoDbClient)
		err = waiter.Wait(context.TODO(), &dynamodb.DescribeTableInput{
			TableName: aws.String(basics.TableName)}, 5*time.Minute)
		if err != nil {
			log.Printf("Wait for table exists failed. Here's why: %v\n", err)
		}
		tableDesc = table.TableDescription
	}
	return tableDesc, err
}
