package types
type MessageOptions struct {
	Echo        bool     `json:"echo,omitempty"`
	BroadcastTo []string `json:"broadcast_to,omitempty"`
}

type ClientMessage struct {
	Type    string                 `json:"type"`
	UserID  string               `json:"peer_id"`
	RoomID  int                `json:"room_id"`
	Payload map[string]interface{} `json:"payload"`
	Options *MessageOptions        `json:"options"`
}
