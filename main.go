package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var p1 *User
var p2 *User
var p3 *User

var users = []*User{}
var hubs = []*Hub{}

func main() {
	fmt.Println("started server on port :1212")

	router := mux.NewRouter()

	addTestHubs()

	// queries
	router.HandleFunc("/hubs", GetHubs).Methods("GET")
	router.HandleFunc("/hubs/{hub_id}", GetHub).Methods("GET")

	router.HandleFunc("/users/{user_id}", GetUser).Methods("GET")
	router.HandleFunc("/users/{user_id}/hubs", GetUserHubs).Methods("GET")
	router.HandleFunc("/users/{user_id}/friends", GetUserFriends).Methods("GET")

	router.HandleFunc("/members/{hub_id}", GetMembers).Methods("GET")
	router.HandleFunc("/messages/{hub_id}", GetMessages).Methods("GET")

	// authentication
	router.HandleFunc("/register", Register).Methods("POST")
	router.HandleFunc("/login", Login).Methods("POST")

	// sockets
	router.HandleFunc("/ws", ConnectionHandler)
	router.HandleFunc("/ws/find-hubs", FuzzyFindHubs)
	router.HandleFunc("/ws/find-users", FuzzyFindUsers)

	log.Fatal(http.ListenAndServe(":1212", router))
}
