package main

import (
	"log"
	"os"
	"sistem-manajemen-armada/service/rabbitmq"
)

func main() {
	amqpURL := os.Getenv("RABBITMQ_URL")
	if amqpURL == "" {
		amqpURL = "amqp://guest:guest@rabbitmq:5672/"
	}

	rabbitMQ, err := rabbitmq.NewRabbitMQ(amqpURL)
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ: %v", err)
	}
	defer rabbitMQ.Close()

	log.Println("Starting RabbitMQ worker...")
	rabbitmq.ConsumeGeofenceAlerts(rabbitMQ.Channel)
}
