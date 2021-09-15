package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type customEvent struct {
	StackName string `json:"stackName"`
}

func HandleRequest(ctx context.Context, event customEvent) (string, error) {
	roleToAssume := os.Getenv("ROLE_ARN")
	templateUrl := os.Getenv("STACK_TEMPLATE_URL")
	allowedIpRangeValue := os.Getenv("ALLOWED_IP_CIDR_RANGE")
	stackName := event.StackName
	allowedIpRangeName := "AllowedIpRange"
	p := types.Parameter{
		ParameterKey:   &allowedIpRangeName,
		ParameterValue: &allowedIpRangeValue,
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		log.Fatal(err)
	}

	// Create the credentials from AssumeRoleProvider to assume the role
	// referenced by the roleToAssume variable
	stsSvc := sts.NewFromConfig(cfg)
	creds := stscreds.NewAssumeRoleProvider(stsSvc, roleToAssume)

	cfg.Credentials = aws.NewCredentialsCache(creds)

	client := cloudformation.NewFromConfig(cfg)

	_, err = client.CreateStack(context.TODO(), &cloudformation.CreateStackInput{
		StackName:   &stackName,
		TemplateURL: &templateUrl,
		Parameters:  []types.Parameter{p},
	})

	if err != nil {
		log.Fatal(err)
		return "error", err
	}

	return "success", nil
}

func main() {
	lambda.Start(HandleRequest)
}
