package main

import (
	"floqars/models"
	"floqars/shared"

	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/mmcloughlin/geohash"
)

func GetPeople(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	latStr, lngStr := req.QueryStringParameters["lat"], req.QueryStringParameters["lng"]
	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, err
	}

	lng, err := strconv.ParseFloat(lngStr, 64)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, err
	}

	hash := geohash.EncodeInt(float64(lat), float64(lng))
	fmt.Println(hash)
	people := []models.Person{}
	res, err := shared.DAL.Scan("people").Execute()
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, err
	}

	for _, p := range res.Items {
		person := models.Person{}
		person.FromDocument(p)
		people = append(people, person)
	}

	b, err := json.Marshal(&people)
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
	lambda.Start(GetPeople)
}
