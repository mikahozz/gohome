package dbsync

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

// SetupDbConn creates a database connection using environment variables
func SetupDbConn() (*sql.DB, error) {
	connStr := fmt.Sprintf("host=localhost port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("DB_APP_USER"),
		os.Getenv("DB_APP_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	return db, nil
}
