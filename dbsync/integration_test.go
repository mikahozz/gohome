package dbsync

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func setupTestDB(t *testing.T) *sql.DB {
	err := godotenv.Load("../.env")
	if err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	connStr := fmt.Sprintf("host=localhost port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("DB_APP_USER"),
		os.Getenv("DB_APP_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("Error connecting to database: %v", err)
	}

	if err := db.Ping(); err != nil {
		t.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	return db
}

func createTestMeasurement(temperature float64) *Measurement {
	return &Measurement{
		Timestamp: time.Now(),
		SensorID:  "test_sensor",
		MainValue: temperature,
		Value: map[string]interface{}{
			"temperature": temperature,
			"humidity":    45.5,
			"battery":     98,
		},
	}
}

func TestMeasurements(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db := setupTestDB(t)
	defer db.Close()

	store := NewMeasurementStore(db)
	ctx := context.Background()

	t.Run("insert and retrieve", func(t *testing.T) {
		// Insert measurement
		temperature := 20.0 + rand.Float64()*10.0
		m := createTestMeasurement(temperature)

		if err := store.Insert(ctx, m); err != nil {
			t.Fatalf("Failed to insert measurement: %v", err)
		}

		// Retrieve and verify
		retrieved, err := store.GetBySensorAndTime(ctx, m.SensorID, m.Timestamp)
		if err != nil {
			t.Fatalf("Failed to get measurement: %v", err)
		}

		if retrieved.MainValue != temperature {
			t.Errorf("MainValue mismatch: got %.1f, want %.1f", retrieved.MainValue, temperature)
		}
	})

	t.Run("update", func(t *testing.T) {
		// Insert initial measurement
		m := createTestMeasurement(21.5)
		if err := store.Insert(ctx, m); err != nil {
			t.Fatalf("Failed to insert initial measurement: %v", err)
		}

		// Update values
		m.MainValue = 22.5
		m.Value["temperature"] = 22.5
		m.Value["humidity"] = 46.0
		m.Value["battery"] = 97

		if err := store.Insert(ctx, m); err != nil {
			t.Fatalf("Failed to update measurement: %v", err)
		}

		// Verify update
		updated, err := store.GetBySensorAndTime(ctx, m.SensorID, m.Timestamp)
		if err != nil {
			t.Fatalf("Failed to get updated measurement: %v", err)
		}

		if updated.MainValue != 22.5 {
			t.Errorf("MainValue not updated: got %.1f, want %.1f", updated.MainValue, 22.5)
		}

		for key, want := range m.Value {
			got := toFloat64(updated.Value[key], t)
			want := toFloat64(want, t)
			if got != want {
				t.Errorf("%s not updated: got %.1f, want %.1f", key, got, want)
			}
		}
	})
}

// Helper function to convert interface{} to float64
func toFloat64(v interface{}, t *testing.T) float64 {
	switch x := v.(type) {
	case float64:
		return x
	case int:
		return float64(x)
	default:
		t.Fatalf("Unexpected type for value: %T", v)
		return 0
	}
}
