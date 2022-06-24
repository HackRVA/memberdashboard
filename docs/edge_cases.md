# Edge Cases
> note: Be extra careful when making direct changes to the DB.  This is typically not recommended.  We humans make errors.

## An existing member starts a subscription for someone else.
We received the subscriptionID from paypal, but we already have that paying member's email in the DB.
The DB has a constraint that prevents it from using the same email address for two different members.
So, the system created a new "member" with only the subscriptionID populated.
(i.e. no name, no email)
i.e.
| id | name | email | rfid | member_tier_id | subscription_id |
| -- | -- | -- | -- | -- | -- |
| <some_id> | null | null | null | inactive | <some_subscription_id> |

I created a new member for them in the dashboard.  This allowed me to input their name and email address and assign them a fob that worked immediately.

| id | name | email | rfid | member_tier_id | subscription_id |
| -- | -- | -- | -- | -- | -- |
| <some_id> | <some_name> | <some_email> | <some_rfid> | <some_membership_level> | null |

On the following day, the member would lose access because they don't have a `subscription_id`

This can't be fixed in the UI because the UI doesn't allow you to delete members and you aren't allowed to give two members the same `subscription_id`.

The only way to fix this is to run a few DB queries.

1. Get the new member's `subscription_id`.  If you don't know how to get this, you may have to contact the treasurer.
2. verify that we don't already have a member with that `subscription_id` (there's a constraint that prevents multiple members having the same `subscription_id`)
```SQL
SELECT id, name, email, rfid, member_tier_id, subscription_id
	FROM membership.members
	WHERE subscription_id = '<new subscription id>';
```
In this case, I found that the new `subscription_id` overwrote the original member's existing `subscription_id`.
So, I will have to go find their existing `subscription_id` and correct this.

```SQL
UPDATE membership.members
	SET subscription_id='<original subscription id>'
	WHERE subscription_id = '<new subscription id>'
```

corrected.

Now that we are sure that no other member has that `subscription_id`, we can update their `subscription_id` via the UI.
