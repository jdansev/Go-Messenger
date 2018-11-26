package main

import (
	"encoding/json"

	"github.com/gorilla/websocket"
)

// User : a user
type User struct {
	ID              string
	Username        string
	password        string
	Friends         []*UserTag
	FriendRequests  []*UserTag
	Hubs            []*UserHub
	JoinInvitations []*HubInvitation

	ws *websocket.Conn // notifications websocket

	main *websocket.Conn // main websocket

	// attach UserTag to user? (not exported)
	// usertags should not be stored in memory?
}

// UserHub is a user's hubs
type UserHub struct {
	Tag        HubTag
	ReadLatest bool
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

// HubPreview : preview of a hub
type HubPreview struct {
	Tag         HubTag
	ReadLatest  bool
	LastMessage Message
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

/* Notification Structs */

// Notification : sending notification
type Notification struct {
	Type string
	Body interface{}
}

// N is receiving notification
type N struct {
	Type string // message, notification, friend request, invitation, etc.

	Body json.RawMessage
}

// HubMessage is a hub message
type HubMessage struct {
	Hub     HubTag
	Sender  UserTag
	Message string
}

// UserMessage is a user message
type UserMessage struct {
	Sender    UserTag
	Recipient UserTag
	Message   string
}

// FriendRequest is a friend request
type FriendRequest struct {
	From UserTag
	To   UserTag
}

// JoinInvitation is a join invitation
type JoinInvitation struct {
	HubID string
	From  UserTag
	To    UserTag
}
