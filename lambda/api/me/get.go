package main

import (
	// models "cdk-infrastructure/database/models"
	gateway_helpers "cdk-infrastructure/gateway/helpers"
	"cdk-infrastructure/lambda/internal/auth/utils"
	"cdk-infrastructure/utils/query_client"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
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

type resp struct {
	Sub   string `json:"sub"`
	Email string `json:"email,omitempty"`
}

// --- adapter: V1 -> V2 ---
func v1ToV2(r events.APIGatewayProxyResponse) events.APIGatewayV2HTTPResponse {
	// Copy headers to V2 shape
	h := map[string]string{}
	for k, v := range r.Headers {
		h[k] = v
	}
	return events.APIGatewayV2HTTPResponse{
		StatusCode: r.StatusCode,
		Headers:    h,
		Body:       r.Body,
		// If you rely on cookies or multi-value headers, mirror them here as needed.
	}
}

func handler(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	claims := map[string]string{}
	if event.RequestContext.Authorizer != nil &&
		event.RequestContext.Authorizer.JWT != nil &&
		event.RequestContext.Authorizer.JWT.Claims != nil {
		claims = event.RequestContext.Authorizer.JWT.Claims
	}

	sub := claims["sub"]
	email := claims["email"]
	if sub == "" {
		// If the authorizer passed the request, sub should be present; defensively handle missing.
		return jsonResp(http.StatusForbidden, map[string]string{"error": "missing sub in JWT claims"})
	}

	// Enforce/ensure student row exists (or create it); fail the request if this fails.
	if err := utils.RequireStudent(ctx, qc, sub, email); err != nil {
		if errors.Is(err, utils.ErrNoSub) {
			v1, _ := gateway_helpers.NewClientErrorGatewayResponse(
				"missing sub in JWT claims",
				map[string]any{"code": "ERR_NO_SUB"},
			)
			return v1ToV2(v1), nil
		}
		// Any DB/other error -> 500
		v1, _ := gateway_helpers.NewServerErrorGatewayResponse(
			"failed to ensure student",
			map[string]any{
				"code":   "ERR_REQUIRE_STUDENT",
				"detail": err.Error(),
			},
		)
		return v1ToV2(v1), nil
	}

	// Success -> 200 with sub (and email if you want)
	// If you prefer "message + fields", use the success helper:
	// v1, _ := gateway_helpers.NewSuccessGatewayResponse("ok", map[string]any{"sub": sub, "email": email})
	// return v1ToV2(v1), nil

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
