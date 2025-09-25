package main

import (
	"context"
	"errors"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"cdk-infrastructure/utils/query_client"
	"os"
	"time"
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

func handler(ctx context.Context, e events.CognitoEventUserPoolsPostConfirmation) (events.CognitoEventUserPoolsPostConfirmation, error) {
	// Only act on the confirmation flow we care about.
	if e.TriggerSource != "PostConfirmation_ConfirmSignUp" {
		return e, nil
	}

	sub := e.Request.UserAttributes["sub"]
	email := e.Request.UserAttributes["email"] // may be empty depending on pool config

	if sub == "" {
		log.Printf("PostConfirmation: missing sub; skipping DB write (email=%s)", email)
		return e, nil
	}

	log.Printf("PostConfirmation received: sub=%s email=%s", sub, email)

	// Short timeout so we never stall the Cognito flow.
	cctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	upsertStudentQuery := query_client.NewQuery("students/UPSERT_user.sql", sub, email)

	// If your QueryClient uses positional params, pass sub, email in order.
	res, err := qc.Exec(upsertStudentQuery)
	if err != nil {
		// Distinguish timeout from other errors.
		if errors.Is(err, context.DeadlineExceeded) {
			log.Printf("students upsert: TIMED OUT (sub=%s, email=%s)", sub, email)
		} else if cctx.Err() != nil {
			// If the context was canceled for another reason, log it explicitly.
			log.Printf("students upsert: context error: %v (sub=%s, email=%s)", cctx.Err(), sub, email)
		} else {
			log.Printf("students upsert: DB error: %v (sub=%s, email=%s)", err, sub, email)
		}
		// Keep signup flow going.
		return e, nil
	}

	// Read the affected row count from the result.
	aff, raErr := res.RowsAffected()
	if raErr != nil {
		log.Printf("students upsert: could not read RowsAffected: %v (sub=%s, email=%s)", raErr, sub, email)
		return e, nil
	}

	// MySQL ON DUPLICATE KEY UPDATE semantics:
	// 1 => inserted new row
	// 2 => existing row updated (email changed)
	// 0 => existing row unchanged (same email)
	switch aff {
	case 1:
		log.Printf("students upsert: INSERTED (first-time) sub=%s email=%s", sub, email)
	case 2:
		log.Printf("students upsert: UPDATED (email changed) sub=%s email=%s", sub, email)
	case 0:
		log.Printf("students upsert: NO-OP (already present, same email) sub=%s email=%s", sub, email)
	default:
		log.Printf("students upsert: unexpected RowsAffected=%d sub=%s email=%s", aff, sub, email)
	}

	// Always return the original event to let Cognito continue.
	return e, nil

	// Always return the original event to let Cognito continue.
	return e, nil
}

func main() { lambda.Start(handler) }
