package main

// GLOBAL helpers

func nukeServerData() {
	users = []*User{}
	hubs = []*Hub{}
	addTestHubs()
}

func addHub(h *Hub) {
	hubs = append(hubs, h)
}

func addUser(u *User) {
	users = append(users, u)
}

func findUserByID(id string) *User {
	for _, u := range users {
		if u.ID == id {
			return u
		}
	}
	return nil
}

func findUserByName(usr string) *User {
	for _, u := range users {
		if u.Username == usr {
			return u
		}
	}
	return nil
}
