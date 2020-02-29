package model

// Message holds the data received from client.
type Message struct {
	MessageType string `json:"messageType"`
	SenderID    string `json:"senderId"`
	ReceiverID  string `json:"receiverId"`
	Message     string `json:"message"`
}
