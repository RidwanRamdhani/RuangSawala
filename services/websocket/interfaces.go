package websocket

// ClientManager defines operations for managing WebSocket clients
type ClientManager interface {
	Register(client *Client)
	Unregister(client *Client)
	Get(userID int) *Client
	GetAll() map[int]*Client
}

// RoomManager defines operations for managing chat rooms
type RoomManager interface {
	Create(roomID string, userA, userB int) *Room
	Get(roomID string) *Room
	Delete(roomID string)
	Exists(roomID string) bool
	GetUsersInRoom(roomID string) (userA, userB int, err error)
}

// Broadcaster defines operations for broadcasting messages
type Broadcaster interface {
	BroadcastToRoom(roomID string, message []byte) error
}

// MatchNotifier defines operations for notifying matched users
type MatchNotifier interface {
	NotifyMatch(userAID, userBID int, score int) error
}
