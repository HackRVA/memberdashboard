package in_memory

import (
	"memberserver/internal/models"
	"memberserver/test/generators"
)

type In_memory struct {
	Members map[string]models.Member
	Tiers   []models.Tier
}

func Setup() (*In_memory, error) {
	db := &In_memory{
		Members: make(map[string]models.Member),
	}

	generators.Seed(db, 20)
	return db, nil
}
