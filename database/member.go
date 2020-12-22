package database

import (
	"context"
	"log"
)

const getMemberQuery = `SELECT id, name, email, rfid, rfid_125, member_tier_id FROM membership.members;`

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
