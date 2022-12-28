package e2e_tests

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

func TestWebping(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../terraform",
		Vars: map[string]interface{}{
			"aws_region":  "us-west-1",
			"endpoints":   []string{"https://dakotadacoda.com", "https://laurensettembrino.com"},
			"environment": "test",
			"stack_name":  "test",
		},
	})

	terraform.WorkspaceSelectOrNew(t, terraformOptions, "test")
	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApply(t, terraformOptions)

	functionName := terraform.Output(t, terraformOptions, "ping_lambda_function_name")

	t.Log("invoking lambda function")
	aws.InvokeFunction(t, "us-west-1", functionName, map[string]interface{}{})
	t.Log("successfully invoked lambda function")
}
