cd sql/migrations
goose postgres postgres://postgres:a@localhost:5432/gator down
goose postgres postgres://postgres:a@localhost:5432/gator up
cd ../..