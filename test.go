package main

import (
	"fmt"
	"time"
	"strconv"
)

// TEST helpers
func addTestHubs() {

	// Manual creation of hubs
	createTestUsers()

	h := createHub("hub1")
	h.addTestUsersToHub()
	h.addTestMessagesToHub()

	addHub(h)

	// user creation of hubs
	h2 := p1.createHub("p1s-hub")
	go h2.MessageHandler()
	// p1.leaveHub(h2)

}

func createTestUsers() {
	p1 = createUser("testuser1", "secret-key")
	p2 = createUser("testuser2", "secret-key")
	p3 = createUser("testuser3", "secret-key")

	// add friends
	p1.sendFriendRequestTo(p3)
	p3.acceptFriendRequest(p1)

	p2.sendFriendRequestTo(p1)
	p1.declineFriendRequest(p2)

	p1.sendFriendRequestTo(p2)

	p3.sendFriendRequestTo(p1)
	p1.sendFriendRequestTo(p2)

	p2.acceptFriendRequest(p3)

	
}

func (h *Hub) addTestUsersToHub() {
	h.addUserToHub(p1)
	h.addUserToHub(p2)
	h.addUserToHub(p3)
}

func (h *Hub) addTestMessagesToHub() {
	h.Messages = append(h.Messages, &Message{"1","hey there"})
	h.Messages = append(h.Messages, &Message{"2","whats up"})
	h.Messages = append(h.Messages, &Message{"3","how's it going"})
}


var count = 1

func registerTestUserLoop() {
	for {
		createUser("new-user-" + strconv.Itoa(count), "my-secret-key")
		fmt.Printf("created new user, %d\n", count)
		count++
		time.Sleep(2 * time.Second)
	}
}