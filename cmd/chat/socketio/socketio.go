package socketio

import (
	"bytes"
	"encoding/json"
	"fmt"
	socketio "github.com/googollee/go-socket.io"
	"github.com/streadway/amqp"
	"github.com/superhero-match/superhero-chat/cmd/chat/model"
	"github.com/superhero-match/superhero-chat/cmd/chat/service"
	"github.com/superhero-match/superhero-chat/internal/config"
	"log"
	"strings"
	"time"
)

// connectedUsers holds all the currently online/connected active users connections/websockets.
// These websockets are going to be used when message is going to be received on RabbitMQ topic.
// Once a message for a specific user will be received, websocket of that specific user will be pulled
// out of the map and the message will be sent ot the user.
var connectedUsers map[string]socketio.Conn
var connectedUsersIDs map[string]string
var connectedUsersQueueNames map[string]string

func init() {
	connectedUsers = make(map[string]socketio.Conn)
	connectedUsersIDs = make(map[string]string)
	connectedUsersQueueNames = make(map[string]string)
}

// SocketIO holds all the data related to Socket.IO.
type SocketIO struct {
	Service *service.Service
}

// NewSocketIO returns new value of type SocketIO.
func NewSocketIO(cfg *config.Config) (*SocketIO, error) {
	srv, err := service.NewService(cfg)
	if err != nil {
		return nil, err
	}

	return &SocketIO{
		Service: srv,
	}, nil
}

// NewSocketIOServer returns Socket.IO server.
func (s *SocketIO) NewSocketIOServer() (*socketio.Server, error) {
	server, err := socketio.NewServer(nil)
	if err != nil {
		return nil, err
	}

	server.OnConnect("/", func(c socketio.Conn) error {
		log.Println("New client connected")

		return nil
	})

	server.OnEvent("/", "onOpen", func(c socketio.Conn, senderID string) {
		log.Println("onOpen event raised...")

		connectedUsers[senderID] = c
		connectedUsersIDs[c.ID()] = senderID

		err = s.Service.SetOnlineUser(senderID)
		if err != nil {
			log.Println(err)
		}

		// User subscribes to RabbitMQ topic message.for.userid.
		q, err := s.Service.RabbitMQ.Channel.QueueDeclare(
			"",    // name, when left empty RabbitMQ generates one automatically.
			true,  // durable means persisted on disk.
			true,  // delete
			true,  // exclusive queue is deleted when connection that declared it is closed.
			false, // no-wait
			nil,   // arguments
		)
		if err != nil {
			log.Println(err)
		}

		err = s.Service.RabbitMQ.Channel.QueueBind(
			q.Name,                          // queue name
			senderID,                // routing key
			s.Service.RabbitMQ.ExchangeName, // exchange
			false,
			nil,
		)
		if err != nil {
			log.Println(err)
		}

		connectedUsersQueueNames[senderID] = q.Name

		msgs, err := s.Service.RabbitMQ.Channel.Consume(
			q.Name, // queue
			"",     // consumer
			false,  // auto ack
			false,  // exclusive
			false,  // no local
			false,  // no wait
			nil,    // args
		)
		if err != nil {
			log.Println(err)
		}

		go func() {
			for d := range msgs {
				log.Printf(" [x] %s", d.Body)

				var m model.Message
				if err := json.Unmarshal(d.Body, &m); err != nil {
					log.Println(err)
					continue
				}

				ws, ok := connectedUsers[m.ReceiverID]
				if !ok {
					// User is not online anymore, that means the offline message needs to be stored in database,
					// cache and Firebase cloud function needs to be run in order to notify user that there is
					// offline message awaiting on the server that needs to be picked up.
					err = s.Service.StoreMessage(m, false, m.CreatedAt)
					if err != nil {
						log.Println(err)
					}

					continue
				}

				ws.Emit("message", m)

				if err = d.Ack(false); err != nil {
					log.Println(err)
				}
			}
		}()

		log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	})

	server.OnEvent("/", "message", func(c socketio.Conn, msg string) {
		log.Println("message event raised...")

		var message model.Message
		if err := json.Unmarshal([]byte(msg), &message); err != nil {
			log.Println(err)
		}

		message.CreatedAt = strings.ReplaceAll(time.Now().UTC().Format(s.Service.TimeFormat), "T", " ")

		// Once a message is received, the check is made whether the receiver is online.
		online, err := s.Service.GetOnlineUser(fmt.Sprintf(s.Service.Cache.OnlineUserKeyFormat, message.ReceiverID))
		if err != nil {
			log.Println(err)
		}

		// This value will be used by Kafka consumer once the message is consumed.
		// If receiver is online, that means that message shouldn't be stored in cache
		// and Firebase cloud function shouldn't be invoked
		// because it was already sent to the receiver via RabbitMQ.
		// If user is offline, then Kafka consumer will store the message in database
		// and in cache and Firebase cloud function must be invoked in order to notify
		// message receiver that message is awaiting on server.
		var isOnline bool

		// If message receiver is online, publish the message on RabbitMQ topic.
		// This message will be emitted to the receiver via websocket stored in connectedUsers map
		// (this is not implemented yet, it is the next step).
		if len(online) > 0 {
			isOnline = true

			messageBytes := new(bytes.Buffer)
			err = json.NewEncoder(messageBytes).Encode(message)
			if err != nil {
				log.Println(err)
			}

			err = s.Service.RabbitMQ.Channel.Publish(
				s.Service.RabbitMQ.ExchangeName,
				message.ReceiverID, // routing key
				false,
				false,
				amqp.Publishing{
					ContentType: s.Service.RabbitMQ.ContentType,
					Body:        messageBytes.Bytes(),
				},
			)
		}

		err = s.Service.StoreMessage(message, isOnline, message.CreatedAt)
		if err != nil {
			log.Println(err)
		}
	})

	server.OnDisconnect("/", func(c socketio.Conn, reason string) {
		log.Println("OnDisconnect event raised...", reason)

		userID, ok := connectedUsersIDs[c.ID()]
		if ok {
			delete(connectedUsers, userID)

			delete(connectedUsersIDs, c.ID())

			queueName, ok := connectedUsersQueueNames[userID]
			if ok {
				err = s.Service.RabbitMQ.Channel.QueueUnbind(
					queueName,
					userID,
					s.Service.RabbitMQ.ExchangeName,
					nil,
				)
				if err != nil {
					log.Println(err)
				}
			}

			if err := s.Service.DeleteOnlineUser(userID); err != nil {
				log.Println(err)
			}
		}
	})

	return server, nil
}
