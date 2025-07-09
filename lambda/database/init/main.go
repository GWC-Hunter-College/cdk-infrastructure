package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

/*
	To test locally, use the following commands (on MacOS with an M-series chip, change `linux/arm64`
		to `linux/amd64` for x86-based chips):

	1. docker buildx build --platform linux/arm64 --provenance=false -t gwc-database-init:test .

	2. docker run --platform linux/arm64 -d -v ~/.aws-lambda-rie:/aws-lambda -p 9000:8080 \
		--entrypoint /aws-lambda/aws-lambda-rie \
		gwc-database-init:test \
		/main

	3. curl "http://localhost:9000/2015-03-31/functions/function/invocations" -d '{}'
*/

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "Database initialization complete",
	}, nil
}

func main() {
	lambda.Start(handler)
}
