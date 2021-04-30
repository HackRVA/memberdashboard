package database

import (
	"context"

	log "github.com/sirupsen/logrus"
)

var tierDbMethod TierDatabaseMethod

// MemberLevel enum
type MemberLevel int

const (
	// Inactive $0
	Inactive MemberLevel = iota + 1
	// Credited $1
	Credited
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
	1:  Credited,
	30: Classic,
	35: Standard,
	50: Premium,
}

// MemberLevelToStr convert MemberLevel to string
var MemberLevelToStr = map[MemberLevel]string{
	Inactive: "Inactive",
	Credited: "Credited",
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
	rows, err := db.getConn().Query(context.Background(), tierDbMethod.getMemberTiers())
	if err != nil {
		log.Errorf("conn.Query failed: %v", err)
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
