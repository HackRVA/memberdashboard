# Membership Server

The membership server is the source of truth for who has access to what at the makerspace.
Membership statuses will be pulled down from Payment Providers on a daily basis.
If a member has made a payment in the past 30 days, they will be considered an active member.

## High level
 * The server pulls payment information from paypal (and stores in the db) so we can tell who is currently an active member
 * the server will maintain access lists and periodically push those access lists to the microcontrollers on the network
 * The microcontroller (aka a resource) stores its access list locally so it's not dependant on the network when someone wants to access the space

## Install prereqs

you need to install at least:
docker-compose
go
go-swagger
typescript
npm

Then follow README from the ui directory to instal the npm modules

maybe do this one separately
```
npm install --global rollup
```

## Start the app

This project uses docker.
A testing environment can be spun up by running the build script.

```
sh buildandrun.sh
```

### Seed the DB with test data

Create a membership database and grant rights to user test

```
cd test/postgres
sh seedLocal.sh
```

### CONFIG file

This app requires a config file.
The path for the config file is set using the `MEMBER_SERVER_CONFIG_FILE` environment variable.

the `sample.config.json` can be used as a template for this file.

```
export MEMBER_SERVER_CONFIG_FILE="/etc/hackrva/config.json"
```

### Generating Swagger Docs
Follow instructions in [./docs/README.md](./docs/README.md)
