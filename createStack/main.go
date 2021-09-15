package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	templateUrl := os.Getenv("STACK_TEMPLATE_URL")
	allowedIpRangeValue := os.Getenv("ALLOWED_IP_CIDR_RANGE")
	stackName := request.QueryStringParameters["stackName"]
	allowedIpRangeName := "AllowedIpRange"
	p := types.Parameter{
		ParameterKey:   &allowedIpRangeName,
		ParameterValue: &allowedIpRangeValue,
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		log.Fatal(err)
	}

	client := cloudformation.NewFromConfig(cfg)

	_, err = client.CreateStack(context.TODO(), &cloudformation.CreateStackInput{
		StackName:   &stackName,
		TemplateURL: &templateUrl,
		Parameters:  []types.Parameter{p},
	})

	if err != nil {
		log.Fatal(err)
		return events.APIGatewayProxyResponse{}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 201,
	}, nil
}

func main() {
	lambda.Start(handler)
}
