package types

// EVENTS
// JOIN_ROOM, LEAVE_ROOM, PEER_CONNECTED, PEER_DISCONNECTED, PEER_MESSAGE

// TODO: remove json tags from this struct
type MessageOptions struct {
	Echo        bool     `json:"echo"`
	BroadcastTo []string `json:"broadcast_to,omitempty"`
}

// TODO: remove json tags from this struct
type ClientMessage struct {
	Type    string                 `json:"type"`
	UserID  string               `json:"user_id"`
	RoomID  int                `json:"room_id"`
	Payload map[string]interface{} `json:"payload"`
	Options *MessageOptions        `json:"options"`
}
