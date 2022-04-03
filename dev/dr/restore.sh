#!/bin/bash

# you can run this script to restore, but you will need to update the 
#   BACKUP_FILE variable to match the path that you want to restore
export BACKUP_FILE=/path/to/your/file/dump_name.tar

pg_restore -d ${POSTGRES_DB} ${BACKUP_FILE} -c -U ${POSTGRES_USER}
