package main

import (
	"fmt"
	"log"
	"net/http"
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


// ConnectionHandler : connects the socket to the requested hub
func ConnectionHandler(w http.ResponseWriter, r *http.Request) {

	// var tokenParam []string
	var hidParam []string
	var ok bool

	// tokenParam, ok = r.URL.Query()["token"]
	// if !ok || len(tokenParam) < 1 {
	// 	http.Error(w, "400 - no token specified!", http.StatusBadRequest)
	// 	return
	// }

	hidParam, ok = r.URL.Query()["hub"]
	if !ok || len(hidParam) < 1 {
		http.Error(w, "400 - no hub specified!", http.StatusBadRequest)
		return
	}

	// tok := tokenParam[0]
	hid := hidParam[0]

	// ok = validateToken(tok)
	// if !ok {
	// 	http.Error(w, "403 - token invalid!", http.StatusForbidden)
	// 	return
	// }

	// // validate hub
	// hub := getHub(hid)
	// if hub == nil {
	// 	// http.Error(w, "400 - hub doesn't exist!", http.StatusForbidden)
	// 	hub = createHub(hid)
	// 	go hub.MessageHandler()
	// 	return
	// }


	hub := createHub(hid)

	fmt.Println("made it here")

	// upgrade http connection to websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	hub.clients[ws] = true

	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(hub.clients, ws)
			break
		}
		hub.broadcast <- msg
	}

}
