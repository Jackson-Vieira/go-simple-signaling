package domain

import (
	"log"
	"sync"
	"time"

	"github.com/Jackson-Vieira/go-simple-signalling/types"
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
			log.Println("Error closing peer connection:", err)
		}
	}
}

// return the users unclocked
func (r *Room) GetUsersUnlocked(except *User) []*User {
	users := make([]*User, 0, len(r.users))
	for _, u := range r.users {
		users = append(users, u)
	}
	return users
}

// return the peers in the room
func (r *Room) GetUsers(except *User) []*User {
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

// add a peer to the room
func (r *Room) AddUser(s *melody.Session) *User {
	r.mu.Lock()
	defer r.mu.Unlock()

	// bind user to room and session
	r.users[s] = &User{
		room: r,
		conn: s,
	}

	return r.users[s]
}

func (r *Room) RemoveUser(s *melody.Session) {
	r.mu.Lock()
	defer r.mu.Unlock()

	u := r.users[s]

	if u == nil {
		log.Println("peer not found")
		return
	}

	// disconnect user
	err := u.Disconnect()
	if err != nil {
		log.Println("Error closing a user connection:", err)
	}

	// remove peer from room
	delete(r.users, s)
}

func (r *Room) Broadcast(msg types.ClientMessage, except *User) {
	r.mu.Lock()
	defer r.mu.Unlock()

	users := r.GetUsersUnlocked(except)

	for _, u := range users {
		err := u.WriteConn(msg)
		if err != nil {
			log.Fatalln("Error writing to peer:", err)
		}
	}
}
