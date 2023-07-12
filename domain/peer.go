package domain

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/Jackson-Vieira/go-simple-signalling/types"
	"github.com/google/uuid"

	"github.com/olahol/melody"
)

type Peer struct {
	id   string
	room *Room
	conn        *melody.Session
	displayName string
	mu          sync.Mutex
}

func (u *Peer) Id() string {
	return u.id
}

func (u *Peer) GetDisplayName() string {
	return u.displayName
}

func (u *Peer) GetRoom() *Room {
	return u.room
}

// write a message to the User
func (u *Peer) WriteConn(m types.ClientMessage) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	messageJSON, err := json.Marshal(m)
	if err != nil {
		log.Println("Error encoding JSON:", err)
		return err
	}

	err = u.conn.Write(messageJSON)
	if err != nil {
		return err
	}

	return nil
}

// disconnect user
func (u *Peer) Disconnect() error {
	u.mu.Lock()
	defer u.mu.Unlock()

	var err error

	m := types.ClientMessage{
		Type: "disconnect",
	}

	err = u.WriteConn(m)

	if err != nil {
		return err
	}

	err = u.conn.Close()

	if err != nil {
		return err
	}

	return nil
}

// FACTORY
func NewUser(room *Room, conn *melody.Session, displayName string) *Peer {
	return &Peer{
		id:          uuid.New().String(),
		room:        room,
		conn:        conn,
		displayName: displayName,
	}
}
