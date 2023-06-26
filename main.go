package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	middleware "github.com/labstack/echo/v4/middleware"
	"github.com/olahol/melody"
)

func isJoinMessage(message ClientMessage) bool {
	return message.Type == "join"
}

type RoomManager struct {
	rooms map[string]*Room
	mu    sync.RWMutex
}

func (rm *RoomManager) CreateRoom(displayName string) string {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	newRoom := &Room{
		Id:          uuid.New().String(),
		DisplayName: displayName,
		Peers:       make([]*Peer, 0),
	}

	rm.rooms[newRoom.Id] = newRoom

	return newRoom.Id
}

func (rm *RoomManager) GetRoom(roomId string) (*Room, bool) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	room, ok := rm.rooms[roomId]
	return room, ok
}

func (rm *RoomManager) GetAllRooms() []*Room {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	rooms := make([]*Room, 0, len(rm.rooms))
	for _, room := range rm.rooms {
		rooms = append(rooms, room)
	}
	return rooms
}

func main() {
	e := echo.New()
	m := melody.New()
	m.Config.MaxMessageSize = 8192

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	roomManager := RoomManager{
		rooms: make(map[string]*Room),
	}

	roomId := roomManager.CreateRoom("Sala Teste 1")
	fmt.Println("RoomID", roomId)

	// http://localhost:1323
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	// ws://localhost:1323
	e.GET("/ws", func(c echo.Context) error {
		m.HandleRequest(c.Response().Writer, c.Request())
		return nil
	})

	// return a peer object
	// m.HandleConnect(func(s *melody.Session){ return })

	m.HandleMessage(func(s *melody.Session, msg []byte) {

		var message ClientMessage
		err := json.Unmarshal(msg, &message)

		if err != nil {
			// returm invalid message
			fmt.Println("Error decoding JSON", m)
			return
		}

		if !isJoinMessage(message) {
			fmt.Println("Is join Message")
			// validate if peer is in a room
		}

		switch message.Type {

		case "join":
			room, found := roomManager.GetRoom(message.RoomID)
			if !found {
				fmt.Println("Room not found")
				return
			}

			peer := &Peer{
				id:   uuid.New().String(),
				room: nil,
				conn: s,
			}

			// room.AddPeer
			room.Peers = append(room.Peers, peer)
			peer.room = room

			response := ClientMessage{
				Type:    "peer_connected",
				PeerID:  peer.Id(),
				RoomID:  room.Id,
				Payload: make(map[string]interface{}),
				Options: message.Options,
			}

			if len(room.Peers) > 0 {
				for _, p := range room.Peers {
					// p.Send | p.WriteConn
					if err := p.WriteConn(response); err != nil {
						fmt.Println("Error sending message:", err)
					}
				}
			}

		case "leave":
			fmt.Println("leave case")
		case "offer":
			fmt.Println("offer case")
		case "answer":
			fmt.Println("answer case")
		case "ice-candidate":
			fmt.Println("ice-candidate-case")
		}

	})
	e.Logger.Fatal(e.Start(":1323"))
}

//  Room.validate
//
