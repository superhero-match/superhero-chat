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
package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/superhero-chat/cmd/chat/model"
	"log"
	"net/http"
)

// We'll need to define an Upgrader
// this will require a Read and Write buffer size
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// connectedUsers holds all the currently online/connected active users connections/websockets.
// These websockets are going to be used when message is going to be received on RabbitMQ topic.
// Once a message for a specific user will be received, websocket of that specific user will be pulled
// out of the map and the message will be sent ot the user.
var connectedUsers map[string]*websocket.Conn

func init()  {
	connectedUsers = make(map[string]*websocket.Conn)
}

// reader is listening for new messages.
func reader(conn *websocket.Conn, c *Controller) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			continue
		}

		fmt.Println(string(p))

		var message model.Message
		if err := json.Unmarshal(p, &message); err != nil {
			panic(err)
		}

		switch message.MessageType {
		case "onOpen":
			fmt.Println("Open connection message has been received...")
			fmt.Printf("%+v\n", message)
			fmt.Println()

			connectedUsers[message.SenderID] = conn

			err = c.Service.SetOnlineUser(message.SenderID)
			if err != nil {
				log.Println(err)
				continue
			}

			// User subscribes to RabbitMQ topic message.for.userid.
			break
		case "message":
			fmt.Println("Text message has been received...")
			fmt.Printf("%+v\n", message)
			fmt.Println()

			// The message is published on the RabbitMQ topic for the receiver to receive the message.
			// Message contains sender id, receiver id, message, something else.
			// Once a message is received on topic, the socket is pulled out the map and the message is sent to the receiver.

			break
		}

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}
}

func (c *Controller) WsEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	// helpful log statement to show connections
	log.Println("Client Connected")

	err = ws.WriteMessage(1, []byte("Hi Client!"))
	if err != nil {
		log.Println(err)
	}

	reader(ws, c)
}

