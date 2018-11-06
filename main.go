package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
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
	router.HandleFunc("/ws", connectionHandler)

	log.Fatal(http.ListenAndServe(":1212", router))
}

func messageHandler(h *Hub) {
	for {
		msg := <-h.broadcast
		for client := range h.clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(h.clients, client)
			}
		}
	}
}

// TODO: resolve data races

func connectionHandler(w http.ResponseWriter, r *http.Request) {

	var tokenParam, hidParam []string
	var ok bool

	tokenParam, ok = r.URL.Query()["token"]
	if !ok || len(tokenParam) < 1 {
		http.Error(w, "400 - token invalid!", http.StatusBadRequest)
		return
	}

	hidParam, ok = r.URL.Query()["hub"]
	if !ok || len(hidParam) < 1 {
		http.Error(w, "400 - no hub specified!", http.StatusBadRequest)
		return
	}

	tok := tokenParam[0]
	hid := hidParam[0]

	ok = validateToken(tok)
	if !ok {
		http.Error(w, "400 - token invalid!", http.StatusBadRequest)
		return
	}

	u := getUserFromToken(tok)
	if u == nil {
		http.Error(w, "500 - error fetching user!", http.StatusInternalServerError)
		return
	}

	// join an existing hub otherwise create it
	h := getHub(hid)
	if h == nil {
		h = u.createHub(hid)
		go messageHandler(h)
		fmt.Println("created new hub")
	} else {
		h.joinUser(u)
		fmt.Println("hub exists, adding new user to it")
	}

	ws, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)

	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()
	h.clients[ws] = true

	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		fmt.Println(msg)

		h.recordMessage(&msg)

		if err != nil {
			log.Printf("error: %v", err)
			delete(h.clients, ws)
			break
		}
		h.broadcast <- msg
	}

}