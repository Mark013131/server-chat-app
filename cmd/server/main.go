package main

import (
	"chat-app/internal/server"
	"log"
)

func main() {
	s := server.NewServer()

	log.Println("Starting the server on :8080...")
	if err := s.Start(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
