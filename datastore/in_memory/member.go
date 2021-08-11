package in_memory

import (
	"errors"
	"memberserver/api/models"
)

var Members = map[string]models.Member{}

func (i In_memory) GetTiers() []models.Tier {
	return []models.Tier{}
}
func (i In_memory) GetMembers() []models.Member {
	return []models.Member{}
}
func (i In_memory) GetMemberByEmail(email string) (models.Member, error) {
	for _, k := range testUsers {
		if k.Email == email {
			return k, nil
		}
	}
	return models.Member{}, errors.New("error getting user: not found")
}
func (i In_memory) AssignRFID(email string, rfid string) (models.Member, error) {
	return models.Member{}, nil
}
func (i In_memory) AddNewMember(newMember models.Member) (models.Member, error) {
	return models.Member{}, nil
}
func (i In_memory) AddMembers(members []models.Member) error {
	return nil
}
func (i In_memory) GetMembersWithCredit() []models.Member {
	return []models.Member{}
}
