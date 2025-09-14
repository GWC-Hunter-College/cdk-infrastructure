// query_client.go
// A client to streamline connecting to  and querying the deployed RDS database.

package query_client

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"strings"

	"github.com/go-sql-driver/mysql"
)

type QueryClient struct {
	Conn *sql.DB

	user     string
	password string
	address  string
}

// The below go:embed directive will embed all .sql files located in the queries directory
// into the Go binary at compile time, allowing for easy access to SQL query files WITHOUT needing to
// create a separate Dockerfile for each Lambda.

// This may seem like a lot of extra bytes added to the binary, but all files are zipped in the binary
// and only read when it's needed.

// The marginal size increase is negligible compared to the scalability of not needing to manage
// separate Dockerfiles for each Lambda function that needs to run SQL queries.
// For more information, please see: https://pkg.go.dev/embed
// - Anthony 09/08/2025

//go:embed queries/**/*.sql
var sqlFiles embed.FS

var queriesDir, err = fs.Sub(sqlFiles, "queries")

// Factory method to create a new QueryClient instance.
// It requires a context and the ARN of the AWS Secrets Manager secret.
//
// Please remember to close the client connection using `defer queryClient.Close()`
//
// Please also note that if using the client in AWS Lambda, the lambda's IAM role must have
// permissions to access the secret in AWS Secrets Manager.
func NewClient(ctx context.Context, dbSecretArn string, dbName string) (*QueryClient, error) {
	dsn := mysql.NewConfig()

	creds, err := loadSecrets(ctx, dbSecretArn)
	if err != nil {
		return nil, err
	}

	dsn.User = creds.user
	dsn.Passwd = creds.password
	dsn.Addr = creds.address
	dsn.Net = "tcp"
	dsn.TLSConfig = "true"

	conn, err := sql.Open("mysql", dsn.FormatDSN())

	if err != nil {
		return nil, err
	}

	// Verify connection
	if err := conn.Ping(); err != nil {
		return nil, err
	}

	return &QueryClient{
		Conn:     conn,
		user:     creds.user,
		password: creds.password,
		address:  creds.address,
	}, nil
}

// Factory method to create a new QueryClient instance with a custom host address.
// It requires a context and the ARN of the AWS Secrets Manager secret, along with the custom host address.
//
// Please remember to close the client connection using `defer queryClient.Close()`
//
// Please also note that if using the client in AWS Lambda, the lambda's IAM role must have
// permissions to access the secret in AWS Secrets Manager.
func NewClientFromHost(ctx context.Context, dbSecretArn string, dbName string, host string) (*QueryClient, error) {
	dsn := mysql.NewConfig()

	creds, err := loadSecrets(ctx, dbSecretArn)
	if err != nil {
		return nil, err
	}

	dsn.User = creds.user
	dsn.Passwd = creds.password
	dsn.Addr = fmt.Sprintf("%s:%s", host, creds.port)
	dsn.Net = "tcp"
	dsn.DBName = dbName
	dsn.TLSConfig = "true"

	conn, err := sql.Open("mysql", dsn.FormatDSN())

	if err != nil {
		return nil, err
	}

	// Verify connection
	if err := conn.Ping(); err != nil {
		return nil, err
	}

	return &QueryClient{
		Conn:     conn,
		user:     creds.user,
		password: creds.password,
		address:  creds.address,
	}, nil
}

func (qc *QueryClient) ChangeDatabase(dbName string) error {
	_, err := qc.Conn.Exec("USE " + dbName)
	if err != nil {
		return err
	}

	return nil
}

// Execute a raw query string with optional args for placeholders.
func (qc *QueryClient) ExecuteQueryRaw(query string, args ...any) (*sql.Rows, error) {
	return qc.Conn.Query(query, args...)
}

// Execute a query from a file located at `filepath` with optional args for placeholders.
//
// Please note that to execute a SQL file that contains multiple statements that does not expect any rows returned
// (e.g. migrations or bulk table creation), you can use `ExecuteFileBulk` instead.
func (qc *QueryClient) ExecuteQuery(filepath string, args ...any) (*sql.Rows, error) {
	query, err := loadSQLFromFile(filepath)
	if err != nil {
		return nil, err
	}

	rows, err := qc.Conn.Query(query, args...)
	if err != nil {
		return nil, err
	}

	log.Printf("Successfully executed query from file: %s.", filepath)

	return rows, nil
}

// Execute a query from a file located at `filepath` that is expected to return a single row.
func (qc *QueryClient) ExecuteQueryRow(filepath string, args ...any) *sql.Row {
	query, err := loadSQLFromFile(filepath)
	if err != nil {
		return nil
	}

	row := qc.Conn.QueryRow(query, args...)

	log.Printf("Successfully executed query from file: %s.", filepath)

	return row
}

// Execute multiple queries in sequence, each from their own file with optional args for placeholders.
// If any query fails, the execution stops and the error is returned.
//
// The verbose flag can be set to true to log each query execution. Plesae note that this will log
// everything, including sensitive information if the query args contain such information.
func (qc *QueryClient) ExecuteMultipleQueries(queries []Query, verbose bool) (rows []*sql.Rows, err error) {
	for _, query := range queries {
		resultingRows, err := qc.ExecuteQuery(query.Filepath, query.Args...)

		if verbose {
			log.Printf("Executed query from file: %s with args: %v", query.Filepath, query.Args)
		}

		if err != nil {
			return rows, err
		}

		rows = append(rows, resultingRows)
	}

	return rows, nil
}

// Execute a SQL file that contains multiple statements (e.g. for migrations or bulk table creation).
// Each statement is executed in sequence.
//
// Please note that this method does NOT return any rows, meaning that this method should only be used
// to execute statements that do not return rows (e.g. CREATE, INSERT, UPDATE, DELETE).
func (qc *QueryClient) ExecuteFileBulk(filepath string) (results []sql.Result, err error) {
	query, err := loadSQLFromFile(filepath)
	if err != nil {
		return nil, err
	}

	statements := strings.Split(query, ";")

	for _, statement := range statements {
		statement = strings.TrimSpace(statement)
		if statement == "" {
			continue
		}

		result, err := qc.Conn.Exec(statement)
		if err != nil {
			return nil, err
		}

		results = append(results, result)

		log.Printf("Executed statement: %s", statement)
	}

	return results, nil
}

func (qc *QueryClient) Close() error {
	return qc.Conn.Close()
}
