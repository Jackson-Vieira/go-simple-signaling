package domain

import (
	"log"
	"sync"
	"time"

	"github.com/Jackson-Vieira/go-simple-signalling/types"
	"github.com/olahol/melody"
)

type Room struct {
	id         	int
	displayName string
	peers       map[*melody.Session]*Peer
	// startAt     time.Time
	createdAt time.Time
	mu        sync.Mutex
}

func (r *Room) Id() int {
	return r.id
}

func (r *Room) GetDisplayName() string {
	return r.displayName
}

// init room
func (r *Room) Init() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.peers = make(map[*melody.Session]*Peer, 0)
	r.createdAt = time.Now()
}

// close room
func (r *Room) Close() {
	log.Println("Closing room", r.Id())
	for _, u := range r.peers {
		log.Println("Disconnect user connection", u.Id())
		err := u.Disconnect()
		if err != nil {
			log.Println("Error closing user connection:", err)
		}
	}
}

// return the users unclocked
func (r *Room) GetPeersUnlocked(except *melody.Session) []*Peer {
	peers := make([]*Peer, 0, len(r.peers))
	for _, u := range r.peers {
		peers = append(peers, u)
	}
	return peers
}

// return the users in the room
func (r *Room) GetUsers(except *melody.Session) []*Peer {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.GetPeersUnlocked(except)
}

// set the room display name
func (r *Room) SetDisplayName(displayName string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.displayName = displayName
}

// add a users to the room
func (r *Room) AddUser(s *melody.Session, message types.ClientMessage) *Peer {
	r.mu.Lock()
	defer r.mu.Unlock()

	var m types.ClientMessage

	users := r.GetPeersUnlocked(nil)

	// add peer to room
	r.peers[s] = &Peer{
		room: r,
		conn: s,
		id:   message.UserID,
	}
	peer := r.peers[s]

	// JOIN MESSAGE
	m = types.ClientMessage{
		Type:   "join",
		UserID: peer.Id(),
		RoomID: r.Id(),
		Payload: map[string]interface{}{
			"room": map[string]interface{}{
				"display_name": r.displayName,
				"id":           r.id,
				"created_at":   r.createdAt,
			},
		},
	}

	err := peer.WriteConn(m)
	if err != nil {
		log.Println("Error writing to user:", err)
	}

	// USER CONNECTED MESSAGE
	m = types.ClientMessage{
		Type:   "peer_connected",
		UserID: peer.Id(),
		RoomID: r.Id(),
	}

	r.Broadcast(m, s)

	for _, u := range users {

		if u.Id() == peer.Id() {
			continue
		}

		// USER CONNECTED MESSAGE
		m := types.ClientMessage{
			Type:   "peer_connected",
			UserID: u.Id(),
			RoomID: r.Id(),
		}

		err := peer.WriteConn(m)

		if err != nil {
			log.Println("Error writing to user:", err)
		}
	}

	return peer
}

func (r *Room) RemoveUser(s *melody.Session) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var m types.ClientMessage

	u := r.peers[s]

	if u == nil {
		log.Println("user not found")
		return
	}

	m = types.ClientMessage{
		Type:    "peer_disconnected",
		UserID:  u.Id(),
		RoomID:  r.Id(),
		Payload: make(map[string]interface{}),
		Options: &types.MessageOptions{},
	}

	r.Broadcast(m, s)

	// send leave room message to user
	m = types.ClientMessage{
		Type:   "leave_room",
		RoomID: r.Id(),
	}

	err := u.WriteConn(m)
	if err != nil {
		log.Println("Error writing to user:", err)
	}

	// remove user from room
	delete(r.peers, s)

	// FIXUP: refactor this for a better solution and remove this for another function wrapper in leaveRoom for exaple
	log.Println("Peer removed successfully")
}

func (r *Room) Broadcast(msg types.ClientMessage, except *melody.Session) {
	peers := r.GetPeersUnlocked(except)

	log.Println("Broadcast message to peers", len(peers))

	for _, u := range peers {

		if u.conn == except {
			continue
		}

		err := u.WriteConn(msg)
		if err != nil {
			log.Fatalln("Error writing to user:", err)
		}
	}
}

// factory
func NewRoom(id int, displayName string) *Room {
	return &Room{
		id:          id,
		displayName: displayName,
		peers:       make(map[*melody.Session]*Peer, 0),
		createdAt:   time.Now(),
	}
}
