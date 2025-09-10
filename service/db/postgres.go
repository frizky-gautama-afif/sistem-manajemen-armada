package db

import (
	"database/sql"
	"fmt"
	"sistem-manajemen-armada/service/model"

	_ "github.com/lib/pq"
)

type DB struct {
	*sql.DB
}

func NewDB(connStr string) (*DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	return &DB{db}, nil
}

// SaveLocation saves a vehicle's location data into the database.
func (d *DB) SaveLocation(loc *model.LocationData) error {
	query := `INSERT INTO vehicle_locations (vehicle_id, latitude, longitude, timestamp) VALUES ($1, $2, $3, $4)`
	_, err := d.Exec(query, loc.VehicleID, loc.Latitude, loc.Longitude, loc.Timestamp)
	if err != nil {
		return fmt.Errorf("failed to save location: %w", err)
	}
	return nil
}

// GetLastLocation retrieves the last known location for a vehicle.
func (d *DB) GetLastLocation(vehicleID string) (*model.LocationData, error) {
	loc := &model.LocationData{}
	query := `SELECT vehicle_id, latitude, longitude, timestamp FROM vehicle_locations WHERE vehicle_id = $1 ORDER BY timestamp DESC LIMIT 1`
	err := d.QueryRow(query, vehicleID).Scan(&loc.VehicleID, &loc.Latitude, &loc.Longitude, &loc.Timestamp)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get last location: %w", err)
	}
	return loc, nil
}

// GetLocationHistory retrieves a vehicle's location history within a specific time range.
func (d *DB) GetLocationHistory(vehicleID string, start, end int64) ([]model.LocationData, error) {
	var locations []model.LocationData
	query := `SELECT vehicle_id, latitude, longitude, timestamp FROM vehicle_locations WHERE vehicle_id = $1 AND timestamp >= $2 AND timestamp <= $3 ORDER BY timestamp ASC`
	rows, err := d.Query(query, vehicleID, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get location history: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var loc model.LocationData
		if err := rows.Scan(&loc.VehicleID, &loc.Latitude, &loc.Longitude, &loc.Timestamp); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		locations = append(locations, loc)
	}
	return locations, nil
}
