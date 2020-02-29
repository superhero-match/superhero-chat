package config

// Producer holds the configuration values for the Kafka producer.
type Producer struct {
	Brokers      []string `env:"KAFKA_BROKERS" default:"[localhost:9092]"`
	Topic        string   `env:"KAFKA_STORE_CHAT_MESSAGE_TOPIC" default:"store.chat.message"`
	BatchSize    int      `env:"KAFKA_BATCH_SIZE" default:"1"`
	BatchTimeout int      `env:"KAFKA_BATCH_TIMEOUT" default:"10"`
}