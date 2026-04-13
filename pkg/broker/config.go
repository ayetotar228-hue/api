package broker

import "os"

type Config struct {
	Brokers       []string
	ConsumerGroup string
}

func LoadConfig() *Config {
	brokers := os.Getenv("KAFKA_BROKERS")
	if brokers == "" {
		brokers = "localhost:9092"
	}
	group := os.Getenv("KAFKA_CONSUMER_GROUP")
	if group == "" {
		group = "default-group"
	}
	return &Config{
		Brokers:       []string{brokers},
		ConsumerGroup: group,
	}
}
