package dbsync

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

// MeasurementStore handles database operations for measurements
type MeasurementStore struct {
	db *sql.DB
}

// Measurement represents a single measurement record
type Measurement struct {
	Timestamp time.Time
	SensorID  string
	MainValue float64
	Value     map[string]interface{}
}

// NewMeasurementStore creates a new MeasurementStore
func NewMeasurementStore(db *sql.DB) *MeasurementStore {
	return &MeasurementStore{db: db}
}

// Insert adds a new measurement to the database or updates if it exists
func (s *MeasurementStore) Insert(ctx context.Context, m *Measurement) error {
	valueJSON, err := json.Marshal(m.Value)
	if err != nil {
		return fmt.Errorf("marshaling value: %w", err)
	}

	query := `
		INSERT INTO measurements (timestamp, sensor_id, main_value, value)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (timestamp, sensor_id) DO UPDATE SET
			main_value = EXCLUDED.main_value,
			value = EXCLUDED.value`

	_, err = s.db.ExecContext(ctx, query,
		m.Timestamp,
		m.SensorID,
		m.MainValue,
		valueJSON,
	)

	if err != nil {
		return fmt.Errorf("upserting measurement: %w", err)
	}

	return nil
}

// GetBySensorAndTime retrieves a measurement for a specific sensor and timestamp
func (s *MeasurementStore) GetBySensorAndTime(ctx context.Context, sensorID string, timestamp time.Time) (*Measurement, error) {
	m := &Measurement{}
	var valueBytes []byte

	query := `
		SELECT timestamp, sensor_id, main_value, value
		FROM measurements
		WHERE sensor_id = $1 AND timestamp = $2`

	err := s.db.QueryRowContext(ctx, query, sensorID, timestamp).Scan(
		&m.Timestamp,
		&m.SensorID,
		&m.MainValue,
		&valueBytes,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("querying measurement: %w", err)
	}

	if err := json.Unmarshal(valueBytes, &m.Value); err != nil {
		return nil, fmt.Errorf("unmarshaling value: %w", err)
	}

	return m, nil
}
