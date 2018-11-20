package main

import (
	"fmt"
	"strconv"
	"time"
)

// TEST helpers
func addTestHubs() {

	createTestUsers()

	spec := Spectrum {
		"#38ef7d",
		"#11998e",
	}

	h2 := p1.createHub("p1-private-hub", "private", spec)

	p3.sendJoinRequest(h2)

	p2.sendJoinRequest(h2)

	h2.acceptJoinRequest(p1, p2)
	h2.declineJoinRequest(p1, p3)

	h2.grantAdmin(p2, p3)
	h2.grantAdmin(p1, p2)

	h2.unjoinUser(p1)

	u1 := createUser("asdf", "asdf")
	u2 := createUser("qwer", "qwer")
	u3 := createUser("zxcv", "zxcv")
	
	u1.sendFriendRequestTo(u2)
	u1.sendFriendRequestTo(u3)
	u2.acceptFriendRequest(u1)
	u3.acceptFriendRequest(u1)

}

var p1 *User
var p2 *User
var p3 *User

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
	h.Messages = append(h.Messages, &Message{"1", "hey there", "one"})
	h.Messages = append(h.Messages, &Message{"2", "whats up", "two"})
	h.Messages = append(h.Messages, &Message{"3", "how's it going", "three"})
}

var count = 1

func registerTestUserLoop() {
	for {
		createUser("new-user-"+strconv.Itoa(count), "my-secret-key")
		fmt.Printf("created new user, %d\n", count)
		count++
		time.Sleep(2 * time.Second)
	}
}
