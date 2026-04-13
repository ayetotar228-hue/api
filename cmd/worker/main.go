package main

import (
	"api/pkg/broker"
	"api/pkg/broker/consumer"
	"api/pkg/broker/producer"
	"api/pkg/worker"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	brokerConfig := broker.LoadConfig()
	producer := producer.NewProducer(brokerConfig.Brokers)
	workerPool := worker.NewPool(10, producer)
	workerPool.Start()
	defer workerPool.Stop()

	consumer := consumer.NewConsumer(
		brokerConfig.Brokers,
		brokerConfig.ConsumerGroup,
		"user-events",
		workerPool,
	)
	defer consumer.Close()

	consumer.RegisterHandler(broker.EventUserCreated, func(ctx context.Context, msg broker.Message) error {
		log.Printf("Processing user.created event: %v", msg.Payload)
		email := msg.Payload["email"].(string)
		log.Printf("Sending welcome email to: %s", email)
		return nil
	})

	consumer.RegisterHandler(broker.EventEmailSend, func(ctx context.Context, msg broker.Message) error {
		log.Printf("Processing email.send event: %v", msg.Payload)
		return nil
	})

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go consumer.Start(ctx)

	log.Println("Worker Service started, waiting for messages...")

	<-quit
	log.Println("Shutting down Worker Service...")

	cancel()
	time.Sleep(2 * time.Second)

	log.Println("Worker Service exited gracefully")
}
