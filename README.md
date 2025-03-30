# aggregator

## Purpose
The purpose of this project is to provide a simple configuration management tool for Go applications. It allows users to read and write configuration settings from a JSON file located in the user's home directory.

## Setup
1. Clone the repository:
   ```
   git clone https://github.com/mrjkey/aggregator.git
   ```
2. Navigate to the project directory:
   ```
   cd aggregator
   ```
3. Install the dependencies:
   ```
   go mod tidy
   ```

## Usage
1. Create a `.gatorconfig.json` file in your home directory with the following content:
   ```json
   {
     "db_url": "your_database_url",
     "current_user_name": "your_username"
   }
   ```
2. Run the application:
   ```
   go run main.go
   ```

## Dependencies
- Go 1.24.0 or higher

## Contact
For any questions or suggestions, please contact the project maintainer at [your_email@example.com].

## Contribution Guidelines
1. Fork the repository.
2. Create a new branch for your feature or bugfix.
3. Commit your changes and push the branch to your fork.
4. Create a pull request with a detailed description of your changes.
=======
psql postgres://postgres:a@localhost:5432/gator

## How to use

The following commands are available

```bash
go run . reset # reset the database, deletes users and feeds
go run . register <username> # register a new user
go run . login <username> # login as a user, must already be registered
go run . users # list all users, shows the current user as well
go run . addfeed <feed name> <feed url> # add an rss feed to the database, and the current user follows this
go run . follow <feed url> # follow a feed
go run . unfollow <feed url> # unfollow a feed
go run . following # list all feeds the current user is following
go run . agg <wait time interval> # gets posts from all the feeds of the current user. time interval is 1s, 1m, 1h, 1d etc.
go run . browse <option number> # browse the posts of the current user, can specify how many
```

## setup

need:

- go
- postgres
- goose
- sqlc

there is an example config file in the top level, but you need to setup the file in ~/.gatorconfig.json

```json
{
  "db_url": "postgres://<user>:<password>@localhost:5432/gator?sslmode=disable",
  "current_user_name": "dave"
}
```

the postgres url should be adjust to the user and password you use

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

### db setup

sudo apt install postgresql postgresql-contrib
psql --version
sudo passwd postgres
sudo -u postgres psql
CREATE DATABASE gator;

ALTER USER postgres PASSWORD 'a';

migrate

cd sql/schema
goose postgres postgres://postgres:a@localhost:5432/gator up

### example input

``` bash
go run . register jared
go run . addfeed "TechCrunch" "https://techcrunch.com/feed/"
go run . addfeed "Hacker News" "https://news.ycombinator.com/rss"
go run . agg 5s
go run . browse 5
```
