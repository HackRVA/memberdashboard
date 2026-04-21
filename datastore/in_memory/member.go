package in_memory

import (
	"context"
	"errors"
	"sort"
	"strconv"

	"github.com/HackRVA/memberserver/models"
)

func (i *In_memory) allocMembers() {
	if i.Members == nil {
		i.Members = map[string]models.Member{}
	}
}

func (i *In_memory) GetTiers(ctx context.Context) []models.Tier {
	return i.Tiers
}

func (i *In_memory) GetMembersPaginated(ctx context.Context, limit int, offset int, active bool) ([]models.Member, error) {
	members := MemberMapToSlice(i.Members)
	var filteredMembers []models.Member

	for _, member := range members {
		if active && member.Level == 0 {
			continue
		}
		filteredMembers = append(filteredMembers, member)
	}

	start := offset
	if start > len(filteredMembers) {
		start = len(filteredMembers)
	}

	end := start + limit
	if end > len(filteredMembers) {
		end = len(filteredMembers)
	}

	return filteredMembers[start:end], nil
}

func (i *In_memory) GetMembers(ctx context.Context) []models.Member {
	return MemberMapToSlice(i.Members)
}

func (i *In_memory) GetMemberByEmail(ctx context.Context, email string) (models.Member, error) {
	println("Member len: ", len(i.Members))
	for _, k := range i.Members {
		println(k.Email)
		if k.Email == email {
			return k, nil
		}
	}
	return models.Member{}, errors.New("error getting user: not found")
}

func (i *In_memory) GetMemberByRFID(ctx context.Context, rfid string) (models.Member, error) {
	return models.Member{}, nil
}

func (i *In_memory) AssignRFID(ctx context.Context, email string, rfid string) (models.Member, error) {
	if len(rfid) == 0 {
		return models.Member{}, errors.New("not a valid rfid")
	}

	for _, member := range i.Members {
		if member.Email != email {
			continue
		}
		member.RFID = rfid
		i.Members[email] = member
		return member, nil
	}
	return models.Member{}, errors.New("user not found")
}

func (i *In_memory) UpdateMember(ctx context.Context, update models.Member) error {
	if len(update.Name) == 0 {
		return errors.New("fullname is required")
	}

	if len(update.Email) == 0 {
		return errors.New("email is required")
	}

	if _, ok := i.Members[update.Email]; !ok {
		return errors.New("not found")
	}

	i.Members[update.Email] = update
	return nil
}

func (i *In_memory) UpdateMemberByID(ctx context.Context, memberID string, update models.Member) error {
	for key, m := range i.Members {
		if m.ID != memberID {
			continue
		}

		m.Email = update.Email
		m.Name = update.Name
		m.SubscriptionID = update.SubscriptionID
		delete(i.Members, key)
		i.Members[m.Email] = m
		return nil
	}
	return errors.New("unable to update member info")
}

func (i *In_memory) UpdateMemberBySubscriptionID(ctx context.Context, subscriptionID string, update models.Member) error {
	for _, m := range i.Members {
		if m.SubscriptionID != subscriptionID {
			continue
		}

		m.Email = update.Email
		m.Name = update.Name
		i.Members[m.Email] = m
		return nil
	}
	return errors.New("unable to update member info")
}

func (i *In_memory) AddNewMember(ctx context.Context, newMember models.Member) (models.Member, error) {
	i.allocMembers()
	if newMember.ID == "" {
		newMember.ID = strconv.Itoa(len(i.Members) + 1)
	}
	i.Members[newMember.Email] = newMember
	return newMember, nil
}

func (i *In_memory) AddMembers(ctx context.Context, members []models.Member) error {
	i.allocMembers()
	for _, m := range members {
		if m.ID == "" {
			m.ID = strconv.Itoa(len(i.Members) + 1)
		}
		i.Members[m.Email] = m
	}
	return nil
}

func (i *In_memory) GetMembersWithCredit(ctx context.Context) []models.Member {
	return []models.Member{}
}

func (i *In_memory) ProcessMember(ctx context.Context, newMember models.Member) error {
	return nil
}

func (i *In_memory) SetMemberLevel(ctx context.Context, memberId string, level models.MemberLevel) error {
	for _, member := range i.Members {
		if member.ID != memberId {
			continue
		}

		update := member
		update.Level = uint8(level)

		delete(i.Members, update.Email)

		i.Members[update.Email] = update
	}
	return nil
}
func (i *In_memory) ApplyMemberCredits(ctx context.Context) {}
func (i *In_memory) UpdateMemberTiers(ctx context.Context)  {}

func (i *In_memory) GetActiveMembersWithoutSubscription(ctx context.Context) []models.Member {
	return nil
}

func (i *In_memory) GetMemberCount(ctx context.Context, isActive bool) (int, error) {
	count := 0
	for _, member := range i.Members {
		if isActive && member.Level != 1 {
			count++
		} else if !isActive && member.Level == 1 {
			count++
		}
	}
	return count, nil
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
