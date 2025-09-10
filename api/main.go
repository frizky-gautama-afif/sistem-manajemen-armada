package main

import (
	"log"
	"os"
	"sistem-manajemen-armada/service/api"
	"sistem-manajemen-armada/service/db"
	"sistem-manajemen-armada/service/mqtt"
	"sistem-manajemen-armada/service/rabbitmq"
	"time"
)

func main() {
	// Database Connection with Retry Logic
	dbConnStr := os.Getenv("DATABASE_URL")
	if dbConnStr == "" {
		dbConnStr = "host=db port=5432 user=user password=password dbname=sistem_manajemen_armada_db sslmode=disable"
	}

	var dbConn *db.DB
	var err error
	for i := 0; i < 5; i++ {
		dbConn, err = db.NewDB(dbConnStr)
		if err == nil {
			log.Println("Database connection successful.")
			break
		}
		log.Printf("Failed to connect to database (attempt %d/5): %v", i+1, err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatalf("Failed to initialize database after multiple attempts: %v", err)
	}
	defer dbConn.Close()

	// RabbitMQ Connection with Retry Logic
	amqpURL := os.Getenv("RABBITMQ_URL")
	if amqpURL == "" {
		amqpURL = "amqp://guest:guest@rabbitmq:5672/"
	}

	var rabbitMQ *rabbitmq.RabbitMQ
	for i := 0; i < 5; i++ {
		rabbitMQ, err = rabbitmq.NewRabbitMQ(amqpURL)
		if err == nil {
			log.Println("RabbitMQ connection successful.")
			break
		}
		log.Printf("Failed to connect to RabbitMQ (attempt %d/5): %v", i+1, err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ after multiple attempts: %v", err)
	}
	defer rabbitMQ.Close()

	// MQTT Subscriber
	mqttBrokerURL := os.Getenv("MQTT_BROKER_URL")
	if mqttBrokerURL == "" {
		mqttBrokerURL = "tcp://mqtt-broker:1883"
	}
	mqttClient := mqtt.NewClient(mqttBrokerURL, "go-subscriber")
	mqtt.SubscribeToLocationTopic(mqttClient, dbConn, rabbitMQ)
	defer mqttClient.Disconnect(250)

	// REST API with Gin
	handler := api.NewHandler(dbConn)
	router := api.SetupRouter(handler)

	log.Println("Starting API server on :9999...")
	if err := router.Run(":9999"); err != nil {
		log.Fatalf("Failed to run API server: %v", err)
	}
}
