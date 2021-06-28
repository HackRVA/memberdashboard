# Generators
Generators create mock data for testing purposes.  

## Members
Generate a test set of members with payment history using members generator by specifying the number of members to create. For example, the following will generate 200 members of various tiers with randomized payment histories.
```
go run members.go 200
```

## Purging test data
The easiest way to purge test data is to just rebuild the database.  A quick method is to completely migrate down, and then migrate back up. The following commands will perform these tasks if run from the root of the workspace.
```
make migrate-down
```
```
make migrate-up
```

