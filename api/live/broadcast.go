package main

import (
	"floqars/models"
	"floqars/shared"
	"log"

	"encoding/json"
	"sync"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/mmcloughlin/geohash"
	"gopkg.in/underarmour/dynago.v2"
)

func BroadcastHandler(req events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	currentConnID := req.RequestContext.ConnectionID
	person := models.Person{}
	if err := json.Unmarshal([]byte(req.Body), &person); err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: err.Error()}, err
	}

	wg := sync.WaitGroup{}
	res, err := shared.DAL.Scan("connected").FilterExpression("connectionid <> :conn", dynago.Param{":conn", currentConnID}).Limit(1000).Execute()
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: err.Error()}, err
	}
	count := len(res.Items)

	// no one else connected :/
	if count == 0 {
		return events.APIGatewayProxyResponse{StatusCode: 200}, err
	}

	// update the current person's location
	// person := models.Person{}
	// person.GetByEmail(person.Email)
	hash := geohash.EncodeInt(float64(person.Lat), float64(person.Lng))
	person.GeoHash = hash
	person.Lat = person.Lat
	person.Lng = person.Lng
	if err := person.Save(); err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: err.Error()}, err
	}

	// get people out of the DB
	// (if we were doing this for real we wouldn't scan on every location update...)
	people := []models.Person{}
	peopleRes, err := shared.DAL.Scan("people").Execute()
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: err.Error()}, err
	}
	for _, p := range peopleRes.Items {
		person := models.Person{}
		person.FromDocument(p)
		people = append(people, person)
	}
	b, err := json.Marshal(people)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: err.Error()}, err
	}

	wg.Add(count)

	// emit to everyone...
	for _, conn := range res.Items {
		go func() {
			defer wg.Done()
			if err := shared.SendWSMessage(conn.GetString("connectionid"), req.RequestContext.APIID, b); err != nil {
				log.Println(err)
				panic(err)
			}
		}()
	}
	wg.Wait()
	return events.APIGatewayProxyResponse{StatusCode: 200}, err
}

func main() {
	shared.Connect()
	lambda.Start(BroadcastHandler)
}
