package main
// MEMBER helpers

func (m *HubMember) setAdmin(a bool) {
	m.IsAdmin = a
}

func (m *HubMember) setOwner(o bool) {
	m.IsOwner = o
}