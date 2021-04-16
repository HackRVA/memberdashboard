package database

import (
	"context"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
)

const getMemberQuery = `SELECT id, name, email, COALESCE(rfid,'notset'), member_tier_id,
ARRAY(
SELECT resource_id
FROM membership.member_resource
LEFT JOIN membership.resources 
ON membership.resources.id = membership.member_resource.resource_id
WHERE member_id = membership.members.id
) as resources
FROM membership.members
ORDER BY name;
`
const getMembersWithCreditQuery = `SELECT id
FROM membership.member_credit
LEFT JOIN membership.members
ON member_id = membership.members.id
ORDER BY name;
`

const getMemberByEmailQuery = `SELECT id, name, email, COALESCE(rfid,'notset'), member_tier_id,
ARRAY(
SELECT resource_id
FROM membership.member_resource
LEFT JOIN membership.resources 
ON membership.resources.id = membership.member_resource.resource_id
WHERE member_id = membership.members.id
) as resources
FROM membership.members
WHERE email = $1;`

const getMemberByIDQuery = `SELECT id, name, email, COALESCE(rfid,'notset'), member_tier_id,
ARRAY(
SELECT resource_id
FROM membership.member_resource
LEFT JOIN membership.resources 
ON membership.resources.id = membership.member_resource.resource_id
WHERE member_id = membership.members.id
) as resources
FROM membership.members
WHERE id = $1;`

const setMemberRFIDTag = `UPDATE membership.members
SET rfid=$2
WHERE email=$1
RETURNING rfid;`

const insertMemberQuery = `INSERT INTO membership.members(
	name, email, rfid, member_tier_id)
	VALUES ($1, $2, null, 1)
RETURNING id, name, email;`

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
	rows, err := db.getConn().Query(db.ctx, getMemberQuery)
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
	rows, err := db.getConn().Query(db.ctx, getMembersWithCreditQuery)
	if err != nil {
		log.Errorf("error getting credited members: %v", err)
	}

	defer rows.Close()

	var members []Member

	for rows.Next() {
		var m Member
		err = rows.Scan(&m.ID)
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

	err := db.getConn().QueryRow(context.Background(), getMemberByEmailQuery, memberEmail).Scan(&m.ID, &m.Name, &m.Email, &m.RFID, &m.Level, &rIDs)
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

// GetMemberByID - lookup a member by their memberID
func (db *Database) GetMemberByID(memberID string) (Member, error) {
	var m Member
	var rIDs []string

	err := db.getConn().QueryRow(context.Background(), getMemberByIDQuery, memberID).Scan(&m.ID, &m.Name, &m.Email, &m.RFID, &m.Level, &rIDs)
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

	err = db.getConn().QueryRow(context.Background(), setMemberRFIDTag, email, encodeRFID(RFIDTag)).Scan(&m.RFID)
	if err != nil {
		return m, fmt.Errorf("conn.Query failed: %v", err)
	}

	return m, err
}

// AddMember adds a member to the database
func (db *Database) AddMember(email string, name string) (Member, error) {
	var m Member

	err := db.getConn().QueryRow(context.Background(), insertMemberQuery, name, email).Scan(&m.ID, &m.Name, &m.Email)
	if err != nil {
		return m, fmt.Errorf("conn.Query failed: %v", err)
	}

	return m, err
}

// AddMember adds multiple members to the database
func (db *Database) AddMembers(members []Member) error {
	sqlStr := `INSERT INTO membership.members(
name, email, member_tier_id)
VALUES `

	var valStr []string
	for _, m := range members {
		// postgres doesn't like apostrophes
		memberName := strings.Replace(m.Name, "'", "''", -1)
		valStr = append(valStr, fmt.Sprintf("('%s', '%s', %d)", memberName, m.Email, 1))
	}

	str := strings.Join(valStr, ",")

	_, err := db.getConn().Query(context.Background(), sqlStr+str+" ON CONFLICT DO NOTHING;")
	if err != nil {
		return fmt.Errorf("add members query failed: %v", err)
	}

	return err
}
