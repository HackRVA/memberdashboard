#!/bin/sh

# generate documentation
sh gendocs.sh

MQTT_BROKER_ADDRESS="tcp://mosquitto:1883" docker-compose -f docker-compose.yml up --build
