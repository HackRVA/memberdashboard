package in_memory

import (
	"context"
	"time"

	"github.com/HackRVA/memberserver/models"
)

func (i *In_memory) UpdateMemberCounts(ctx context.Context) {}
func (i *In_memory) GetMemberCounts(ctx context.Context) ([]models.MemberCount, error) {
	return []models.MemberCount{}, nil
}

func (i *In_memory) GetMemberCountByMonth(ctx context.Context, month time.Time) (models.MemberCount, error) {
	return models.MemberCount{}, nil
}

func (i *In_memory) GetAccessStats(ctx context.Context, date time.Time, resourceName string) ([]models.AccessStats, error) {
	return []models.AccessStats{}, nil
}

func (i *In_memory) GetMemberChurn(ctx context.Context) (int, error) {
	return 0, nil
}
