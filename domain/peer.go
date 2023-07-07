package domain

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/olahol/melody"

	"github.com/Jackson-Vieira/go-simple-signalling/types"
)

type Peer struct {
	ID   string `json:"id"`
	Room *Room  `json:"-"`
	// FIXME: dont export this field in future
	Conn        *melody.Session `json:"-"`
	DisplayName string          `json:"displayName"`
	mu          sync.RWMutex    `json:"-"`
}

// return the peer id
func (p *Peer) Id() string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.ID
}

// return the display name
func (p *Peer) GetDisplayName() string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.DisplayName
}

// return the room id
func (p *Peer) GetRoom() *Room {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.Room
}

// write a message to the peer
func (p *Peer) WriteConn(m types.ClientMessage) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	messageJSON, err := json.Marshal(m)
	if err != nil {
		log.Println("Error encoding JSON:", err)
		return err
	}

	err = p.Conn.Write(messageJSON)
	if err != nil {
		return err
	}

	return nil
}

// Close peer
func (p *Peer) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	err := p.Conn.CloseWithMsg([]byte(`{"type":"peer_disconnected","payload":{}}`))
	if err != nil {
		return err
	}

	return nil
}
