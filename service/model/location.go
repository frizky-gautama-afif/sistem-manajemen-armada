package model

// Location struct for latitude and longitude
type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// LocationData represents the vehicle location data from MQTT.
type LocationData struct {
	VehicleID string  `json:"vehicle_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timestamp int64   `json:"timestamp"`
}

// GeofenceEvent represents the event message sent to RabbitMQ.
type GeofenceEvent struct {
	VehicleID string   `json:"vehicle_id"`
	Event     string   `json:"event"`
	Location  Location `json:"location"` // Menggunakan struct yang telah dideklarasikan
	Timestamp int64    `json:"timestamp"`
}
