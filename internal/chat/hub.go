package chat

import (
	"net/http"

	"github.com/gorilla/websocket"
)

func (h *Hub) Run() {
	for {
		select {

		case client := <-h.register:
			h.clients[client] = true

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}

		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					delete(h.clients, client)
					close(client.send)
				}
			}
		}
	}
}

// WebSocketアップグレーダー
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// 全てのリクエストを許可（セキュリティ調整が必要な場合もあります）
		return true
	},
}

// Hub はWebSocket接続を管理するための構造体です
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

// NewHub は新しいHubを作成します
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// HandleWebSocket はWebSocket接続をアップグレードし、クライアントを管理するメソッド
func (h *Hub) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// WebSocket接続にアップグレード
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		// エラー処理
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}

	// 新しいクライアントを作成し、ハブに登録
	client := &Client{
		hub:  h,
		conn: conn,
		send: make(chan []byte, 256),
	}
	h.register <- client

	// クライアントの読み書きを処理
	go client.ReadPump()
	go client.WritePump()
}
