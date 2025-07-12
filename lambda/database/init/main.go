package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"strings"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"

	"github.com/aws/aws-lambda-go/events"
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

/*
	TODO:
	- Implement database initialization logic here.
	- Steps to do this:
		1. Create migration files - Done
		2. Remember to copy them to the image in the Dockerfile
		3. Create a new test database stack to instantiate the RDS instance
		4. In this handler, connect to the RDS instance (probably using AWS Secrets Manager for creds)
		5. Run the migrations to set up the schema

	https://aws.amazon.com/blogs/infrastructure-and-automation/use-aws-cdk-to-initialize-amazon-rds-instances/#:~:text=with%20initialization%20support%3A-,Create%20the,folder%2C%20and%20paste%20the%20following%20content%20inside%3A,-import%20*%20as
*/

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "Database initialization complete",
	}, nil
}

func connectToMySQL(user string, password string, dbName string) (*sql.DB, error) {
	dsn := mysql.NewConfig()
	dsn.User = user
	dsn.Passwd = password
	dsn.DBName = dbName

	db, err := sql.Open("mysql", dsn.FormatDSN())

	if err != nil {
		db.Close()
		log.Print("failed to connect to MySQL:", err)
	}

	log.Println("Connected to MySQL successfully")

	return db, nil
}

// runMigration executes the SQL statements in the given migration file against the provided database connection.
//
// It assumes that your migration files are under a folder "migrations" in the current working directory.
func runMigration(db *sql.DB, filename string) error {
	fileBytes, err := os.ReadFile("migrations/" + filename)
	if err != nil {
		log.Fatalf("Failed to read migration file %s: %v", filename, err)
	}

	statements := strings.Split(string(fileBytes), ";")

	for _, statement := range statements {
		statement = strings.TrimSpace(statement)
		if statement == "" {
			continue
		}

		_, err := db.Exec(statement)
		if err != nil {
			log.Fatalf("Failed to execute statement in %s: %v", filename, err)
		}

		log.Printf("Executed statement: %s", statement)
	}

	return nil
}

func main() {
	// lambda.Start(handler)

	user := "root"
	password := ""
	targetDBName := "PROD"

	initTableMigrationFiles := []string{
		"07_11_2025_create_core_tables_up.sql",
		"07_11_2025_create_member_form_migration_table_up.sql",
		"07_11_2025_create_core_tables_down.sql",
	}

	initDatabaseMigrationFile := "07_11_2025_create_databases_up.sql"

	// Run migrations to create staging and production databases
	createDatabasesDb, err := connectToMySQL(user, password, "")
	if err != nil {
		log.Fatalf("Failed to connect to MySQL to init databases: %v", err)
	}
	defer createDatabasesDb.Close()

	if err = runMigration(createDatabasesDb, initDatabaseMigrationFile); err != nil {
		log.Fatalf("Failed to initialize databases: %v", err)
	}

	// Run migrations to create all tables
	initTablesDb, err := connectToMySQL(user, password, targetDBName)
	if err != nil {
		log.Fatalf("Failed to connect to MySQL to init tables: %v", err)
	}
	defer initTablesDb.Close()

	for _, file := range initTableMigrationFiles {
		err := runMigration(initTablesDb, file)
		if err != nil {
			log.Fatalf("Failed to run migration %s: %v", file, err)
		}
		log.Printf("Migration %s completed successfully", file)
	}
}
