package in_memory

import (
	"context"
	"errors"
	"math"
	"time"

	"github.com/HackRVA/memberserver/models"
)

func (i *In_memory) UpdateMemberCounts(ctx context.Context) {
	i.mu.Lock()
	defer i.mu.Unlock()

	current := memberCountForMonth(time.Now(), i.Members)
	for index := range i.memberCounts {
		if sameMonth(i.memberCounts[index].Month, current.Month) {
			i.memberCounts[index] = current
			return
		}
	}

	i.memberCounts = append(i.memberCounts, current)
}

func (i *In_memory) GetMemberCounts(ctx context.Context) ([]models.MemberCount, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	counts := make([]models.MemberCount, len(i.memberCounts))
	copy(counts, i.memberCounts)
	return counts, nil
}

func (i *In_memory) GetMemberCountByMonth(ctx context.Context, month time.Time) (models.MemberCount, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	for _, count := range i.memberCounts {
		if sameMonth(count.Month, month) {
			return count, nil
		}
	}

	return models.MemberCount{}, errors.New("member count not found")
}

func (i *In_memory) GetAccessStats(ctx context.Context, date time.Time, resourceName string) ([]models.AccessStats, error) {
	return []models.AccessStats{}, nil
}

func (i *In_memory) GetMemberChurn(ctx context.Context) (int, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	if len(i.memberCounts) < 2 {
		return 0, nil
	}

	latest := activeMemberTotal(i.memberCounts[len(i.memberCounts)-1])
	previous := activeMemberTotal(i.memberCounts[len(i.memberCounts)-2])
	return latest - previous, nil
}

// demoMemberCounts creates enough history to exercise the reports UI while
// keeping the latest snapshot consistent with the demo's generated members.
func demoMemberCounts(now time.Time, members map[string]models.Member) []models.MemberCount {
	latest := memberCountForMonth(now, members)
	growth := []float64{0.72, 0.75, 0.78, 0.76, 0.81, 0.84, 0.86, 0.89, 0.91, 0.94, 0.97, 1}
	counts := make([]models.MemberCount, 0, len(growth))

	for index, factor := range growth {
		month := now.AddDate(0, index-len(growth)+1, 0)
		counts = append(counts, models.MemberCount{
			Month:    time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, time.UTC),
			Classic:  scaledCount(latest.Classic, factor),
			Standard: scaledCount(latest.Standard, factor),
			Premium:  scaledCount(latest.Premium, factor),
			Credited: scaledCount(latest.Credited, factor),
		})
	}

	return counts
}

func memberCountForMonth(month time.Time, members map[string]models.Member) models.MemberCount {
	count := models.MemberCount{
		Month: time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, time.UTC),
	}

	for _, member := range members {
		switch models.MemberLevel(member.Level) {
		case models.Classic:
			count.Classic++
		case models.Standard:
			count.Standard++
		case models.Premium:
			count.Premium++
		case models.Credited:
			count.Credited++
		}
	}

	return count
}

func scaledCount(value int, factor float64) int {
	if value == 0 {
		return 0
	}
	return max(1, int(math.Round(float64(value)*factor)))
}

func activeMemberTotal(count models.MemberCount) int {
	return count.Classic + count.Standard + count.Premium
}

func sameMonth(a, b time.Time) bool {
	return a.Year() == b.Year() && a.Month() == b.Month()
}
