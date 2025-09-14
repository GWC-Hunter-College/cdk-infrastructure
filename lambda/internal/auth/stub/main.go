package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type resp struct {
	Sub   string `json:"sub"`
	Email string `json:"email,omitempty"`
}

func handler(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	claims := map[string]string{}
	if event.RequestContext.Authorizer != nil &&
		event.RequestContext.Authorizer.JWT != nil &&
		event.RequestContext.Authorizer.JWT.Claims != nil {
		claims = event.RequestContext.Authorizer.JWT.Claims
	}

	// claims["sub"] and claims["email"] should be strings if present
	sub := claims["sub"]
	email := claims["email"]

	if sub == "" {
		// If the authorizer passed the request, sub should be present; defensively handle missing.
		return jsonResp(http.StatusForbidden, map[string]string{"error": "missing sub in JWT claims"})
	}

	return jsonResp(http.StatusOK, resp{Sub: sub, Email: email})
}

func jsonResp(status int, v any) (events.APIGatewayV2HTTPResponse, error) {
	b, _ := json.Marshal(v)
	return events.APIGatewayV2HTTPResponse{
		StatusCode: status,
		Headers: map[string]string{
			"content-type":                 "application/json",
			"access-control-allow-origin":  "*", // dev CORS; tighten in prod
			"access-control-allow-headers": "authorization,content-type",
		},
		Body: string(b),
	}, nil
}

func main() { lambda.Start(handler) }
