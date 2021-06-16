## ----------------------------------------------------------------------
## The purpose of this Makefile is to simplify common development tasks.
## ----------------------------------------------------------------------
##

.PHONY:help
help: ## Show this help.
##
	@sed -ne '/@sed/!s/##//p' $(MAKEFILE_LIST)

.PHONY:migration
migration: ## Create a new migration
##   Create a new migration with name
##   Usage:
##     make migration name=<NAME>
##   Arguments:
##     <NAME>	The name for the migration
##   Example:
##     make migration name=communication_log
##
ifdef name
	migrate create -ext sql -dir migrations -seq $(name)
else
	@echo No name provided.  Provide a name to make a migration
	@echo Example:
	@echo "  make migration name=communication_log"
endif

.PHONY:migrate-up
migrate-up: ## Apply database migrations
##   Apply all or N up migrations
##   Usage:
##     make migrate-up [n=<N>]
##   Arguments:
##     <N>	The number of migrations to apply
##   Example:
##     make migrate-up
##    or
##     make migrate-up n=2
##
	migrate -database ${DB_CONNECTION_STRING}?sslmode=disable -verbose -path migrations up $(n)

.PHONY:migrate-down
migrate-down: ## Revert database migrations
##   Apply all or N down migrations
##   Usage:
##     make migrate-down [n=<N>]
##   Arguments:
##     <N>	The number of migrations to apply
##   Example:
##     make migrate-to version=2
##    or
##     make migrate-down n=2
##
	migrate -database ${DB_CONNECTION_STRING}?sslmode=disable -verbose -path migrations down $(n)

.PHONY:migrate-to
migrate-to: ## Migrate to the specified version
##   Migrate to the specified database version
##   Usage:
##     make migrate-to version=<VERSION>
##   Arguments:
##     <VERSION>	The version to migrate to
##   Example:
##     make migrate-to version=2
##
ifdef version
	migrate -database ${DB_CONNECTION_STRING}?sslmode=disable -verbose -path migrations goto $(version)
else
	@echo No version provided. Provide a version to migrate to
	@echo Example:
	@echo "  make migrate-to version=2"
endif

.PHONY:migrate-force
migrate-force: ## Force migrate to a database version
##   Force migrate to the specified version to fix a dirty database version.
##   Migrate to the specified database version
##   Usage:
##     make migrate-force version=<VERSION>
##   Arguments:
##     <VERSION>	The version to migrate to
##   Example:
##     make migrate-force version=2
##
ifdef version
	migrate -database ${DB_CONNECTION_STRING}?sslmode=disable -verbose -path migrations force $(version)
else
	@echo No version provided. Provide a version to force migrate
	@echo Example:
	@echo "  make migrate-force version=2"
endif

.PHONY:run-sql
run-sql: ## Run a sql script
##   Run a sql script file against postgres by specifying a script path
##   Usage:
##     make run-sql script=<SCRIPT>
##   Arguments:
##     <SCRIPT>	The path to the sql script to run
##   Example:
##     make run-sql script=test/postgres/myscript.sql
##
ifdef script
	psql -Atx $(DB_CONNECTION_STRING) $(script)
else
	@echo No script provided.  Provide a script to run sql
	@echo Example:
	@echo "  make run-sql script=test/postgres/myscript.sql"
endif

.PHONY:run-sql-command
run-sql-command: ## Run a sql command
##   Run a sql command against postgres by specifying a command
##    Usage:
##      make run-sql-command command="<COMMAND>"
##    Arguments:
##      <COMMAND>	The sql command to execute
##    Example:
##      make run-sql-command command="\copy membership.member_credit FROM 'test/postgres/seedData/member_credit.csv' DELIMITER ',' CSV HEADER;"
##
ifdef command
	psql -Atx $(DB_CONNECTION_STRING) -c "$(command)"
else
	@echo No command provided.  Provide a command to run
	@echo Example:
	@echo "  make run-sql-command command=\"\\\copy membership.member_credit FROM 'test/postgres/seedData/member_credit.csv' DELIMITER ',' CSV HEADER;\""
endif

.PHONY:swagger
swagger: ## Generate swagger doc
	swagger generate spec -o ./docs/swaggerui/swagger.json --scan-models
