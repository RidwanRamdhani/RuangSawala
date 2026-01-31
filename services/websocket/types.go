package websocket

import (
	"github.com/gorilla/websocket"
)

// MessageType represents the type of WebSocket message
type MessageType string

const (
	// Message types
	TypeMatchFound MessageType = "match_found"
	TypeMessage    MessageType = "message"
	TypeTyping     MessageType = "typing"
	TypeError      MessageType = "error"
)

// Message represents a WebSocket message structure
type Message struct {
	Type     MessageType `json:"type"`
	RoomID   string      `json:"room_id,omitempty"`
	Content  string      `json:"content,omitempty"`
	From     string      `json:"from,omitempty"`
	Partner  PartnerInfo `json:"partner,omitempty"`
	IsTyping bool        `json:"is_typing,omitempty"`
	Error    string      `json:"error,omitempty"`
}

// PartnerInfo represents information about a matched partner
type PartnerInfo struct {
	ID    int `json:"id"`
	Score int `json:"score"`
}

// RoomMessage represents a message to be broadcast to a room
type RoomMessage struct {
	RoomID  string
	Message []byte
}

// Client represents a WebSocket client connection
type Client struct {
	UserID int
	Hub    *Hub
	Conn   *websocket.Conn
	Send   chan []byte
}

// Room represents a chat room between two users
type Room struct {
	ID    string
	UserA int
	UserB int
}
