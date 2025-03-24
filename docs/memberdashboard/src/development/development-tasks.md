# Development Tasks

## Database Migrations

Database migrations are managed using [golang-migrate/migrate](https://github.com/golang-migrate/migrate). The CLI is already available when using Remote-Containers, otherwise it will need to be installed.

## How to add a migration

A migration can be added by creating an up and down migration file in the migrations folder. The migration file names should be named according to {sequentialNumber}_{description}.up.sql and {sequentialNumber}_{description}.down.sql. Migrations can be created with the following command.

```
make migration name=<NAME>
```

Populate the up and down scripts. Up scripts should be idempotent. Down scripts should revert all changes made by up script

## How to run a migration

Migrations can be applied using the CLI

```
make migrate-up
```

## Querying Postgres

The following options are available if using Remote-Containers.

- Use the [PostgreSQL extension](https://marketplace.visualstudio.com/items?itemName=ckolkman.vscode-postgres) in VS Code
- Use pgAdmin on [localhost:8080](http://localhost:8080) with info@hackrva.org/test
- Use make run-sql or make run-sql-command

## Troubleshooting devContainer

Sometimes i'm getting an error when tryign to build the containers.

```
Command failed: docker inspect --type image memberdashboard_dev
```

This happened after I ran a `docker system prune -a`, but I think I've seen it before.

To fix this, build the containers with the following command then reopen vscode in Container

```
docker-compose -f docker-compose.yml -f .devcontainer/docker-compose.yml build
```
