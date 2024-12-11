package dbsync

import (
	"context"
	"database/sql"
	"encoding/json"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// setupTestDB creates and returns a PostgreSQL database connection for testing
func setupTestDB(t *testing.T) *sql.DB {
	err := godotenv.Load("../.env")
	if err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	host := "localhost"
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("DB_APP_USER")
	password := os.Getenv("DB_APP_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")

	connStr := "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable"
	connStr = "host=" + host + " port=" + port + " user=" + user + " password=" + password + " dbname=" + dbname + " sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("Error connecting to database: %v", err)
	}

	return db
}

func TestPostgresHealth(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db := setupTestDB(t)
	defer db.Close()

	err := db.Ping()
	if err != nil {
		t.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
}

func TestMeasurementWrite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db := setupTestDB(t)
	defer db.Close()

	ctx := context.Background()

	// Generate test data
	temperature := 20.0 + rand.Float64()*10.0
	timestamp := time.Now()
	value := map[string]interface{}{
		"temperature": temperature,
		"humidity":    45.5,
		"battery":     98,
		"type":        "temperature_sensor",
		"location":    "living_room",
	}

	// Convert value map to JSON
	valueJSON, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("Failed to marshal value to JSON: %v", err)
	}

	// Test writing to database
	_, err = db.ExecContext(ctx,
		`INSERT INTO measurements (timestamp, sensor_id, main_value, value) 
		 VALUES ($1, $2, $3, $4)`,
		timestamp,
		"test_sensor",
		temperature,
		valueJSON,
	)
	if err != nil {
		t.Fatalf("Failed to write test measurement: %v", err)
	}

	// Test reading from database
	var resultTimestamp time.Time
	var resultMainValue float64
	var valueBytes []byte
	err = db.QueryRowContext(ctx,
		`SELECT timestamp, main_value, value FROM measurements 
		 WHERE sensor_id = $1 AND timestamp = $2`,
		"test_sensor",
		timestamp,
	).Scan(&resultTimestamp, &resultMainValue, &valueBytes)

	if err != nil {
		t.Fatalf("Failed to query measurement: %v", err)
	}

	// Parse JSON value
	var resultValue map[string]interface{}
	err = json.Unmarshal(valueBytes, &resultValue)
	if err != nil {
		t.Fatalf("Failed to parse JSON value: %v", err)
	}

	// Verify results
	if !resultTimestamp.Equal(timestamp) {
		t.Errorf("Retrieved timestamp %v doesn't match inserted timestamp %v", resultTimestamp, timestamp)
	}

	if resultMainValue != temperature {
		t.Errorf("Retrieved main_value %.2f doesn't match inserted temperature %.2f", resultMainValue, temperature)
	}

	if resultValue["temperature"].(float64) != temperature {
		t.Errorf("Retrieved temperature %.2f from value JSON doesn't match inserted temperature %.2f",
			resultValue["temperature"].(float64), temperature)
	}
}
