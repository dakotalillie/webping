package tests

import (
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func TestWebping(t *testing.T) {
	if err := godotenv.Load(); err != nil {
		// This should not fail the test, because in CI, these values aren't derived from .env
		logger.Log(t, "failed to load environment variables from .env:", err)
	}

	awsRegion := "us-west-1"
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "terraform",
		Vars: map[string]interface{}{
			"aws_region":       awsRegion,
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
	time.Sleep(5 * time.Second)

	functionName := terraform.Output(t, terraformOptions, "ping_lambda_function_name")
	logger.Log(t, "invoking lambda function")
	invokeLambda(t, awsRegion, functionName)
	logger.Log(t, "successfully invoked lambda function")

	failureEndpoint := "https://dakotadacoda.com/404"
	terraformOptions.Vars["endpoints"] = []string{failureEndpoint}
	terraform.Apply(t, terraformOptions)

	// The sleeps here are important because otherwise the records from the previous invocations will not be
	// successfully retrieved from dynamodb
	logger.Log(t, "invoking lambda repeatedly to meet error threshold")
	invokeLambda(t, awsRegion, functionName)
	time.Sleep(5 * time.Second)
	invokeLambda(t, awsRegion, functionName)
	time.Sleep(5 * time.Second)
	invokeLambda(t, awsRegion, functionName)
	logger.Log(t, "successfully invoked lambda repeatedly to meet error threshold")

	logger.Log(t, "fetching message from Twilio")
	twilioMsg := getTwilioMessage(t, failureEndpoint)
	require.NotNil(t, twilioMsg, "target message must not be nil")
	require.NotNil(t, twilioMsg.Status, "target message status must not be nil")
	require.Equal(t, "delivered", *twilioMsg.Status, "target message must have status of delivered")
}
