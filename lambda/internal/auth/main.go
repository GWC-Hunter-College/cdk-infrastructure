// lambdas/db-test/main.go
package main

//=============================================
// Imports
//=============================================
import (
	"context"       // carries cancellation / deadlines across calls
	"database/sql"  // std-lib DB abstraction
	"encoding/json" // (un)marshal Secrets Manager JSON payload
	"fmt"           // string formatting
	"log"           // structured Lambda logging
	"os"            // read environment variables
	"sync"          // sync.Once for cold-start cache

	"github.com/aws/aws-lambda-go/events" // API Gateway V2 types
	"github.com/aws/aws-lambda-go/lambda" // Lambda bootstrap
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"

	cip "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"

	// loads AWS creds/region
	_ "github.com/go-sql-driver/mysql" // MySQL driver; blank import means “register”
)

//=============================================
// Global cache (warm-start optimisation)
//=============================================

// ❯ sync.Once guarantees loadSecret() runs exactly once per warm container.
var (
	once   sync.Once
	dbUser string
	dbPass string

	dbPort int
	dbName string

	secErr error // remember the first error so later calls return it
)

// =============================================
// loadSecret: fetches & caches username/password from Secrets Manager
// =============================================
func loadSecret(ctx context.Context, arn string) error {
	once.Do(func() { // executes only the first time
		// AWS SDK config (region/creds from env/IAM role)
		cfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			secErr = err
			return
		}
		sm := secretsmanager.NewFromConfig(cfg)

		// Get secret value by ARN
		out, err := sm.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
			SecretId: &arn,
		})
		if err != nil {
			secErr = err
			return
		}

		// Parse
		var tmp struct {
			User string `json:"username"`
			Pass string `json:"password"`

			Port int    `json:"port"`
			Name string `json:"dbname"`
		}
		if err := json.Unmarshal([]byte(*out.SecretString), &tmp); err != nil {
			secErr = err
			return
		}

		// Cache for the life of the container
		dbUser, dbPass, dbPort, dbName =
			tmp.User, tmp.Pass, tmp.Port, tmp.Name
	})
	return secErr // nil on success, first error otherwise
}

// =============================================
//
//	Response helper (serialises to JSON)
//
// =============================================
type resp struct {
	Success bool   `json:"success"`
	Result  int    `json:"result,omitempty"` // present only on OK
	Error   string `json:"error,omitempty"`  // present only on error
}

// =============================================
// EventBridge (CloudTrail) details from rule set
// =============================================

type trailDetail struct {
	EventName         string `json:"eventName"`
	RequestParameters struct {
		UserPoolId string `json:"userPoolId"`
		Username   string `json:"username"`
	} `json:"requestParameters"`
}

// =============================================
// Lambda handler – runs every invocation
// =============================================
func handler(ctx context.Context, ev events.CloudWatchEvent) (any, error) {

	// arn of secret
	arn := os.Getenv("DB_SECRET_ARN") // secret reference

	host := os.Getenv("DB_HOST") // secret reference

	// fet credentionals from secret
	if err := loadSecret(ctx, arn); err != nil {
		return fail(err) // early return on error
	}

	// Parse the CloudTrail detail from EventBridge
	var d trailDetail
	if err := json.Unmarshal(ev.Detail, &d); err != nil {
		log.Println("detail unmarshal:", err)
		return nil, nil
	}

	userPoolId := d.RequestParameters.UserPoolId
	username := d.RequestParameters.Username
	if userPoolId == "" || username == "" {
		log.Printf("missing userPoolId/username; event=%s", d.EventName)
		return nil, nil
	}

	// Look up user attributes to get sub + email
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Println("aws config:", err)
		return nil, nil
	}
	c := cip.NewFromConfig(cfg)

	u, err := c.AdminGetUser(ctx, &cip.AdminGetUserInput{
		UserPoolId: aws.String(userPoolId),
		Username:   aws.String(username),
	})
	if err != nil {
		log.Println("AdminGetUser:", err)
		return nil, nil
	}

	var sub, email string
	for _, a := range u.UserAttributes {
		switch aws.ToString(a.Name) {
		case "sub":
			sub = aws.ToString(a.Value)
		case "email":
			email = aws.ToString(a.Value)
		}
	}
	if sub == "" {
		// Many pools use UUID as username; safe fallback
		sub = username
	}
	// tls stuff needed

	// Build MySQL DSN and connect
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?tls=true", dbUser, dbPass, host, dbPort, dbName)
	db, err := sql.Open("mysql", dsn) // creates connection pool
	if err != nil {
		return fail(err)
	}
	defer db.Close()

	// Idempotent upsert by sub
	const stmt = `
		INSERT INTO students (id, email)
		VALUES (?, ?)
		ON DUPLICATE KEY UPDATE email = VALUES(email)
	`
	if _, err := db.ExecContext(ctx, stmt, sub, email); err != nil {
		log.Println("upsert:", err)
		return nil, nil
	}

	log.Printf("upsert OK: sub=%s email=%s", sub, email)
	return nil, nil

	// // Marshal success JSON
	// ok, _ := json.Marshal(resp{Success: true, Result: result})
	// return events.APIGatewayV2HTTPResponse{
	// 	StatusCode: 200,
	// 	Body:       string(ok),
	// }, nil
}

// =============================================
// error wrapper
// =============================================
func fail(e error) (events.APIGatewayV2HTTPResponse, error) {
	log.Println("handler error:", e) // appears in CloudWatch Logs
	body, _ := json.Marshal(resp{Success: false, Error: e.Error()})
	return events.APIGatewayV2HTTPResponse{
		StatusCode: 500,
		Body:       string(body),
	}, nil
}

// =============================================
// Bootstrap
// =============================================
func main() {
	lambda.Start(handler)
}
