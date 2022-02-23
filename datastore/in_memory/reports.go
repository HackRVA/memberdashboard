package in_memory

import (
	"memberserver/api/models"
	"time"
)

func (i In_memory) UpdateMemberCounts() {}
func (i In_memory) GetMemberCounts() ([]models.MemberCount, error) {
	return []models.MemberCount{}, nil
}
func (i In_memory) GetMemberCountByMonth(month time.Time) (models.MemberCount, error) {
	return models.MemberCount{}, nil
}
