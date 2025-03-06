package dbstore

var tierDbMethod TierDatabaseMethod

// TierDatabaseMethod -- method container that holds the extension methods to query the tier table
type TierDatabaseMethod struct{}

func (tier *TierDatabaseMethod) getMemberTiers() string {
	const getMemberTiersQuery = `SELECT id, description FROM membership.member_tiers;`

	return getMemberTiersQuery
}
