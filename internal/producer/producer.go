/*
  Copyright (C) 2019 - 2022 MWSOFT
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
	"time"

	"github.com/segmentio/kafka-go"

	"github.com/superhero-match/superhero-chat/internal/config"
	"github.com/superhero-match/superhero-chat/internal/producer/model"
)

// Producer interface defines producer methods.
type Producer interface {
	StoreMessage(c model.Message) error
}

// producer holds Kafka producer related data.
type producer struct {
	Producer *kafka.Writer
}

// NewProducer configures Kafka producer that produces to configured topic.
func NewProducer(cfg *config.Config) Producer {
	w := &kafka.Writer{
		Addr:         kafka.TCP(cfg.Producer.Brokers),
		Topic:        cfg.Producer.Topic,
		BatchSize:    cfg.Producer.BatchSize,
		BatchTimeout: time.Duration(cfg.Producer.BatchTimeout) * time.Millisecond,
		Balancer:     &kafka.LeastBytes{},
	}

	return &producer{
		Producer: w,
	}
}
