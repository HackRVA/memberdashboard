# ResourceDummy
Resource Dummy is intended to be a fake implementation of a resource (for testing purposes).  i.e. a device that is capable of receiving updates from the membership server -- typically this will be a rfid reader.

## run resourceDummy
```
go run resourcedummy.go
```


## MQTT CLI
For MQTT cli, I'm using [hivemq](https://hivemq.github.io/mqtt-cli/docs/quick_start.html)

After installing hivemq, you should be able to run the shell scripts in this dir to simulate pub/sub to and from a resource.

## Authenticating

https://www.hivemq.com/blog/mqtt-security-fundamentals-authentication-username-password/
