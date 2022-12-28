package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func InsertRecordsIntoDynamoDB(records []PingRecord) error {
	requests := make([]types.WriteRequest, len(records))
	for i, record := range records {
		requests[i] = types.WriteRequest{
			PutRequest: &types.PutRequest{
				Item: map[string]types.AttributeValue{
					"Endpoint":       &types.AttributeValueMemberS{Value: record.Endpoint},
					"ExpirationTime": &types.AttributeValueMemberN{Value: fmt.Sprint(record.ExpirationTime)},
					"Result":         &types.AttributeValueMemberS{Value: record.Result},
					"Timestamp":      &types.AttributeValueMemberN{Value: fmt.Sprint(record.Timestamp)},
				},
			},
		}
	}

	log.Println("writing records to dynamodb")
	_, err := dynamodbClient.BatchWriteItem(
		context.TODO(),
		&dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{
				dbTable: requests,
			},
		},
	)
	if err != nil {
		log.Println("failed to write records to dynamodb:", err)
		return err
	}
	return nil
}

func GetPreviousRecords(endpoint string) ([]PingRecord, error) {
	log.Println("retrieving previous records for", endpoint)
	out, err := dynamodbClient.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String(dbTable),
		KeyConditionExpression: aws.String("Endpoint = :endpoint"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":endpoint": &types.AttributeValueMemberS{Value: endpoint},
		},
		Limit:            aws.Int32(3),
		ScanIndexForward: aws.Bool(false), // Descending order
	})
	if err != nil {
		log.Printf("failed to retrieve previous records for %s: %s", endpoint, err)
		return nil, err
	}

	log.Println("unmarshaling previous records for", endpoint)
	var records []PingRecord
	err = attributevalue.UnmarshalListOfMaps(out.Items, &records)
	if err != nil {
		log.Printf("failed to unmarshal previous records for %s: %s", endpoint, err)
		return nil, err
	}

	return records, nil
}
