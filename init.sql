CREATE TABLE IF NOT EXISTS vehicle_locations (
    vehicle_id VARCHAR(255) NOT NULL,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    timestamp BIGINT NOT NULL,
    PRIMARY KEY (vehicle_id, timestamp)
);

CREATE INDEX IF NOT EXISTS idx_vehicle_timestamp ON vehicle_locations (vehicle_id, timestamp);