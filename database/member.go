package database

import (
	"context"
	"fmt"
	"log"
)

const getMemberQuery = `SELECT id, name, email, rfid, member_tier_id FROM membership.members;`
const getMemberByEmailQuery = `SELECT id, name, email, rfid, member_tier_id
FROM membership.members
WHERE email = $1;`
const setMemberRFIDTag = `UPDATE membership.members
rfid=$2
WHERE email = $1;`

// Member -- a member of the makerspace
type Member struct {
	ID    uint8  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	RFID  string `json:"rfid"`
	Level uint8  `json:"memberLevel"`
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
		var m Member
		err = rows.Scan(&m.ID, &m.Name, &m.Email, &m.RFID, &m.Level)
		members = append(members, m)
	}

	return members
}

// GetMemberByEmail - lookup a member by their email address
func (db *Database) GetMemberByEmail(memberEmail string) (Member, error) {
	var m Member
	err := db.pool.QueryRow(context.Background(), getMemberByEmailQuery, memberEmail).Scan(&m.ID, &m.Name, &m.Email, &m.RFID, &m.Level)
	if err != nil {
		return m, fmt.Errorf("conn.Query failed: %v", err)
	}

	return m, err
}

// SetRFIDTag sets the rfid tag as
func (db *Database) SetRFIDTag(email string, RFIDTag string) (Member, error) {
	var m Member
	rows, err := db.pool.Query(context.Background(), setMemberRFIDTag, email, RFIDTag)
	if err != nil {
		return m, fmt.Errorf("conn.Query failed: %v", err)
	}

	defer rows.Close()

	err = rows.Scan(&m.ID, &m.Name, &m.Email, &m.Level, &m.RFID)
	if err != nil {
		return m, fmt.Errorf("SetRFIDTag failed: %s", err)
	}

	return m, err
}
