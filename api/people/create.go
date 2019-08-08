package main

import (
	"floqars/models"
	"floqars/shared"
	"os"

	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func CreatePerson(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	person := models.Person{}
	if err := json.Unmarshal([]byte(req.Body), &person); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, err
	}
	if err := person.Save(); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, err
	}

	b, err := json.Marshal(&person)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       err.Error(),
		}, err
	}
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(b),
		Headers: map[string]string{
			"floqars-region":              os.Getenv("AWS_REGION"),
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
	}, nil
}

func main() {
	shared.Connect()
	lambda.Start(CreatePerson)
}
