package main

import "github.com/gorilla/websocket"

// User : a user
type User struct {
	ID             string
	Username       string
	password       string
	Friends        []*UserTag
	FriendRequests []*UserTag
	Hubs           []*HubTag

	ws *websocket.Conn // notifications websocket

	// attach UserTag to user? (not exported)
	// usertags should not be stored in memory?
}

// UserTag : ID in friend list
type UserTag struct {
	ID       string
	Username string
}

// Hub : collection of users
type Hub struct {
	ID           string
	Visibility   string
	Members      []*HubMember
	Messages     []*Message
	JoinRequests []*UserTag

	clients   map[*websocket.Conn]bool
	broadcast chan Message
}

// HubTag : hubs a user has joined
type HubTag struct {
	ID         string
	Visibility string
}

// HubMember : joined users of a hub
type HubMember struct {

	Member  *User

	/* Admin
	- invite users (private and secret hubs only)
	- accept join requests (private hubs only)
	- remove members
	- change hub visibility
	- change hub details
	*/
	IsAdmin bool

	/* Owner
	- has admin privileges plus
	- can assign other members as admin
	- can delete the hub
	*/
	IsOwner bool

}

// Message : message struct
type Message struct {
	// Sender UserTag
	ID       string
	Username string

	Message string
}

// Notification : holds notification data
type Notification struct {
	// Sender *UserTag
	Recipient *User
	Type      string
}
