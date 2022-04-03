# Resource Bridge

The Resource Bridge is intended to run onsite with limited compute power.
This device will proxy requests from the local `resources` to do things like update access lists on the local `resources` and push updates (e.g. logs) from the `resources` to the servers if need be.

This will allow the membership server to exist in the cloud somewhere and safe from the dangers of local mishaps.

## Listen for broadcasts

The resource bridge listens for broadcasts on the network.
Unprovisioned `resources` will broadcast to the network to try to be found.

When the resource bridge sees a broadcast we will communicate to the server that a resource is available to be provisioned.
From the server, we can label this resource


# MQTT User Management
## Generate a Password file
Create a password file and add a user:
```
mosquitto_passwd -c /etc/hackrva/mqtt.passwordfile user
```

## Add Additional users to the password file
Add a user to an existing password file:
```
mosquitto_passwd -b /etc/hackrva/mqtt.passwordfile user password
```

 ## Point to your password file

You must set an environment variable to tell the docker container where you password file is
```
export MOSQUITTO_PASSWORD_FILE=/etc/hackrva/mqtt.passwordfile
```

## Delete User

```
mosquitto_passwd -D /etc/hackrva/mqtt.passwordfile user
```
