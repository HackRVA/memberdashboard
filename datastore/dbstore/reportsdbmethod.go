package dbstore

import (
	"fmt"
	"time"
)

// ReportsDatabaseMethod -- method container that holds the extension methods to query the members, credit, and resource tables
type ReportsDatabaseMethod struct{}

func (report *ReportsDatabaseMethod) updateMemberCounts() string {
	return `INSERT INTO membership.member_counts(
		month, credited, classic, standard, premium)
		SELECT 
		date_trunc('month', NOW()) as month,
		SUM(case when member_tiers.id = 2 then 1 else 0 end) credited,
		SUM(case when member_tiers.id = 3 then 1 else 0 end) classic,
		SUM(case when member_tiers.id = 4 then 1 else 0 end) standard,
		SUM(case when member_tiers.id = 5 then 1 else 0 end) premium
		FROM membership.members
		JOIN membership.member_tiers
		ON member_tiers.id = members.member_tier_id
	ON CONFLICT (month)
	DO UPDATE SET credited = EXCLUDED.credited, classic = EXCLUDED.classic, standard = EXCLUDED.standard, premium = EXCLUDED.premium;
	`
}

func (report *ReportsDatabaseMethod) getMemberCounts() string {
	return `SELECT 
	month,
	classic,
	standard,
	premium,
	credited
	FROM membership.member_counts
	ORDER BY month;`
}

func (report *ReportsDatabaseMethod) getMemberCountByMonth() string {
	return `SELECT
	month,
	classic,
	standard,
	premium,
	credited
	FROM membership.member_counts
		WHERE month = date_trunc('month', '$1');;`
}

func (report *ReportsDatabaseMethod) getAccessStats(day time.Time, resourceName string) string {
	var dayFilter string
	var resourceFilter string

	if len(resourceName) > 0 {
		resourceFilter = fmt.Sprintf("WHERE door = '%s'", resourceName)
	}

	if !day.IsZero() {
		dayFilter = fmt.Sprintf("WHERE date_trunc('day', event_time) = date_trunc('day', %s)", day)
		if len(resourceName) > 0 {
			resourceFilter = fmt.Sprintf("AND door = '%s'", resourceName)
		}
	}

	return fmt.Sprintf(`SELECT date_trunc('day', event_time) as day, door as resource, COUNT(*)
	FROM membership.access_events
	%s
	%s
	GROUP BY date_trunc('day', event_time), door
	ORDER BY day;`, dayFilter, resourceFilter)
}
