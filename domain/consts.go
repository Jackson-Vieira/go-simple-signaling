package domain

const (

	// -- Connection events --

	// PeerConnected represents the event received by a peer when a new peer has joined the room.
	PeerConnected = "peer_connected"

	// PeerDisconnected represents the event received by a peer when a peer has left the room.
	PeerDisconnected = "peer_disconnected"

	// PeerMessage represents the event received when a new message is received from another client.
	PeerMessage = "peer_message"

	// JoinRoom represents the message received by a client when the client joins the room.
	JoinRoom = "room_join"

	// LeaveRoom represents the message received by a client when the client leaves the room.
	LeaveRoom = "room_leave"

	// ConnectionClosed represents the message received by a client when the client has disconnected.
	ConnectionClosed = "disconnected"

	// ConnectionOpen represents the message received by a client when the client has connected.
	ConnectionOpen = "connected"
)