#!/bin/bash

# add timestamp as filename
basefilename=memberserver_postgres_backup
filename=${basefilename}_$(date +%Y-%m-%d_%H-%M-%S).tar.gz

echo "perform some kind of backup here"

# dump the db and send to NAS
# pg_dump -U ${POSTGRES_USER} -W -F t POSTGRES_DB > ${BACKUP_DIR}${filename}.tar
