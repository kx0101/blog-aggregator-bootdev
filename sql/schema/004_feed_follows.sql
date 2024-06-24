-- +goose Up
create table feed_follows(
    id UUID primary key,
    feed_id UUID unique not null,
    user_id UUID unique not null,
    created_at timestamp not null,
    updated_at timestamp not null,
    foreign key(feed_id) references feeds(id),
    foreign key(user_id) references users(id)
);

-- +goose Down
drop table feed_follows;
