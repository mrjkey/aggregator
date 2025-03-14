-- +goose Up
create table feed_follows (
    created_at timestamp not null,
    updated_at timestamp not null,
    user_id uuid,
    constraint fk_user_id
        foreign key (user_id)
        references public.users(id)
        on delete cascade,
    feed_id uuid,
    constraint fk_feed_id
        foreign key (feed_id)
        references public.feeds(id) 
        on delete cascade,
    primary key (user_id, feed_id)
);

-- +goose Down
drop table feed_follows;