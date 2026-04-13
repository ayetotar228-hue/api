package producer

import (
	"api/pkg/broker"
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

func NewProducer(brokers []string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Balancer:     &kafka.LeastBytes{},
			RequiredAcks: kafka.RequireOne,
		},
	}
}

func (p *Producer) SendMessage(ctx context.Context, topic string, msg broker.Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	return p.writer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Key:   []byte(msg.ID),
		Value: data,
		Time:  time.Now(),
	})
}

func (p *Producer) SendUserCreated(ctx context.Context, userID int, email string) error {
	return p.SendMessage(ctx, "user-events", broker.Message{
		ID:        uuid.New().String(),
		Type:      broker.EventUserCreated,
		Timestamp: time.Now(),
		Payload:   map[string]interface{}{"user_id": userID, "email": email},
	})
}

func (p *Producer) SendEmailNotification(ctx context.Context, to, subject, body string) error {
	return p.SendMessage(ctx, "email-notifications", broker.Message{
		ID:        uuid.New().String(),
		Type:      broker.EventEmailSend,
		Timestamp: time.Now(),
		Payload:   map[string]interface{}{"to": to, "subject": subject, "body": body},
	})
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
