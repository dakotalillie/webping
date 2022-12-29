package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

func PublishToSNS(record PingRecord) error {
	log.Println("publishing error message to sns")
	_, err := snsClient.Publish(context.TODO(), &sns.PublishInput{
		Message:  aws.String(record.Endpoint + " has failed multiple sequential health checks, indicating it is currently down."),
		Subject:  aws.String(record.Endpoint + " is currently down"),
		TopicArn: aws.String(snsTopic),
	})
	return err
}
