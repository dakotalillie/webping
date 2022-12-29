package tests

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/require"
)

func TestWebping(t *testing.T) {
	logger.Log(t, "running build script")
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal("failed to get working directory:", err)
	}
	rootDir := filepath.Dir(cwd)

	cmd := exec.Command("scripts/build.sh")
	cmd.Dir = rootDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatal("failed to run build script:", string(out))
	}
	logger.Log(t, "successfully ran build script")

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../terraform",
		Vars: map[string]interface{}{
			"aws_region":       "us-west-1",
			"enable_ping_cron": false,
			"endpoints":        []string{"https://dakotadacoda.com", "https://laurensettembrino.com"},
			"environment":      "test",
			"stack_name":       "test",
		},
	})

	terraform.Init(t, terraformOptions)
	terraform.WorkspaceSelectOrNew(t, terraformOptions, "test")
	defer terraform.Destroy(t, terraformOptions)
	terraform.Apply(t, terraformOptions)
	logger.Log(t, "waiting to allow all changes to take effect")
	time.Sleep(10 * time.Second)

	functionName := terraform.Output(t, terraformOptions, "ping_lambda_function_name")
	logger.Log(t, "invoking lambda function")
	aws.InvokeFunction(t, "us-west-1", functionName, map[string]interface{}{})
	logger.Log(t, "successfully invoked lambda function")

	terraformOptions.Vars["endpoints"] = []string{"https://dakotadacoda.com/404"}
	terraform.Apply(t, terraformOptions)

	logger.Log(t, "setting up aws clients")
	sqsClient := aws.NewSqsClient(t, "us-west-1")
	snsClient := aws.NewSnsClient(t, "us-west-1")
	logger.Log(t, "successfully set up aws clients")

	logger.Log(t, "creating sqs queue")
	sqsQueueName := "webping-test-queue"
	sqsQueueArn := fmt.Sprintf("arn:aws:sqs:%s:%s:%s", "us-west-1", aws.GetAccountId(t), sqsQueueName)
	snsTopicArn := terraform.Output(t, terraformOptions, "sns_topic_arn")
	sqsQueuePolicy := fmt.Sprintf(`{
	 "Version": "2012-10-17",
	 "Id": "QueuePolicy",
	 "Statement": [
	   {
	     "Effect": "Allow",
	     "Principal": {
	       "Service": "sns.amazonaws.com"
	     },
	     "Action": "sqs:SendMessage",
	     "Resource": "%s",
	     "Condition": {
	       "ArnEquals": {
	         "aws:SourceArn": "%s"
	       }
	     }
	   }
	 ]
	}`, sqsQueueArn, snsTopicArn)
	sqsQueueOut, err := sqsClient.CreateQueue(&sqs.CreateQueueInput{
		QueueName: &sqsQueueName,
		Attributes: map[string]*string{
			"Policy": &sqsQueuePolicy,
		},
	})
	if err != nil {
		t.Fatal("failed to create sqs queue:", err)
	}
	defer func(sqsClient *sqs.SQS, input *sqs.DeleteQueueInput) {
		_, err := sqsClient.DeleteQueue(input)
		if err != nil {
			t.Fatal("failed to delete sqs queue:", err)
		}
	}(sqsClient, &sqs.DeleteQueueInput{QueueUrl: sqsQueueOut.QueueUrl})
	logger.Log(t, "successfully created sqs queue")

	logger.Log(t, "subscribing sqs queue to sns topic")
	snsProtocol := "sqs"
	returnSubscriptionArn := true
	subscribeOut, err := snsClient.Subscribe(&sns.SubscribeInput{
		Endpoint:              &sqsQueueArn,
		Protocol:              &snsProtocol,
		TopicArn:              &snsTopicArn,
		ReturnSubscriptionArn: &returnSubscriptionArn,
	})
	if err != nil {
		t.Fatal("failed to subscribe sqs queue to sns topic:", err)
	}
	defer func(snsClient *sns.SNS, input *sns.UnsubscribeInput) {
		_, err := snsClient.Unsubscribe(input)
		if err != nil {
			t.Fatal("failed to unsubscribe sqs queue from sns topic:", err)
		}
	}(snsClient, &sns.UnsubscribeInput{SubscriptionArn: subscribeOut.SubscriptionArn})
	logger.Log(t, "successfully subscribed sqs queue to sns topic")

	// The sleeps here are important because otherwise the records from the previous invocations will not be
	// successfully retrieved from dynamodb
	logger.Log(t, "invoking lambda repeatedly to meet error threshold")
	aws.InvokeFunction(t, "us-west-1", functionName, map[string]interface{}{})
	time.Sleep(5 * time.Second)
	aws.InvokeFunction(t, "us-west-1", functionName, map[string]interface{}{})
	time.Sleep(5 * time.Second)
	aws.InvokeFunction(t, "us-west-1", functionName, map[string]interface{}{})
	logger.Log(t, "successfully invoked lambda repeatedly to meet error threshold")

	logger.Log(t, "waiting to receive message from sqs queue")
	var waitSeconds int64 = 20 // This is the max value the SDK will allow
	receiveMessageOut, err := sqsClient.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl:        sqsQueueOut.QueueUrl,
		WaitTimeSeconds: &waitSeconds,
	})
	if err != nil {
		t.Fatal("failed to receive message from sqs queue:", err)
	}
	require.Len(t, receiveMessageOut.Messages, 1, "expected sqs queue to receive message from sns")
	logger.Log(t, "successfully received message from sqs queue")
}
