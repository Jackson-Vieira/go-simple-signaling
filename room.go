package main

import (
	"fmt"
	"sync"
	"time"
)

type Room struct {
	ID          string       `json:"ID"`
	DisplayName string       `json:"displayName"`
	Peers       []*Peer      `json:"-"`
	StartAt     time.Time    `json:"start_at"`
	CreatedAt   time.Time    `json:"created_at,omitempty"`
	mu          sync.RWMutex `json:"-"`
}

// return the room id
func (r *Room) Id() string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.ID
}

// init room
func (r *Room) Init() {
	r.mu.Lock()
	defer r.mu.Unlock()

	// r.lobby = make([]*melody.Session, 0)
	r.Peers = make([]*Peer, 0)
	r.CreatedAt = time.Now()
}

// close room
func (r *Room) Close() {
	r.mu.Lock()
	defer r.mu.Unlock()

	fmt.Println("Closing room", r.Id())
	for _, peer := range r.Peers {
		fmt.Println("Closing peer", peer.Id())
		err := peer.Close()
		if err != nil {
			fmt.Println("Error closing peer connection:", err)
		}
	}
}

// return the peers in the room
func (r *Room) GetPeers() []*Peer {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.Peers
}

// set the room peer list
func (r *Room) SetPeers(peers []*Peer) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.Peers = peers
}

// set the room display name
func (r *Room) SetDisplayName(displayName string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.DisplayName = displayName
}

// return the room display name
func (r *Room) GetDisplayName() string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.DisplayName
}

// add a peer to the room
func (r *Room) AddPeer(peer *Peer) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.Peers = append(r.Peers, peer)
}

// find peer index by ID
func (r *Room) FindPeerIndexById(peerId string) int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for i, peer := range r.Peers {
		if peer.Id() == peerId {
			return i
		}
	}
	return -1
}

func (r *Room) RemovePeer(peerId string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	peerIndex := r.FindPeerIndexById(peerId)

	if peerIndex == -1 {
		// return peer not found
		fmt.Println("peer not found")
		return
	}

	// remove peer from room
	r.Peers = append(r.Peers[:peerIndex], r.Peers[peerIndex+1:]...)

	// if room is empty, remove (close) room from room manager
	return
}

func (r *Room) Broadcast(msg ClientMessage) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	peers := r.Peers
	for _, peer := range peers {
		err := peer.WriteConn(msg)
		if err != nil {
			fmt.Println("Error writing to peer:", err)
		}
	}
}
