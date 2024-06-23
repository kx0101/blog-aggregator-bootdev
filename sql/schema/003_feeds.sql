-- +goose Up
create table feeds(
    id UUID primary key,
    created_at timestamp not null,
    updated_at timestamp not null,
    name varchar(255) not null,
    url varchar(255) unique not null,
    user_id UUID not null,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down
drop table feeds;
