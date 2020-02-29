package mapper

import (
	"github.com/superhero-chat/cmd/chat/model"
	pm "github.com/superhero-chat/internal/producer/model"
)

// MapAPIMessageToProducer maps API Message model to Producer Message model.
func MapAPIMessageToProducer(m model.Message) (message pm.Message) {
	return pm.Message{
		MessageType: m.MessageType,
		SenderID:    m.SenderID,
		ReceiverID:  m.ReceiverID,
		Message:     m.Message,
	}
}
