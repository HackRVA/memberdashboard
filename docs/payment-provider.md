# Payment Provider

## Paypal
### New Member
When a member creates a subscription, Paypal sends us a message via a webhook.  There is no guaranty that this will be immediate.
The message might not include all the expected information.

Hopefully we receive their `subscription_id` and their email address.
If the information isn't correct, we can modify it in the UI.

### Evaluating Membership
Memberships are evaluated when the server starts up and every day at the same time.

