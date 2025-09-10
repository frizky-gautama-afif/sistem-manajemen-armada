package rabbitmq

import (
	"encoding/json"
	"log"
	"sistem-manajemen-armada/service/model"

	amqp "github.com/rabbitmq/amqp091-go"
)

func ConsumeGeofenceAlerts(ch *amqp.Channel) {
	msgs, err := ch.Consume(
		"geofence_alerts",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("Received geofence message: %s", d.Body)
			var event model.GeofenceEvent
			if err := json.Unmarshal(d.Body, &event); err != nil {
				log.Printf("Failed to parse JSON message: %v", err)
				continue
			}
			log.Printf("✅ ALERT: Vehicle %s entered geofence at %d", event.VehicleID, event.Timestamp)
		}
	}()

	log.Println("RabbitMQ worker is waiting for messages.")
	<-forever
}
