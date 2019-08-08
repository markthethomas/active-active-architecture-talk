package main

import (
	"floqars/shared"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"gopkg.in/underarmour/dynago.v2"
)

func DisconnectHandler(req events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	connID := req.RequestContext.ConnectionID
	if _, err := shared.DAL.DeleteItem("connected", dynago.Document{"connectionid": connID}).Execute(); err != nil {
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
	lambda.Start(DisconnectHandler)
}
