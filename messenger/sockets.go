package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/lithammer/fuzzysearch/fuzzy"
)

// NotificationHandler : handles all notifications across a single socket
func NotificationHandler(w http.ResponseWriter, r *http.Request) {

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

	// 4. Upgrade http to websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	defer ws.Close()

	// 5. Set user's main socket
	u.main = ws

	// Begin notification handler for this user
	for {

		var r N
		err := ws.ReadJSON(&r)

		// detect when user breaks socket connection
		if err != nil {
			log.Printf("error: %v", err)
			u.main = nil
			break
		}

		// Parse notification body

		// Handle hub messages
		if r.Type == "hubMessage" {

			var hm HubMessage
			var h *Hub

			json.Unmarshal(r.Body, &hm)

			if h = findHubByID(hm.Hub.ID); h == nil {
				fmt.Println("hub doesn't exist")
				break // hub doesn't exist
			}

			h.notifyMembers(hm)

		}

		// Handle direct friend messages
		if r.Type == "userMessage" {

			var msg UserMessage
			var recipient *User

			json.Unmarshal(r.Body, &msg)

			if recipient = findUserByID(msg.Recipient.ID); recipient == nil {
				fmt.Println("user doesn't exist")
				break
			}

			if ok := recipient.notify(msg); !ok {
				u.main.WriteJSON([]string{
					"message sent, recipient is offline!",
				})
			}
		}

		// Handle friend requests
		if r.Type == "friendRequest" {

			var fr FriendRequest
			var recipient *User

			json.Unmarshal(r.Body, &fr)
			fmt.Println(fr)

			if recipient = findUserByID(fr.To.ID); recipient == nil {
				fmt.Println("user doesn't exist")
				break
			}

			if ok := recipient.notify(fr); !ok {
				u.main.WriteJSON([]string{
					"friend request sent, recipient is offline!",
				})
			}
		}

		// Handler join invitations
		if r.Type == "JoinInvitation" {
			var ji JoinInvitation
			json.Unmarshal(r.Body, &ji)
			fmt.Println(ji)
		}

	}

}

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

		http.Error(w, "400 - hub doesn't exist!", http.StatusBadRequest)
		return

	}

	h.joinUser(u)
	fmt.Println("hub exists, adding new user to it")

	ws, err := upgrader.Upgrade(w, r, nil)

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

		h.saveMessage(&msg)

		// set all hub user's read latest to false
		for _, hubMember := range h.Members {

			for _, memberHub := range hubMember.Member.Hubs {

				// if websocket isn't open
				if memberHub.Tag.ID == h.ID && !h.clients[hubMember.Member.ws] {
					memberHub.ReadLatest = false;
				}

			}
		}

		// notify hub members a message was received
		hm := constructHubMessage(h, u, msg)
		n := constructNotification("hubMessage", hm)
		h.notifyMembers(n)

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
			for _, h := range hubs {
				if h.Visibility != "secret" {
					matches = append(matches, h)
				}
			}
		} else if q != "" {
			for _, h := range hubs {
				if fuzzy.Match(strings.ToLower(q), strings.ToLower(h.ID)) && h.Visibility != "secret" {
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
				if fuzzy.Match(strings.ToLower(q), strings.ToLower(u.Username)) {
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
