package dbsync

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// setupTestDB creates and returns a PostgreSQL database connection for testing
func setupTestDB(t *testing.T) *sql.DB {
	// Load environment variables
	err := godotenv.Load("../.env")
	if err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	// Get environment variables
	host := "localhost"
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("DB_APP_USER")
	password := os.Getenv("DB_APP_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")

	// Create connection string
	connStr := "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable"
	connStr = "host=" + host + " port=" + port + " user=" + user + " password=" + password + " dbname=" + dbname + " sslmode=disable"

	// Open database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("Error connecting to database: %v", err)
	}

	return db
}

func TestPostgresHealth(t *testing.T) {
	// Skip if not running integration tests
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db := setupTestDB(t)
	defer db.Close()

	// Test connection with a simple query
	err := db.Ping()
	if err != nil {
		t.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
}

func TestPostgresWriteAndQuery(t *testing.T) {
	// Skip if not running integration tests
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db := setupTestDB(t)
	defer db.Close()

	ctx := context.Background()

	// Test if we can write to the database
	timestamp := time.Now()
	_, err := db.ExecContext(ctx,
		`INSERT INTO measurements (timestamp, sensor_id, sensor_type, location, value, metadata) 
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		timestamp,
		"test_sensor",
		"temperature",
		"test_location",
		21.5,
		`{"test": "integration"}`,
	)
	if err != nil {
		t.Errorf("Failed to write test measurement: %v", err)
	}

	// Test if we can query the database
	var count int
	err = db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM measurements 
		 WHERE sensor_id = $1 AND timestamp > $2`,
		"test_sensor",
		timestamp.Add(-1*time.Hour),
	).Scan(&count)

	if err != nil {
		t.Errorf("Failed to query measurements: %v", err)
	}

	if count == 0 {
		t.Error("No measurements found after insertion")
	}
}
