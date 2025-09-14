package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, e events.CognitoEventUserPoolsPostConfirmation) (events.CognitoEventUserPoolsPostConfirmation, error) {
	attrs := e.Request.UserAttributes
	sub := attrs["sub"]
	email := attrs["email"] // may be empty if not collected

	// TODO: write (sub, email) into your DB
	log.Printf("Confirm signup for sub=%s email=%s", sub, email)

	return e, nil
}

func main() { lambda.Start(handler) }
