package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	stackName := request.QueryStringParameters["stackName"]

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	client := cloudformation.NewFromConfig(cfg)

	params := cloudformation.DeleteStackInput{
		StackName: &stackName,
	}

	_, err = client.DeleteStack(context.TODO(), &params)

	if err != nil {
		log.Fatal(err)
		return events.APIGatewayProxyResponse{}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 202,
	}, nil
}

func main() {
	lambda.Start(handler)
}
