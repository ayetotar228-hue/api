package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

type ProducerInterface interface {
	SendMessage(ctx context.Context, topic string, msg Message) error
	SendUserCreated(ctx context.Context, userID int, email string) error
	SendEmailNotification(ctx context.Context, to, subject, body string) error
	Close() error
}

func NewProducer(brokers []string) *Producer {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireOne,
	}

	return &Producer{writer: writer}
}

func (p *Producer) SendMessage(ctx context.Context, topic string, msg Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	kafkaMsg := kafka.Message{
		Topic: topic,
		Key:   []byte(msg.ID),
		Value: data,
		Time:  time.Now(),
	}

	if err := p.writer.WriteMessages(ctx, kafkaMsg); err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	return nil
}

func (p *Producer) SendUserCreated(ctx context.Context, userID int, email string) error {
	msg := Message{
		ID:        uuid.New().String(),
		Type:      EventUserCreated,
		Timestamp: time.Now(),
		Payload: map[string]interface{}{
			"user_id": userID,
			"email":   email,
		},
	}

	return p.SendMessage(ctx, "user-events", msg)
}

func (p *Producer) SendEmailNotification(ctx context.Context, to, subject, body string) error {
	msg := Message{
		ID:        uuid.New().String(),
		Type:      EventEmailSend,
		Timestamp: time.Now(),
		Payload: map[string]interface{}{
			"to":      to,
			"subject": subject,
			"body":    body,
		},
	}

	return p.SendMessage(ctx, "email-notifications", msg)
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
