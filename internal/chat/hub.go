package chat

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)


var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
	
		return true
	},
}


type Hub struct {
	clients    map[string]*Client 
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}


func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}


func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client.ID] = client 
			log.Printf("Client registered: %s", client.ID)

		case client := <-h.unregister:
			if _, ok := h.clients[client.ID]; ok {
				delete(h.clients, client.ID) 
				close(client.send)
				log.Printf("Client unregistered: %s", client.ID)
			}

		case message := <-h.broadcast:
			for id, client := range h.clients {
				select {
				case client.send <- message:
					log.Printf("Message sent to client: %s", id)
				default:
					close(client.send)
					delete(h.clients, id)
					log.Printf("Failed to send message, client removed: %s", id)
				}
			}
		}
	}
}


func (h *Hub) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		id = uuid.New().String() 
	}

	client := &Client{
		hub:  h,
		conn: conn,
		send: make(chan []byte, 256),
		ID:   id, 
	}
	h.register <- client

	go client.ReadPump()
	go client.WritePump()

	log.Printf("New client connected: %s", client.ID)
}
