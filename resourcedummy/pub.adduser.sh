#!/bin/bash

mqtt pub -t frontdoor -h localhost -p 1883 --message '{"doorip": "192.168.1.211", "cmd": "adduser", "user": "dustin", "uid": "4755ca35", "acctype":1,"validuntil":-86400}'
