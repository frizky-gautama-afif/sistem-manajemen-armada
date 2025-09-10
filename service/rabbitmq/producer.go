package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"sistem-manajemen-armada/service/model"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

func NewRabbitMQ(amqpURL string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	err = ch.ExchangeDeclare("fleet.events", "direct", true, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}
	_, err = ch.QueueDeclare("geofence_alerts", true, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}
	err = ch.QueueBind("geofence_alerts", "geofence_entry", "fleet.events", false, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to bind queue: %w", err)
	}

	return &RabbitMQ{Conn: conn, Channel: ch}, nil
}

func (r *RabbitMQ) PublishEvent(event model.GeofenceEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	return r.Channel.PublishWithContext(ctx,
		"fleet.events",
		"geofence_entry",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
}

func (r *RabbitMQ) Close() {
	if r.Channel != nil {
		r.Channel.Close()
	}
	if r.Conn != nil {
		r.Conn.Close()
	}
}
