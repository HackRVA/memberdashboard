package database

import (
	"context"
	"log"
)

const getMemberTiersQuery = `SELECT id, description FROM membership.member_tiers;`

// MemberLevel enum
type MemberLevel int

const (
	// Inactive $0
	Inactive MemberLevel = iota + 1
	// Student $15
	Student
	// Classic $30
	Classic
	// Standard $35
	Standard
	// Premium $50
	Premium
)

// MemberLevelFromAmount convert amount to MemberLevel
var MemberLevelFromAmount = map[int64]MemberLevel{
	0:  Inactive,
	15: Student,
	30: Classic,
	35: Standard,
	50: Premium,
}

// MemberLevelToStr convert MemberLevel to string
var MemberLevelToStr = map[MemberLevel]string{
	Inactive: "Inactive",
	Student:  "Student",
	Classic:  "Classic",
	Standard: "Standard",
	Premium:  "Premium",
}

// Tier - level of membership
type Tier struct {
	ID   uint8  `json:"id"`
	Name string `json:"level"`
}

// GetMemberTiers - gets the member tiers from DB
func (db *Database) GetMemberTiers() []Tier {
	rows, err := db.pool.Query(context.Background(), getMemberTiersQuery)
	if err != nil {
		log.Fatalf("conn.Query failed: %v", err)
	}

	defer rows.Close()

	var tiers []Tier

	for rows.Next() {
		var t Tier
		err = rows.Scan(&t.ID, &t.Name)
		tiers = append(tiers, t)
	}

	return tiers
}
