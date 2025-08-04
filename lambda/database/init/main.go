package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/go-sql-driver/mysql"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

/*
	To test locally, use the following commands (on MacOS with an M-series chip, change `linux/arm64`
		to `linux/amd64` for x86-based chips):

	1. docker-compose build

	2. docker-compose start

	3. curl "http://localhost:9000/2015-03-31/functions/function/invocations" -d '{}'
*/

var databaseNames = []string{"STAGING", "PROD"}

var initTableMigrationFiles = []string{
	"07_11_2025_create_core_tables_up.sql",
	"07_11_2025_create_member_form_migration_table_up.sql",
}

var initDatabaseMigrationFile = "07_11_2025_create_databases_up.sql"

var (
	once          sync.Once
	user          string
	password      string
	host          string
	databaseName  string
	secretLoadErr error
)

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	err := initDatabase(ctx)
	if err != nil {
		log.Printf("Error initializing database: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("Failed to initialize database: %v", err),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "Database initialization complete",
	}, nil
}

func initDatabase(ctx context.Context) error {
	secretArn, success := os.LookupEnv("DB_SECRET_ARN")
	if !success {
		log.Printf("DB_SECRET_ARN environment variable not set")
		return os.ErrInvalid
	}

	loadSecrets(ctx, secretArn)
	if secretLoadErr != nil {
		log.Printf("Failed to load secrets: %v", secretLoadErr)
		return secretLoadErr
	}

	databaseName = ""

	// Run migrations to create staging and production databases
	mysqlConn, err := connectToMySQL(user, password, databaseName, host)
	if err != nil {
		log.Printf("Failed to connect to MySQL to init databases: %v", err)
		return err
	}
	defer mysqlConn.Close()

	if err = runMigration(mysqlConn, initDatabaseMigrationFile); err != nil {
		log.Printf("Failed to initialize databases: %v", err)
		return err
	}

	// Run migrations to create all tables in staging
	for _, file := range initTableMigrationFiles {
		err := runMigration(mysqlConn, file)
		if err != nil {
			log.Printf("Failed to run migration %s: %v", file, err)
			return err
		}
		log.Printf("Migration %s completed successfully", file)
	}

	return nil
}

func loadSecrets(ctx context.Context, arn string) {
	once.Do(func() {
		cfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			secretLoadErr = err
			return
		}

		sm := secretsmanager.NewFromConfig(cfg)

		out, err := sm.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
			SecretId: &arn,
		})

		if err != nil {
			secretLoadErr = err
			return
		}

		var creds struct {
			User string `json:"username"`
			Pass string `json:"password"`
			Host string `json:"host"`
		}

		if err = json.Unmarshal([]byte(*out.SecretString), &creds); err != nil {
			secretLoadErr = err
			return
		}

		user, password, host = creds.User, creds.Pass, creds.Host
	})

	return
}

func connectToMySQL(user string, password string, dbName string, address string) (*sql.DB, error) {
	dsn := mysql.NewConfig()
	dsn.User = user
	dsn.Passwd = password
	dsn.DBName = dbName
	dsn.Addr = address + ":3306"
	dsn.Net = "tcp"

	db, err := sql.Open("mysql", dsn.FormatDSN())
	if err != nil {
		log.Printf("sql.Open failed: %v", err)
		return nil, err
	}

	log.Printf("sql.Open succeeded, now testing connection with Ping...")
	err = db.Ping()
	if err != nil {
		log.Printf("db.Ping failed: %v", err)
		db.Close()
		return nil, err
	}

	log.Printf("Connected to MySQL successfully (user: %s, db: %s)", user, dbName)
	return db, nil
}

// runMigration executes the SQL statements in the given migration file against the provided database connection.
//
// It assumes that your migration files are under a folder "migrations" in the current working directory.
func runMigration(db *sql.DB, filename string) error {
	fileBytes, err := os.ReadFile("migrations/" + filename)
	if err != nil {
		log.Printf("Failed to read migration file %s: %v", filename, err)
		return err
	}

	statements := strings.Split(string(fileBytes), ";")

	for _, statement := range statements {
		statement = strings.TrimSpace(statement)
		if statement == "" {
			continue
		}

		_, err := db.Exec(statement)
		if err != nil {
			log.Printf("Failed to execute statement in %s: %v", filename, err)
			return err
		}

		log.Printf("Executed statement: %s", statement)
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
