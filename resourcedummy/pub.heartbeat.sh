#!/bin/bash

mqtt pub -t frontdoor/sync -h localhost -p 1883 --message '{"type":"heartbeat","time":1616731044,"ip":"192.168.1.211","door":"esp-rfid"}'
