package producer

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"

	"github.com/superhero-chat/internal/producer/model"
)

// StoreMessage publishes new choice on Kafka topic for it to be
// consumed by consumer and stored in Cache.
func (p *Producer) StoreMessage(c model.Message) error {
	var sb bytes.Buffer

	err := json.NewEncoder(&sb).Encode(c)
	if err != nil {
		return err
	}

	err = p.Producer.WriteMessages(
		context.Background(),
		kafka.Message{
			Value: sb.Bytes(),
		},
	)
	if err != nil {
		return err
	}

	return nil
}