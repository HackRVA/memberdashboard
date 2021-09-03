package in_memory

import "memberserver/api/models"

type In_memory struct {
	Members map[string]models.Member
	Tiers   []models.Tier
}
