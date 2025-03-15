-- +goose Up
ALTER TABLE feeds ADD COLUMN titties timestamp;

-- +goose Down
alter table feeds drop column titties;