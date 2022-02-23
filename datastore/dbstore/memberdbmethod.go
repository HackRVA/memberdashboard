package dbstore

// MemberDatabaseMethod -- method container that holds the extension methods to query the members, credit, and resource tables
type MemberDatabaseMethod struct{}

func (member *MemberDatabaseMethod) getMember() string {
	const getMemberQuery = `SELECT id, name, email, COALESCE(rfid,'notset'), member_tier_id,
	ARRAY(
	SELECT resource_id
	FROM membership.member_resource
	LEFT JOIN membership.resources 
	ON membership.resources.id = membership.member_resource.resource_id
	WHERE member_id = membership.members.id
	) as resources, COALESCE(subscription_id,'none')
	FROM membership.members
	ORDER BY name;
	`

	return getMemberQuery
}

func (member *MemberDatabaseMethod) getMembersWithCredit() string {
	const getMembersWithCreditQuery = `SELECT id, name, email, COALESCE(rfid,'notset'), member_tier_id
	FROM membership.members
	RIGHT JOIN membership.member_credit
	ON membership.member_credit.member_id = id
	ORDER BY name;
	`

	return getMembersWithCreditQuery
}

func (member *MemberDatabaseMethod) getMemberByEmail() string {
	const getMemberByEmailQuery = `SELECT id, name, LOWER(email), COALESCE(rfid,'notset'), member_tier_id,
	ARRAY(
	SELECT resource_id
	FROM membership.member_resource
	LEFT JOIN membership.resources 
	ON membership.resources.id = membership.member_resource.resource_id
	WHERE member_id = membership.members.id
	) as resources
	FROM membership.members
	WHERE LOWER(email) = LOWER($1);`

	return getMemberByEmailQuery
}

func (member *MemberDatabaseMethod) updateMemberByEmail() string {
	const updateMemberByEmailQuery = `UPDATE membership.members SET name=$1 WHERE email=$2;`

	return updateMemberByEmailQuery
}

func (member *MemberDatabaseMethod) getMemberByRFID() string {
	const getMemberByEmailQuery = `SELECT id, name, LOWER(email), COALESCE(rfid,'notset'), member_tier_id,
	ARRAY(
	SELECT resource_id
	FROM membership.member_resource
	LEFT JOIN membership.resources 
	ON membership.resources.id = membership.member_resource.resource_id
	WHERE member_id = membership.members.id
	) as resources
	FROM membership.members
	WHERE rfid = $1;`

	return getMemberByEmailQuery
}

func (member *MemberDatabaseMethod) setMemberRFIDTag() string {
	const setMemberRFIDTagQuery = `UPDATE membership.members
	SET rfid=$2
	WHERE email=$1
	RETURNING rfid;`

	return setMemberRFIDTagQuery
}

func (member *MemberDatabaseMethod) updateMemberName() string {
	return `UPDATE membership.members
	SET name=$2
	WHERE id=$1
	RETURNING name;`
}

func (member *MemberDatabaseMethod) updateMemberSubscriptionID() string {
	return `UPDATE membership.members
	SET subscription_id=$2
	WHERE id=$1
	RETURNING name;`
}
