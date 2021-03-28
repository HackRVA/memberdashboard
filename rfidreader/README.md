
> :warning: **We are no longer using this code for the rfid reader**: We are now using [esp-rfid](https://github.com/esprfid/esp-rfid)!


# HackRVA RFID Reader

HackRVA RFID reader for access control on doors and tools.

See [developer setup](./docs/SETUP.md) for how to get started.

See [board pinout](./docs/board_pinout.md) for information on how to setup a prototype board.

- An endpoint exists that will accept an access list via an http request
- An access list will be saved into flash
- When a user swipes their badge, it will determine if their ID is in the access list
- an endpoint exists that returns a hash of the existing access list.
  - The [member server](https://github.com/Ranthalion/rfidLock) will use this hash to determine if the access list is up to date.
- The [member server](https://github.com/Ranthalion/rfidLock) will push updates to the access list on this device as needed and keep track of the device status.

This device should continue to function even if it does not have network access or access to the member server.

## Configure wifi

When you first flash the device, it doesn't know how to connect to your network.
It will come online as an access point with an ssid of `RFIDReaderConfig`. You should be able to connect to this without a password.

> note: this `RFIDReaderConfig` wifi not run after the device is setup

When you connect, find your default gateway address and send an http post request to `<default-gateway-address>/setup-wifi` with the following json body:

```json
{
  "ssid": "your wifi networks ssid",
  "password": "the password to connect to your wifi"
}
```

The device will save this to a config file in flash and will skip this setup on the next boot. If you want to reset the configuration, you will have to reset the device and reflash the firmware.

## Updating the readers Access Control List

`POST REQUEST: http://<rfid-reader-ip-address>/update`

```json
{
  "acl": ["2755459513", "848615840"]
}
```

> note: this will replace the existing access control list
