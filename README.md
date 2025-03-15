# aggregator

psql postgres://postgres:a@localhost:5432/gator

## build up

```bash
cd sql/schema
goose postgres postgres://postgres:a@localhost:5432/gator up
cd ../..
```

## tear down

```bash
cd sql/schema
goose postgres postgres://postgres:a@localhost:5432/gator down
cd ../..
```

# things installed

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

```bash
go get github.com/google/uuid
go get github.com/lib/pq
```

# troubleshooting

Edit: the solution is you need to capitilize Up and Down in the goose migration file, or it doesn't see it

```sql
-- +goose Up
...
-- +goose Down
```

I ran into a problem generating the sqlc files. this happened after I added a new migration sql file

```sql
-- +goose up
alter table feeds
add column last_fetched_at timestamp;

-- +goose down
alter table feeds
drop column last_fetched_at;

-- name: MarkFeedFetched :exec
update feeds
set last_fetched_at = $2,
    updated_at = $3
where id = $1;
```

when I ran goose, the migration worked and the table was updated. but when I ran sqlc, I got the following error

```bash
sqlc generate
# package
sql/queries/feeds.sql:47:5: column "last_fetched_at" does not exist
```

So to fix this, I dumbed the schema, moved the old migrations to a back directory and then ran sqlc again

```bash
pg_dump --schema-only --no-owner --file=sql/schema/schema.sql postgres://postgres:a@localhost:5432/gator
sqlc generate
```

just re-running the migration did not work.
