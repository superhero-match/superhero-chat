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
package config

// RabbitMQ holds the configuration values for the RabbitMQ.
type RabbitMQ struct {
	Host               string `env:"RABBITMQ_HOST" yaml:"host" default:"amqp://%s:%s@%s:%s/"`
	User               string `env:"RABBITMQ_USER" yaml:"user" default:"dev"`
	Password           string `env:"RABBITMQ_PASSWORD" yaml:"password" default:"password"`
	Address            string `env:"RABBITMQ_ADDRESS" yaml:"address" default:"192.168.1.229"`
	Port               string `env:"RABBITMQ_PORT" yaml:"port" default:"5672"`
	ExchangeName       string `env:"RABBITMQ_EXCHANGE_NAME" yaml:"exchange_name" default:"message.for.*"`
	ExchangeType       string `env:"RABBITMQ_EXCHANGE_TYPE" yaml:"exchange_type" default:"topic"`
	ExchangeDurable    bool   `env:"RABBITMQ_EXCHANGE_DURABLE" yaml:"exchange_durable" default:"true"`
	ExchangeAutoDelete bool   `env:"RABBITMQ_EXCHANGE_AUTO_DELETE" yaml:"exchange_auto_delete" default:"false"`
	ExchangeInternal   bool   `env:"RABBITMQ_EXCHANGE_INTERNAL" yaml:"exchange_internal" default:"false"`
	ExchangeNoWait     bool   `env:"RABBITMQ_EXCHANGE_NO_WAIT" yaml:"exchange_no_wait" default:"false"`
	TopicMandatory     bool   `env:"RABBITMQ_TOPIC_MANDATORY" yaml:"topic_mandatory" default:"false"`
	TopicImmediate     bool   `env:"RABBITMQ_TOPIC_IMMEDIATE" yaml:"topic_immediate" default:"false"`
	ContentType        string `env:"RABBITMQ_CONTENT_TYPE" yaml:"content_type" default:"application/json"`
}
