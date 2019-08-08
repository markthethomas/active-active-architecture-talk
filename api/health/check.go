package main

import (
	"encoding/json"
	"os"
	"time"

	"floqars/shared"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"gopkg.in/underarmour/dynago.v2"
)

type HealthCheckReply struct {
	Region string    `json:"region"`
	Time   time.Time `json:"time"`
}

func checker() (events.APIGatewayProxyResponse, error) {
	reg := os.Getenv("AWS_REGION")
	res, err := shared.DAL.GetItem("region-health", dynago.HashKey("region", reg)).Execute()
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, err
	}

	up := res.Item.GetBool("healthy")
	if !up {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, nil
	}

	body := HealthCheckReply{
		Region: reg,
		Time:   time.Now().UTC(),
	}
	b, err := json.Marshal(body)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, err
	}
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(b),
	}, nil
}

func main() {
	shared.Connect()
	lambda.Start(checker)
}
