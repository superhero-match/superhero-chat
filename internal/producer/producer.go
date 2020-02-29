package producer

import (
	"time"

	"github.com/segmentio/kafka-go"

	"github.com/superhero-chat/internal/config"
)

// Producer holds Kafka producer related data.
type Producer struct {
	Producer *kafka.Writer
}

// NewProducer configures Kafka producer that produces to configured topic.
func NewProducer(cfg *config.Config) *Producer {
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      cfg.Producer.Brokers,
		Topic:        cfg.Producer.Topic,
		BatchSize:    cfg.Producer.BatchSize,
		BatchTimeout: time.Duration(cfg.Producer.BatchTimeout) * time.Millisecond,
		Balancer:     &kafka.LeastBytes{},
	})

	return &Producer{
		Producer: w,
	}
}
