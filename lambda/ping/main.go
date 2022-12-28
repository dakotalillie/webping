package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

var endpoints []string
var dbTable string
var dynamodbClient *dynamodb.Client
var snsClient *sns.Client
var snsTopic string

func HandleRequest(context.Context, any) error {
	records := PingAllEndpoints(endpoints)

	for _, record := range records {
		if record.Result == "FAIL" {
			prevRecords, err := GetPreviousRecords(record.Endpoint)
			if err != nil {
				return err
			}
			if HasTransitionedIntoErrorState(prevRecords) {
				err = PublishToSNS(record)
				if err != nil {
					return err
				}
			}
		}
	}

	return InsertRecordsIntoDynamoDB(records)
}

func init() {
	endpoints = strings.Split(os.Getenv("ENDPOINTS"), ",")
	dbTable = os.Getenv("DB_TABLE")
	snsTopic = os.Getenv("SNS_TOPIC")

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	dynamodbClient = dynamodb.NewFromConfig(cfg)
	snsClient = sns.NewFromConfig(cfg)
}

func main() {
	lambda.Start(HandleRequest)
}
