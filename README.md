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
