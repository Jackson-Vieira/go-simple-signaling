package types

type MessageOptions struct {
	Echo        bool     `json:"echo"`
	BroadcastTo []string `json:"broadcast_to,omitempty"`
}

type ClientMessage struct {
	Type    string                 `json:"type"`
	PeerID  string                 `json:"peer_id"`
	RoomID  string                 `json:"room_id"`
	Payload map[string]interface{} `json:"payload"`
	Options *MessageOptions        `json:"options"`
}
