package chat

import (
	"encoding/json"
	"log"
)

type Chat struct {
	clients map[uint]*Client

	Leave chan *Client
	Join  chan *Client

	foward chan []byte
}

type Message struct {
	Id   uint   `json:"id"`
	Text string `json:"text"`
}

func NewChat() *Chat {
	return &Chat{
		foward:  make(chan []byte),
		Leave:   make(chan *Client),
		Join:    make(chan *Client),
		clients: make(map[uint]*Client),
	}
}

func (c *Chat) Run() {
	for {
		select {
		case Client := <-c.Join:
			c.clients[Client.Id] = Client
		case Client := <-c.Leave:
			delete(c.clients, Client.Id)
			close(Client.Receive)
		case msgJson := <-c.foward:
			var msg Message
			if err := json.Unmarshal(msgJson, &msg); err != nil {
				log.Println("not supported json message")
				continue
			}

			client, exist := c.clients[msg.Id]
			if !exist {
				log.Println("can't fiend client in map ")
				continue
			}
			client.Receive <- []byte(msg.Text)
		}
	}
}
