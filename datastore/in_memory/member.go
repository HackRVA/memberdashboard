package in_memory

import (
	"errors"
	"memberserver/api/models"
	"sort"
)

var Members = map[string]models.Member{}

func (i In_memory) GetTiers() []models.Tier {
	return i.Tiers
}
func (i In_memory) GetMembers() []models.Member {
	return MemberMapToSlice(i.Members)
}
func (i In_memory) GetMemberByEmail(email string) (models.Member, error) {
	for _, k := range i.Members {
		if k.Email == email {
			return k, nil
		}
	}
	return models.Member{}, errors.New("error getting user: not found")
}
func (i In_memory) AssignRFID(email string, rfid string) (models.Member, error) {
	if len(rfid) == 0 {
		return models.Member{}, errors.New("not a valid rfid")
	}

	for _, member := range i.Members {
		if member.Email != email {
			continue
		}
		return member, nil
	}
	return models.Member{}, errors.New("user not found")
}
func (i In_memory) AddNewMember(newMember models.Member) (models.Member, error) {
	return newMember, nil
}
func (i In_memory) AddMembers(members []models.Member) error {
	return nil
}
func (i In_memory) GetMembersWithCredit() []models.Member {
	return []models.Member{}
}
func (i In_memory) ProcessMember(newMember models.Member) error {
	return nil
}

func MemberMapToSlice(m map[string]models.Member) []models.Member {
	var members []models.Member

	for _, member := range m {
		members = append(members, member)
	}

	sort.Sort(ByID(members))

	return members
}

// ByID implements sort.Interface based on the ID field.
type ByID []models.Member

func (a ByID) Len() int           { return len(a) }
func (a ByID) Less(i, j int) bool { return a[i].ID < a[j].ID }
func (a ByID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
