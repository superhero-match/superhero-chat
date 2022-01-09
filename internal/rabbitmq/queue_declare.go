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
	"github.com/streadway/amqp"
)

// QueueDeclare declares queue.
func (r *rabbitMQ) QueueDeclare() (amqp.Queue, error) {
	q, err := r.Channel.QueueDeclare(
		"",    // name, when left empty RabbitMQ generates one automatically.
		true,  // durable means persisted on disk.
		true,  // delete
		true,  // exclusive queue is deleted when connection that declared it is closed.
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return q, err
	}

	return q, nil
}
