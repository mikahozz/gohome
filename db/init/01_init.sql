-- Create the measurements table
CREATE TABLE measurements (
    id BIGSERIAL PRIMARY KEY,
    timestamp TIMESTAMPTZ NOT NULL,
    sensor_id VARCHAR(50) NOT NULL,
    sensor_type VARCHAR(50) NOT NULL,
    location VARCHAR(50),
    value NUMERIC,
    metadata JSONB,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for common query patterns
CREATE INDEX idx_measurements_timestamp_sensor ON measurements(timestamp DESC, sensor_id);
CREATE INDEX idx_measurements_sensor_type ON measurements(sensor_type);
CREATE INDEX idx_measurements_metadata ON measurements USING GIN (metadata);
  