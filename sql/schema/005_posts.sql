-- +goose Up
create table posts (
    id uuid primary key, 
    created_at timestamp not null, 
    updated_at timestamp not null, 
    title text not null, 
    url text unique not null,
    description text,
    published_at timestamp not null, 
    feed_id uuid not null,
    constraint fk_feed_id 
        foreign key (feed_id)
        references public.feeds(id)
        on delete cascade
);

-- +goose Down
drop table posts;