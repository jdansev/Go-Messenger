package main

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

// GLOBAL helpers

func addHub(h *Hub) {
	hubs = append(hubs, h)
}

func addUser(u *User) {
	users = append(users, u)
}

func findUserByID(id string) *User {
	for _, u := range users {
		if u.ID == id {
			return u
		}
	}
	return nil
}

func findUserByName(usr string) *User {
	for _, u := range users {
		if u.Username == usr {
			return u
		}
	}
	return nil
}

// TOKEN helpers

func (u *User) generateToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
		"user_id": u.ID,
	})
	tokenString, _ := token.SignedString(mySigningKey)
	return tokenString
}

func validateToken(tokenString string) bool {
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	})
	return err == nil
}

func getUserFromToken(t string) *User {
	token, _ := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	})
	claims := token.Claims.(jwt.MapClaims)
	uid := claims["user_id"].(string)
	return findUserByID(uid)
}

// MEMBER helpers

func (m *Member) setAdmin(a bool) {
	m.IsAdmin = a
}

// HUB helpers

func removeHub(h *Hub) bool {
	for i, hub := range hubs {
		if hub == h {
			*hubs[i] = *hubs[len(hubs)-1]
			hubs = hubs[:len(hubs)-1]
			return true
		}
	}
	return false
}

func createHub(id string) *Hub {
	return &Hub{
		ID:        id,
		broadcast: make(chan Message),
		clients:   make(map[*websocket.Conn]bool),
		Members:   []*Member{},
		Messages:  []*Message{},
	}
}

func (j *JoinedHubs) getJoinedHub() *Hub {
	for _, h := range hubs {
		if j.ID == h.ID {
			return h
		}
	}
	return nil
}

func getHub(id string) *Hub {
	for _, h := range hubs {
		if h.ID == id {
			return h
		}
	}
	return nil
}

func (h *Hub) joinUser(u *User) {
	h.addUserToHub(u)
	u.joinHub(h)
}

func (h *Hub) unjoinUser(u *User) {
	h.removeUserFromHub(u)
	u.leaveHub(h)
}

func (h *Hub) setAdmin(u *User, isAdmin bool) bool {
	for _, m := range h.Members {
		if m.Member == u {
			m.setAdmin(isAdmin)
			return true
		}
	}
	return false
}

func (h *Hub) addUserToHub(p *User) {
	h.Members = append(h.Members, &Member{p, false})
}

func (h *Hub) removeUserFromHub(p *User) bool {
	m := h.Members
	for i, member := range m {
		if p == member.Member {
			*m[i] = *m[len(m)-1]
			h.Members = m[:len(m)-1]
			return true
		}
	}
	return false
}

func (h *Hub) recordMessage(m *Message) {
	h.Messages = append(h.Messages, m)
}

// FRIEND helpers
func (f *Friend) getFriendUser() *User {
	for _, u := range users {
		if f.ID == u.ID {
			return u
		}
	}
	return nil
}

// User helpers

func createUser(username, password string) *User {
	uid, _ := uuid.NewV4()
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	u := &User{
		uid.String(),
		username,
		string(passwordHash),
		[]*Friend{},
		[]*JoinedHubs{},
	}
	addUser(u)
	return u
}

// TODO: deactivate user

func (u *User) isFriendsWith(f *User) bool {
	for _, friend := range u.Friends {
		if friend.ID == f.ID {
			return true
		}
	}
	return false
}

func (u *User) addFriend(f *User) {
	if u.isFriendsWith(f) || u == f {
		return
	}
	u.Friends = append(u.Friends, &Friend{f.ID})
}

func (u *User) removeFriend(f *Friend) bool {
	m := u.Friends
	for i, friend := range m {
		if f == friend {
			*m[i] = *m[len(m)-1]
			u.Friends = m[:len(m)-1]
			return true
		}
	}
	return false
}

func (u *User) createHub(id string) *Hub {
	if getHub(id) != nil { // hub already exists, so don't create it
		return nil
	}
	h := createHub(id)
	addHub(h)
	h.joinUser(u)
	h.setAdmin(u, true)
	return h
}

func (u *User) joinHub(h *Hub) {
	u.Hubs = append(u.Hubs, &JoinedHubs{h.ID})
}

func (u *User) leaveHub(h *Hub) bool {
	jhs := u.Hubs
	for i, jh := range jhs {
		hub := jh.getJoinedHub()
		if hub == h {
			*jhs[i] = *jhs[len(jhs)-1]
			u.Hubs = jhs[:len(jhs)-1]
			return true
		}
	}
	return false
}