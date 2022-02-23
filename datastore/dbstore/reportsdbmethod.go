package dbstore

// ReportsDatabaseMethod -- method container that holds the extension methods to query the members, credit, and resource tables
type ReportsDatabaseMethod struct{}

func (member *ReportsDatabaseMethod) updateMemberCounts() string {
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

func (member *ReportsDatabaseMethod) getMemberCounts() string {
	return `SELECT 
	month,
	classic,
	standard,
	premium,
	credited
	 FROM membership.member_counts;`
}

func (member *ReportsDatabaseMethod) getMemberCountByMonth() string {
	return `SELECT
	month,
	classic,
	standard,
	premium,
	credited
	FROM membership.member_counts
		WHERE month = date_trunc('month', '$1');;`
}
