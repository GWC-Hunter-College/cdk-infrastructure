// main.go
package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// HTTP API (v2) Lambda authorizer with "Simple response" enabled.
func handler(ctx context.Context, req events.APIGatewayV2CustomAuthorizerV2Request) (events.APIGatewayV2CustomAuthorizerSimpleResponse, error) {
	return events.APIGatewayV2CustomAuthorizerSimpleResponse{
		IsAuthorized: true, // always allow
		Context: map[string]interface{}{
			"user": "test-user",
			"role": "allow-all",
		},
	}, nil
}

func main() {
	lambda.Start(handler)
}
