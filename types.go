package main

type MessageOptions struct {
	Echo bool `json:"echo"`
}

type ClientMessage struct {
	Type    string                 `json:"type"`
	PeerID  string                 `json:"peer_id"`
	RoomID  string                 `json:"room_id"`
	Payload map[string]interface{} `json:"payload"`
	Options *MessageOptions        `json:"options"`
}
