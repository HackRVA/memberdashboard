# Disaster Recovery

## If this system goes down, How do we recover?
We take periodic snapshots and backups of the entire DB.  It's probably preferable to restore the entire DB.  
However, if that's not an option, read on.

Most of the things in the database can be calculated or retrieved.

The bare minimum you need to restore the system is returned from the following query.

```SQL
SELECT name, email, rfid, subscription_id
	FROM membership.members;
```

### Member Status
Member Status is evaluated daily.  The `subscription_id` is used to lookup the member's subscription in the payment provider.

Occasionally, we have members that have `Credited` memberships.  This can happen if someone wants to pay for only a month of membership or if someone gifts a few months of membership to someone.
With the bare minimum restore, `Credited` members would not persisted their `Credited` status.  This should be tracked out of band from this application (i.e. the treasurer keeps track of who has this and when they need to be removed from `Credited` status).

### RFID fobs
In a pinch, you could get a list of membership details (name, email, subscription_id) from the treasurer and then get a list of the active RFIDs off of the the rfid device (it store them on local flash).

### Worst case scenario
If you've lost everything, you'll just have to get people to register their fobs again.

The dashboard does have an option for self service.  A member can create an account and assign themselves a new fob.  Normal members don't have access to change other members fobs and they don't have access to assign new access to themselves.

