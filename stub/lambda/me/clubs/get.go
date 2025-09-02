package main

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// ----- local structs (only for this handler) -----

type Club struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type MeClubsResponse struct {
	Clubs []Club `json:"clubs"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// NOTE: This uses HTTP API (payload v2). If you're on REST API (v1),
// swap types to APIGatewayProxyRequest/APIGatewayProxyResponse.
func handle(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	// Stubbed data (e.g., for kyle@example.com)
	resp := MeClubsResponse{
		Clubs: []Club{
			{ID: 2, Name: "Girls Who Code @ Hunter"},
			{ID: 3, Name: "Hunter CS Club"},
		},
	}

	body, err := json.Marshal(resp)
	if err != nil {
		errBody, _ := json.Marshal(ErrorResponse{Error: "failed_to_marshal"})
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       string(errBody),
		}, err
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: 200,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(body),
	}, nil
}

func main() {
	lambda.Start(handle)
}
