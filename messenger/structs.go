package main

import "github.com/gorilla/websocket"

// User : a user
type User struct {
	ID              string
	Username        string
	password        string
	Friends         []*UserTag
	FriendRequests  []*UserTag
	Hubs            []*HubTag
	JoinInvitations []*HubInvitation

	ws *websocket.Conn // notifications websocket

	// attach UserTag to user? (not exported)
	// usertags should not be stored in memory?
}

// UserTag : ID in friend list
type UserTag struct {
	ID       string
	Username string
}

// Spectrum : a hub's designated color gradient
type Spectrum struct {
	Start string
	End   string
}

// Hub : collection of users
type Hub struct {
	ID           string
	Visibility   string
	Members      []*HubMember
	Messages     []*Message
	JoinRequests []*UserTag

	Spectrum Spectrum

	clients   map[*websocket.Conn]bool
	broadcast chan Message
}

// HubInvitation : an invitation to join a hub
type HubInvitation struct {
	Hub  HubTag
	From UserTag
}

// HubTag : hubs a user has joined
type HubTag struct {
	ID         string
	Visibility string
	Spectrum   Spectrum
}

// HubMember : joined users of a hub
type HubMember struct {
	Member *User

	/* Admins
	- can invite users (private and secret hubs only)
	- can accept join requests (private hubs only)
	- can remove members
	- can change hub visibility
	- can change hub details
	*/
	IsAdmin bool

	/* Owners
	- have admin privileges plus
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
	Message  string
}

// Notification : holds notification data
type Notification struct {
	// Sender *UserTag
	recipient *User
	Type      string
	Body      interface{}
}
