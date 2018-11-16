package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func runHubTests(h *Hub) {

	// TEST 1: make first member not admin
	if len(h.Members) > 1 {
		h.Members[0].setAdmin(true)
	}

	// TEST 2: change second sender id
	if len(h.Members) > 2 {
		h.Members[1].Member.ID = "2345"
	}

	// TEST 3: remove the last member
	if len(h.Members) > 1 {
		h.removeUserFromHub(h.Members[len(h.Members)-1].Member)
	}

	// TEST 4: remove two friends from p1
	if len(p1.Friends) > 0 {
		p1.removeFriend(p1.Friends[0])
	}
	if len(p1.Friends) > 0 {
		p1.removeFriend(p1.Friends[0])
	}

}

// NukeServer : deletes all existing users and hubs and repopulates test data
func NukeServer(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")

	nukeServerData()
}

// GetUser : returns a user from an id
func GetUser(w http.ResponseWriter, r *http.Request) {

	// 1. Get user id from path
	if u, ok := validateUserIDFromPath(w, r); ok {
		json.NewEncoder(w).Encode(u)
	}

}

// GetUsers : returns all users
func GetUsers(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(users)
}

// GetMessages : returns all messages for a hub
func GetMessages(w http.ResponseWriter, r *http.Request) {

	// 1. Get hub id from path
	if h, ok := validateHubIDFromPath(w, r); ok {
		json.NewEncoder(w).Encode(h.Messages)
	}

}

// GetMembers : returns all members in a hub
func GetMembers(w http.ResponseWriter, r *http.Request) {

	// 1. Get hub id from path
	if h, ok := validateHubIDFromPath(w, r); ok {
		json.NewEncoder(w).Encode((h.Members))
	}

}

// GetMyFriends : returns requesting user's friends
func GetMyFriends(w http.ResponseWriter, r *http.Request) {

	var tok string
	var ok bool
	var u *User

	w.Header().Set("Access-Control-Allow-Origin", "*")

	// 1. Validate token from url
	if tok, ok = validateURLToken(w, r); !ok {
		return
	}

	// 2. Get the user's profile
	if u, ok = validateUserFromToken(tok, w); !ok {
		return
	}

	// 3. Return user hubs
	json.NewEncoder(w).Encode(u.Friends)
}

/* HUB APIs */

// GetHubs : returns all hubs
func GetHubs(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(hubs)
}

// GetHub : returns hub with id
func GetHub(w http.ResponseWriter, r *http.Request) {

	// 1. Get hub id from path
	if h, ok := validateHubIDFromPath(w, r); ok {
		json.NewEncoder(w).Encode((h))
	}

}

// GetUserFriends : returns all user friends
func GetUserFriends(w http.ResponseWriter, r *http.Request) {

	// 1. Get user id from path
	if u, ok := validateUserIDFromPath(w, r); ok {
		json.NewEncoder(w).Encode(u.Friends)
	}

}

// GetUserHubs : returns all user hubs
func GetUserHubs(w http.ResponseWriter, r *http.Request) {

	// 1. Get user id from path
	if u, ok := validateUserIDFromPath(w, r); ok {
		json.NewEncoder(w).Encode(u.Hubs)
	}
}

// GetMyHubs : returns requesting user's hubs
func GetMyHubs(w http.ResponseWriter, r *http.Request) {

	var tok string
	var ok bool
	var u *User

	w.Header().Set("Access-Control-Allow-Origin", "*")

	// 1. Validate token from url
	if tok, ok = validateURLToken(w, r); !ok {
		return
	}

	// 2. Get the user's profile
	if u, ok = validateUserFromToken(tok, w); !ok {
		return
	}

	// 3. Return user hubs
	json.NewEncoder(w).Encode(u.Hubs)
}

// CreateHub : allow users to start a new hub
func CreateHub(w http.ResponseWriter, r *http.Request) {

	var tok string
	var ok bool
	var u *User

	w.Header().Set("Access-Control-Allow-Origin", "*")

	// 1. Validate token from url
	if tok, ok = validateURLToken(w, r); !ok {
		return
	}

	// 2. Get the user's profile
	if u, ok = validateUserFromToken(tok, w); !ok {
		return
	}

	vis := r.FormValue("hub_visibility")

	specStart := r.FormValue("hub_spec_start")
	specEnd := r.FormValue("hub_spec_end")

	spec := Spectrum{
		specStart,
		specEnd,
	}


	// 3. Create the new hub
	hid := r.FormValue("hub_id")
	h := u.createHub(hid, vis, spec)
	if h == nil {
		http.Error(w, "400 - hub already exists!", http.StatusBadRequest)
		return
	}

	fmt.Println(h)

	// 4. Return it
	json.NewEncoder(w).Encode(h)
}

// TODO: change all post requests in template to not use forms

// SendFriendRequest : sends a request to the user with id
func SendFriendRequest(w http.ResponseWriter, r *http.Request) {

	var tok string
	var ok bool
	var u, fu *User

	w.Header().Set("Access-Control-Allow-Origin", "*")

	// 1. Validate token
	if tok, ok = validateFormToken(w, r); !ok {
		return
	}

	// 2. Get the user's profile
	if u, ok = validateUserFromToken(tok, w); !ok {
		return
	}

	// 3. Get the user who will receive the friend request
	if fu, ok = validateUserIDFromForm(w, r); !ok || u == fu {
		return
	}

	// 4. Send the friend request
	if ok = u.sendFriendRequestTo(fu); !ok {
		http.Error(w, "400 - cannot do that!", http.StatusBadRequest)
		return
	}

	// 5. Send the notification
	n := &Notification{
		fu,
		"friendRequestReceived",
		UserTag{u.ID, u.Username},
	}

	ok = n.Notify()
	if !ok {
		u.ws.WriteJSON([]string{
			"request sent, recipient is offline!",
		})
	}
}

// AcceptFriendRequest : accepts friend request from user with id
func AcceptFriendRequest(w http.ResponseWriter, r *http.Request) {

	var tok string
	var ok bool
	var u, fu *User

	w.Header().Set("Access-Control-Allow-Origin", "*")

	// 1. Validate token
	if tok, ok = validateFormToken(w, r); !ok {
		return
	}

	// 2. Get the user's profile
	if u, ok = validateUserFromToken(tok, w); !ok {
		return
	}

	// 3. Get the user who's request will be accepted
	if fu, ok = validateUserIDFromForm(w, r); !ok || u == fu {
		return
	}

	// 4. Accept the request
	if ok = u.acceptFriendRequest(fu); !ok {
		http.Error(w, "400 - cannot accept request from this user!", http.StatusBadRequest)
		return
	}

	// 5. Notify the accepting user
	n1 := &Notification{
		u,
		"youAcceptedFriendRequest",
		UserTag{fu.ID, fu.Username},
	}
	ok = n1.Notify()
	if !ok {
		u.ws.WriteJSON([]string{
			"user is not connected!",
		})
	}

	// 6. Notify the requesting user
	n2 := &Notification{
		fu,
		"requestAccepted",
		UserTag{u.ID, u.Username},
	}
	fmt.Println(n2)
	ok = n2.Notify()
	if !ok {
		u.ws.WriteJSON([]string{
			"request accepted, recipient is offline!",
		})
	}

}

// DeclineFriendRequest : declines friend request from user with id
func DeclineFriendRequest(w http.ResponseWriter, r *http.Request) {

	var tok string
	var ok bool
	var u, fu *User

	w.Header().Set("Access-Control-Allow-Origin", "*")

	// 1. Validate token
	if tok, ok = validateFormToken(w, r); !ok {
		return
	}

	// 2. Get the user's profile
	if u, ok = validateUserFromToken(tok, w); !ok {
		return
	}

	// 3. Get the user who's request will be declined
	if fu, ok = validateUserIDFromForm(w, r); !ok || u == fu {
		return
	}

	// 4. Decline the request
	if ok = u.declineFriendRequest(fu); !ok {
		http.Error(w, "400 - cannot decline request from this user!", http.StatusBadRequest)
		return
	}

	// 6. Notify the declining user
	n := &Notification{
		u,
		"youDeclinedFriendRequest",
		UserTag{fu.ID, fu.Username},
	}
	n.Notify()

}

// GetMyFriendRequests : returns requesting user's friend requests
func GetMyFriendRequests(w http.ResponseWriter, r *http.Request) {

	var tok string
	var ok bool
	var u *User

	w.Header().Set("Access-Control-Allow-Origin", "*")

	// 1. Validate token from url
	if tok, ok = validateURLToken(w, r); !ok {
		return
	}

	// 2. Get the user's profile
	if u, ok = validateUserFromToken(tok, w); !ok {
		return
	}

	// 3. Return user friend requests
	json.NewEncoder(w).Encode(u.FriendRequests)
}

// GetHubMessages : returns hub message history
func GetHubMessages(w http.ResponseWriter, r *http.Request) {

	var tok string
	var ok bool
	var u *User
	var h *Hub

	w.Header().Set("Access-Control-Allow-Origin", "*")

	// 1. Validate token from url
	if tok, ok = validateURLToken(w, r); !ok {
		return
	}

	// 2. Get the user's profile
	if u, ok = validateUserFromToken(tok, w); !ok {
		return
	}

	// 3. Validate hub id from path
	if h, ok = validateHubIDFromPath(w, r); !ok {
		return
	}

	// Check that user is a member of this hub
	if !u.isMemberOf(h) {
		return
	}

	// 4. Return user friend requests
	json.NewEncoder(w).Encode(h.Messages)

}