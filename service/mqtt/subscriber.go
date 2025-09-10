package mqtt

import (
	"encoding/json"
	"log"
	"math"
	"sistem-manajemen-armada/service/db"
	"sistem-manajemen-armada/service/model"
	"sistem-manajemen-armada/service/rabbitmq"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var GeofenceCenter = struct {
	Latitude  float64
	Longitude float64
}{
	Latitude:  -6.1952,
	Longitude: 106.8236,
}

const GeofenceRadius = 50

func NewClient(broker, clientID string) mqtt.Client {
	opts := mqtt.NewClientOptions().AddBroker(broker).SetClientID(clientID).SetCleanSession(true)
	opts.SetConnectionLostHandler(func(c mqtt.Client, err error) {
		log.Printf("MQTT connection lost: %v. Reconnecting...", err)
	})
	opts.SetOnConnectHandler(func(c mqtt.Client) {
		log.Println("Connected to MQTT broker.")
	})
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to connect to MQTT broker: %v", token.Error())
	}
	return client
}

func Haversine(lat1, lon1, lat2, lon2 float64) float64 {
	var R = 6371e3
	var φ1 = lat1 * math.Pi / 180
	var φ2 = lat2 * math.Pi / 180
	var Δφ = (lat2 - lat1) * math.Pi / 180
	var Δλ = (lon2 - lon1) * math.Pi / 180

	var a = math.Sin(Δφ/2)*math.Sin(Δφ/2) + math.Cos(φ1)*math.Cos(φ2)*math.Sin(Δλ/2)*math.Sin(Δλ/2)
	var c = 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	var d = R * c
	return d
}

func SubscribeToLocationTopic(client mqtt.Client, dbConn *db.DB, rabbitMQ *rabbitmq.RabbitMQ) {
	topic := "/fleet/vehicle/+/location"
	token := client.Subscribe(topic, 1, func(c mqtt.Client, msg mqtt.Message) {
		log.Printf("Received message on topic: %s", msg.Topic())

		parts := strings.Split(msg.Topic(), "/")
		if len(parts) != 5 || parts[3] == "" {
			log.Printf("Invalid topic format or missing vehicle ID: %s", msg.Topic())
			return
		}
		vehicleIDFromTopic := parts[3]

		var loc model.LocationData
		if err := json.Unmarshal(msg.Payload(), &loc); err != nil {
			log.Printf("Failed to unmarshal JSON: %v", err)
			return
		}

		if loc.VehicleID == "" || loc.Latitude == 0 || loc.Longitude == 0 || loc.Timestamp == 0 {
			log.Println("Data validation failed: empty fields")
			return
		}

		if loc.VehicleID != vehicleIDFromTopic {
			log.Printf("Vehicle ID in payload (%s) does not match ID in topic (%s)", loc.VehicleID, vehicleIDFromTopic)
			return
		}

		if err := dbConn.SaveLocation(&loc); err != nil {
			log.Printf("Failed to save location to DB: %v", err)
		} else {
			log.Printf("Location for vehicle %s saved: Lat=%.4f, Lon=%.4f", loc.VehicleID, loc.Latitude, loc.Longitude)
		}

		// ... (Logika Geofence tetap sama) ...
		distance := Haversine(GeofenceCenter.Latitude, GeofenceCenter.Longitude, loc.Latitude, loc.Longitude)
		if distance <= GeofenceRadius {
			log.Printf("Vehicle %s entered geofence. Distance: %.2f meters", loc.VehicleID, distance)
			event := model.GeofenceEvent{
				VehicleID: loc.VehicleID,
				Event:     "geofence_entry",
				Location: model.Location{
					Latitude:  loc.Latitude,
					Longitude: loc.Longitude,
				},
				Timestamp: loc.Timestamp,
			}
			if err := rabbitMQ.PublishEvent(event); err != nil {
				log.Printf("Failed to publish geofence event to RabbitMQ: %v", err)
			}
		}
	})
	if token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to subscribe to topic: %v", token.Error())
	}
	log.Printf("Successfully subscribed to topic: %s", topic)
}
