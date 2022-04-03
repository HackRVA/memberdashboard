# Disaster Recovery
This is a backup procedure to periodically backup our DB

## Backups happen automatically
Backups should run every day.
Backups will be stored on the NAS, but the host will need a few environment variables to function properly.

```
BACKUP_DIR=<some value>
POSTGRES_HOST=<some value>
POSTGRES_DB=<some value>
POSTGRES_USER=<some value>
POSTGRES_PASSWORD=<some value>
```

## Restore
```
pg_restore -d ${POSTGRES_DB} /path/to/your/file/dump_name.tar -c -U ${POSTGRES_USER}
```