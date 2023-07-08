package domain

import (
	"log"
	"sync"
	"time"

	"github.com/Jackson-Vieira/go-simple-signalling/types"
	"github.com/google/uuid"
	"github.com/olahol/melody"
)

type Room struct {
	id          string
	displayName string
	users       map[*melody.Session]*User
	// startAt     time.Time
	createdAt time.Time
	mu        sync.Mutex
}

func (r *Room) Id() string {
	return r.id
}

func (r *Room) GetDisplayName() string {
	return r.displayName
}

// init room
func (r *Room) Init() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.users = make(map[*melody.Session]*User, 0)
	r.createdAt = time.Now()
}

// close room
func (r *Room) Close() {
	log.Println("Closing room", r.Id())
	for _, u := range r.users {
		log.Println("Disconnect user connection", u.Id())
		err := u.Disconnect()
		if err != nil {
			log.Println("Error closing user connection:", err)
		}
	}
}

// return the users unclocked
func (r *Room) GetUsersUnlocked(except *melody.Session) []*User {
	users := make([]*User, 0, len(r.users))
	for _, u := range r.users {
		users = append(users, u)
	}
	return users
}

// return the users in the room
func (r *Room) GetUsers(except *melody.Session) []*User {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.GetUsersUnlocked(except)
}

// set the room display name
func (r *Room) SetDisplayName(displayName string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.displayName = displayName
}

// add a users to the room
func (r *Room) AddUser(s *melody.Session) *User {
	r.mu.Lock()
	defer r.mu.Unlock()

	users := r.GetUsersUnlocked(nil)

	// add user to room
	r.users[s] = &User{
		room: r,
		conn: s,
	}
	user := r.users[s]

	m := types.ClientMessage{
		Type:    "user_connected",
		UserID:  user.Id(),
		RoomID:  r.Id(),
		Payload: make(map[string]interface{}),
		Options: &types.MessageOptions{},
	}

	r.Broadcast(m, s)

	// TODO: Exchange peer information with new user
	for _, u := range users {

		log.Println("Exchange user information with new user", u.Id())
		if u.Id() == user.Id() {
			continue
		}

		m := types.ClientMessage{
			Type:    "user_connected",
			UserID:  u.Id(),
			RoomID:  r.Id(),
			Payload: make(map[string]interface{}),
			Options: &types.MessageOptions{},
		}

		err := user.WriteConn(m)

		if err != nil {
			log.Println("Error writing to user:", err)
		}
	}

	return user
}

func (r *Room) RemoveUser(s *melody.Session) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var m types.ClientMessage

	u := r.users[s]

	if u == nil {
		log.Println("user not found")
		return
	}

	m = types.ClientMessage{
		Type:    "user_disconnected",
		UserID:  u.Id(),
		RoomID:  r.Id(),
		Payload: make(map[string]interface{}),
		Options: &types.MessageOptions{},
	}

	r.Broadcast(m, s)

	// remove user from room
	delete(r.users, s)

	// FIXUP: refactor this for a better solution and remove this for another function wrapper in leaveRoom for exaple

	// send leave room message to user
	m = types.ClientMessage{
		Type:   "leave_room",
		RoomID: r.Id(),
	}

	err := u.WriteConn(m)
	if err != nil {
		log.Println("Error writing to user:", err)
	}

	log.Println("User removed successfully")
}

func (r *Room) Broadcast(msg types.ClientMessage, except *melody.Session) {
	r.mu.Lock()
	defer r.mu.Unlock()

	users := r.GetUsersUnlocked(except)

	for _, u := range users {
		err := u.WriteConn(msg)
		if err != nil {
			log.Fatalln("Error writing to user:", err)
		}
	}
}

// factory
func NewRoom(displayName string) *Room {
	return &Room{
		id:          uuid.New().String(),
		displayName: displayName,
		users:       make(map[*melody.Session]*User, 0),
		createdAt:   time.Now(),
	}
}
