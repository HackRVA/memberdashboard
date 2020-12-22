# HackRVA RFID Reader

HackRVA RFID reader for access control on doors and tools.

See [developer setup](./docs/SETUP.md) for how to get started.

See [board pinout](./docs/board_pinout.md) for information on how to setup a prototype board.


* An endpoint exists that will accept an access list via an http request
* An access list will be saved into flash
* When a user swipes their badge, it will determine if their ID is in the access list
* an endpoint exists that returns a hash of the existing access list.
    * The [member server](https://github.com/Ranthalion/rfidLock) will use this hash to determine if the access list is up to date.
* The [member server](https://github.com/Ranthalion/rfidLock) will push updates to the access list on this device as needed and keep track of the device status.

This device should continue to function even if it does not have network access or access to the member server.


## Updating the readers Access Control List

`POST REQUEST: http://<rfid-reader-ip-address>/update`
```json
{
    "acl": [ 2755459513, 848615840 ]
}
```

> note: this will replace the existing access control list