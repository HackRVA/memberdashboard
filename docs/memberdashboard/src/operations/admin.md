# Admin Access
Admin access is the ability to assign rfid fobs to members and change what they have access to.

## Granting access
To give someone access, you assign them to the `admin` resource.

The `admin` resource is a special resource that doesn't broadcast MQTT messages.

If a member has the `admin` resource, they will see `reports`, `members`, and `resources` tabs in the UI.
