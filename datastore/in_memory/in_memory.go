package in_memory

import (
	"sync"

	"github.com/HackRVA/memberserver/models"
	"github.com/HackRVA/memberserver/test/generators"
)

type In_memory struct {
	mu        sync.RWMutex
	Members   map[string]models.Member
	Tiers     []models.Tier
	resources map[string]models.Resource
}

func Setup() (*In_memory, error) {
	db := &In_memory{
		Members:   make(map[string]models.Member),
		resources: make(map[string]models.Resource),
	}

	generators.Seed(db, 20)
	return db, nil
}
