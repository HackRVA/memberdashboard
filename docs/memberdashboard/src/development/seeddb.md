# Seed the DB with test data

> note: before seeding the db, it has to be up and running with the current schema
> run the db migrations

```
make migrate-up
```

seed the db

```
make seed
```

This will create some random members and a test user.

| username      | password |
| ------------- | -------- |
| test@test.com | test     |
