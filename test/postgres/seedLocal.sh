#!/bin/sh

#######################################################################################################
#
#   Note: this script deletes all contents of the db and reseeds with new test data.
#         This should probably only be used for testing/development
#
#######################################################################################################

export PGUSER=test
export PGPASSWORD=test
export PGHOST=postgres
export PGDB=membership

psql -U $PGUSER -h $PGHOST -d $PGDB -f drop.sql
psql -U $PGUSER -h $PGHOST -d $PGDB -f postgres.sql
psql -U $PGUSER -h $PGHOST -d $PGDB -c "\copy  membership.member_tiers FROM './seedData/tiers.csv' DELIMITER ',' CSV HEADER;"
psql -U $PGUSER -h $PGHOST -d $PGDB -c "\copy  membership.members FROM './seedData/members.csv' DELIMITER ',' CSV HEADER;"
psql -U $PGUSER -h $PGHOST -d $PGDB -c "\copy  membership.resources FROM './seedData/resources.csv' DELIMITER ',' CSV HEADER;"
psql -U $PGUSER -h $PGHOST -d $PGDB -c "\copy  membership.member_resource FROM './seedData/member_to_resource_relations.csv' DELIMITER ',' CSV HEADER;"
psql -U $PGUSER -h $PGHOST -d $PGDB -c "\copy  membership.member_credit FROM './seedData/member_credit.csv' DELIMITER ',' CSV HEADER;"
