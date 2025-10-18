package mq

import (
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func NewMQConnection(url string) (*amqp.Connection, error) {
	if url == "" {
		return nil, fmt.Errorf("message queue URL is empty")
	}

	var conn *amqp.Connection
	var err error
	maxRetries := 10
	retryDelay := 3 * time.Second

	for i := range maxRetries {
		conn, err = amqp.Dial(url)
		if err == nil {
			log.Printf("Connected to message queue")
			return conn, nil
		}

		if i < maxRetries-1 {
			log.Printf("Failed to connect to message queue (attempt %d/%d): %v. Retrying in %v...", i+1, maxRetries, err, retryDelay)
			time.Sleep(retryDelay)
		}
	}

	return nil, fmt.Errorf("failed to connect to message queue after %d attempts: %w", maxRetries, err)
}

func NewChannel(conn *amqp.Connection) (*amqp.Channel, error) {
	if conn == nil {
		return nil, fmt.Errorf("MQ connection is nil")
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	return ch, nil
}

func DeclareQueue(ch *amqp.Channel, queueName string) error {
	if ch == nil {
		return fmt.Errorf("channel is nil")
	}

	_, err := ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue %s: %w", queueName, err)
	}

	return nil
}
