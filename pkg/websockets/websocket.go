package websocket

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	Clients   map[*websocket.Conn]bool
	Broadcast chan interface{}
	Lock      sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		Clients:   make(map[*websocket.Conn]bool),
		Broadcast: make(chan interface{}),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case message := <-h.Broadcast:
			h.Lock.Lock()
			for client := range h.Clients {
				if err := client.WriteJSON(message); err != nil {
					client.Close()
					delete(h.Clients, client)
				}
			}
			h.Lock.Unlock()
		}
	}
}

func (h *Hub) RegisterClient(client *websocket.Conn) {
	h.Lock.Lock()
	defer h.Lock.Unlock()
	h.Clients[client] = true
}
