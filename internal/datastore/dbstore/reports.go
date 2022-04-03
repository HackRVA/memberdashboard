package dbstore

import (
	"context"
	"memberserver/internal/models"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
)

var reportsDbMethod ReportsDatabaseMethod

func (db *DatabaseStore) UpdateMemberCounts() {
	dbPool, err := pgxpool.Connect(db.ctx, db.connectionString)
	if err != nil {
		log.Printf("got error: %v\n", err)
	}
	defer dbPool.Close()

	err = dbPool.QueryRow(context.Background(), reportsDbMethod.updateMemberCounts()).Scan()
	if err != nil {
		if err.Error() != "no rows in result set" {
			log.Errorf("updateMemberCounts failed: %v", err)
		}
	}
}

func (db *DatabaseStore) GetMemberCounts() ([]models.MemberCount, error) {
	var memberCounts []models.MemberCount

	dbPool, err := pgxpool.Connect(db.ctx, db.connectionString)
	if err != nil {
		log.Printf("got error: %v\n", err)
	}
	defer dbPool.Close()

	rows, err := dbPool.Query(db.ctx, reportsDbMethod.getMemberCounts())
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

func (db *DatabaseStore) GetMemberCountByMonth(month time.Time) (models.MemberCount, error) {
	var memberCount models.MemberCount

	dbPool, err := pgxpool.Connect(db.ctx, db.connectionString)
	if err != nil {
		log.Printf("got error: %v\n", err)
	}
	defer dbPool.Close()

	err = dbPool.QueryRow(context.Background(), reportsDbMethod.getMemberCountByMonth(), month).Scan(&memberCount.Classic, &memberCount.Standard, &memberCount.Premium, &memberCount.Credited)
	if err != nil {
		log.Errorf("etMemberCountByMonth failed: %v", err)
	}

	return memberCount, nil
}
