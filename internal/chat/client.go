package chat

import (
	"encoding/json"

	"github.com/gorilla/websocket"
)

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
	ID  string
}

func (c *Client) ReadPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	for {
		_, rawMessage, err := c.conn.ReadMessage()
		if err != nil {
			break
		}

		message := Message{
			Sender:  c.ID,              
			Content: string(rawMessage), 
		}

		encodedMessage, err := json.Marshal(message)
		if err != nil {
			continue
		}

		c.hub.broadcast <- encodedMessage
	}
}

func (c *Client) WritePump() {
	defer c.conn.Close()
	for message := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
			break
		}
	}
}
