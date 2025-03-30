-- +goose Up
ALTER TABLE feeds ADD COLUMN last_fetched_at timestamp;

-- +goose Down
alter table feeds drop column last_fetched_at;