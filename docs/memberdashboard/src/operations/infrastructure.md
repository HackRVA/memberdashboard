# Infrastructure
```
    cloud
{ database } <--> [firewall] <-[proxy]-> [server] <--> [mqtt broker] <--> [...resources/rfid devices]
```

* The db runs on a hosting provider in the cloud. 
* The server is internal.  The server could also run in the cloud, but that would require us to expose the MQTT broker through the firewall.
* The web ui (which runs on the server) is exposed through a proxy.
* The server communicates to the rfid endpoints (aka resources) through the MQTT protocol.

## Database
The database runs postgres.


## Server
Responsibilities:
* handle backend requests for the "Member Dashboard" (aka the web ui)
* serve up the web ui
* publish/subscribe MQTT messages
* run tasks at scheduled intervals (e.g. reach out to payment providers to verify that members are up to date)

## Web UI
The web ui is a dashboard that makes managing the space more convenient.  It controls who has access to what and shows some membership stats.

## MQTT Broker
The MQTT Broker allows the server to communicate to resources on the network (e.g. the frontdoor rfid device).


## Resources
Resources is a term used for RFID devices.  However, they don't have to be an rfid device.  Essentially, a resource is just something that we can communicate with via MQTT.

For example, you could potentially build out an LED light controller that publishes and subscribes to MQTT events and then build some module on this member dashboard to control it.
