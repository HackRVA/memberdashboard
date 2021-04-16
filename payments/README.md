# Payments
This is a package that fetches payment information from our payment provider (i.e. Paypal).
There are a few edge cases involved which I will detail below.

HackRVA memberships are established by making a subscription to our paypal.

## Evaluate Membership
We evaluate a member's status everyday.

### Active
If a member has made a payment in the last 30 days.  They are considered an `Active Member`

### Grace Period
If a member hasn't made a payment in the last 30 days, we will offer a grace period for their membership.

We will send the member a notification, but their `member_status` won't change

### Revoked
If a member hasn't made a payment in the last 45 days, the membership will be revoked.

We will send the member an email stating that their membership has been revoked and we will update their `member_status` to `Inactive`.

## Membership Levels

| Level    | price |
|----------|-------|
| Inactive | $0    |
| Credited | $1    | 
| *Classic | $30   |
| Standard | $35   |
| Permiere | $50   |
| Donation | >$50  |


> Classic Memberships - some members are grandfathered in at the previous rate

## Gifting a membership
> TODO: This is a case we need to handle.
