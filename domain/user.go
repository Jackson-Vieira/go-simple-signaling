package domain

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/Jackson-Vieira/go-simple-signalling/types"

	"github.com/olahol/melody"
)

type User struct {
	id   string
	room *Room
	// FIXME: dont export this field in future
	conn        *melody.Session
	displayName string
	mu          sync.Mutex
}

func (u *User) Id() string {
	return u.id
}

func (u *User) GetDisplayName() string {
	return u.displayName
}

func (u *User) GetRoom() *Room {
	return u.room
}

// write a message to the User
func (u *User) WriteConn(m types.ClientMessage) error {
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
func (u *User) Disconnect() error {
	u.mu.Lock()
	defer u.mu.Unlock()

	err := u.conn.Close()
	if err != nil {
		return err
	}

	return nil
}
