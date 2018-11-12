package main

import (
	"testing"
)

var u1 *User
var u2 *User
var u3 *User

func TestFriends(t *testing.T) {

	u1 = createUser("testuser1", "secret-key")
	u2 = createUser("testuser2", "secret-key")
	u3 = createUser("testuser3", "secret-key")

	// u2 accepts u1's friend request
	u1.sendFriendRequestTo(u2)

	if !u1.hasRequested(u2) {
		t.Error("Expected u1 has requested u2")
	}

	u2.acceptFriendRequest(u1)

	// u1 is friends with u2
	if u1.isFriendsWith(u2) == false {
		t.Error("Expected u1 is friends with u2")
	}

	// u1 declines u3's friend request
	u3.sendFriendRequestTo(u1)
	u1.declineFriendRequest(u2)

	if u1.isFriendsWith(u2) == false {
		t.Error("Expected u3 is not friends with u1")
	}

	// u1 can't send a request to a friend
	if u1.sendFriendRequestTo(u2) == true {
		t.Error("Excepted u1 cannot send friend request")
	}

}


func TestHubs(t *testing.T) {

	pub := u1.createHub("public-hub", "public")
	priv := u2.createHub("private-hub", "private")
	secr := u3.createHub("secret-hub", "secret")

	// created hubs are not nil
	if pub == nil || priv == nil || secr == nil {
		t.Error("Expected created hubs to not be nil")
	}

	// can't send join request when you're already a member
	if u1.sendJoinRequest(pub) == true {
		t.Error("Expected can't send join request when you're already a member")
	}

	// can't send join request to public hub
	if u2.sendJoinRequest(pub) == true {
		t.Error("Expected can't send join request to public hub")
	}

	// can send join request to private hub
	if u3.sendJoinRequest(priv) == false {
		t.Error("Expected can send join request to private hub")
	}

	// can't send join request to secret hub
	if u1.sendJoinRequest(secr) == true {
		t.Error("Expected can't send join request to secret hub")
	}

	// non members or non admins can't accept join requests
	if priv.acceptJoinRequest(u1, u3) == true {
		t.Error("Expected non members or non admins can't accept join requests")
	}

	// multiple requests by the same user are not allowed
	if u3.sendJoinRequest(priv) == true {
		t.Error("Expected non members or non admins can't accept join requests")
	}

	// once accepted user becomes a member
	if !priv.acceptJoinRequest(u2, u3) || !u3.isMemberOf(priv) {
		t.Error("Expected non members or non admins can't accept join requests")
	}

	u1.sendJoinRequest(priv)

	// when declined user doesn't become a member
	if !priv.declineJoinRequest(u2, u1) || u1.isMemberOf(priv) {
		t.Error("Expected non members or non admins can't accept join requests")
	}

}
