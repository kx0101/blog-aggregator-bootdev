-- name: CreateFeedFollow :one
insert into feed_follows (id, feed_id, user_id, created_at, updated_at)
values ($1, $2, $3, $4, $5)
returning *;

-- name: GetFeedFollows :many
select *
from feed_follows;

-- name: GetFeedFollowsForUser :many
select *
from feed_follows
where user_id = $1;

-- name: GetFeedFollow :one
select *
from feed_follows
where id = $1;

-- name: DeleteFeedFollow :exec
delete from feed_follows where id = $1;
