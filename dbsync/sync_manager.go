package dbsync

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type SyncType string
type SyncFrequency string
type SyncStatus string

const (
	SyncTypeSpotPrice SyncType = "SPOT_PRICE"
	// Add more sync types as needed
	// SyncTypeWeather  SyncType = "WEATHER"
	// SyncTypeUsage    SyncType = "USAGE"

	FrequencyHourly  SyncFrequency = "HOURLY"
	FrequencyDaily   SyncFrequency = "DAILY"
	FrequencyWeekly  SyncFrequency = "WEEKLY"
	FrequencyMonthly SyncFrequency = "MONTHLY"

	StatusSynced    SyncStatus = "SYNCED"
	StatusNotSynced SyncStatus = "NOT_SYNCED"
	StatusError     SyncStatus = "ERROR"
)

type SyncManager struct {
	db *sql.DB
}

func (sm *SyncManager) MarkSynced(ctx context.Context, syncType SyncType, date time.Time, metadata map[string]interface{}) error {
	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	_, err = sm.db.ExecContext(ctx, `
		INSERT INTO sync_status 
			(sync_type, target_date, frequency, status, last_attempt, metadata, updated_at)
		VALUES 
			($1, $2, $3, $4, NOW(), $5, NOW())
		ON CONFLICT (sync_type, target_date) 
		DO UPDATE SET 
			status = $4,
			last_attempt = NOW(),
			metadata = $5,
			updated_at = NOW()
	`, syncType, date, FrequencyDaily, StatusSynced, metadataJSON)

	return err
}

func (sm *SyncManager) MarkError(ctx context.Context, syncType SyncType, date time.Time, errMsg string) error {
	_, err := sm.db.ExecContext(ctx, `
		INSERT INTO sync_status 
			(sync_type, target_date, frequency, status, last_attempt, error_message, retry_count)
		VALUES 
			($1, $2, $3, $4, NOW(), $5, 1)
		ON CONFLICT (sync_type, target_date) 
		DO UPDATE SET 
			status = $4,
			last_attempt = NOW(),
			error_message = $5,
			retry_count = sync_status.retry_count + 1,
			updated_at = NOW()
	`, syncType, date, FrequencyDaily, StatusError, errMsg)

	return err
}

func (sm *SyncManager) GetUnsynced(ctx context.Context, syncType SyncType, limit int) ([]time.Time, error) {
	rows, err := sm.db.QueryContext(ctx, `
		SELECT target_date 
		FROM sync_status 
		WHERE sync_type = $1 
		  AND status = $2
		ORDER BY target_date ASC 
		LIMIT $3
	`, syncType, StatusNotSynced, limit)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dates []time.Time
	for rows.Next() {
		var date time.Time
		if err := rows.Scan(&date); err != nil {
			return nil, err
		}
		dates = append(dates, date)
	}

	return dates, rows.Err()
}

// Helper to initialize sync status for a date range
func (sm *SyncManager) InitializeDateRange(ctx context.Context, syncType SyncType, startDate, endDate time.Time) error {
	_, err := sm.db.ExecContext(ctx, `
		INSERT INTO sync_status 
			(sync_type, target_date, frequency, status)
		SELECT 
			$1, 
			generate_series::date, 
			$2,
			$3
		FROM generate_series($4::date, $5::date, '1 day'::interval)
		ON CONFLICT (sync_type, target_date) DO NOTHING
	`, syncType, FrequencyDaily, StatusNotSynced, startDate, endDate)

	return err
}

// Get sync status for a specific date
func (sm *SyncManager) GetStatus(ctx context.Context, syncType SyncType, date time.Time) (SyncStatus, error) {
	var status SyncStatus
	err := sm.db.QueryRowContext(ctx, `
		SELECT status 
		FROM sync_status 
		WHERE sync_type = $1 
		  AND target_date = $2
	`, syncType, date).Scan(&status)

	if err == sql.ErrNoRows {
		return StatusNotSynced, nil
	}
	return status, err
}
