package chat

type Message struct {
	Sender  string `json:"sender"`
	Content string `json:"content"`
}
