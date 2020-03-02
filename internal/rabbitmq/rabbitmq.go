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
package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
	"github.com/superhero-chat/internal/config"
)

// NewRabbitMQChannel connects to RabbitMQ, creates channel, declares queue and returns channel.
func NewRabbitMQChannel(cfg *config.Config) (*amqp.Channel, error) {
	// "amqp://guest:guest@localhost:5672/
	conn, err := amqp.Dial(fmt.Sprintf(cfg.RabbitMQ.Host, cfg.RabbitMQ.User, cfg.RabbitMQ.Password, cfg.RabbitMQ.Address, cfg.RabbitMQ.Port))
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

	return ch, nil
}
