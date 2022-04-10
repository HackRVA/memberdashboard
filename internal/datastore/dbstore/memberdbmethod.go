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
	return `UPDATE membership.members SET name=$1, subscription_id=$2 WHERE email=$3;`
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

func (member *MemberDatabaseMethod) getPayments() string {
	const getPaymentsQuery = `
	SELECT id, date, amount
	FROM membership.payments
	ORDER BY date;`

	return getPaymentsQuery
}

func (member *MemberDatabaseMethod) insertPayment() string {
	const insertPaymentQuery = `
	INSERT INTO membership.payments(
	date, amount, member_id)
	VALUES ($1, $2, $3)
	RETURNING *;`

	return insertPaymentQuery
}

func (member *MemberDatabaseMethod) updateMembershipLevel() string {
	const updateMembershipLevelQuery = `
	UPDATE membership.members
	SET member_tier_id=$2
	WHERE id=$1
	RETURNING *;`

	return updateMembershipLevelQuery
}

func (member *MemberDatabaseMethod) pastDuePayments() string {
	const sql = `
	SELECT m.id, m.name, m.email, COALESCE(max(p.date), '0001-01-01') as lastPaymentDate,
		current_date - COALESCE(max(p.date), '0001-01-01') as daysSinceLastPayment
	FROM membership.members m
	INNER JOIN membership.member_tiers t
	on m.member_tier_id = t.id
	LEFT JOIN membership.payments p
	on m.id = p.member_id
	WHERE t.description not in ('Inactive', 'Credited')
	GROUP BY m.id, m.name, m.email
	HAVING MAX(p.date) is null or MAX(p.date) < current_date - interval '1 month';`
	return sql
}

func (member *MemberDatabaseMethod) updateMemberTiers() string {
	const sql = `
	with cte as (
		SELECT m.id as MemberId, p.amount,
			ROW_NUMBER() over (
				Partition By m.id
				order by p.date DESC
			) row_num
		FROM membership.members m
		INNER JOIN membership.payments p
		ON m.id = p.member_id
			AND p.amount > 0
		WHERE p.date > current_date - interval '1 month'
	)
	UPDATE membership.members m
	SET member_tier_id = t.id
	FROM cte c
	INNER JOIN membership.member_tiers t
	ON c.amount = t.price
	WHERE c.memberid = m.id
		AND c.row_num = 1
		AND m.member_tier_id != t.id;
	`
	return sql
}
