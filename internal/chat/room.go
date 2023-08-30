package chat

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: messageBufferSize}

type Room struct {
	clients map[*client]struct{}

	join chan *client

	leave chan *client

	foward chan []byte
}

func NewRoom() *Room {
	return &Room{
		foward:  make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]struct{}),
	}
}

func (r *Room) ServeHTTP(c *gin.Context) {
	socket, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "can't upgrade to websocket",
		})
		return
	}

	client := &client{
		socket:  socket,
		receive: make(chan []byte, messageBufferSize),
		room:    r,
	}
	r.join <- client

	defer func() { r.leave <- client }()
	go client.write()
	client.read()

}

func (r *Room) Run() {
	for {
		select {
		case client := <-r.join:
			r.clients[client] = struct{}{}
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.receive)
		case msg := <-r.foward:
			for client := range r.clients {
				client.receive <- msg
			}
		}
	}
}
