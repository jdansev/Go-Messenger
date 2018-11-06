package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
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

// GetUser : returns a user from an id
func GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	user := findUserByID(params["user_id"])
	if user == nil {
		http.Error(w, "400 - user doesn't exist!", http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(user)
}

// GetUsers : returns all users
func GetUsers(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(users)
}

// GetMessages : returns all messages for a hub
func GetMessages(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	hub := getHub(params["hub_id"])
	if hub == nil {
		http.Error(w, "400 - hub doesn't exist!", http.StatusBadRequest)
		return
	}

	runHubTests(hub)

	json.NewEncoder(w).Encode(hub.Messages)
}

// GetMembers : returns all members in a hub
func GetMembers(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	hub := getHub(params["hub_id"])
	if hub == nil {
		http.Error(w, "400 - hub doesn't exist!", http.StatusBadRequest)
		return
	}

	runHubTests(hub)

	json.NewEncoder(w).Encode((hub.Members))
}

// GetHub : returns hub with id
func GetHub(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	hub := getHub(params["hub_id"])
	if hub == nil {
		http.Error(w, "400 - hub doesn't exist!", http.StatusBadRequest)
		return
	}

	runHubTests(hub)

	json.NewEncoder(w).Encode(hub)
}

// GetHubs : returns all hubs
func GetHubs(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(hubs)
}

// GetUserHubs : returns all user hubs
func GetUserHubs(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	user := findUserByID(params["user_id"])
	if user == nil {
		http.Error(w, "400 - user doesn't exist!", http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(user.Hubs)
}

// GetUserFriends : returns all user friends
func GetUserFriends(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	user := findUserByID(params["user_id"])
	if user == nil {
		http.Error(w, "400 - user doesn't exist!", http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(user.Friends)
}

// CreateHub : allow users to start a new hub
func CreateHub(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")

	// 1. Validate token

	var tokenParam []string
	var ok bool

	tokenParam, ok = r.URL.Query()["token"]

	if !ok || len(tokenParam) < 1 || tokenParam[0] == "undefined" {
		http.Error(w, "400 - token invalid!", http.StatusBadRequest)
		return
	}

	tok := tokenParam[0]
	ok = validateToken(tok)
	if !ok {
		http.Error(w, "403 - you are not authorized to create hubs!", http.StatusForbidden)
		return
	}

	// 2. Get the user's profile
	u := getUserFromToken(tok)
	if u == nil {
		http.Error(w, "500 - error fetching user!", http.StatusInternalServerError)
		return
	}

	// // 2. Create the new hub

	hid := r.FormValue("hub_id")

	h := u.createHub(hid)
	if h == nil {
		http.Error(w, "400 - hub already exists!", http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(h)
}
