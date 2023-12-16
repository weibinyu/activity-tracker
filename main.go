package main

import (
    "time"
    "context"
	"log"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/credentials"
)

type Activity_Types int64

const (
    Code Activity_Types = iota
    Rest
    BJJ
    Entertainment
)

type activity struct {
    ID string  `json:"id"`
    Name string  `json:"name"`
    Date time.Time   `json:"date"`
    TotalMinutes int     `json:"totalMinutes"`
    CreatorID string     `json:"creatorID"`
    Type Activity_Types  `json:"type"`
}

var activities = []activity{
    { 
        ID:"1", Name: "Programming", Date: time.Now(),
        CreatorID: "user",TotalMinutes: 60,
    },
    { 
        ID:"2", Name: "Programming", Date: time.Now(),
        CreatorID: "user",TotalMinutes: 60,
    },
    { 
        ID:"3", Name: "Programming", Date: time.Now(),
        CreatorID: "user",TotalMinutes: 60,
    },
}

func getActivities (c *gin.Context){
    c.HTML(http.StatusOK,"index.html",activities)
}

func postActivity (c *gin.Context){
    var newActivity activity 
    if err := c.BindJSON(&newActivity); err != nil {
        return
    }

    activities = append(activities, newActivity)
    c.JSON(http.StatusCreated, activities)
}

func createDatabase(){
   localdb,err := newclient("local")

   if err != nil {
    log.Printf("Wait for table exists failed. Here's why: %v\n", err)
    return
    }
    
   local := TableBasics{
    localdb,
    "Activities", 
   }

   local.CreateActivitiesTable()
}

func main(){
    router := gin.Default()
    router.LoadHTMLGlob("html/*")
    router.GET("/activities", getActivities)
    router.POST("/activities",postActivity)
    router.Run("localhost:8080")
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

func (basics TableBasics) CreateActivitiesTable() (*types.TableDescription, error) {
	var tableDesc *types.TableDescription
	table, err := basics.DynamoDbClient.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
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
