# aggregator

psql postgres://postgres:a@localhost:5432/gator

cd sql/schema

goose postgres postgres://postgres:a@localhost:5432/gator up
