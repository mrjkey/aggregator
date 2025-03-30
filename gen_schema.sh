pg_dump --schema-only --no-owner --file=sql/schema/schema.sql postgres://postgres:a@localhost:5432/gator
sqlc generate