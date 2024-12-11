package chat

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
	ID   string
}

func (c *Client) ReadPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	for {
		_, rawMessage, err := c.conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		var message map[string]interface{}
		if err := json.Unmarshal(rawMessage, &message); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			continue
		}

		messageType, ok := message["type"].(string)
		if !ok {
			log.Printf("Error: Message type is not a string or missing")
			continue
		}

		if messageType == "CONNECT" {
			id, ok := message["id"].(string)
			if !ok {
				log.Printf("Error: ID is not a string or missing")
				continue
			}
			c.ID = id
			log.Printf("My ID is: %s", c.ID)

			continue
		} else if messageType == "MESSAGE" {
			sender, ok := message["sender"].(string)
			if !ok {
				log.Printf("Error: Sender is not a string or missing")
				continue
			}
			content, ok := message["content"].(string)
			if !ok {
				log.Printf("Error: Content is not a string or missing")
				continue
			}
			log.Printf("Received message from %s: %s", sender, content)

			encodedMessage, err := json.Marshal(message)
			if err != nil {
				log.Printf("Error marshaling message: %v", err)
				continue
			}
			c.hub.broadcast <- encodedMessage
		} else {
			log.Printf("Unknown message type: %s", messageType)
		}
	}
}

func (c *Client) WritePump() {
	defer c.conn.Close()
	for message := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Printf("Error writing message: %v", err)
			break
		}
	}
}
