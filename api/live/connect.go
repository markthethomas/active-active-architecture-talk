package main

import (
	"floqars/shared"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"gopkg.in/underarmour/dynago.v2"
)

func ConnectHandler(req events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	connID := req.RequestContext.ConnectionID
	if _, err := shared.DAL.PutItem("connected", dynago.Document{
		"connectionid": connID,
		"createdat":    time.Now().UTC().Add(1 * time.Hour),
	}).Execute(); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       err.Error(),
		}, err
	}
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
	}, nil
}

func main() {
	shared.Connect()
	lambda.Start(ConnectHandler)
}
