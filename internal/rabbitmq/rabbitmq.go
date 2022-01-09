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
package rabbitmq

import (
	"fmt"

	"github.com/streadway/amqp"

	"github.com/superhero-match/superhero-chat/internal/config"
)

// RabbitMQ interface defines RabbitMQ methods.
type RabbitMQ interface {
	QueueDeclare() (amqp.Queue, error)
	QueueBind(queueName string, senderID string) error
	Consume(queueName string) (<-chan amqp.Delivery, error)
	Publish(receiverID string, message []byte) error
	QueueUnbind(queueName string, userID string) error
}

type rabbitMQ struct {
	Channel      *amqp.Channel
	ExchangeName string
	ContentType  string
}

// NewRabbitMQ connects to RabbitMQ, creates channel, declares queue and returns channel.
func NewRabbitMQ(cfg *config.Config) (RabbitMQ, error) {
	// amqp://guest:guest@localhost:5672/
	rabbitMQURL := fmt.Sprintf(
		cfg.RabbitMQ.Host,
		cfg.RabbitMQ.User,
		cfg.RabbitMQ.Password,
		cfg.RabbitMQ.Address,
		cfg.RabbitMQ.Port,
	)

	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	err = ch.ExchangeDeclare(
		cfg.RabbitMQ.ExchangeName,
		cfg.RabbitMQ.ExchangeType,
		cfg.RabbitMQ.ExchangeDurable,
		cfg.RabbitMQ.ExchangeAutoDelete,
		cfg.RabbitMQ.ExchangeInternal,
		cfg.RabbitMQ.ExchangeNoWait,
		nil,
	)

	return &rabbitMQ{
		Channel:      ch,
		ExchangeName: cfg.RabbitMQ.ExchangeName,
		ContentType:  cfg.RabbitMQ.ContentType,
	}, nil
}
