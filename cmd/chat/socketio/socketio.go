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
)

// connectedUsers holds all the currently online/connected active users connections/websockets.
// These websockets are going to be used when message is going to be received on RabbitMQ topic.
// Once a message for a specific user will be received, websocket of that specific user will be pulled
// out of the map and the message will be sent ot the user.
var connectedUsers map[string]socketio.Conn
var connectedUsersIDs map[string]string

func init() {
	connectedUsers = make(map[string]socketio.Conn)
	connectedUsersIDs = make(map[string]string)
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

	server.OnEvent("/", "onOpen", func(c socketio.Conn, msg string) {
		log.Println("onOpen event raised...")

		var message model.Message
		if err := json.Unmarshal([]byte(msg), &message); err != nil {
			log.Println(err)
		}

		connectedUsers[message.SenderID] = c
		connectedUsersIDs[c.ID()] = message.SenderID

		err = s.Service.SetOnlineUser(message.SenderID)
		if err != nil {
			log.Println(err)
		}

		// User subscribes to RabbitMQ topic message.for.userid.
		q, err := s.Service.RabbitMQ.Channel.QueueDeclare(
			"",    // name, when left empty RabbitMQ generates one automatically.
			false, // durable means persisted on disk.
			false, // delete
			false, // exclusive queue when connections is closed.
			false, // no-wait
			nil,   // arguments
		)
		if err != nil {
			log.Println(err)
		}

		err = s.Service.RabbitMQ.Channel.QueueBind(
			q.Name,                          // queue name
			message.SenderID,                // routing key
			s.Service.RabbitMQ.ExchangeName, // exchange
			false,
			nil,
		)
		if err != nil {
			log.Println(err)
		}

		msgs, err := s.Service.RabbitMQ.Channel.Consume(
			q.Name, // queue
			"",     // consumer
			true,   // auto ack
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
					err = s.Service.StoreMessage(m, false)
					if err != nil {
						log.Println(err)
					}

					continue
				}

				ws.Emit("message", m)
			}
		}()

		log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	})

	server.OnEvent("/", "message", func(c socketio.Conn, msg string) {
		log.Println("message event raised...")
		fmt.Printf("%s\n", msg)
		fmt.Println()

		var message model.Message
		if err := json.Unmarshal([]byte(msg), &message); err != nil {
			log.Println(err)
		}

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

		err = s.Service.StoreMessage(message, isOnline)
		if err != nil {
			log.Println(err)
		}
	})

	server.OnError("/", func(e error) {
		log.Println("OnError event raised...")
		log.Println(e)
	})

	server.OnDisconnect("/", func(c socketio.Conn, reason string) {
		log.Println("OnDisconnect event raised...", reason)

		userID, ok := connectedUsersIDs[c.ID()]
		if ok {
			delete(connectedUsers, userID)
			delete(connectedUsersIDs, c.ID())
		}
	})

	return server, nil
}
