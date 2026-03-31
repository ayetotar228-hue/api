package broker

import "os"

type Config struct {
	Brokers       []string
	ConsumerGroup string
	Topics        []string
}

func LoadConfig() *Config {
	brokers := os.Getenv("KAFKA_BROKERS")
	if brokers == "" {
		brokers = "localhost:9092"
	}

	return &Config{
		Brokers:       []string{brokers},
		ConsumerGroup: os.Getenv("KAFKA_CONSUMER_GROUP"),
		Topics:        []string{"user-events", "email-notifications"},
	}
}
