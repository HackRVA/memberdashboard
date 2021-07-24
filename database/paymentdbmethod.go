package database

// PaymentDatabaseMethod -- method container that holds the extension methods to query the payments table
type PaymentDatabaseMethod struct{}

func (payment *PaymentDatabaseMethod) getPayments() string {
	const getPaymentsQuery = `
	SELECT id, date, amount
	FROM membership.payments
	ORDER BY date;`

	return getPaymentsQuery
}

func (payment *PaymentDatabaseMethod) insertPayment() string {
	const insertPaymentQuery = `
	INSERT INTO membership.payments(
	date, amount, member_id)
	VALUES ($1, $2, $3)
	RETURNING *;`

	return insertPaymentQuery
}

func (payment *PaymentDatabaseMethod) updateMembershipLevel() string {
	const updateMembershipLevelQuery = `
	UPDATE membership.members
	SET member_tier_id=$2
	WHERE id=$1
	RETURNING *;`

	return updateMembershipLevelQuery
}

func (payment *PaymentDatabaseMethod) pastDuePayments() string {
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

func (payment *PaymentDatabaseMethod) updateMemberTiers() string {
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
