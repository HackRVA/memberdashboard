#!/bin/sh

# point to a config file on the system
# this directory is mapped as a volume
export MEMBER_SERVER_CONFIG_FILE='/etc/hackrva/config.json'

docker-compose up -d --build
