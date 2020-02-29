package service

import (
	"github.com/superhero-chat/cmd/chat/model"
	"github.com/superhero-chat/cmd/chat/service/mapper"
)

// StoreMessage publishes new message on Kafka topic for it to be
// consumed by consumer and stored in Cache and DB.
func (s *Service) StoreMessage(m model.Message) error {
	return s.Producer.StoreMessage(mapper.MapAPIMessageToProducer(m))
}