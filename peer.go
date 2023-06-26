package main

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/olahol/melody"
)

type Peer struct {
	id          string
	room        *Room
	conn        *melody.Session
	displayName string
	mu          sync.RWMutex
}

func (p *Peer) Id() string {
	return p.id
}

func (p *Peer) DisplayName() string {
	return p.displayName
}

func (p *Peer) Room() *Room {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.room
}

func (p *Peer) WriteConn(m ClientMessage) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	messageJSON, err := json.Marshal(m)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return err
	}

	err = p.conn.Write(messageJSON)
	if err != nil {
		return err
	}

	return nil
}
