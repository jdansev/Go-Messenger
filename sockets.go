package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/lithammer/fuzzysearch/fuzzy"
)

// MessageHandler : relays messages to users in a hub
func (h *Hub) MessageHandler() {
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

// ConnectionHandler : connects the socket to the requested hub
func ConnectionHandler(w http.ResponseWriter, r *http.Request) {

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
		http.Error(w, "403 - you are not authorized!", http.StatusForbidden)
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
		fmt.Println("created new hub")
	} else {
		h.joinUser(u)
		fmt.Println("hub exists, adding new user to it")
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	fmt.Println(ws)

	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()
	h.clients[ws] = true
	u.ws = ws

	for {
		var msg Message
		err := ws.ReadJSON(&msg)

		if err != nil {
			log.Printf("error: %v", err)
			delete(h.clients, ws)
			u.ws = nil
			break
		}

		h.recordMessage(&msg)

		h.broadcast <- msg
	}

}

// TODO: combine searching hubs/users into one and distinguish them by a socket layer

// FuzzyFindHubs : returns a list of hubs with matching ids
func FuzzyFindHubs(w http.ResponseWriter, r *http.Request) {

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	for {
		mt, query, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		matches := []*Hub{}
		q := string(query)

		if q == "*" {
			matches = hubs
		} else if q != "" {
			for _, h := range hubs {
				if fuzzy.Match(q, h.ID) {
					matches = append(matches, h)
				}
			}
		}

		js, _ := json.Marshal(matches)

		err = c.WriteMessage(mt, js)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}

}

// FuzzyFindUsers : returns a list of users with matching ids
func FuzzyFindUsers(w http.ResponseWriter, r *http.Request) {

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	for {
		mt, query, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		matches := []*User{}
		q := string(query)

		if q == "*" {
			matches = users
		} else if q != "" {
			for _, u := range users {
				if fuzzy.Match(q, u.Username) {
					matches = append(matches, u)
				}
			}
		}
		js, _ := json.Marshal(matches)

		err = c.WriteMessage(mt, js)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}

}

// Notifications : basic echo socket with json
func Notifications(w http.ResponseWriter, r *http.Request) {

	var tok string
	var ok bool
	var u *User
	// var recipient *User

	// 1. Validate token from url
	if tok, ok = validateURLToken(w, r); !ok {
		return
	}

	// 2. Get the user's profile
	if u, ok = validateUserFromToken(tok, w); !ok {
		return
	}

	// 4. upgrade http to websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	defer ws.Close()

	// 5. set user's notification socket
	u.ws = ws

	// Begin notification handler for this user
	for {

		// detect when user breaks socket connection
		if _, _, err := ws.ReadMessage(); err != nil {
			log.Printf("error: %v", err)
			u.ws = nil
			break
		}

	}
}

// Notify : relays a notification to a specific user
func (n *Notification) Notify() bool {
	var err error

	notification := map[string]string{
		"Type": n.Type,
	}

	if n.Recipient.ws == nil {
		// user has no notifications socket
		return false
	}

	err = n.Recipient.ws.WriteJSON(notification)

	return err == nil
}