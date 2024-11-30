package dbsync

import (
	"context"
	"database/sql"
	"math/rand"
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

	// Generate random temperature between 20.0 and 30.0
	temperature := 20.0 + rand.Float64()*10.0
	timestamp := time.Now()

	// Test if we can write to the database
	_, err := db.ExecContext(ctx,
		`INSERT INTO measurements (timestamp, sensor_id, sensor_type, location, value, metadata) 
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		timestamp,
		"test_sensor",
		"temperature",
		"test_location",
		temperature,
		`{"test": "integration"}`,
	)
	if err != nil {
		t.Fatalf("Failed to write test measurement: %v", err)
	}

	// Test if we can query the database and verify the values
	var resultTimestamp time.Time
	var resultTemperature float64
	err = db.QueryRowContext(ctx,
		`SELECT timestamp, value FROM measurements 
		 WHERE sensor_id = $1 AND timestamp = $2`,
		"test_sensor",
		timestamp,
	).Scan(&resultTimestamp, &resultTemperature)

	if err != nil {
		t.Fatalf("Failed to query measurement: %v", err)
	}

	if !resultTimestamp.Equal(timestamp) {
		t.Errorf("Retrieved timestamp %v doesn't match inserted timestamp %v", resultTimestamp, timestamp)
	}

	if resultTemperature != temperature {
		t.Errorf("Retrieved temperature %.2f doesn't match inserted temperature %.2f", resultTemperature, temperature)
	}
}
