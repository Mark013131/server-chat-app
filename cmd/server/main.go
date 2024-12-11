package main

import (
	"chat-app/internal/server"
	"log"
)

func main() {
	// サーバーのインスタンスを作成
	s := server.NewServer()

	// サーバーを起動して、エラーが発生した場合にはログに出力
	log.Println("Starting the server on :8080...")
	if err := s.Start(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
