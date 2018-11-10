package main

import "github.com/gorilla/websocket"

// User : a user
type User struct {
	ID             string
	Username       string
	Password       string
	Friends        []*UserTag
	FriendRequests []*UserTag
	Hubs           []*HubTag

	ws *websocket.Conn // notifications websocket
}

// HubTag : hubs a user has joined
type HubTag struct {
	ID         string
	Visibility string
}

// UserTag : ID in friend list
type UserTag struct {
	ID       string
	Username string
}

// Hub : collection of users
type Hub struct {
	ID         string
	Visibility string
	Members    []*HubMember
	Messages   []*Message

	clients   map[*websocket.Conn]bool
	broadcast chan Message
}

// HubMember : of a Hub
type HubMember struct {
	Member  *User
	IsAdmin bool
}

// Message : message struct
type Message struct {
	ID           string
	Username     string
	Message      string
	JoinRequests []*UserTag
}

// Notification : holds notification data
type Notification struct {
	Recipient *User
	Type      string
}
