package chat

import (
	"log"

	"github.com/gorilla/websocket"
)

type client struct {
	socket  *websocket.Conn
	room    *Room
	receive chan []byte
}

func (c *client) read() {
	defer c.socket.Close()

	for {
		_, msg, err := c.socket.ReadMessage()
		if err != nil {
			return
		}
		log.Println(msg)
		c.room.foward <- msg
	}
}

func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.receive {
		err := c.socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return
		}
	}
}
