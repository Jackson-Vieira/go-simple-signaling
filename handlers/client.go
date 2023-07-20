package handlers

import (
	"log"

	"github.com/Jackson-Vieira/go-simple-signalling/controllers"
	"github.com/Jackson-Vieira/go-simple-signalling/types"
	"github.com/olahol/melody"
)

const (
	// EventJoin is the event type used when a client wants to join a room.
	EventJoin = "join"

	// EventLeave is the event type used when a client wants to leave a room.
	EventLeave = "leave"

	// EventOffer is the event type used to send a WebRTC offer to another client.
	EventOffer = "offer"

	// EventIceCandidate is the event type used to send information about an ICE candidate to another client.
	EventIceCandidate = "ice-candidate"

	// EventAnswer is the event type used to send a WebRTC answer to another client.
	EventAnswer = "answer"

	// EventUnknown is the event type used when a received event is not recognized or not supported.
	EventUnknown = "unknown"
)

type EventHandler interface {
	Handle(s *melody.Session, msg types.ClientMessage, roomManager *controllers.RoomManager)
}

type JoinHandler struct{}

func (j *JoinHandler) Handle(s *melody.Session, msg types.ClientMessage, roomManager *controllers.RoomManager) {
	room, found := roomManager.GetRoom(msg.RoomID)
	if !found {
		log.Println("Room not found")
		return
	}
	room.AddUser(s, msg)

	// send room id to client
	s.Set("room_id", room.Id())
}

type LeaveHandler struct{}

func (l *LeaveHandler) Handle(s *melody.Session, msg types.ClientMessage, roomManager *controllers.RoomManager) {
	room, found := roomManager.GetRoom(msg.RoomID)
	if !found {
		log.Println("Room not found")
		return
	}
	room.RemoveUser(s)
}

type OfferHandler struct{}

func (o *OfferHandler) Handle(s *melody.Session, msg types.ClientMessage, roomManager *controllers.RoomManager) {
	room, found := roomManager.GetRoom(msg.RoomID)
	if !found {
		log.Println("Room not found")
		return
	}
	room.Broadcast(msg, s)
}

type IceCandidateHandler struct{}

func (i *IceCandidateHandler) Handle(s *melody.Session, msg types.ClientMessage, roomManager *controllers.RoomManager) {
	room, found := roomManager.GetRoom(msg.RoomID)
	if !found {
		log.Println("Room not found")
		return
	}
	room.Broadcast(msg, s)
}

type AnswerHandler struct{}

func (a *AnswerHandler) Handle(s *melody.Session, msg types.ClientMessage, roomManager *controllers.RoomManager) {
	room, found := roomManager.GetRoom(msg.RoomID)
	if !found {
		log.Println("Room not found")
		return
	}
	room.Broadcast(msg, s)
}

type UnknownHandler struct{}

func (u *UnknownHandler) Handle(s *melody.Session, msg types.ClientMessage, roomManager *controllers.RoomManager) {
	log.Println("Unknown message type", msg.Type)
}


func CreateEventHandler(eventType string) EventHandler {
	switch eventType {
	case EventJoin:
		return &JoinHandler{}
	case EventLeave:
		return &LeaveHandler{}
	case EventOffer:
		return &OfferHandler{}
	case EventIceCandidate:
		return &IceCandidateHandler{}
	case EventAnswer:
		return &AnswerHandler{}
	default:
		return &UnknownHandler{}
	}
}