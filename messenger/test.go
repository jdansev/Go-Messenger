package main

import (
	"fmt"
	"strconv"
	"time"
)

var u1 *User
var u2 *User
var u3 *User

func seedDatabase() {
	u1 = createUser("Derek Smith", "derek")
	u2 = createUser("Lisa Ingham", "lisa")
	u3 = createUser("Fred Mercury", "fred")

	h1 := u1.createHub("Gang of One", "private", Spectrum{
		// "#38ef7d",
		// "#11998e",

		"#ff6a00",
		"#ee0979",
	})

	u1.sendFriendRequestTo(u2)
	u2.acceptFriendRequestFrom(u1)
	u2.sendFriendRequestTo(u3)
	u3.acceptFriendRequestFrom(u2)
	u3.sendFriendRequestTo(u1)
	u1.acceptFriendRequestFrom(u3)

	u2.sendJoinRequest(h1)
	u3.sendJoinRequest(h1)
	h1.acceptJoinRequest(u2, u1)
	h1.grantAdmin(u1, u2)
	h1.acceptJoinRequest(u3, u2)

	h1.saveMessage(&Message{"3", "Fred Mercury", "hey there"})
	h1.saveMessage(&Message{"3", "Fred Mercury", "how's it going everyone?"})
	h1.saveMessage(&Message{"3", "Fred Mercury", "nice to meet you all"})

	h1.saveMessage(&Message{"2", "Lisa Ingham", "You too!"})
	h1.saveMessage(&Message{"2", "Lisa Ingham", "The pleasure is mine"})

	h1.saveMessage(&Message{"3", "Fred Mercury", "What's up you guys?"})
	h1.saveMessage(&Message{"3", "Fred Mercury", "Love the new style"})
	
	h1.saveMessage(&Message{"2", "Lisa Ingham", "You really outdid yourself this time huh"})

	h1.saveMessage(&Message{"3", "Fred Mercury", "Sure did looks like it"})
	h1.saveMessage(&Message{"3", "Fred Mercury", "What an ass."})

	h1.saveMessage(&Message{"2", "Lisa Ingham", "Seconded."})
}

// TEST helpers
func addTestHubs() {

	createTestUsers()

	spec := Spectrum{
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

	// u1 := createUser("asdf", "asdf")
	// u2 := createUser("qwer", "qwer")
	// u3 := createUser("zxcv", "zxcv")

	// u1.sendFriendRequestTo(u2)
	// u1.sendFriendRequestTo(u3)
	// u2.acceptFriendRequest(u1)
	// u3.acceptFriendRequest(u1)

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
	p3.acceptFriendRequestFrom(p1)

	p2.sendFriendRequestTo(p1)
	p1.declineFriendRequestFrom(p2)

	p1.sendFriendRequestTo(p2)

	p3.sendFriendRequestTo(p1)
	p1.sendFriendRequestTo(p2)

	p2.acceptFriendRequestFrom(p3)

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
