package controllers

import (
	"sync"

	"github.com/Jackson-Vieira/go-simple-signalling/domain"
)

type RoomManager struct {
	rooms map[int]*domain.Room
	mu    sync.RWMutex
}

func (rm *RoomManager) CreateRoom(displayName string) int {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	
	rid := len(rm.rooms) + 1

	newRoom := domain.NewRoom(rid, displayName)

	rm.rooms[rid] = newRoom
	return rid
}

func (rm *RoomManager) GetRoom(roomId int) (*domain.Room, bool) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	room, ok := rm.rooms[roomId]
	return room, ok
}

func (rm *RoomManager) GetAllRooms() []*domain.Room {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	rooms := make([]*domain.Room, 0, len(rm.rooms))
	for _, room := range rm.rooms {
		rooms = append(rooms, room)
	}
	return rooms
}

func NewRoomManager() *RoomManager {
	return &RoomManager{
		rooms: make(map[int]*domain.Room),
	}
}