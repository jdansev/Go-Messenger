package main

import (
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

/* User Tag Helpers */
func (ut *UserTag) getUserFromTag() *User {
	for _, u := range users {
		if ut.ID == u.ID {
			return u
		}
	}
	return nil
}

/* User Helpers */

// TODO: convert all for loops into contains and remove


// Hub join Invites
func (u *User) acceptJoinInviteFrom(h *Hub) bool {
	if !u.hasJoinInviteFrom(h) {
		return false
	}
	// remove the request from user's invitations
	u.removeJoinInviteFrom(h)
	h.joinUser(u)
	return true
}

func (u *User) declineJoinInviteFrom(h *Hub) bool {
	if !u.hasJoinInviteFrom(h) {
		return false
	}
	// remove the request from user's invitations
	u.removeJoinInviteFrom(h)
	return true
}

func (u *User) removeJoinInviteFrom(h *Hub) bool {
	r := u.JoinInvitations
	for i, invitation := range r {
		if invitation.Hub.ID == h.ID {
			*r[i] = *r[len(r)-1]
			u.JoinInvitations = r[:len(r)-1]
			return true
		}
	}
	return false
}

func (u *User) hasJoinInviteFrom(h *Hub) bool {
	for _, invite := range u.JoinInvitations {
		if invite.Hub.ID == h.ID {
			return true
		}
	}
	return false
}

// Join Requests
func (u *User) sendJoinRequest(h *Hub) bool {
	// only if the hub is private
	if h.Visibility != "private" {
		return false
	}
	// and the user hasn't already requested
	if h.hasJoinRequestFrom(u) {
		return false
	}
	r := &UserTag{u.ID, u.Username}
	h.JoinRequests = append(h.JoinRequests, r)
	return true
}

// Create Users
func createUser(username, password string) *User {
	uid, _ := uuid.NewV4()
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	u := &User{
		uid.String(),
		username,
		string(passwordHash),
		[]*UserTag{},
		[]*UserTag{},
		[]*HubTag{},
		[]*HubInvitation{},
		nil,
	}
	addUser(u)
	return u
}

// Friend Requests
func (u *User) hasRequested(f *User) bool {
	for _, friend := range f.FriendRequests {
		if friend.ID == u.ID {
			return true
		}
	}
	return false
}

func (u *User) sendFriendRequestTo(f *User) bool {
	if u.isFriendsWith(f) || u.hasRequested(f) {
		return false // already friends or already requested
	}
	request := &UserTag{u.ID, u.Username}
	f.FriendRequests = append(f.FriendRequests, request)
	return true
}

func (u *User) acceptFriendRequest(f *User) bool {
	if !f.hasRequested(u) { // didn't send a request
		return false
	}
	u.addFriend(f)
	f.addFriend(u)
	u.removeFriendRequest(f)
	return true
}

func (u *User) declineFriendRequest(f *User) bool {
	if u.isFriendsWith(f) || !f.hasRequested(u) {
		return false // already friends or haven't requested
	}
	u.removeFriendRequest(f)
	return true
}

func (u *User) removeFriendRequest(f *User) bool {
	r := u.FriendRequests
	for i, friend := range r {
		if f.ID == friend.ID {
			*r[i] = *r[len(r)-1]
			u.FriendRequests = r[:len(r)-1]
			return true
		}
	}
	return false
}

// TODO: deactivate user

// Friends
func (u *User) isFriendsWith(f *User) bool {
	for _, friend := range u.Friends {
		if friend.ID == f.ID {
			return true
		}
	}
	return false
}

func (u *User) addFriend(f *User) {
	if u.isFriendsWith(f) || u == f {
		return
	}
	u.Friends = append(u.Friends, &UserTag{f.ID, f.Username})
}

func (u *User) removeFriend(f *UserTag) bool {
	m := u.Friends
	for i, friend := range m {
		if f == friend {
			*m[i] = *m[len(m)-1]
			u.Friends = m[:len(m)-1]
			return true
		}
	}
	return false
}

// Hubs
func (u *User) createHub(id, vis string, spec Spectrum) *Hub {
	if getHub(id) != nil { // hub already exists, so don't create it
		return nil
	}
	h := createHub(id, vis, spec)
	addHub(h)
	h.joinUser(u)
	h.setAdmin(u, true)
	h.setOwner(u, true)
	return h
}

func (u *User) joinHub(h *Hub) {
	u.Hubs = append(u.Hubs, &HubTag{h.ID, h.Visibility, h.Spectrum})
}

func (u *User) leaveHub(h *Hub) bool {
	htags := u.Hubs
	for i, htag := range htags {
		hub := htag.getHubFromTag()
		if hub == h {
			*htags[i] = *htags[len(htags)-1]
			u.Hubs = htags[:len(htags)-1]
			return true
		}
	}
	return false
}

func (u *User) isMemberOf(h *Hub) bool {
	for _, m := range h.Members {
		if u.ID == m.Member.ID {
			return true
		}
	}
	return false
}
