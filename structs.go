package main

import "github.com/gorilla/websocket"

// User : a user
type User struct {
	ID       string
	Username string
	Password string
	Friends  []*Friend
	Requests []*Friend
	Hubs     []*JoinedHubs
}

// JoinedHubs : hubs a user has joined
type JoinedHubs struct {
	ID string
}

// Friend : ID in friend list
type Friend struct {
	ID       string
	Username string
}

// TODO: public and private hubs

// Hub : collection of users
type Hub struct {
	ID       string
	Members  []*Member
	Messages []*Message

	clients   map[*websocket.Conn]bool
	broadcast chan Message
}

// Member : of a Hub
type Member struct {
	Member  *User
	IsAdmin bool
}

// Message : message struct
type Message struct {
	Sender  string
	Message string
}
