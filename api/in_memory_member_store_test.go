package api

import (
	"errors"
	"memberserver/api/models"
	"sync"
)

func NewInMemoryMemberStore() *InMemoryMemberStore {
	return &InMemoryMemberStore{}
}

type InMemoryMemberStore struct {
	store map[string]models.Member
	// A mutex is used to synchronize read/write access to the map
	lock sync.RWMutex
}

func (m *InMemoryMemberStore) GetMembers() []models.Member {
	m.lock.Lock()
	defer m.lock.Unlock()
	var members []models.Member

	for _, member := range m.store {
		members = append(members, member)
	}

	return members
}

type StubMemberStore struct {
	members map[string]models.Member
	tiers   []models.Tier
}

var testMemberStore = StubMemberStore{
	map[string]models.Member{
		"test@test.com": {
			ID:        "0",
			Name:      "testuser",
			Email:     "test@test.com",
			RFID:      "rfid1",
			Level:     0,
			Resources: []models.MemberResource{},
		},
	},
	[]models.Tier{
		{ID: 0,
			Name: "fake-inactive"},
		{ID: 1,
			Name: "fake-active"},
		{ID: 2,
			Name: "fake-premium"},
	},
}

func (s *StubMemberStore) GetMembers() []models.Member {
	return memberMapToSlice(s.members)
}

func (s *StubMemberStore) GetTiers() []models.Tier {
	return s.tiers
}

func (s *StubMemberStore) GetMemberByEmail(email string) (models.Member, error) {
	if val, ok := s.members[email]; !ok {
		return val, errors.New(email)
	}
	return s.members[email], nil
}

func (s *StubMemberStore) AssignRFID(email string, rfid string) (models.Member, error) {
	if val, ok := s.members[email]; !ok {
		return val, errors.New("not a valid member email")
	}
	member := s.members[email]

	member.RFID = rfid

	s.members[email] = member
	return member, nil
}

func (s *StubMemberStore) AddNewMember(newMember models.Member) (models.Member, error) {
	s.members[newMember.Email] = newMember
	return s.members[newMember.Email], nil
}

func memberMapToSlice(m map[string]models.Member) []models.Member {
	var members []models.Member

	for _, member := range m {
		members = append(members, member)
	}
	return members
}
