package socketio

import (
	"bytes"
	"encoding/json"
	"fmt"
	gosocketio "github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
	"github.com/streadway/amqp"
	"github.com/superhero-chat/cmd/chat/model"
	"github.com/superhero-chat/cmd/chat/service"
	"github.com/superhero-chat/internal/config"
	"log"
)

// connectedUsers holds all the currently online/connected active users connections/websockets.
// These websockets are going to be used when message is going to be received on RabbitMQ topic.
// Once a message for a specific user will be received, websocket of that specific user will be pulled
// out of the map and the message will be sent ot the user.
var connectedUsers map[string]*gosocketio.Channel

func init() {
	connectedUsers = make(map[string]*gosocketio.Channel)
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
func (s *SocketIO) NewSocketIOServer() (*gosocketio.Server, error) {
	//create
	server := gosocketio.NewServer(transport.GetDefaultWebsocketTransport())

	//handle connected
	err := server.On(gosocketio.OnConnection, func(c *gosocketio.Channel) {
		log.Println("New client connected")
	})
	if err != nil {
		return nil, err
	}

	err = server.On("onOpen", func(c *gosocketio.Channel, message model.Message) string {
		log.Println("onOpen event raised...")
		fmt.Printf("%+v\n", message)
		fmt.Println()

		connectedUsers[message.SenderID] = c

		err = s.Service.SetOnlineUser(message.SenderID)
		if err != nil {
			log.Println(err)
		}

		// User subscribes to RabbitMQ topic message.for.userid.
		q, err := s.Service.RabbitMQ.Channel.QueueDeclare(
			"",    // name
			false, // durable
			false, // delete when unused
			true,  // exclusive
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
			true,   // exclusive
			false,  // no local
			false,  // no wait
			nil,    // args
		)
		if err != nil {
			log.Println(err)
		}

		forever := make(chan bool)

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

				if err = ws.Emit("message", m); err != nil {
					log.Println(err)
					continue
				}
			}
		}()

		log.Printf(" [*] Waiting for logs. To exit press CTRL+C")

		<-forever

		return "OK"
	})
	if err != nil {
		return nil, err
	}

	err = server.On("message", func(c *gosocketio.Channel, message model.Message) string {
		//send event to all in room
		// c.Emit("my event", "my data")
		log.Println("message event raised...")
		fmt.Printf("%+v\n", message)
		fmt.Println()

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

			reqBodyBytes := new(bytes.Buffer)
			err = json.NewEncoder(reqBodyBytes).Encode(message)
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
					Body:        reqBodyBytes.Bytes(),
				},
			)
		}

		err = s.Service.StoreMessage(message, isOnline)
		if err != nil {
			log.Println(err)
		}

		return "OK"
	})
	if err != nil {
		return nil, err
	}

	return server, nil
}
