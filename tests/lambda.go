package tests

import (
	"context"
	"encoding/base64"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/stretchr/testify/require"
)

func invokeLambda(t *testing.T, region, name string) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	require.NoError(t, err, "failed to load aws config")

	client := lambda.NewFromConfig(cfg)
	res, err := client.Invoke(context.TODO(), &lambda.InvokeInput{
		FunctionName: aws.String(name),
		LogType:      types.LogTypeTail,
		Payload:      []byte("{}"),
	})
	require.NoError(t, err, "failed to invoke lambda function")

	logs, err := base64.StdEncoding.DecodeString(*res.LogResult)
	require.NoError(t, err, "failed to decode lambda logs")
	logger.Log(t, "logs from lambda invocation:", string(logs))
}
