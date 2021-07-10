#!/bin/sh

# generate documentation
sh gendocs.sh

GIT_COMMIT=$(git rev-parse --short HEAD) MQTT_BROKER_ADDRESS="tcp://mosquitto:1883" docker-compose -f docker-compose.yml up --build
