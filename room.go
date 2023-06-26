package main

import "sync"

type Room struct {
	Id          string
	DisplayName string
	Peers       []*Peer
	// lobby       []*melody.Session
	// start_at    time.Time
	// created_at  time.Time
	mu sync.RWMutex
}

// TODO:
// AddPeer
// RemovePeer
// Id()
// GetDisplayName()
// GetPeers()
// Broadcast
