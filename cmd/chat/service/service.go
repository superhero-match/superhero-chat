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
	"go.uber.org/zap"

	"github.com/superhero-match/superhero-chat/cmd/chat/model"
	"github.com/superhero-match/superhero-chat/internal/cache"
	"github.com/superhero-match/superhero-chat/internal/config"
	"github.com/superhero-match/superhero-chat/internal/producer"
	"github.com/superhero-match/superhero-chat/internal/rabbitmq"
)

// Service interface defines service methods.
type Service interface {
	SetOnlineUser(key string, userID string) error
	GetOnlineUser(key string) (string, error)
	DeleteOnlineUser(keys []string, userID string) error
	StoreMessage(m model.Message, isOnline bool, createdAt string) error
	Consume(queueName string) (<-chan amqp.Delivery, error)
	Publish(receiverID string, message []byte) error
	QueueBind(queueName string, senderID string) error
	QueueDeclare() (amqp.Queue, error)
	QueueUnbind(queueName string, userID string) error
}

// service holds all the different services that are used when handling request.
type service struct {
	Producer producer.Producer
	Cache    cache.Cache
	RabbitMQ rabbitmq.RabbitMQ
	Logger   *zap.Logger
}

// NewService creates value of type Service.
func NewService(cfg *config.Config) (Service, error) {
	c, err := cache.NewCache(cfg)
	if err != nil {
		return nil, err
	}

	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}

	defer logger.Sync()

	rabbit, err := rabbitmq.NewRabbitMQ(cfg)
	if err != nil {
		return nil, err
	}

	return &service{
		Producer: producer.NewProducer(cfg),
		Cache:    c,
		RabbitMQ: rabbit,
		Logger:   logger,
	}, nil
}
