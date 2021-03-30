#!/bin/bash
mqtt pub -t frontdoor -h localhost -p 1883 --message '{"doorip": "192.168.1.211", "cmd": "deletusers"}'
