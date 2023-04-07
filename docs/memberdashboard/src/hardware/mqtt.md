# mqtt examples 

|Mqtt Command|Topic|Host|Port|Message|
|-|-|-|-|-|
|sub|frontdoor/sync|localhost|1883||
|sub|frontdoor/send|localhost|1883||
|sub|frontdoor|localhost|1883||
|pub|frontdoor|localhost|1883|{"doorip": "192.168.1.211", "cmd": "listusr"}|
|pub|frontdoor/sync|localhost|1883|{"type":"heartbeat","time":1616731044,"ip":"192.168.1.211","door":"esp-rfid"}|
|pub|frontdoor/send|localhost|1883|{"cmd":"log","type":"access","time":'$(date +%s)',"isKnown":"true","access":"Always","username":"Fake User","uid":"not an rfid tag","door":"frontdoor"}|
|pub|frontdoor|localhost|1883|{"doorip": "192.168.1.211", "cmd": "deletusers"}|
|pub|frontdoor|localhost|1883|{"doorip": "192.168.1.211", "cmd": "adduser", "user": "dustin", "uid": "4755ca35", "acctype":1,"validuntil":-86400}|
