#!/bin/sh

# generate documentation
sh gendocs.sh

docker-compose up --build
