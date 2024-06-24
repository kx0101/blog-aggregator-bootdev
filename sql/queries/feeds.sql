-- name: CreateFeed :one
insert into feeds (id, created_at, updated_at, name, url, user_id)
values ($1, $2, $3, $4, $5, $6)
returning *;

-- name: UpdateLastFetchedAt :exec
update feeds
set last_fetched_at = now()
where id = ANY($1::uuid[]);

-- name: GetFeeds :many
select *
from feeds;

-- name: GetFeed :one
select *
from feeds
where id = $1;
