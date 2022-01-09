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
package service

import (
	"github.com/streadway/amqp"
)

// Consume consumes messages from RabbitMQ.
func (s *service) Consume(queueName string) (<-chan amqp.Delivery, error) {
	return s.RabbitMQ.Consume(queueName)
}

// Publish publishes RabbitMQ message.
func (s *service) Publish(receiverID string, message []byte) error {
	return s.RabbitMQ.Publish(receiverID, message)
}

// QueueBind binds newly created queue.
func (s *service) QueueBind(queueName string, senderID string) error {
	return s.RabbitMQ.QueueBind(queueName, senderID)
}

// QueueDeclare declares queue.
func (s *service) QueueDeclare() (amqp.Queue, error) {
	return s.RabbitMQ.QueueDeclare()
}

// QueueUnbind unbinds queue.
func (s *service) QueueUnbind(queueName string, userID string) error {
	return s.RabbitMQ.QueueUnbind(queueName, userID)
}
