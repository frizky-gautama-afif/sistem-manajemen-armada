package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sistem-manajemen-armada/service/model"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	broker := os.Getenv("MQTT_BROKER_URL")
	if broker == "" {
		broker = "tcp://mqtt-broker:1883"
	}

	opts := mqtt.NewClientOptions().AddBroker(broker).SetClientID("mock-publisher")
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to connect to MQTT broker: %v", token.Error())
	}
	defer client.Disconnect(250)
	log.Println("Publisher connected to MQTT broker.")

	vehicleIDs := []string{"B1234XYZ", "D5678ABC", "E9012JKL"}
	geofenceLat, geofenceLon := -6.1952, 106.8236 // Bundaran HI

	for {
		for _, id := range vehicleIDs {
			var lat, lon float64
			if id == "B1234XYZ" {
				lat = geofenceLat + (rand.Float64()-0.5)*0.0005
				lon = geofenceLon + (rand.Float64()-0.5)*0.0005
			} else {
				lat = -6.2088 + (rand.Float64()-0.5)*0.05
				lon = 106.8456 + (rand.Float64()-0.5)*0.05
			}

			loc := model.LocationData{
				VehicleID: id,
				Latitude:  lat,
				Longitude: lon,
				Timestamp: time.Now().Unix(),
			}

			payload, err := json.Marshal(loc)
			if err != nil {
				log.Printf("Failed to marshal JSON: %v", err)
				continue
			}

			topic := fmt.Sprintf("/fleet/vehicle/%s/location", id)
			token := client.Publish(topic, 1, false, payload)
			if token.Wait() && token.Error() != nil {
				log.Printf("Failed to publish message: %v", token.Error())
			} else {
				log.Printf("Message published to %s", topic)
			}
		}

		time.Sleep(2 * time.Second)
	}
}
