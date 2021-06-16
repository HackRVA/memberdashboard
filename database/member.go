package database

import (
	"context"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
)

var memberDbMethod MemberDatabaseMethod

// MemberResource a resource that a member belongs to
type MemberResource struct {
	ResourceID string `json:"resourceID"`
	Name       string `json:"name"`
}

// Member -- a member of the makerspace
type Member struct {
	ID        string           `json:"id"`
	Name      string           `json:"name"`
	Email     string           `json:"email"`
	RFID      string           `json:"rfid"`
	Level     uint8            `json:"memberLevel"`
	Resources []MemberResource `json:"resources"`
}

// AssignRFIDRequest -- request to associate an rfid to a member
type AssignRFIDRequest struct {
	Email string `json:"email"`
	RFID  string `json:"rfid"`
}

// GetMembers - gets the status from DB
func (db *Database) GetMembers() []Member {
	rows, err := db.getConn().Query(db.ctx, memberDbMethod.getMember())
	if err != nil {
		log.Errorf("conn.Query failed: %v", err)
	}

	defer rows.Close()

	var members []Member

	resourceMemo := make(map[string]MemberResource)

	for rows.Next() {
		var rIDs []string
		var m Member
		err = rows.Scan(&m.ID, &m.Name, &m.Email, &m.RFID, &m.Level, &rIDs)
		if err != nil {
			log.Errorf("error scanning row: %s", err)
		}

		// having issues with unmarshalling a jsonb object array from pgx
		// using a less efficient approach for now
		// TODO: fix this on the query level
		for _, rID := range rIDs {
			if _, exist := resourceMemo[rID]; exist {
				m.Resources = append(m.Resources, MemberResource{ResourceID: rID, Name: resourceMemo[rID].Name})
				continue
			}

			resource, err := db.GetResourceByID(rID)
			if err != nil {
				log.Debugf("error getting resource by id in memberResource lookup: %s %s_\n", err.Error(), rID)
				continue
			}

			resourceMemo[rID] = MemberResource{
				ResourceID: resource.ID,
				Name:       resource.Name,
			}

			m.Resources = append(m.Resources, MemberResource{ResourceID: rID, Name: resource.Name})
		}

		members = append(members, m)
	}

	return members
}

// GetMembersWithCredit - gets members that have been credited a membership
//  if a member exists in the member_credits table
//  they are credited a membership
func (db *Database) GetMembersWithCredit() []Member {
	rows, err := db.getConn().Query(db.ctx, memberDbMethod.getMembersWithCredit())
	if err != nil {
		log.Errorf("error getting credited members: %v", err)
	}

	defer rows.Close()

	var members []Member

	for rows.Next() {
		var m Member
		err = rows.Scan(&m.ID, &m.Name, &m.Email, &m.RFID, &m.Level)
		if err != nil {
			log.Errorf("error scanning row: %s", err)
		}

		members = append(members, m)
	}

	return members
}

// GetMemberByEmail - lookup a member by their email address
func (db *Database) GetMemberByEmail(memberEmail string) (Member, error) {
	var m Member
	var rIDs []string

	err := db.getConn().QueryRow(context.Background(), memberDbMethod.getMemberByEmail(), memberEmail).Scan(&m.ID, &m.Name, &m.Email, &m.RFID, &m.Level, &rIDs)
	if err != nil {
		log.Errorf("error getting member by email: %v", memberEmail)
		return m, fmt.Errorf("conn.Query failed: %w", err)
	}

	resourceMemo := make(map[string]MemberResource)

	// having issues with unmarshalling a jsonb object array from pgx
	// using a less efficient approach for now
	// TODO: fix this on the query level
	for _, rID := range rIDs {
		if _, exist := resourceMemo[rID]; exist {
			m.Resources = append(m.Resources, MemberResource{ResourceID: rID, Name: resourceMemo[rID].Name})
			continue
		}
		resource, err := db.GetResourceByID(rID)
		if err != nil {
			log.Debugf("error getting resource by id in memberResource lookup: %s %s\n", err.Error(), rID)
		}

		resourceMemo[rID] = MemberResource{
			ResourceID: resource.ID,
			Name:       resource.Name,
		}
		m.Resources = append(m.Resources, MemberResource{ResourceID: rID, Name: resource.Name})
	}

	return m, err
}

// GetMemberByID - lookup a member by their memberID
func (db *Database) GetMemberByID(memberID string) (Member, error) {
	var m Member
	var rIDs []string

	err := db.getConn().QueryRow(context.Background(), memberDbMethod.getMemberByID(), memberID).Scan(&m.ID, &m.Name, &m.Email, &m.RFID, &m.Level, &rIDs)
	if err != nil {
		return m, fmt.Errorf("conn.Query failed: %v", err)
	}

	resourceMemo := make(map[string]MemberResource)

	// having issues with unmarshalling a jsonb object array from pgx
	// using a less efficient approach for now
	// TODO: fix this on the query level
	for _, rID := range rIDs {
		if _, exist := resourceMemo[rID]; exist {
			m.Resources = append(m.Resources, MemberResource{ResourceID: rID, Name: resourceMemo[rID].Name})
			continue
		}
		resource, err := db.GetResourceByID(rID)
		if err != nil {
			log.Debugf("error getting resource by id in memberResource lookup: %s %s\n", err.Error(), rID)
		}

		resourceMemo[rID] = MemberResource{
			ResourceID: resource.ID,
			Name:       resource.Name,
		}
		m.Resources = append(m.Resources, MemberResource{ResourceID: rID, Name: resource.Name})
	}

	return m, err
}

// SetRFIDTag sets the rfid tag as
func (db *Database) SetRFIDTag(email string, RFIDTag string) (Member, error) {
	m, err := db.GetMemberByEmail(email)
	if err != nil {
		log.Errorf("error retrieving a member with that email address %s", err.Error())
		return m, err
	}

	err = db.getConn().QueryRow(context.Background(), memberDbMethod.setMemberRFIDTag(), email, encodeRFID(RFIDTag)).Scan(&m.RFID)
	if err != nil {
		return m, fmt.Errorf("conn.Query failed: %v", err)
	}

	return m, err
}

// AddMembers adds multiple members to the database
func (db *Database) AddMembers(members []Member) error {
	sqlStr := `INSERT INTO membership.members(
name, email, member_tier_id)
VALUES `

	var valStr []string
	for _, m := range members {
		// postgres doesn't like apostrophes
		memberName := strings.Replace(m.Name, "'", "''", -1)

		// if member level isn't set them to inactive,
		//   otherwise, use the level they already have.
		if m.Level == 0 {
			m.Level = uint8(Inactive)
		}

		valStr = append(valStr, fmt.Sprintf("('%s', '%s', %d)", memberName, m.Email, m.Level))

		db.AddUserToDefaultResources(m.Email)
	}

	str := strings.Join(valStr, ",")

	_, err := db.getConn().Query(context.Background(), sqlStr+str+" ON CONFLICT DO NOTHING;")
	if err != nil {
		return fmt.Errorf("add members query failed: %v", err)
	}

	return err
}
