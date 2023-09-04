package chat

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	Id      uint
	Socket  *websocket.Conn
	Chat    *Chat
	Receive chan []byte
}

func (c *Client) Read() {
	defer c.Socket.Close()

	for {
		_, msg, err := c.Socket.ReadMessage()
		if err != nil {
			return
		}
		log.Println(msg)
		c.Chat.foward <- msg
	}
}

func (c *Client) Write() {
	defer c.Socket.Close()
	for msg := range c.Receive {
		err := c.Socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return
		}
	}
}
