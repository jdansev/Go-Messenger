package main

func (u *User) notify(n interface{}) bool {
	if u.main == nil {
		return false // user has no main socket
	}
	u.main.WriteJSON(n)
	return true
}

func (h *Hub) notifyMembers(n interface{}) {
	for _, m := range h.Members {
		m.Member.notify(n)
	}
}

func constructFriendRequest(from, to *User) *FriendRequest {
	return &FriendRequest{
		from.getTagFromUser(),
		to.getTagFromUser(),
	}
}

func constructHubMessage(hub *Hub, sender *User, msg Message) *HubMessage {
	return &HubMessage{
		hub.getTagFromHub(),
		sender.getTagFromUser(),
		msg.Message,
	}
}

func constructNotification(t string, n interface{}) *Notification {
	return &Notification{t, n}
}
