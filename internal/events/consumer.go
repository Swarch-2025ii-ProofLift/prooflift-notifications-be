package events

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"prooflift-notifications-be/internal/configs"
	"prooflift-notifications-be/internal/mq"
	"prooflift-notifications-be/internal/services"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   string
	service *services.NotificationService
}

func NewConsumer(conn *amqp.Connection, queue string, service *services.NotificationService) (*Consumer, error) {
	ch, err := mq.NewChannel(conn)
	if err != nil {
		return nil, fmt.Errorf("failed to create MQ channel: %w", err)
	}

	err = mq.DeclareQueue(ch, queue)
	if err != nil {
		_ = ch.Close()
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	return &Consumer{
		conn:    conn,
		channel: ch,
		queue:   queue,
		service: service,
	}, nil
}

func (c *Consumer) Start(ctx context.Context) error {
	msgs, err := c.channel.Consume(
		c.queue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to start consuming messages: %w", err)
	}

	log.Printf("Listening for messages on queue: %s", c.queue)

	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping consumer...")
			return nil

		case d, ok := <-msgs:
			if !ok {
				return fmt.Errorf("message channel closed")
			}

			go c.handleMessage(ctx, d)
		}
	}
}

func (c *Consumer) handleMessage(ctx context.Context, d amqp.Delivery) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic recovered in message handler: %v", r)
			d.Nack(false, false)
		}
	}()

	var event Event
	if err := json.Unmarshal(d.Body, &event); err != nil {
		log.Printf("invalid JSON message: %v", err)
		_ = d.Nack(false, false)
		return
	}

	if !event.IsValidType() {
		log.Printf("ignoring unsupported event type: %s", event.Type)
		_ = d.Ack(false)
		return
	}

	userID, actorID, postID, commentID, err := event.ParseIDs()
	if err != nil {
		log.Printf("invalid event UUIDs: %v", err)
		_ = d.Nack(false, false)
		return
	}

	if userID == actorID {
		log.Printf("skipping self-notification: user %s tried to notify themselves", userID)
		_ = d.Ack(false)
		return
	}

	var svcError error
	switch event.Type {
	case EventCommentCreated:
		_, svcError = c.service.CreateCommentNotification(ctx, userID, actorID, postID, commentID, event.Message)
	case EventReactionAdded:
		_, svcError = c.service.CreateReactionNotification(ctx, userID, actorID, postID, event.Message)
	default:
		log.Printf("no handler for event type: %s", event.Type)
		_ = d.Ack(false)
		return
	}

	if svcError != nil {
		if configs.IsValidationError(svcError) {
			log.Printf("validation error, not requeuing: %v", svcError)
			_ = d.Nack(false, false)
		} else {
			log.Printf("transient error, requeuing: %v", svcError)
			_ = d.Nack(false, true)
		}
		return
	}

	_ = d.Ack(false)
	log.Printf("notification created for user %s (type: %s)", event.UserID, event.Type)
}

func (c *Consumer) Close() {
	if err := c.channel.Close(); err != nil {
		log.Printf("failed to close MQ channel: %v", err)
	}
}
