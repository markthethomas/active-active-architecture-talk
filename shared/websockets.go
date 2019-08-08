package shared

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
)

func SendWSMessage(connID, apiID string, data []byte) error {
	sess, err := session.NewSession(&aws.Config{Region: aws.String(os.Getenv("AWS_REGION"))})
	if err != nil {
		panic(err)
	}
	WSClient := apigatewaymanagementapi.New(sess)
	WSClient.Endpoint = "https://6hy0ivbdee.execute-api.us-west-1.amazonaws.com/prod/"
	payload := apigatewaymanagementapi.PostToConnectionInput{
		Data:         data,
		ConnectionId: aws.String(connID),
	}
	_, err = WSClient.PostToConnection(&payload)
	return err
}
