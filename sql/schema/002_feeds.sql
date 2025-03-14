-- +goose Up
create table feeds (
    id uuid primary key,
    created_at timestamp not null,
    updated_at timestamp not null,
    name text not null,
    url  text unique not null
);

-- +goose Down
drop table feeds;