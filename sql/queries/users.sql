-- name: CreateUser :one
insert into users (id, created_at, updated_at, name)
values ($1, $2, $3, $4)
returning *;

-- name: GetUser :one
select * from users
where users.name = $1;

-- name: RemoveAllUsers :exec
delete from users;

-- name: GetUsers :many
select * from users;