BEGIN;

ALTER TABLE membership.members
REMOVE COLUMN subscription_id;

COMMIT;
