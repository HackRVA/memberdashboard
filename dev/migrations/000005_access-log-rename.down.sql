BEGIN;

ALTER TABLE membership.access_events
RENAME TO event_log;


COMMIT;
