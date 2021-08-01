#!/bin/bash

mqtt pub -t frontdoor/send -h localhost -p 1883 --message '{"cmd":"log","type":"access","time":'$(date +%s)',"isKnown":"true","access":"Always","username":"Fake User","uid":"not an rfid tag","door":"frontdoor"}'
