-- +goose up
ALTER TABLE feeds ADD COLUMN last_fetched_at timestamp;

-- +goose down
alter table feeds drop column last_fetched_at;