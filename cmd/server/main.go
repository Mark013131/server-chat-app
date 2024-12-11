package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type Message struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

func main() {
	serverURL := "ws://localhost:8080/ws?id=your-client-id"
	conn, _, err := websocket.DefaultDialer.Dial(serverURL, nil)
	if err != nil {
		log.Fatal("Error connecting to WebSocket server:", err)
	}
	defer conn.Close()

	_, message, err := conn.ReadMessage()
	if err != nil {
		log.Fatal("Error reading message:", err)
	}
	fmt.Printf("Received: %s\n", message)
	
	msg := Message{
		Type:    "MESSAGE",          
		Content: "Hello, everyone!", 
	}

	encodedMessage, err := json.Marshal(msg)
	if err != nil {
		log.Fatal("Error marshaling message:", err)
	}

	err = conn.WriteMessage(websocket.TextMessage, encodedMessage)
	if err != nil {
		log.Fatal("Error sending message:", err)
	}

	_, response, err := conn.ReadMessage()
	if err != nil {
		log.Fatal("Error reading response:", err)
	}
	fmt.Printf("Received response: %s\n", response)
}
