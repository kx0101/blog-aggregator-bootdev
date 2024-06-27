-- +goose Up
create table posts(
    id uuid primary key,
    created_at timestamp not null,
    updated_at timestamp not null,
    title varchar(255) not null,
    url varchar(255) not null,
    description varchar(255),
    published_at timestamp not null,
    feed_id uuid not null,
    foreign key(feed_id) references feeds(id)
);

-- +goose Down
drop table posts;
