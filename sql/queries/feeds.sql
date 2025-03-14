-- name: AddFeed :one
insert into feeds (id, created_at, updated_at, name, url)
values ($1, $2, $3, $4, $5)
returning *;

-- -- name: GetFeeds :many
-- select feeds.name, feeds.url, users.name
-- from feeds
-- left join users on feeds.user_id = users.id;


-- name: GetFeedByUrl :one
select *
from feeds
where feeds.url = $1;

-- name: CreateFeedFollow :one
with inserted_feed_follow as (
    insert into feed_follows (created_at, updated_at, user_id, feed_id)
    values ($1, $2, $3, $4)
    returning *
)
select inserted_feed_follow.* , feeds.name as feed_name, users.name as user_name
from inserted_feed_follow
inner join users on inserted_feed_follow.user_id = users.id 
inner join feeds on inserted_feed_follow.feed_id = feeds.id;

-- name: GetFeedFollowersForUser :many
select feed_follows.*, feeds.name as feed_name, users.name as user_name
from feed_follows
inner join users on feed_follows.user_id = users.id
inner join feeds on feed_follows.feed_id = feeds.id
where feed_follows.user_id = $1;

-- name: RemoveAllFeeds :exec
delete from feeds;

-- name: DeleteFeedFollow :exec
delete from feed_follows
using feeds
where feed_follows.feed_id = feeds.id
and feed_follows.user_id = $1
and feeds.url = $2;