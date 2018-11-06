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
		go h.MessageHandler()
		fmt.Println("created new hub")
	} else {
		h.joinUser(u)
		fmt.Println("hub exists, adding new user to it")
	}

	ws, err := upgrader.Upgrade(w, r, nil)

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

// FuzzyFindHubs : returns a list of matching hubs
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

		for _, h := range hubs {
			if fuzzy.Match(string(query), h.ID) {
				matches = append(matches, h)
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
