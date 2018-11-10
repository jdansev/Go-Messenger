package main

import "github.com/gorilla/websocket"

// HUB Tag helpers
func (ht *HubTag) getHubFromTag() *Hub {
	for _, h := range hubs {
		if h.ID == ht.ID {
			return h
		}
	}
	return nil
}

// HUB helpers

func (h *Hub) grantAdmin(a, b *User) bool {

	// one of them is not a member
	if !a.isMemberOf(h) || !b.isMemberOf(h)  {
		return false
	}

	// member granting admin status is not themself an admin
	if !h.getHubMemberFromUser(a).IsOwner {
		return false
	}

	h.getHubMemberFromUser(b).setAdmin(true)
	return true
}

func (h *Hub) declineJoinRequest(m, r *User) {

	if !h.hasJoinRequestFrom(r) {
		return // didn't request to join
	}

	if !m.isMemberOf(h) || r.isMemberOf(h) {
		return // user to accept is not a member or requeter is already a member
	}

	var hm *HubMember

	if hm = h.getHubMemberFromUser(m); hm == nil {
		return // hub member doesn't exist
	}

	if !hm.IsAdmin { // not an admin
		return
	}

	h.removeJoinRequest(r)
}

func (h *Hub) acceptJoinRequest(m, r *User) { // member to accept the requesting user

	if !h.hasJoinRequestFrom(r) {
		return // didn't request to join
	}

	if !m.isMemberOf(h) {
		return // user to accept is not a member
	}

	var hm *HubMember

	if hm = h.getHubMemberFromUser(m); hm == nil {
		return // hub member doesn't exist
	}

	if !hm.IsAdmin { // not an admin
		return
	}

	h.joinUser(r)
	h.removeJoinRequest(r)
}

func (h *Hub) removeJoinRequest(u *User) bool {
	jrs := h.JoinRequests
	for i, jr := range jrs {
		if jr.ID == u.ID {
			*jrs[i] = *jrs[len(jrs)-1]
			h.JoinRequests = jrs[:len(jrs)-1]
			return true
		}
	}
	return false

}

func (h *Hub) getHubMemberFromUser(u *User) *HubMember {
	for _, m := range h.Members {
		if u.ID == m.Member.ID {
			return m
		}
	}
	return nil
}

func (h *Hub) hasJoinRequestFrom(u *User) bool {
	for _, r := range h.JoinRequests {
		if r.ID == u.ID {
			return true
		}
	}
	return false
}

func removeHub(h *Hub) bool {
	for i, hub := range hubs {
		if hub == h {
			*hubs[i] = *hubs[len(hubs)-1]
			hubs = hubs[:len(hubs)-1]
			return true
		}
	}
	return false
}

func createHub(id, vis string) *Hub {
	h := &Hub{
		ID:           id,
		Visibility:   vis,
		Members:      []*HubMember{},
		Messages:     []*Message{},
		JoinRequests: []*UserTag{},

		broadcast: make(chan Message),
		clients:   make(map[*websocket.Conn]bool),
	}
	go h.MessageHandler()
	return h
}

func getHub(id string) *Hub {
	for _, h := range hubs {
		if h.ID == id {
			return h
		}
	}
	return nil
}

func (h *Hub) joinUser(u *User) {
	if u.isMemberOf(h) { // already a member
		return
	}
	h.addUserToHub(u)
	u.joinHub(h)
}

func (h *Hub) unjoinUser(u *User) {
	if !u.isMemberOf(h) { // not a member
		return
	}
	h.removeUserFromHub(u)
	u.leaveHub(h)
}

func (h *Hub) setAdmin(u *User, isAdmin bool) bool {
	for _, m := range h.Members {
		if m.Member == u {
			m.setAdmin(isAdmin)
			return true
		}
	}
	return false
}

func (h *Hub) setOwner(u *User, isOwner bool) bool {
	for _, m := range h.Members {
		if m.Member == u {
			m.setOwner(isOwner)
			return true
		}
	}
	return false
}

func (h *Hub) addUserToHub(p *User) {
	h.Members = append(h.Members, &HubMember{p, false, false})
}

func (h *Hub) removeUserFromHub(p *User) bool {
	m := h.Members
	for i, member := range m {
		if p == member.Member {
			*m[i] = *m[len(m)-1]
			h.Members = m[:len(m)-1]
			return true
		}
	}
	return false
}

func (h *Hub) recordMessage(m *Message) {
	h.Messages = append(h.Messages, m)
}
