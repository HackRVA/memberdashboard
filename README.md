# Member Dashboard

Member Dashboard is the source of truth for who has access to what at the makerspace.
Membership statuses will be pulled down from Payment Providers on a daily basis.
If a member has made a payment in the past 30 days, they will be considered an active member.

## High level

- The server pulls payment information from paypal (and stores in the db) so we can tell who is currently an active member
- the server will maintain access lists and periodically push those access lists to the microcontrollers on the network
- The microcontroller (aka a resource) stores its access list locally so it's not dependant on the network when someone wants to access the space

## Install prereqs

You need to install at least:

- [docker](https://docs.docker.com/get-docker/)
- [docker-compose](https://docs.docker.com/compose/install/)
- [go](https://golang.org/doc/install)
- [node](https://nodejs.org/en/)

Follow the [README](/ui/README.md) in the ui directory to install the npm modules

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

## Database Migrations
Database migrations are managed using [golang-migrate/migrate](https://github.com/golang-migrate/migrate).  Migrations can be applied using the migrate CLI or through a docker container.

### How to install golang-migrate
```
go install github.com/golang-migrate/migrate/v4/cmd/migrate
```

### How to add a migration
A migration can be added by creating an up and down migration file in the migrations folder.  The migration file names should be named according to {sequentialNumber}_{description}.up.sql and {sequentialNumber}_{description}.down.sql.  These files can be created manaully, or by using the migrate CLI.
```
migrate create -ext sql -dir migrations -seq <description>
```
Populate the up and down scripts.  Up scripts should be idempotent.  Down scripts should revert all changes made by up script

### How to run a migration

Migrations can be applied using the CLI
```
migrate -database postgres://test:test@localhost:5432/membership?sslmode=disable -verbose -path migrations up
```
or via docker
```
docker run -v migrations:/migrations --network host migrate/migrate -path=/migrations -database postgres://test:test@localhost:5432/membership?sslmode=disable up
```
