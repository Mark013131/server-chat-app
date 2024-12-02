package server

import (
	"chat-app/internal/chat"
	"net/http"
)

type Server struct {
	hub *chat.Hub
}

func NewServer() *Server {
	hub := chat.NewHub()
	go hub.Run()
	return &Server{hub: hub}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/ws" {
		s.hub.HandleWebSocket(w, r)
		return
	}

	http.NotFound(w, r)
}

func (s *Server) Start() error {

	server := &http.Server{
		Addr:    ":8080",
		Handler: s,
	}

	return server.ListenAndServe()
}
