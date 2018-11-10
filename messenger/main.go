package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var users = []*User{}
var hubs = []*Hub{}

func main() {
	fmt.Println("started server on port :1212")

	router := mux.NewRouter()

	addTestHubs()

	// test APIs
	router.HandleFunc("/nuke", NukeServer).Methods("GET")

	// public queries
	router.HandleFunc("/hubs", GetHubs).Methods("GET")
	router.HandleFunc("/hubs/{hub_id}", GetHub).Methods("GET")

	router.HandleFunc("/users/{user_id}", GetUser).Methods("GET")
	router.HandleFunc("/users/{user_id}/hubs", GetUserHubs).Methods("GET")
	router.HandleFunc("/users/{user_id}/friends", GetUserFriends).Methods("GET")

	router.HandleFunc("/members/{hub_id}", GetMembers).Methods("GET")
	router.HandleFunc("/messages/{hub_id}", GetMessages).Methods("GET")

	// actions (secure APIs) must include valid token
	router.HandleFunc("/create-hub", CreateHub).Methods("POST")
	router.HandleFunc("/send-friend-request", SendFriendRequest).Methods("POST")
	router.HandleFunc("/accept-friend-request", AcceptFriendRequest).Methods("POST")
	router.HandleFunc("/decline-friend-request", DeclineFriendRequest).Methods("POST")
	router.HandleFunc("/my-hubs", GetMyHubs).Methods("GET")
	router.HandleFunc("/my-friends", GetMyFriends).Methods("GET")
	router.HandleFunc("/my-friend-requests", GetMyFriendRequests).Methods("GET")

	// authentication
	router.HandleFunc("/register", Register).Methods("POST")
	router.HandleFunc("/login", Login).Methods("POST")

	// query sockets
	router.HandleFunc("/ws/find-hubs", FuzzyFindHubs)
	router.HandleFunc("/ws/find-users", FuzzyFindUsers)

	// notifications socket
	router.HandleFunc("/ws/notifications", Notifications)

	// chat socket
	router.HandleFunc("/ws", ConnectionHandler)

	log.Fatal(http.ListenAndServe(":1212", router))
}
