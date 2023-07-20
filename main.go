package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Jackson-Vieira/go-simple-signalling/controllers"
	"github.com/Jackson-Vieira/go-simple-signalling/handlers"
	"github.com/Jackson-Vieira/go-simple-signalling/types"

	"github.com/labstack/echo/v4"
	middleware "github.com/labstack/echo/v4/middleware"
	"github.com/olahol/melody"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	m := melody.New()
	m.Config.MaxMessageSize = 8192

	// TODO: room manager package 
	roomManager := controllers.NewRoomManager()
	
	
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "/")
	})

	e.GET("/ws", func(c echo.Context) error {
		m.HandleRequest(c.Response().Writer, c.Request())
		return nil
	})

	m.HandleConnect(func(s *melody.Session) {
		// this maybe is useful for log or make some operation when client connect 
		// set a timer for the client join in a room for example, ... 
		log.Println("New Client connected")
	})
	
	m.HandleDisconnect(func(s *melody.Session) {
		log.Println("Client Disconnected")
		
		roomId, exist := s.Get("room_id")
		
		if !exist {
			return
		}
		
		room, found := roomManager.GetRoom(roomId.(int))
		
		if !found {
			// reset room id
			s.Set("room_id", nil)
			
			return
		}
		
		room.RemoveUser(s)
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
		
		eventHadler := handlers.CreateEventHandler(message.Type)
		eventHadler.Handle(s, message, roomManager)
	})
	
	rid := roomManager.CreateRoom("Mindmeet-01")
	log.Println(rid)
	e.Logger.Fatal(e.Start(":8000"))
}
