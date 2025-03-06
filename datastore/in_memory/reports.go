package in_memory

import (
	"time"

	"github.com/HackRVA/memberserver/models"
)

func (i *In_memory) UpdateMemberCounts() {}
func (i *In_memory) GetMemberCounts() ([]models.MemberCount, error) {
	return []models.MemberCount{}, nil
}

func (i *In_memory) GetMemberCountByMonth(month time.Time) (models.MemberCount, error) {
	return models.MemberCount{}, nil
}

func (i *In_memory) GetAccessStats(date time.Time, resourceName string) ([]models.AccessStats, error) {
	return []models.AccessStats{}, nil
}

func (i *In_memory) GetMemberChurn() (int, error) {
	return 0, nil
}
