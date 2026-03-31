package broker

import (
	"api/internal/worker"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader     *kafka.Reader
	workerPool *worker.Pool
	handlers   map[EventType]func(context.Context, Message) error
	topic      string
}

func NewConsumer(brokers []string, groupID string, topic string, pool *worker.Pool) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})

	return &Consumer{
		reader:     reader,
		workerPool: pool,
		handlers:   make(map[EventType]func(context.Context, Message) error),
		topic:      topic,
	}
}

func (c *Consumer) RegisterHandler(eventType EventType, handler func(context.Context, Message) error) {
	c.handlers[eventType] = handler
}

func (c *Consumer) Start(ctx context.Context) {
	log.Printf("Kafka consumer started on topic '%s' with group '%s'", c.topic, c.reader.Config().GroupID)

	for {
		select {
		case <-ctx.Done():
			log.Println("Kafka consumer stopping...")
			return
		default:
			msg, err := c.reader.FetchMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				log.Printf("Failed to fetch message: %v", err)
				continue
			}

			c.workerPool.Submit(worker.Job{
				ID:      len(msg.Value),
				Type:    "kafka_message",
				Payload: msg,
				Context: ctx,
			})

			if err := c.reader.CommitMessages(ctx, msg); err != nil {
				log.Printf("Failed to commit message: %v", err)
			}
		}
	}
}

func (c *Consumer) ProcessMessage(ctx context.Context, kafkaMsg kafka.Message) error {
	var msg Message
	if err := json.Unmarshal(kafkaMsg.Value, &msg); err != nil {
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}

	handler, exists := c.handlers[msg.Type]
	if !exists {
		return fmt.Errorf("no handler for event type: %s", msg.Type)
	}

	return handler(ctx, msg)
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
