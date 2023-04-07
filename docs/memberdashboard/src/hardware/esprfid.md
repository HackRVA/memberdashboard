# esprfid
[esprfid github](https://github.com/esprfid/esp-rfid)

[buy](https://www.tindie.com/products/nardev/esp-rfid-relay-blue-board/)


The "esprfid" blue board works well for our purposes for the most part.  

Some issues:
* after several months of operation, we've experienced memory corruption (this could have to do with how frequently we update it).
* we currently don't have a verification step between the rfid reader and the member server
* in the event of a power outage, we seem to lose connection (or at least we don't receive mqtt messages from the rfid device).  This could be related to the rfid device coming online before the mqtt broker is available.  Usually power cycling the rfid device gets restored to a state that can send/receive mqtt messages.