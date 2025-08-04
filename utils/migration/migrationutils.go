package migrationutils

import (
	"database/sql"
	"log"
	"os"
	"strings"
)

// RunMigration executes the SQL statements in the given migration file against the provided database connection.
//
// It assumes that your migration files are under a folder "migrations" in the current working directory.
func RunMigration(db *sql.DB, filename string) error {
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
