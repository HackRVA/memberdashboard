#!/bin/sh

# temporarily putting this here because there's a weird problem with docker building 
#  when the ui build folders don't exist
cd ui
npm run rollup
rm tsconfig.tsbuildinfo 2> /dev/null
cd ..

# generate documentation
sh gendocs.sh

# point to a config file on the system
# this directory is mapped as a volume
export MEMBER_SERVER_CONFIG_FILE='/etc/hackrva/config.json'

docker-compose up --build
