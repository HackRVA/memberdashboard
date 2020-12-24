package database

import (
	"context"
	"fmt"
	"log"
)

const getMemberQuery = `SELECT id, name, email, rfid, rfid_125, member_tier_id FROM membership.members;`
const getMemberByNameQuery = `SELECT id, name, email, rfid, rfid_125, member_tier_id
FROM membership.members
WHERE name = $1;`

// Member -- a member of the makerspace
type Member struct {
	ID      uint8  `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	RFID    string `json:"rfid"`
	RFID125 string `json:"rfid_125"`
	Level   uint8  `json:"memberLevel"`
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
		err = rows.Scan(&m.ID, &m.Name, &m.Email, &m.RFID, &m.RFID125, &m.Level)
		members = append(members, m)
	}

	return members
}

// GetMemberByName - lookup a resource by it's name
func (db *Database) GetMemberByName(memberName string) (Member, error) {
	var m Member
	rows, err := db.pool.Query(context.Background(), getMemberByNameQuery, memberName)
	if err != nil {
		return m, fmt.Errorf("conn.Query failed: %v", err)
	}

	defer rows.Close()

	err = rows.Scan(&m.ID, &m.Name, &m.Email, &m.Level, &m.RFID, &m.RFID125)
	if err != nil {
		return m, fmt.Errorf("getResourceByName failed: %s", err)
	}

	return m, err
}
