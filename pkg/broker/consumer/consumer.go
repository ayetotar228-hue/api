package consumer

import (
	"api/pkg/broker"
	"api/pkg/worker"
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader   *kafka.Reader
	pool     *worker.Pool
	topic    string
	handlers map[broker.EventType]func(context.Context, broker.Message) error
}

func NewConsumer(brokers []string, groupID, topic string, pool *worker.Pool) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  brokers,
			GroupID:  groupID,
			Topic:    topic,
			MinBytes: 10e3,
			MaxBytes: 10e6,
		}),
		pool:     pool,
		topic:    topic,
		handlers: make(map[broker.EventType]func(context.Context, broker.Message) error),
	}
}

func (c *Consumer) RegisterHandler(eventType broker.EventType, handler func(context.Context, broker.Message) error) {
	c.handlers[eventType] = handler
}

func (c *Consumer) Start(ctx context.Context) {
	log.Printf("Kafka consumer started: topic=%s, group=%s", c.topic, c.reader.Config().GroupID)

	for {
		select {
		case <-ctx.Done():
			log.Println("Consumer stopping...")
			return
		default:
			msg, err := c.reader.FetchMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				log.Printf("Fetch error: %v", err)
				continue
			}

			var eventMsg broker.Message
			if err := json.Unmarshal(msg.Value, &eventMsg); err != nil {
				log.Printf("Unmarshal error: %v", err)
				continue
			}

			resultCh := make(chan error, 1)
			c.pool.Submit(worker.Job{
				ID:      int(msg.Offset),
				Type:    string(eventMsg.Type),
				Payload: eventMsg,
				Context: ctx,
				Result:  resultCh,
			})

			procErr := <-resultCh
			if procErr == nil {
				if err := c.reader.CommitMessages(ctx, msg); err != nil {
					log.Printf("Commit error: %v", err)
				}
			} else {
				log.Printf("Processing failed (offset NOT committed): %v", procErr)
			}
		}
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
