package database

import (
	"context"
	"fmt"

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
	rows, err := db.pool.Query(context.Background(), getMemberQuery)
	if err != nil {
		log.Fatalf("conn.Query failed: %v", err)
	}

	defer rows.Close()

	var members []Member

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
			resource, err := db.GetResourceByID(rID)
			if err != nil {
				log.Debugf("error getting resource by id in memberResource lookup: %s\n", err.Error())
			}
			m.Resources = append(m.Resources, MemberResource{ResourceID: rID, Name: resource.Name})
		}

		members = append(members, m)
	}

	return members
}

// GetMemberByEmail - lookup a member by their email address
func (db *Database) GetMemberByEmail(memberEmail string) (Member, error) {
	var m Member
	var rIDs []string

	err := db.pool.QueryRow(context.Background(), getMemberByEmailQuery, memberEmail).Scan(&m.ID, &m.Name, &m.Email, &m.RFID, &m.Level, &rIDs)
	if err != nil {
		return m, fmt.Errorf("conn.Query failed: %v", err)
	}

	// having issues with unmarshalling a jsonb object array from pgx
	// using a less efficient approach for now
	// TODO: fix this on the query level
	for _, rID := range rIDs {
		resource, err := db.GetResourceByID(rID)
		if err != nil {
			log.Debugf("error getting resource by id in memberResource lookup: %s\n", err.Error())
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

	err = db.pool.QueryRow(context.Background(), setMemberRFIDTag, email, RFIDTag).Scan(&m.RFID)
	if err != nil {
		return m, fmt.Errorf("conn.Query failed: %v", err)
	}

	return m, err
}

// AddMember adds a member to the database
func (db *Database) AddMember(email string, name string) (Member, error) {
	var m Member

	err := db.pool.QueryRow(context.Background(), insertMemberQuery, name, email).Scan(&m.ID, &m.Name, &m.Email)
	if err != nil {
		return m, fmt.Errorf("conn.Query failed: %v", err)
	}

	return m, err
}
