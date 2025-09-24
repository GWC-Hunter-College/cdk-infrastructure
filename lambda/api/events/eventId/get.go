package main

import (
	"cdk-infrastructure/database/models"
	gateway_helpers "cdk-infrastructure/gateway/helpers"
	"cdk-infrastructure/utils/query_client"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	qc *query_client.QueryClient
)

func init() {
	dbName := os.Getenv("DB_NAME")
	arn := os.Getenv("DB_SECRET_ARN")
	host := os.Getenv("DB_HOST")

	client, err := query_client.NewClientFromHost(context.Background(), arn, dbName, host)

	if err != nil {
		log.Printf("Error creating query client: %v", err)

		panic(err)
	}

	qc = client
}

type SQLSchema struct {
	ClubID    string  `db:"club_id"`
	IsOwner   bool    `db:"is_owner"`
	ObjectKey *string `db:"object_key"`
	models.Event
}

type Club struct {
	ClubID       string `json:"id"`
	ThumbnailUrl string `json:"thumbnailUrl"`
}

type OwnerSchema struct {
	Owner      Club   `json:"owner"`
	Associates []Club `json:"associates"`
}

type ResponseSchema struct {
	models.Event
	OwnerSchema
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	eventId := request.PathParameters["eventId"]

	if eventId == "" {
		return gateway_helpers.NewClientErrorGatewayResponse("Missing eventId path parameter", nil)
	}

	events := []SQLSchema{}
	selectEventsQuery := query_client.NewQuery("events/SELECT_event_by_id.sql", "posted", eventId)
	err := qc.Select(&events, selectEventsQuery)

	if err != nil {
		log.Printf("Error executing query: %v", err)

		return gateway_helpers.NewServerErrorGatewayResponse("Could not fetch events: "+err.Error(), nil)
	}

	if len(events) == 0 {
		return gateway_helpers.NewNotFoundGatewayResponse(fmt.Sprintf("No events found with eventId %s", eventId), nil)
	}

	// Populate resulting events
	responseEvent := ResponseSchema{
		Event: events[0].Event,
		OwnerSchema: OwnerSchema{
			Owner:      Club{},
			Associates: []Club{},
		},
	}

	for _, event := range events {
		var isOwner bool = event.IsOwner

		club := Club{
			ClubID: event.ClubID,
			// This is a placeholder thumbnail URL
			// TODO: Get a new presigned URL when club profile images are implemented
			ThumbnailUrl: "https://media.istockphoto.com/id/1495088043/vector/user-profile-icon-avatar-or-person-icon-profile-picture-portrait-symbol-default-portrait.jpg?s=612x612&w=0&k=20&c=dhV2p1JwmloBTOaGAtaA3AW1KSnjsdMt7-U_3EZElZ0=",
		}

		if isOwner {
			responseEvent.Owner = club
		} else {
			responseEvent.Associates = append(responseEvent.Associates, club)
		}
	}

	response := map[string]any{
		"event": responseEvent,
	}

	return gateway_helpers.NewSuccessGatewayResponse(fmt.Sprintf("Succesfully fetched %d events", len(events)), response)
}

func main() {
	lambda.Start(handler)
}
