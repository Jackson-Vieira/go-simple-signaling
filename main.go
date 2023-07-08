package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/Jackson-Vieira/go-simple-signalling/domain"
	"github.com/Jackson-Vieira/go-simple-signalling/types"

	"github.com/labstack/echo/v4"
	middleware "github.com/labstack/echo/v4/middleware"
	"github.com/olahol/melody"
)

type RoomManager struct {
	rooms map[string]*domain.Room
	mu    sync.RWMutex
}

func (rm *RoomManager) CreateRoom(displayName string) string {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	newRoom := domain.NewRoom(displayName)
	roomId := newRoom.Id()

	rm.rooms[roomId] = newRoom
	return roomId
}

func (rm *RoomManager) GetRoom(roomId string) (*domain.Room, bool) {
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

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	m := melody.New()
	m.Config.MaxMessageSize = 8192

	roomManager := RoomManager{
		rooms: make(map[string]*domain.Room),
	}

	roomId := roomManager.CreateRoom("Mindmeet-01")
	log.Println("Created room test", roomId)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.GET("/ws", func(c echo.Context) error {
		m.HandleRequest(c.Response().Writer, c.Request())
		return nil
	})

	m.HandleConnect(func(s *melody.Session) {
		log.Println("New Client connected")
	})

	m.HandleDisconnect(func(s *melody.Session) {
		log.Println("Client Disconnected")
	})

	m.HandleClose(func(s *melody.Session, i int, s2 string) error {
		log.Println("Client Closed")
		return nil
	})

	m.HandleError(func(s *melody.Session, e error) {
		log.Println("Error", e)
	})

	// websocket event handlers
	m.HandleMessage(func(s *melody.Session, msg []byte) {

		var message types.ClientMessage
		err := json.Unmarshal(msg, &message)

		if err != nil {
			log.Println("Error decoding JSON", m)
			return
		}

		switch message.Type {

		case "join":
			room, found := roomManager.GetRoom(message.RoomID)
			if !found {
				log.Println("Room not found")
				return
			}
			room.AddUser(s)

		case "leave":
			log.Println("leave case")

		case "offer":
			log.Println("offer case")

		case "answer":
			log.Println("answer case")

		case "ice-candidate":
			log.Println("ice-candidate case")

		default:
			// handle unknown message type
			log.Println("Unknown message type", message.Type)
		}

	})
	e.Logger.Fatal(e.Start(":1323"))
}
