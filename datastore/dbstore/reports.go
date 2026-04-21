package dbstore

import (
	"context"
	"fmt"
	"time"

	"github.com/HackRVA/memberserver/models"

	log "github.com/sirupsen/logrus"
)

var reportsDbMethod ReportsDatabaseMethod

func (db *DatabaseStore) UpdateMemberCounts(ctx context.Context) {
	err := db.pool.QueryRow(ctx, reportsDbMethod.updateMemberCounts()).Scan()
	if err != nil {
		if err.Error() != "no rows in result set" {
			log.Errorf("updateMemberCounts failed: %v", err)
		}
	}
}

func (db *DatabaseStore) GetMemberCounts(ctx context.Context) ([]models.MemberCount, error) {
	var memberCounts []models.MemberCount

	rows, err := db.pool.Query(ctx, reportsDbMethod.getMemberCounts())
	if err != nil {
		log.Errorf("error getting member counts: %v", err)
		return memberCounts, err
	}

	defer rows.Close()

	for rows.Next() {
		var m models.MemberCount
		err = rows.Scan(&m.Month, &m.Classic, &m.Standard, &m.Premium, &m.Credited)
		if err != nil {
			log.Errorf("error scanning row: %s", err)
		}

		memberCounts = append(memberCounts, m)
	}

	return memberCounts, nil
}

func (db *DatabaseStore) GetMemberCountByMonth(ctx context.Context, month time.Time) (models.MemberCount, error) {
	var memberCount models.MemberCount

	err := db.pool.QueryRow(ctx, reportsDbMethod.getMemberCountByMonth(), month).Scan(&memberCount.Classic, &memberCount.Standard, &memberCount.Premium, &memberCount.Credited)
	if err != nil {
		log.Errorf("etMemberCountByMonth failed: %v", err)
	}

	return memberCount, nil
}

func (db *DatabaseStore) GetAccessStats(ctx context.Context, date time.Time, resourceName string) ([]models.AccessStats, error) {
	var stats []models.AccessStats

	rows, err := db.pool.Query(ctx, reportsDbMethod.getAccessStats(date, resourceName))
	if err != nil {
		log.Errorf("error getting member counts: %v", err)
		return stats, err
	}

	defer rows.Close()

	for rows.Next() {
		var m models.AccessStats
		err = rows.Scan(&m.Date, &m.ResourceName, &m.AccessCount)
		if err != nil {
			log.Errorf("error scanning row: %s", err)
		}

		stats = append(stats, m)
	}

	return stats, nil
}

func (db *DatabaseStore) GetMemberChurn(ctx context.Context) (int, error) {
	rows, err := db.pool.Query(ctx, reportsDbMethod.getMemberChurn())
	if err != nil {
		return -1, fmt.Errorf("error running query: %s", err)
	}

	defer rows.Close()

	counts := []struct {
		Month       time.Time
		MemberCount int
	}{}
	for rows.Next() {
		var count struct {
			Month       time.Time
			MemberCount int
		}
		err = rows.Scan(&count.Month, &count.MemberCount)
		if err != nil {
			log.Errorf("error scanning row: %s", err)
		}

		counts = append(counts, count)
	}

	return counts[0].MemberCount - counts[1].MemberCount, nil
}
