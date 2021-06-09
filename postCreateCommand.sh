#!/bin/sh

# PostCreateCommand.sh is intended to run the first time the dev container is created.
# Perform initial ui dependency installaion and apply all migrations to database.

npm --prefix /workspace/ui ci 
wait-for-it.sh postgres:5432 -- make migrate-up