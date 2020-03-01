/*
  Copyright (C) 2019 - 2020 MWSOFT
  This program is free software: you can redistribute it and/or modify
  it under the terms of the GNU General Public License as published by
  the Free Software Foundation, either version 3 of the License, or
  (at your option) any later version.
  This program is distributed in the hope that it will be useful,
  but WITHOUT ANY WARRANTY; without even the implied warranty of
  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
  GNU General Public License for more details.
  You should have received a copy of the GNU General Public License
  along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/
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