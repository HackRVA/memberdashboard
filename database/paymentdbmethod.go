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

func (payment *PaymentDatabaseMethod) checkLastPayment() string {
	// checkRecentPayment - if the member doesn't have a recent payment,
	//    we will revoke their membership
	const checkLastPaymentQuery = `
	SELECT current_date - date as last_payment, amount, email
	FROM membership.payments
	LEFT JOIN membership.members
	ON membership.payments.member_id = membership.members.id
	WHERE member_id = $1
	ORDER BY date DESC
	limit 2;`

	return checkLastPaymentQuery
}

func (payment *PaymentDatabaseMethod) countPaymentsOfMemberSince() string {
	const countPaymentsOfMemberSinceQuery = `
	SELECT COUNT(*) as num_payments
	FROM membership.payments
	LEFT JOIN membership.members
	ON membership.payments.member_id = membership.members.id
	WHERE member_id = $1
	AND date >= current_date - $2;`

	return countPaymentsOfMemberSinceQuery
}

func (payment *PaymentDatabaseMethod) updateMembershipLevel() string {
	const updateMembershipLevelQuery = `
	UPDATE membership.members
	SET member_tier_id=$2
	WHERE id=$1
	RETURNING *;`

	return updateMembershipLevelQuery
}
